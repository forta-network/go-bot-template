package publish_metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forta-network/forta-core-go/ipfs"
	manifest2 "github.com/forta-network/forta-core-go/manifest"
	"github.com/tidwall/sjson"
	"os"
	"time"

	"github.com/forta-network/forta-core-go/security"
)

type Params struct {
	KeyDirPath      string
	Passphrase      string
	BotManifestPath string
	DocFilePath     string
	IPFSGatewayPath string
	Environment     string
	Image           string
	BotID           string
}

func compactJson(s string) string {
	dst := &bytes.Buffer{}
	if err := json.Compact(dst, []byte(s)); err != nil {
		panic(err)
	}
	return dst.String()
}

func getID() (string, error) {
	b, err := os.ReadFile("./.settings/botId")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func Run(ctx context.Context, cfg *Params) error {
	b, err := os.ReadFile(cfg.BotManifestPath)
	if err != nil {
		return err
	}
	var mf manifest2.AgentManifest
	if err := json.Unmarshal(b, &mf); err != nil {
		return err
	}
	if mf.Name == nil || len(*mf.Name) == 0 {
		return errors.New("set the .name value in the manifest.json (required)")
	}

	ic, err := ipfs.NewClient(cfg.IPFSGatewayPath)
	if err != nil {
		return err
	}
	db, err := os.ReadFile(cfg.DocFilePath)
	if err != nil {
		return err
	}

	docIpfsRef, err := ic.AddFile(db)
	if err != nil {
		return err
	}

	k, err := security.LoadKeyWithPassphrase(cfg.KeyDirPath, cfg.Passphrase)
	if err != nil {
		fmt.Println("Please run `./publish-cli initialize` first to generate a private key")
		return err
	}

	m, err := sjson.SetRaw(`{}`, "manifest", string(b))
	if err != nil {
		return err
	}

	botID := cfg.BotID
	if botID == "" {
		botID, err = getID()
		if err != nil {
			return err
		}
	}

	m, _ = sjson.Set(m, "manifest.imageReference", cfg.Image)
	m, _ = sjson.Set(m, "manifest.agentIdHash", botID)
	m, _ = sjson.Set(m, "manifest.agentId", botID)
	m, _ = sjson.Set(m, "manifest.from", k.Address.Hex())
	m, _ = sjson.Set(m, "manifest.timestamp", time.Now().UTC().Format(time.RFC3339))
	m, _ = sjson.Set(m, "manifest.documentation", docIpfsRef)
	m = compactJson(m)

	// sign manifest
	sig, err := security.SignString(k, m)
	if err != nil {
		return err
	}
	manifest, err := sjson.Set(m, "signature", sig.Signature)
	if err != nil {
		return err
	}

	manifestRef, err := ic.AddFile([]byte(manifest))
	if err != nil {
		return err
	}

	fmt.Println(manifestRef)

	return nil
}
