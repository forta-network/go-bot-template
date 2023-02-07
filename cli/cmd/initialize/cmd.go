package initialize

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"os"
)

type Params struct {
	KeyDirPath string
	Passphrase string
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Run(ctx context.Context, cfg *Params) error {
	if err := os.MkdirAll(cfg.KeyDirPath, 0777); err != nil {
		return err
	}
	files, err := os.ReadDir(cfg.KeyDirPath)
	if err != nil {
		return err
	}
	if len(files) > 0 {
		fmt.Println("directory not empty, not creating key")
		return nil
	}

	ks := keystore.NewKeyStore(cfg.KeyDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	acct, err := ks.NewAccount(cfg.Passphrase)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Address Initialized: %s", acct.Address.Hex()))
	return nil
}
