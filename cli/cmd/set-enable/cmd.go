package set_enable

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"math/big"
	"os"
)

type Params struct {
	KeyDirPath  string
	Passphrase  string
	Environment string
	Enable      bool
	BotID       string
}

func getID() (string, error) {
	b, err := os.ReadFile("./.settings/botId")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Run(ctx context.Context, cfg *Params) error {
	k, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		fmt.Println("Please run `./publish-cli initialize` first to generate a private key")
		return err
	}
	agentID := cfg.BotID
	if agentID == "" {
		agentID, err = getID()
		if err != nil {
			fmt.Println("bad state, ID not initialized in ./.settings/botId")
			return err
		}
	}

	var r registry.Client
	var chainID *big.Int
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

	agt, err := r.GetAgent(agentID)
	if err != nil {
		return nil
	}

	if agt == nil {
		fmt.Println(fmt.Sprintf("%s not found", agentID))
		return nil
	}

	opts, err := bind.NewKeyedTransactorWithChainID(k.PrivateKey, chainID)
	if err != nil {
		return err
	}

	if cfg.Enable {
		fmt.Println(fmt.Sprintf("enabling %s", agentID))
		res, err := r.Contracts().AgentRegTx.EnableAgent(
			opts,
			utils.AgentHexToBigInt(agentID),
			1,
		)
		if err != nil {
			return err
		}
		fmt.Println(res.Hash().Hex())
	} else {
		fmt.Println(fmt.Sprintf("disabling %s", agentID))
		res, err := r.Contracts().AgentRegTx.DisableAgent(
			opts,
			utils.AgentHexToBigInt(agentID),
			1,
		)
		if err != nil {
			return err
		}
		fmt.Println(res.Hash().Hex())
	}

	return nil
}
