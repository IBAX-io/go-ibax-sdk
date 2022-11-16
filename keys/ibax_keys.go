package keys

import (
	"encoding/hex"
	"errors"
	as "github.com/IBAX-io/go-ibax-sdk/pkg/accounts"
	hd "github.com/IBAX-io/go-ibax-sdk/pkg/accounts/hdwallet"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto/hashalgo"
	"github.com/tyler-smith/go-bip39"
)

type wordsLenType int

const (
	defaultBitSizeOfEntropy = 128
	interval                = 32

	WordsLenTwelve    wordsLenType = 12
	WordsLenFifteen   wordsLenType = 15
	WordsLenEighteen  wordsLenType = 18
	WordsTwentyOne    wordsLenType = 21
	WordsLenTwentyTwo wordsLenType = 24
)

// NewMnemonic
// The number of length words should be 12, 15, 18, 21 or 24
func (p *accountClient) NewMnemonic(length wordsLenType) (mnemonic string, err error) {
	bitSize := defaultBitSizeOfEntropy
	var increment int
	switch length {
	case WordsLenTwelve:
	case WordsLenFifteen:
		increment = interval
	case WordsLenEighteen:
		increment = interval * 2
	case WordsTwentyOne:
		increment = interval * 3
	case WordsLenTwentyTwo:
		increment = interval * 4
	default:
		return "", errors.New("the number length of words should be 12, 15, 18, 21 or 24")
	}
	bitSize += increment
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func (p *accountClient) NewWallet(length wordsLenType) (*hd.Wallet, error) {
	mnemonic, err := p.NewMnemonic(length)
	if err != nil {
		return nil, err
	}
	wallet, err := p.NewWalletFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (p *accountClient) NewWalletFromMnemonic(mnemonic string) (*hd.Wallet, error) {
	wallet, err := hd.NewWalletFromMnemonic(mnemonic, "")
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// NewAccountFromPath
// deriving a new account at the specific derivation path. If pin is set to true, the account will be added to the list of tracked accounts.
func (p *accountClient) NewAccountFromPath(wallet *hd.Wallet, path string, pin bool) (as.Account, error) {
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

func (p *accountClient) GetPrivateKey(wallet *hd.Wallet, account as.Account) ([]byte, error) {
	if account.URL.Path == "" {
		return nil, errors.New("empty derivation path")
	}
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func (p *accountClient) GetPublicKey(wallet *hd.Wallet, account as.Account) ([]byte, error) {
	if account.URL.Path == "" {
		return nil, errors.New("empty derivation path")
	}
	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func (p *accountClient) FormatAddress(account as.Account) string {
	return string(account.Address)
}

func (p *accountClient) GetAddress(publicKey []byte) string {
	return crypto.KeyToAddress(publicKey)
}

func (p *accountClient) GetKeyId(publicKey []byte) int64 {
	return crypto.Address(publicKey)
}

func (p *accountClient) GetETHAddress(publicKey []byte) string {
	if len(publicKey) == 0 {
		return ""
	}
	keccak := &hashalgo.Keccak256{}
	hash256 := keccak.GetHash(crypto.CutPub(publicKey))
	return "0x" + hex.EncodeToString(hash256[len(hash256)-20:])
}
