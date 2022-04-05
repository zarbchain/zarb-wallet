package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

/// AllAddresses lists all the wallet addresses
func AllAddresses() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			addrs := w.Addresses()
			for addr, label := range addrs {
				PrintInfoMsg("%s %s", addr, label)
			}
		}
	}
}

/// NewAddress creates a new address
func NewAddress() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			label := PromptInput("Label: ")
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			addr, err := w.NewAddress(passphrase, label)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("%s", addr)
		}
	}
}

// GetPrivateKey returns the private key of an address
func GetPrivateKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addressArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "address string",
		})

		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			prv, err := w.PrivateKey(passphrase, *addressArg)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintWarnMsg("Private Key: \"%v\"", prv)
		}
	}
}

// GetPrivateKey returns the public key of an address
func GetPublicKey() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		addressArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "address string",
		})

		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			pub, err := w.PublicKey(passphrase, *addressArg)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("Public Key: \"%v\"", pub)
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
