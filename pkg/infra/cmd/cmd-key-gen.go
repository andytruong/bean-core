package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"bean/components/claim"
	"bean/pkg/infra"
)

func KeyGenCommand(_ *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "gen-key",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verify",
				Value: false,
				Usage: "Set to true to verify generated key pair.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("verify") {
				return verify()
			}

			return generate()
		},
	}
}

func verify() error {
	// public key
	{
		pub, err := claim.ParseRsaPublicKeyFromFile("resources/keys/id_rsa.pub")
		if nil != pub && err == nil {
			fmt.Println("Public key: 🆗 ")
		} else {
			fmt.Println("Public key: 🙅 ")
		}
	}

	// private key
	{
		// …
	}

	return nil
}

func generate() error {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return errors.Wrap(err, "failed to generate private key")
	} else if err := writePrivateKey(pk); nil != err {
		return errors.Wrap(err, "failed to write private key")
	} else if err := writePublicKey(pk.PublicKey); nil != err {
		return errors.Wrap(err, "failed to write publi key")
	}

	return nil
}

func writePrivateKey(pk *rsa.PrivateKey) error {
	file, err := os.Create("resources/keys/id_rsa")
	if nil != err {
		return err
	}

	defer file.Close()

	return pem.Encode(
		file,
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Headers: map[string]string{
				"Note": "For QA usage only, don't use on production",
			},
			Bytes: x509.MarshalPKCS1PrivateKey(pk),
		},
	)
}

func writePublicKey(key rsa.PublicKey) error {
	bytes, err := asn1.Marshal(key)
	if nil != err {
		return err
	}

	file, err := os.Create("resources/keys/id_rsa.pub")
	if nil != err {
		return err
	}

	defer file.Close()
	return pem.Encode(file, &pem.Block{Type: "PUBLIC KEY", Bytes: bytes})
}
