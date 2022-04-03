package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

/// Recover recovers a wallet from mnemonic (seed phrase)
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(ZARB) }
		c.Action = func() {
			mnemonic := PromptInput("Mnemonic: ")

			fmt.Println()

			_, err := wallet.RecoverWallet(*path, mnemonic, 0)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

		}
	}
}
