package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/forta-network/forta-core-go/ipfs"
	manifest2 "github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"math/big"
)

type Params struct {
	KeyDirPath      string
	Passphrase      string
	IPFSGatewayPath string
	Environment     string
	Manifest        string
}

func Run(ctx context.Context, cfg *Params) error {
	ic, err := ipfs.NewClient(cfg.IPFSGatewayPath)
	if err != nil {
		return err
	}
	k, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		fmt.Println("Please run `./publish-cli initialize` first to generate a private key")
		return err
	}

	b, err := ic.GetBytes(ctx, cfg.Manifest)
	if err != nil {
		return err
	}
	var smf manifest2.SignedAgentManifest
	if err := json.Unmarshal(b, &smf); err != nil {
		return err
	}

	var r registry.Client
	chainID := big.NewInt(137)
	if cfg.Environment == "dev" {
		r, err = registry.NewClient(ctx, registry.ClientConfig{
			JsonRpcUrl: "https://rpc-mumbai.matic.today",
			ENSAddress: "0x5f7c5bbBa72e1e1fae689120D76D2f334A390Ae9",
			PrivateKey: k.PrivateKey,
		})
		if err != nil {
			return err
		}
		r.SetRegistryChainID(80001)
		chainID = big.NewInt(80001)
	} else {
		r, err = registry.NewClient(ctx, registry.ClientConfig{
			JsonRpcUrl: "https://polygon-rpc.com",
			ENSAddress: "0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7",
			PrivateKey: k.PrivateKey,
		})
		if err != nil {
			return err
		}
		r.SetRegistryChainID(137)
		chainID = big.NewInt(137)
	}

	agt, err := r.GetAgent(*smf.Manifest.AgentIDHash)
	if err != nil {
		return nil
	}

	opts, err := bind.NewKeyedTransactorWithChainID(k.PrivateKey, chainID)
	if err != nil {
		return err
	}

	var chainIDs []*big.Int
	for _, cID := range smf.Manifest.ChainIDs {
		chainIDs = append(chainIDs, big.NewInt(cID))
	}

	agentID := utils.AgentHexToBigInt(*smf.Manifest.AgentID)

	if agt == nil {
		res, err := r.Contracts().AgentRegTx.CreateAgent(
			opts,
			agentID,
			k.Address,
			cfg.Manifest,
			chainIDs,
			k.Address,
		)
		if err != nil {
			return err
		}
		fmt.Println(res.Hash().Hex())
	} else {
		res, err := r.Contracts().AgentRegTx.UpdateAgent(
			opts,
			agentID,
			cfg.Manifest,
			chainIDs,
		)
		if err != nil {
			return err
		}
		fmt.Println(res.Hash().Hex())
	}

	return nil
}
