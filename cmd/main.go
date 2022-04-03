package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
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

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
