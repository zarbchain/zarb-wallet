package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-wallet/wallet"
)

/// Addresses lists the wallet addresses
func Addresses() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			addrs := wallet.Addresses()
			for _, addr := range addrs {
				PrintInfoMsg("%s", addr.String())
			}

		}
	}
}
