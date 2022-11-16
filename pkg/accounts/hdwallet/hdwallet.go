/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package hdwallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	as "github.com/IBAX-io/go-ibax-sdk/pkg/accounts"
	"net/url"
	"strconv"
	"sync"

	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto/asymalgo"
	"github.com/IBAX-io/go-ibax-sdk/pkg/smart"
	"github.com/IBAX-io/go-ibax-sdk/pkg/utils"

	// "github.com/IBAX-io/go-ibax/packages/smart"
	// "github.com/IBAX-io/go-ibax/packages/utils"

	// "github.com/IBAX-io/go-ibax/packages/smart"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/shopspring/decimal"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	mnemonic  string
	masterKey *hdkeychain.ExtendedKey
	seed      []byte
	url       as.URL
	paths     map[as.Address]as.DerivationPath
	accounts  []as.Account
	stateLock sync.RWMutex
}

// BitSizeOfEntropy can be {128, 256}.
const BitSizeOfEntropy = 128

/* -------------------New Seed and Mnemonic functions--------------- */

// NewMnemonic returns a randomly generated BIP-39 mnemonic.
func NewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(BitSizeOfEntropy)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// NewSeedFromMnemonic returns a BIP-39 seed based on a BIP-39 mnemonic and password.
// DK = PBKDF2(PRF, Password, Salt, c, dkLen) and HMAC-SHA512
func NewSeedFromMnemonic(mnemonic string, password string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}

// NewWalletFromMnemonic returns a new wallet from a BIP-39 mnemonic.
func NewWalletFromMnemonic(mnemonic string, password string) (*Wallet, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("mnemonic is invalid")
	}

	seed, err := NewSeedFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, err
	}

	wallet, err := newWallet(seed)
	if err != nil {
		return nil, err
	}
	wallet.mnemonic = mnemonic

	return wallet, nil
}

func NewWallet(password string) (*Wallet, error) {
	newMnemonic, err := NewMnemonic()
	if err != nil {
		return nil, err
	}
	return NewWalletFromMnemonic(newMnemonic, password)
}

// NewWalletFromSeed returns a new wallet from a BIP-39 seed.
func NewWalletFromSeed(seed []byte) (*Wallet, error) {
	if len(seed) == 0 {
		return nil, errors.New("seed is required")
	}
	return newWallet(seed)
}

// newWallet is a helper function to create a new wallet instance.
func newWallet(seed []byte) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		masterKey: masterKey,
		seed:      seed,
		accounts:  []as.Account{},
		paths:     map[as.Address]as.DerivationPath{},
	}, nil
}

/* -------------------Derive functions--------------- */

// Derive implements accounts.Wallet, deriving a new account at the specific
// derivation path. If pin is set to true, the account will be added to the list
// of tracked accounts.
func (w *Wallet) Derive(path as.DerivationPath, pin bool) (as.Account, error) {
	// Try to derive the actual account and update its URL if successful
	w.stateLock.RLock() // Avoid device disappearing during derivation
	address, err := w.deriveAddress(path)
	w.stateLock.RUnlock()

	// If an error occurred or no pinning was requested, return
	if err != nil {
		return as.Account{}, err
	}

	account := as.Account{
		Address: address,
		URL: as.URL{
			Scheme: "",
			Path:   path.String(),
		},
	}
	if !pin {
		return account, nil
	}

	// Pinning needs to modify the state
	w.stateLock.Lock()
	defer w.stateLock.Unlock()

	if _, ok := w.paths[address]; !ok {
		w.accounts = append(w.accounts, account)
		w.paths[address] = path
	}

	return account, nil
}

// DerivePrivateKey derives the private key of the derivation path.
func (w *Wallet) derivePrivateKey(path as.DerivationPath) ([]byte, error) {
	var err error
	key := w.masterKey
	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}
	// fmt.Println("path: ", path, "privateKey: ", privateKeyECDSA.D.Bytes())
	return privateKeyECDSA.D.Bytes(), nil
}

