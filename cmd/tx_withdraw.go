package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func WithdrawTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromOpt := c.String(cli.StringOpt{
			Name: "from",
			Desc: "withdraw from Validator address",
		})

		toOpt := c.String(cli.StringOpt{
			Name: "to",
			Desc: "Deposit to address",
		})

		amountOpt := c.String(cli.StringOpt{
			Name: "amount",
			Desc: "The amount to be transferred",
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
			if *amountOpt == "" {
				cmd.PrintWarnMsg("Stake is not defined.")
				c.PrintHelp()
				return
			}
			if *fromOpt == "" {
				cmd.PrintWarnMsg("Validator address is not defined.")
				c.PrintHelp()
				return
			}
			if *toOpt == "" {
				cmd.PrintWarnMsg("Account address is not defined.")
				c.PrintHelp()
				return
			}

			trx, err := w.MakeWithdrawTx(*stampOpt, *seqOpt, *fromOpt, *toOpt, *amountOpt, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("You are going to sign and broadcast a withdraw transition to the network.")
			PrintInfoMsg("Validator: %s", *fromOpt)
			PrintInfoMsg("Account: %s", *toOpt)
			PrintInfoMsg("Amount: %s", *amountOpt)

			signAndPublishTx(w, trx)
		}
	}
}
