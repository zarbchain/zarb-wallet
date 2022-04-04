package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func UnbondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		valOpt := c.String(cli.StringOpt{
			Name: "val",
			Desc: "Validator's address",
		})
		stampOpt, seqOpt, memoOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			// ---
			if *valOpt == "" {
				cmd.PrintWarnMsg("Sender address is not defined.")
				c.PrintHelp()
				return
			}

			trx, err := w.MakeUnbondTx(*stampOpt, *seqOpt, *valOpt, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("You are going to sign and broadcast an unbond transition to the network.")
			PrintInfoMsg("Validator: %s", *valOpt)

			signAndPublishTx(w, trx)

		}
	}
}
