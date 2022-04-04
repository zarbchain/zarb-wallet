package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func UnbondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		valArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "Validator's address",
		})
		stampOpt, seqOpt, memoOpt, _ := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			trx, err := w.MakeUnbondTx(*stampOpt, *seqOpt, *valArg, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("You are going to sign and broadcast an Unbond transition to the network:")
			PrintInfoMsg("Validator: %s", *valArg)

			signAndPublishTx(w, trx)

		}
	}
}
