package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

// GetPrivateKey returns the private key of an address
func GetPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addressArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "Address string",
		})

		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			passphrase := PromptPassphrase("Passphrase: ", false)
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			prv, err := w.PrivateKey(passphrase, *addressArg)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintDangerMsg("Private Key: \"%v\"", prv)
		}
	}
}

// GetPrivateKey returns the public key of an address
func GetPublicKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addressArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "Address string",
		})

		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			passphrase := PromptPassphrase("Passphrase: ", false)
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			pub, err := w.PublicKey(passphrase, *addressArg)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintDangerMsg("Public Key: \"%v\"", pub)
		}
	}
}

// ImportPrivateKey imports a private key into the wallet
func ImportPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			prv := PromptInput("Private Key: ")

			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			err = w.ImportPrivateKey(passphrase, prv)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintSuccessMsg("Private Key imported")
		}
	}
}
