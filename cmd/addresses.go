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
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			addrs := w.Addresses()
			for _, addr := range addrs {
				PrintInfoMsg("%s", addr)
			}
		}
	}
}
