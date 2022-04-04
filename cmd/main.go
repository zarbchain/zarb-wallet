package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-wallet/wallet"
)

var path *string

func main() {
	app := cli.App("zarb-wallet", "Zarb wallet")

	path = app.String(cli.StringOpt{
		Name:  "w wallet file",
		Desc:  "A path to the wallet file",
		Value: ZarbWalletsDir() + "default_wallet",
	})

	app.Command("create", "Create a new wallet", Generate())
	app.Command("recover", "Recover waller from mnemonic (seed phrase)", Recover())
	app.Command("list_addresses", "List of wallet addresses", Addresses())
	app.Command("get_privkey", "Get private key of an address", GetPrivateKey())
	app.Command("import_privkey", "Import a private key into wallet", ImportPrivateKey())
	app.Command("tx", "Create, sign and publish a transaction", func(k *cli.Cmd) {
		k.Command("bond", "Create, sign and publish a bond transaction", BondTx())
		k.Command("send", "Create, sign and publish a send transaction", SendTx())
		k.Command("unbond", "Create, sign and publish an unbond transaction", UnbondTx())
		k.Command("withdraw", "Create, sign and publish a withdraw transaction", WithdrawTx())
	})

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func addCommonTxOptions(c *cli.Cmd) (*string, *string, *string, *string) {
	stampOpt := c.String(cli.StringOpt{
		Name: "stamp",
		Desc: "Transaction stamp, if not specified will query from gRPC server",
	})

	seqOpt := c.String(cli.StringOpt{
		Name: "seq",
		Desc: "Transaction sequence, if not specified will query from gRPC server",
	})
	memoOpt := c.String(cli.StringOpt{
		Name:  "memo",
		Desc:  "Transaction memo, maximum should be 64 character (optional)",
		Value: "",
	})
	feeOpt := c.String(cli.StringOpt{
		Name:  "fee",
		Desc:  "Transaction fee, if not specified will calculate automatically",
		Value: "",
	})

	return stampOpt, seqOpt, memoOpt, feeOpt
}

func signAndPublishTx(w *wallet.Wallet, trx *tx.Tx) {
	PrintWarnMsg("THIS ACTION IS NOT REVERSABLE")
	PromptConfirm("Do you want to continue? ")

	passphrase := ""
	if w.IsEncrypted() {
		passphrase = PromptPassphrase("Wallet password: ", false)
	}

	res, err := w.SignAndBroadcast(passphrase, trx)
	if err != nil {
		PrintDangerMsg("An error occurred: %s", err.Error())
	}
	PrintInfoMsg(res)
}
