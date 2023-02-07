package generate_id

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"os"
)

type Params struct{}

func Run(ctx context.Context, cfg *Params) error {
	_, err := os.Stat("./.settings/botId")
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll("./.settings", 0777); err != nil {
			return err
		}
		file, err := os.OpenFile("./.settings/botId", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		id := crypto.Keccak256Hash([]byte(uuid.Must(uuid.NewUUID()).String())).Hex()
		_, err = file.WriteString(id)
		return err
	}
	return nil
}
