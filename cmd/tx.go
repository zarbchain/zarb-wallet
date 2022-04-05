package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-wallet/wallet"
)

func SendTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "sender address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "receiver address",
		})

		amountArg := c.String(cli.StringArg{
			Name: "AMOUNT",
			Desc: "the amount to be transferred",
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

func BondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		senderArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "sender account address",
		})

		pubArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "validator public key",
		})

		stakeArg := c.String(cli.StringArg{
			Name: "STAKE",
			Desc: "stake amount",
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

func UnbondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		valArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "validator's address",
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

func WithdrawTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "withdraw from Validator address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "deposit to account address",
		})

		amountArg := c.String(cli.StringArg{
			Name: "AMOUNT",
			Desc: "the amount to be transferred",
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

func addCommonTxOptions(c *cli.Cmd) (*string, *string, *string, *string) {
	stampOpt := c.String(cli.StringOpt{
		Name: "stamp",
		Desc: "transaction stamp, if not specified will query from gRPC server",
	})

	seqOpt := c.String(cli.StringOpt{
		Name: "seq",
		Desc: "transaction sequence, if not specified will query from gRPC server",
	})
	memoOpt := c.String(cli.StringOpt{
		Name:  "memo",
		Desc:  "transaction memo, maximum should be 64 character (optional)",
		Value: "",
	})
	feeOpt := c.String(cli.StringOpt{
		Name:  "fee",
		Desc:  "transaction fee, if not specified will calculate automatically",
		Value: "",
	})

	return stampOpt, seqOpt, memoOpt, feeOpt
}

func signAndPublishTx(w *wallet.Wallet, trx *tx.Tx) {
	PrintWarnMsg("THIS ACTION IS NOT REVERSIBLE")
	confirmed := PromptConfirm("Do you want to continue? ")
	if !confirmed {
		return
	}

	passphrase := getPassphrase(w)
	res, err := w.SignAndBroadcast(passphrase, trx)
	if err != nil {
		PrintDangerMsg(err.Error())
		return
	}
	PrintInfoMsg(res)
}

func getPassphrase(w *wallet.Wallet) string {
	passphrase := ""
	if w.IsEncrypted() {
		passphrase = PromptPassphrase("Wallet password: ", false)
	}
	return passphrase
}
