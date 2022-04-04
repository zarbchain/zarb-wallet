package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func WithdrawTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "withdraw from Validator address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "Deposit to account address",
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

			trx, err := w.MakeWithdrawTx(*stampOpt, *seqOpt, *fromArg, *toArg, *amountArg, *feeOpt, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("You are going to sign and broadcast a Withdraw transition to the network.")
			PrintInfoMsg("Validator: %s", *fromArg)
			PrintInfoMsg("Account: %s", *toArg)
			PrintInfoMsg("Amount: %s", *amountArg)

			signAndPublishTx(w, trx)
		}
	}
}
