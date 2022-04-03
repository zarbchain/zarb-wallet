package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

// Generate creates a new wallet
func Generate() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(ZARB) }
		c.Action = func() {
			passphrase := PromptPassphrase("Passphrase: ", true)

			fmt.Println()

			wallet, err := wallet.CreateWallet(*path, passphrase, 0)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			mnemonic := wallet.Mnemonic(passphrase)

			PrintSuccessMsg("mnemonic: \"%v\"", mnemonic)
		}
	}
}
