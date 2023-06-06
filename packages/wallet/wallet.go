package wallet

import (
	"encoding/hex"
	"errors"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	as "github.com/IBAX-io/go-ibax-sdk/packages/pkg/accounts"
	hd "github.com/IBAX-io/go-ibax-sdk/packages/pkg/accounts/hdwallet"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto/hashalgo"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/tyler-smith/go-bip39"
)

type walletClient struct{}

func New(base modus.Base) modus.Wallet {
	return &walletClient{}
}

// NewMnemonic
// The number of length words should be 12, 15, 18, 21 or 24
func (p *walletClient) NewMnemonic(w request.WordsLenType) (mnemonic string, err error) {
	bitSize := w.GetBitSize()
	if bitSize == 0 {
		return "", errors.New("the number length of words should be 12, 15, 18, 21 or 24")
	}
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func (p *walletClient) NewWallet(length request.WordsLenType) (*hd.Wallet, error) {
	var (
		wallet   = &hd.Wallet{}
		err      error
		mnemonic string
	)
	mnemonic, err = p.NewMnemonic(length)
	if err != nil {
		return wallet, err
	}
	wallet, err = p.NewWalletFromMnemonic(mnemonic)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}

func (p *walletClient) NewWalletFromMnemonic(mnemonic string) (*hd.Wallet, error) {
	var (
		wallet = &hd.Wallet{}
		err    error
	)
	wallet, err = hd.NewWalletFromMnemonic(mnemonic, "")
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}

// NewAccountFromPath
// deriving a new account at the specific derivation path. If pin is set to true, the account will be added to the list of tracked accounts.
func (p *walletClient) NewAccountFromPath(wallet *hd.Wallet, path string, pin bool) (as.Account, error) {
	divPath, err := hd.ParseDerivationPath(path)
	if err != nil {
		return as.Account{}, err
	}

	//new account by Derivation path
	account, err := wallet.Derive(divPath, pin)
	if err != nil {
		return as.Account{}, err
	}
	return account, nil
}

func (p *walletClient) GetPrivateKey(wallet *hd.Wallet, account as.Account) ([]byte, error) {
	if account.URL.Path == "" {
		return nil, errors.New("empty derivation path")
	}
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func (p *walletClient) GetPublicKey(wallet *hd.Wallet, account as.Account) ([]byte, error) {
	if account.URL.Path == "" {
		return nil, errors.New("empty derivation path")
	}
	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func (p *walletClient) FormatAddress(account as.Account) string {
	return string(account.Address)
}

func (p *walletClient) GetAddress(publicKey []byte) string {
	return crypto.KeyToAddress(publicKey)
}

func (p *walletClient) GetKeyId(publicKey []byte) int64 {
	return crypto.Address(publicKey)
}

func (p *walletClient) GetETHAddress(publicKey []byte) string {
	if len(publicKey) == 0 {
		return ""
	}
	keccak := &hashalgo.Keccak256{}
	hash256 := keccak.GetHash(crypto.CutPub(publicKey))
	return "0x" + hex.EncodeToString(hash256[len(hash256)-20:])
}
