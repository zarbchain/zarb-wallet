package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func SendTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "Sender address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "Receiver address",
		})

		amountArg := c.String(cli.StringArg{
			Name: "AMOUNT",
			Desc: "The amount to be transferred",
		})
		stampOpt, seqOpt, memoOpt, feeOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			trx, err := w.MakeSendTx(*stampOpt, *seqOpt, *fromArg, *toArg, *amountArg, *feeOpt, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("You are going to sign and broadcast a Send transition to the network:")
			PrintInfoMsg("From: %s", *fromArg)
			PrintInfoMsg("To: %s", *toArg)
			PrintInfoMsg("Amount: %s", *amountArg)

			signAndPublishTx(w, trx)
		}
	}
}
