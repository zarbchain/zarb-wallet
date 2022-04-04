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
			Name: "SENDER",
			Desc: "Sender address",
		})

		toArg := c.String(cli.StringArg{
			Name: "RECEIVER",
			Desc: "Receiver address",
		})

		amountArg := c.String(cli.StringArg{
			Name: "AMOUNT",
			Desc: "The amount to be transferred",
		})
		stampOpt, seqOpt, memoOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			PrintLine()

			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}
			if *amountArg == "" {
				cmd.PrintWarnMsg("Stake is not defined.")
				c.PrintHelp()
				return
			}
			if *fromArg == "" {
				cmd.PrintWarnMsg("Sender address is not defined.")
				c.PrintHelp()
				return
			}
			if *toArg == "" {
				cmd.PrintWarnMsg("Receiver is not defined.")
				c.PrintHelp()
				return
			}

			trx, err := w.MakeSendTx(*stampOpt, *seqOpt, *fromArg, *toArg, *amountArg, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintInfoMsg("You are going to sign and broadcast a send transition to the network.")
			PrintInfoMsg("From: %s", *fromArg)
			PrintInfoMsg("To: %s", *toArg)
			PrintInfoMsg("Amount: %s", *amountArg)

			signAndPublishTx(w, trx)
		}
	}
}
