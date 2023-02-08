package main

import (
	"fmt"
	"forta-network/publish/cmd/publish"
	publish_metadata "forta-network/publish/cmd/publish-metadata"
	set_enable "forta-network/publish/cmd/set-enable"
	"os"

	generate_id "forta-network/publish/cmd/generate-id"
	"forta-network/publish/cmd/initialize"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const defaultIpfsUrl = "https://ipfs.forta.network"

func deployKeyPath() string {
	uhd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s/.deploy-keys", uhd, ".forta")
}

func main() {
	app := &cli.App{
		Name:  "forta-publish",
		Usage: "a cli for publishing a bot to the network",
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "keydir",
						Value: deployKeyPath(),
					},
					&cli.StringFlag{
						Name: "passphrase",
					},
				},
				Action: func(c *cli.Context) error {
					kd := c.String("keydir")
					if kd == "" {
						uhd, err := os.UserHomeDir()
						if err != nil {
							return err
						}
						kd = fmt.Sprintf("%s/%s/.deploy-keys", uhd, ".forta")
					}

					return initialize.Run(c.Context, &initialize.Params{
						KeyDirPath: c.String("keydir"),
						Passphrase: c.String("passphrase"),
					})
				},
			},
			{
				Name: "publish-metadata",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "prod",
					},
					&cli.StringFlag{
						Name:  "manifest",
						Value: "manifest.json",
					},
					&cli.StringFlag{
						Name:     "image",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "doc-file",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "keydir",
						Value: deployKeyPath(),
					},
					&cli.StringFlag{
						Name:  "ipfs-gateway",
						Value: defaultIpfsUrl,
					},
					&cli.StringFlag{
						Name: "passphrase",
					},
					&cli.StringFlag{
						Name: "bot-id",
					},
				},
				Action: func(c *cli.Context) error {
					return publish_metadata.Run(c.Context, &publish_metadata.Params{
						Environment:     c.String("env"),
						KeyDirPath:      c.String("keydir"),
						Passphrase:      c.String("passphrase"),
						BotManifestPath: c.String("manifest"),
						IPFSGatewayPath: c.String("ipfs-gateway"),
						DocFilePath:     c.String("doc-file"),
						Image:           c.String("image"),
						BotID:           c.String("bot-id"),
					})
				},
			},
			{
				Name: "publish",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "prod",
					},
					&cli.StringFlag{
						Name:     "manifest",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "keydir",
						Value: deployKeyPath(),
					},
					&cli.StringFlag{
						Name:  "ipfs-gateway",
						Value: defaultIpfsUrl,
					},
					&cli.StringFlag{
						Name: "passphrase",
					},
				},
				Action: func(c *cli.Context) error {
					return publish.Run(c.Context, &publish.Params{
						Environment:     c.String("env"),
						KeyDirPath:      c.String("keydir"),
						Passphrase:      c.String("passphrase"),
						Manifest:        c.String("manifest"),
						IPFSGatewayPath: c.String("ipfs-gateway"),
					})
				},
			},
			{
				Name: "enable",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "prod",
					},
					&cli.StringFlag{
						Name:  "keydir",
						Value: deployKeyPath(),
					},
					&cli.StringFlag{
						Name: "passphrase",
					},
					&cli.StringFlag{
						Name: "bot-id",
					},
				},
				Action: func(c *cli.Context) error {
					return set_enable.Run(c.Context, &set_enable.Params{
						Environment: c.String("env"),
						KeyDirPath:  c.String("keydir"),
						Passphrase:  c.String("passphrase"),
						BotID:       c.String("bot-id"),
						Enable:      true,
					})
				},
			},
			{
				Name: "disable",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "env",
						Value: "prod",
					},
					&cli.StringFlag{
						Name:  "keydir",
						Value: deployKeyPath(),
					},
					&cli.StringFlag{
						Name: "passphrase",
					},
				},
				Action: func(c *cli.Context) error {
					return set_enable.Run(c.Context, &set_enable.Params{
						Environment: c.String("env"),
						KeyDirPath:  c.String("keydir"),
						Passphrase:  c.String("passphrase"),
						Enable:      false,
					})
				},
			},
			{
				Name:  "generate-id",
				Flags: []cli.Flag{},
				Action: func(c *cli.Context) error {
					return generate_id.Run(c.Context, &generate_id.Params{})
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