// DerivePublicKey derives the public key of the derivation path.
func (w *Wallet) derivePublicKey(path as.DerivationPath) ([]byte, error) {
	var err error
	key := w.masterKey

	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}
	publicKey := append(asymalgo.FillLeft(privateKeyECDSA.X.Bytes()), asymalgo.FillLeft(privateKeyECDSA.Y.Bytes())...)
	return append([]byte{04}, publicKey...), nil
}

// DeriveAddress derives the account address of the derivation path.
func (w *Wallet) deriveAddress(path as.DerivationPath) (as.Address, error) {
	publicKey, err := w.derivePublicKey(path)
	if err != nil {
		return "", err
	}

	address := crypto.KeyToAddress(publicKey)

	return as.Address(address), nil
}

// ! @todo [ ] ParseDerivationPath parses the derivation path in string format into []uint32.
func ParseDerivationPath(path string) (as.DerivationPath, error) {
	return as.ParseDerivationPath(path)
}

// MustParseDerivationPath parses the derivation path in string format into
// []uint32 but will panic if it can't parse it.
func MustParseDerivationPath(path string) as.DerivationPath {
	parsed, err := as.ParseDerivationPath(path)
	if err != nil {
		panic(err)
	}

	return parsed
}

/* -------------------Account Malipulate functions--------------- */

// Accounts implements accounts.Wallet, returning the list of accounts pinned to
// the wallet. If self-derivation was enabled, the account list is
// periodically expanded based on current chain state.
func (w *Wallet) Accounts() []as.Account {
	// Attempt self-derivation if it's running
	// Return whatever account list we ended up with
	// @todo [X] why Rlock : many readers and one writer
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	cpy := make([]as.Account, len(w.accounts))
	copy(cpy, w.accounts)
	return cpy
}

// Contains implements accounts.Wallet, returning whether a particular account is
// or is not pinned into this wallet instance.
func (w *Wallet) Contains(account as.Account) bool {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	_, exists := w.paths[account.Address]
	return exists
}

// Unpin unpins account from list of pinned accounts.
func (w *Wallet) Unpin(account as.Account) error {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	for i, acct := range w.accounts {
		// @todo [ ] if acct.Address.String() == account.Address.String()
		if acct.Address == account.Address {
			w.accounts = removeAtIndex(w.accounts, i)
			delete(w.paths, account.Address)
			return nil
		}
	}

	return errors.New("account not found")
}

// Path return the derivation path of the account.
func (w *Wallet) Path(account as.Account) (string, error) {
	return account.URL.Path, nil
}

// Address returns the address of the account.
func (w *Wallet) Address(account as.Account) (as.Address, error) {
	return account.Address, nil
}

// PrivateKey returns the ECDSA private key of the account.
func (w *Wallet) PrivateKey(account as.Account) ([]byte, error) {
	path, err := ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}
	return w.derivePrivateKey(path)
}

// PublicKey returns the ECDSA public key of the account.
func (w *Wallet) PublicKey(account as.Account) ([]byte, error) {
	path, err := ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}

	return w.derivePublicKey(path)
}

// @todo [ ] do nothing for now.
func (w *Wallet) URL() as.URL {
	return w.url
}
func (w *Wallet) Open(password string) error {
	return nil
}
func (w *Wallet) Status() (string, error) {
	return "ok", nil
}
func (w *Wallet) Close() error {
	return nil
}

// SelfDerive implements accounts.Wallet, trying to discover accounts that the
// user used previously (based on the chain state), but ones that he/she did not
// explicitly pin to the wallet manually. To avoid chain head monitoring, self
// derivation only runs during account listing (and even then throttled).
//@todo [ ]
// func (w *Wallet) SelfDerive(base []accounts.DerivationPath, chain ) {
// 	// TODO: self derivation
// }

/* -------------------Sign functions--------------- */

// SignHash implements accounts.Wallet, which allows signing arbitrary data.
func (w *Wallet) SignHash(account as.Account, hash []byte) ([]byte, error) {
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, errors.New("account not found")
	}

	privateKey, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(privateKey, hash)

}

