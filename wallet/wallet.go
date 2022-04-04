package wallet

import (
	_ "embed"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
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

	/// ErrWalletExits describes an error in which the address doesn't
	/// exist in wallet
	ErrAddressNotFound = errors.New("address not found")

	/// ErrWalletExits describes an error in which the address already
	/// exist in wallet
	ErrAddressExists = errors.New("address already exists")
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

	bs, err := json.MarshalIndent(w.store, "  ", "  ")
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

func (w *Wallet) Addresses() []string {
	return w.store.Addresses()
}

/// MakeBondTx creates a new bond transaction based on the given parameters
func (w *Wallet) MakeBondTx(stampStr, seqStr, senderStr, valPubStr, stakeStr, memo string) (*tx.Tx, error) {
	sender, err := crypto.AddressFromString(senderStr)
	if err != nil {
		return nil, err
	}
	valPub, err := bls.PublicKeyFromString(valPubStr)
	if err != nil {
		return nil, err
	}
	stake, err := strconv.ParseInt(stakeStr, 10, 64)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsAccSeq(sender, seqStr)
	if err != nil {
		return nil, err
	}

	// TODO
	fee := stake / 10000

	tx := tx.NewBondTx(stamp, seq, sender, valPub, stake, fee, memo)
	return tx, nil
}

/// MakeUnbondTx creates a new unbond transaction based on the given parameters
func (w *Wallet) MakeUnbondTx(stampStr, seqStr, addrStr, memo string) (*tx.Tx, error) {
	addr, err := crypto.AddressFromString(addrStr)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsValSeq(addr, seqStr)
	if err != nil {
		return nil, err
	}

	tx := tx.NewUnbondTx(stamp, seq, addr, memo)
	return tx, nil
}

/// MakeWithdrawTx creates a new unbond transaction based on the given parameters
func (w *Wallet) MakeWithdrawTx(stampStr, seqStr, valAddrStr, accAddrStr, amountStr, memo string) (*tx.Tx, error) {
	valAddr, err := crypto.AddressFromString(valAddrStr)
	if err != nil {
		return nil, err
	}
	accAddr, err := crypto.AddressFromString(accAddrStr)
	if err != nil {
		return nil, err
	}
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsValSeq(valAddr, seqStr)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return nil, err
	}
	// TODO
	fee := amount / 10000
	tx := tx.NewWithdrawTx(stamp, seq, valAddr, accAddr, amount, fee, memo)
	return tx, nil
}

/// MakeSendTx creates a new send transaction based on the given parameters
func (w *Wallet) MakeSendTx(stampStr, seqStr, senderStr, receiverStr, amountStr, memo string) (*tx.Tx, error) {
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
	stamp, err := w.parsStamp(stampStr)
	if err != nil {
		return nil, err
	}
	seq, err := w.parsAccSeq(sender, seqStr)
	if err != nil {
		return nil, err
	}

	// TODO
	fee := amount / 10000

	tx := tx.NewSendTx(stamp, seq, sender, receiver, amount, fee, memo)
	return tx, nil
}

func (w *Wallet) parsAccSeq(signer crypto.Address, seqStr string) (int32, error) {
	if seqStr != "" {
		seq, err := strconv.ParseInt(seqStr, 10, 32)
		if err != nil {
			return -1, err
		}
		return int32(seq), nil
	}

	return w.client.GetAccountSequence(signer)
}

func (w *Wallet) parsValSeq(signer crypto.Address, seqStr string) (int32, error) {
	if seqStr != "" {
		seq, err := strconv.ParseInt(seqStr, 10, 32)
		if err != nil {
			return -1, err
		}
		return int32(seq), nil
	}

	return w.client.GetValidatorSequence(signer)
}

func (w *Wallet) parsStamp(stampStr string) (hash.Stamp, error) {
	if stampStr != "" {
		stamp, err := hash.StampFromString(stampStr)
		if err != nil {
			return hash.UndefHash.Stamp(), err
		}
		return stamp, nil
	}
	return w.client.GetStamp()
}

func (w *Wallet) SignAndBroadcast(passphrase string, tx *tx.Tx) (string, error) {
	prv, err := w.PrivateKey(passphrase, tx.Payload().Signer().String())
	if err != nil {
		return "", err
	}

	signer := crypto.NewSigner(prv)
	signer.SignMsg(tx)
	b, err := tx.Bytes()
	if err != nil {
		return "", err
	}

	return w.client.SendTx(b)
}
