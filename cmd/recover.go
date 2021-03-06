package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

/// Recover recovers a wallet from mnemonic (seed phrase)
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			mnemonic := PromptInput("Seed: ")
			w, err := wallet.RecoverWallet(*path, mnemonic, 0)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("Wallet recovered successfully at: %s", w.Path())
			PrintWarnMsg("Never share your private key.")
		}
	}
}

/// GetSeed prints the seed phrase (mnemonics)
func GetSeed() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(header) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			passphrase := getPassphrase(w)
			mnemonic, err := w.Mnemonic(passphrase)

			PrintLine()
			PrintInfoMsg("Seed: \"%v\"", mnemonic)
		}
	}
}
