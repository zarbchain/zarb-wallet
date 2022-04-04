package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func BondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		senderArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "Sender account address",
		})

		pubArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "Validator public key",
		})

		stakeArg := c.String(cli.StringArg{
			Name: "STAKE",
			Desc: "Stake amount",
		})
		stampOpt, seqOpt, memoOpt, feeOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			trx, err := w.MakeBondTx(*stampOpt, *seqOpt, *senderArg, *pubArg, *stakeArg, *feeOpt, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}

			PrintLine()
			PrintInfoMsg("You are going to sign and broadcast a bond transition to the network.")
			PrintInfoMsg("Account: %s", *senderArg)
			PrintInfoMsg("Validator: %s", trx.Payload().(*payload.BondPayload).PublicKey.Address())
			PrintInfoMsg("Stake: %s", *stakeArg)

			signAndPublishTx(w, trx)
		}
	}
}
