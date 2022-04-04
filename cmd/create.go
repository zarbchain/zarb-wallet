package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

// Generate creates a new wallet
func Generate() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			passphrase := PromptPassphrase("Passphrase: ", true)
			w, err := wallet.CreateWallet(*path, passphrase, 0)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			mnemonic := w.Mnemonic(passphrase)

			PrintLine()
			PrintSuccessMsg("Wallet created successfully at: %s", w.Path())
			PrintInfoMsg("Seed: \"%v\"", mnemonic)
			PrintWarnMsg("Please keep your seed in a safe place; if you lose it, you will not be able to restore your wallet.")
		}
	}
}