func (wallet *Wallet) MultiSignForm(
	Data string,
	Quorum int,
	Threshold int) (url.Values, error) {
	/* MultiSign.sim
	   data{
	       Signatures array //H0, H1
	       Pubkeys array
	       Hash string
	       Data string //signdata, amount, toid
	       Quorum int
	       Threshold int
	   }
	*/

	// Data = []string{"123456"}
	// Quorum = 5
	// Thresshold = 2
	var Signatures []string
	var Pubkeys []string

	// data, err := json.Marshal(Data)
	// 	if err != nil {
	// 	return nil, err
	// }
	data := Data

	// @todo [ ] use private keys to sign data
	for i := 0; i < 5; i++ {
		// privateKey, err := wallet.PrivateKey(wallet.accounts[i])
		// if err != nil {
		// 	return nil, err
		// }
		// hash := crypto.Hash([]byte(data))
		signature, err := wallet.SignHash(wallet.accounts[i], []byte(data))
		if err != nil {
			return nil, err
		}
		publicKey, err := wallet.PublicKey(wallet.accounts[i])
		if err != nil {
			return nil, err
		} else {
			// ChekSign offline first
			_, err = crypto.Verify(publicKey, []byte(data), signature)
			if err != nil {
				fmt.Println("Multi CheckSign error（crypto.Verify）:", err)
				return nil, err
			}

			_, err = utils.CheckSign([][]byte{publicKey}, []byte(data), signature, true)
			if err != nil {
				fmt.Println("Multi CheckSign error（utils.CheckSign）:", err)
				return nil, err
			}

			Pubkeys = append(Pubkeys, hex.EncodeToString(publicKey))
			Signatures = append(Signatures, hex.EncodeToString(signature))

			// fmt.Println("crypto.PubToHex(publicKey)")
			// fmt.Println(crypto.PubToHex(publicKey))
			// fmt.Println(publicKey)
			_, err = smart.CheckSign(crypto.PubToHex(publicKey), data, Signatures[i])
			if err != nil {
				fmt.Println("Multi CheckSign error(smart.CheckSign):", err)
				return nil, err
			}
		}

	}
	sigs, err := json.Marshal(Signatures)
	if err != nil {
		return nil, err
	}
	pubs, err := json.Marshal(Pubkeys)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 5; i++ {
	}

	// @todo [ ] construct Signatures, Data, Pubkeys, Hash, Quorum, Threshold
	multiSignForm := url.Values{"Signatures": {string(sigs)}, "Pubkeys": {string(pubs)}, "Data": {string(data)}, "Quorum": {strconv.Itoa(Quorum)}, "Threshold": {strconv.Itoa(Threshold)}}
	return multiSignForm, nil

}

func (wallet *Wallet) CreateWalletForm(
	Owners string,
	Threshold int) ([]byte, error) {
	/* contract CreateWallet {
	   data {
	       Owners array
	       Threshold int
	   }
	   } */

	// Data = []string{"123456"}
	// Quorum = 5
	// Thresshold = 2
	return nil, nil
}

func (wallet *Wallet) CreateProposeForm(WalletToPropose string, Proposal string, Deadline int, Postscript string) {
}

// @todo [ ] MultiSend implements accounts.Wallet, which allows signing arbitrary data.

func (w *Wallet) SendTokens(recipient string, amount decimal.Decimal) error {
	accounts := w.Accounts()
	var totalBalance decimal.Decimal
	for _, account := range accounts {
		// keyLogin
		// getBal
		// MultiSign
		fmt.Println(totalBalance, account)
	}

	return nil
}

// @todo [ ] SignHash signs the given hash with the requested account.
func (w *Wallet) SignData(account as.Account, mimeType string, data []byte) ([]byte, error) {
	return nil, nil
}

func (w *Wallet) SignDataWithPassword(account as.Account, password, mimeType string, data []byte) ([]byte, error) {
	return nil, nil
}

func (w *Wallet) SignText(account as.Account, text []byte) ([]byte, error) {
	return nil, nil
}

func (w *Wallet) SignTextWithPassword(account as.Account, password string, text []byte) ([]byte, error) {
	return nil, nil
}

// removeAtIndex removes an account at index.
// @todo [X]
func removeAtIndex(accts []as.Account, index int) []as.Account {
	return append(accts[:index], accts[index+1:]...)
}
