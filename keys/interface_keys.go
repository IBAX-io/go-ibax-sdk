package keys

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	as "github.com/IBAX-io/go-ibax-sdk/pkg/accounts"
	hd "github.com/IBAX-io/go-ibax-sdk/pkg/accounts/hdwallet"
)

type Keys interface {
	NewMnemonic(wordsLenType) (mnemonic string, err error)
	NewWallet(length wordsLenType) (*hd.Wallet, error)
	NewWalletFromMnemonic(mnemonic string) (*hd.Wallet, error)
	NewAccountFromPath(wallet *hd.Wallet, path string, pin bool) (as.Account, error)
	GetPrivateKey(wallet *hd.Wallet, account as.Account) ([]byte, error)
	GetPublicKey(wallet *hd.Wallet, account as.Account) ([]byte, error)
	FormatAddress(account as.Account) string
	GetAddress(publicKey []byte) string
	GetKeyId(publicKey []byte) int64
	GetETHAddress(publicKey []byte) string
}

type accountClient struct {
	base.Base
}

func NewClient(b base.Base) Keys {
	return &accountClient{Base: b}
}
