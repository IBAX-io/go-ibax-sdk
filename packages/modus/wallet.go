package modus

import (
	as "github.com/IBAX-io/go-ibax-sdk/packages/pkg/accounts"
	hd "github.com/IBAX-io/go-ibax-sdk/packages/pkg/accounts/hdwallet"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
)

type Wallet interface {
	NewMnemonic(length request.WordsLenType) (mnemonic string, err error)
	NewWallet(length request.WordsLenType) (*hd.Wallet, error)
	NewWalletFromMnemonic(mnemonic string) (*hd.Wallet, error)
	NewAccountFromPath(wallet *hd.Wallet, path string, pin bool) (as.Account, error)
	GetPrivateKey(wallet *hd.Wallet, account as.Account) ([]byte, error)
	GetPublicKey(wallet *hd.Wallet, account as.Account) ([]byte, error)
	FormatAddress(account as.Account) string
	GetAddress(publicKey []byte) string
	GetKeyId(publicKey []byte) int64
	GetETHAddress(publicKey []byte) string
}
