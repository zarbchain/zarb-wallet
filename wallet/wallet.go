package wallet

import (
	_ "embed"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

var (
	/// ErrWalletExits describes an error in which there is a wallet
	/// exists in the given path
	ErrWalletExits = errors.New("wallet exists")

	/// ErrWalletExits describes an error in which the wallet CRC is
	/// invalid
	ErrInvalidCRC = errors.New("invalid CRC")

	/// ErrWalletExits describes an error in which the network is not
	/// valid
	ErrInvalidNetwork = errors.New("invalid network")
)

type Wallet struct {
	path   string
	store  *Store
	client *GrpcClient
}

type serverInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}
type servers = map[string][]serverInfo

//go:embed servers.json
var serversJSON []byte

/// OpenWallet generates an empty wallet and save the seed string
func OpenWallet(path string) (*Wallet, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := new(Store)
	err = json.Unmarshal(data, s)
	exitOnErr(err)

	if s.VaultCRC != s.calcVaultCRC() {
		exitOnErr(ErrInvalidCRC)
	}

	return newWallet(path, s, true)
}

/// Recover recovers a wallet from mnemonic (seed phrase)
func RecoverWallet(path, mnemonic string, net int) (*Wallet, error) {
	path = util.MakeAbs(path)
	if util.PathExists(path) {
		return nil, ErrWalletExits
	}
	s := RecoverStore(mnemonic, net)
	w, err := newWallet(path, s, false)
	if err != nil {
		return nil, err
	}

	err = w.saveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

/// CreateWallet generates an empty wallet and save the seed string
func CreateWallet(path, passphrase string, net int) (*Wallet, error) {
	path = util.MakeAbs(path)
	if util.PathExists(path) {
		return nil, ErrWalletExits
	}
	s := NewStore(passphrase, net)
	w, err := newWallet(path, s, false)
	if err != nil {
		return nil, err
	}

	err = w.saveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newWallet(path string, store *Store, online bool) (*Wallet, error) {
	w := &Wallet{
		store: store,
		path:  path,
	}

	err := w.connectToRandomServer()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wallet) connectToRandomServer() error {
	serversInfo := servers{}
	err := json.Unmarshal(serversJSON, &serversInfo)
	exitOnErr(err)

	var netServers []serverInfo
	switch w.store.Network {
	case 0:
		{ // mainnet
			netServers = serversInfo["mainnet"]
		}
	case 1:
		{ // testnet
			netServers = serversInfo["testnet"]
		}

	default:
		{
			return ErrInvalidNetwork
		}
	}

	for i := 0; i < 3; i++ {
		n := rand.Intn(len(netServers))
		serverInfo := netServers[n]
		client, err := MewGRPCClient(serverInfo.IP)
		if err == nil {
			w.client = client
			return nil
		}
	}

	return errors.New("unable to connect to the servers")
}
func (w *Wallet) Path() string {
	return w.path
}

func (w *Wallet) IsEncrypted() bool {
	return w.store.Encrypted
}

func (w *Wallet) saveToFile() error {
	w.store.VaultCRC = w.store.calcVaultCRC()

	bs, err := json.Marshal(w.store)
	exitOnErr(err)

	return util.WriteFile(w.path, bs)
}

func (w *Wallet) ImportPrivateKey(passphrase string, prv *bls.PrivateKey) error {
	err := w.store.ImportPrivateKey(passphrase, prv)
	if err != nil {
		return err
	}
	return w.saveToFile()
}

func (w *Wallet) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	return w.store.PrivateKey(passphrase, addr)
}

func (w *Wallet) Mnemonic(passphrase string) string {
	return w.store.Mnemonic(passphrase)
}

func (w *Wallet) Addresses() []crypto.Address {
	return w.store.Addresses()
}

func (w *Wallet) MakeSendTx(senderStr, receiverStr, amountStr string) (*tx.Tx, error) {
	sender, err := crypto.AddressFromString(senderStr)
	if err != nil {
		return nil, err
	}
	receiver, err := crypto.AddressFromString(senderStr)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return nil, err
	}

	stamp, err := w.client.GetStamp()
	if err != nil {
		return nil, err
	}

	seq, err := w.client.GetSequence(sender)
	if err != nil {
		return nil, err
	}

	// TODO
	fee := amount / 10000

	tx := tx.NewSendTx(stamp, seq, sender, receiver, amount, fee, "")

	return tx, nil
}

func (w *Wallet) BroadcastSendTx(tx *tx.Tx) (string, error) {
	b, err := tx.Bytes()
	if err != nil {
		return "", err
	}

	return w.client.SendTx(b)
}
