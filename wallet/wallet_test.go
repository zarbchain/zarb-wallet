package wallet

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

func TempDirPath() string {
	p, err := ioutil.TempDir("", "zarb*")
	if err != nil {
		panic(err)
	}
	return p
}

func TempFilePath() string {
	return path.Join(TempDirPath(), "file")
}

var tWallet *Wallet
var tPassphrase string

func setup(t *testing.T) {
	tPassphrase = ""
	w, err := CreateWallet(TempFilePath(), tPassphrase, 2) // 2 for testing
	assert.NoError(t, err)

	tWallet = w
}

func reopenWallet(t *testing.T) {
	w, err := OpenWallet(tWallet.path)
	assert.NoError(t, err)

	tWallet = w
}

func TestRecoverWallet(t *testing.T) {
	setup(t)

	mnemonic := tWallet.Mnemonic(tPassphrase)
	recovered, err := RecoverWallet(TempFilePath(), mnemonic, 2)
	assert.NoError(t, err)

	reopenWallet(t)
	assert.Equal(t, tWallet.store.ParentKey(tPassphrase), recovered.store.ParentKey(""))
}

func TestGetPrivateKey(t *testing.T) {
	setup(t)

	addrs := tWallet.Addresses()
	assert.NotEmpty(t, addrs)
	for _, addr := range addrs {
		prv, err := tWallet.PrivateKey(tPassphrase, addr.String())
		assert.NoError(t, err)
		assert.Equal(t, prv.PublicKey().Address().String(), addr.String())
	}
}

func TestInvalidAddress(t *testing.T) {
	setup(t)

	_, err := tWallet.PrivateKey(tPassphrase, crypto.GenerateTestAddress().String())
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	setup(t)

	_, prv1 := bls.GenerateTestKeyPair()
	assert.NoError(t, tWallet.ImportPrivateKey(tPassphrase, prv1))
	reopenWallet(t)

	assert.True(t, tWallet.store.Contains(prv1.PublicKey().Address()))
	prv2, err := tWallet.PrivateKey(tPassphrase, prv1.PublicKey().Address().String())
	assert.NoError(t, err)
	assert.Equal(t, prv1.Bytes(), prv2.Bytes())

	// Import again
	assert.Error(t, tWallet.ImportPrivateKey(tPassphrase, prv1))
}
