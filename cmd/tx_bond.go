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
		senderOpt := c.String(cli.StringOpt{
			Name: "sender",
			Desc: "Sender account address",
		})

		pubOpt := c.String(cli.StringOpt{
			Name: "pub",
			Desc: "Validator public key",
		})

		stakeOpt := c.String(cli.StringOpt{
			Name: "stake",
			Desc: "Stake amount",
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

			// ---
			if *stakeOpt == "" {
				cmd.PrintWarnMsg("Stake is not defined.")
				c.PrintHelp()
				return
			}

			if *senderOpt == "" {
				cmd.PrintWarnMsg("Sender address is not defined.")
				c.PrintHelp()
				return
			}
			if *pubOpt == "" {
				cmd.PrintWarnMsg("Public key is not defined.")
				c.PrintHelp()
				return
			}

			trx, err := w.MakeBondTx(*stampOpt, *seqOpt, *senderOpt, *pubOpt, *stakeOpt, *memoOpt)
			if err != nil {
				PrintDangerMsg(err.Error())
				return
			}
			PrintInfoMsg("You are going to sign and broadcast a bond transition to the network.")
			PrintInfoMsg("Account: %s", *senderOpt)
			PrintInfoMsg("Validator: %s", trx.Payload().(*payload.BondPayload).PublicKey.Address())
			PrintInfoMsg("Stake: %s", *stakeOpt)

			signAndPublishTx(w, trx)
		}
	}
}
