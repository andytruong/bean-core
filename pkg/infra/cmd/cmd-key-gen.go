package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"

	"github.com/urfave/cli/v2"

	"bean/pkg/infra"
)

func KeyGenCommand(container *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "gen-key",
		Action: func(ctx *cli.Context) error {
			// Generate RSA Keys
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return err
			} else {
				fmt.Println("key.private: ", privateKey)
				fmt.Println("key.public ", privateKey.Public())

				pem.Encode(nil, privateKey)
			}

			return nil
		},
	}
}
