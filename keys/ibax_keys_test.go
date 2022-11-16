package keys

import (
	"encoding/hex"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto/hashalgo"
	"os"
	"strings"
	"testing"
)

//IBAX HDWallet Generator https://ibax-io.github.io/ibax-hdwallet-generator

var cnf *config.IbaxConfig

func init() {
	cnf = &config.IbaxConfig{}
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = os.Getenv("IBAX_Private_Key1")
	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7179"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"
}

func TestNewMnemonic(t *testing.T) {
	b := base.NewClient(cnf)
	c := NewClient(b)

	var le wordsLenType
	le = WordsLenTwelve
	//le = WordsLenFifteen
	//le = WordsLenEighteen
	//le = WordsTwentyOne
	//le = WordsLenTwentyTwo
	mnemonic, err := c.NewMnemonic(le)
	if err != nil {
		t.Errorf("new mnemonic failed :%s", err.Error())
		return
	}
	length := len(strings.Fields(mnemonic))

	words := wordsLenType(length)
	if words != le {
		t.Errorf("mnemonic words length expected %d,got:%d", le, words)
		return
	}
	fmt.Println("mnemonic:", mnemonic)
}

func TestNewWallet(t *testing.T) {
	b := base.NewClient(cnf)
	c := NewClient(b)

	var le wordsLenType
	le = WordsLenTwelve
	_, err := c.NewWallet(le)
	if err != nil {
		t.Errorf("New Wallet Failed:%s", err.Error())
		return
	}
}

func TestNewWalletFromMnemonic(t *testing.T) {
	b := base.NewClient(cnf)
	c := NewClient(b)

	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	_, err := c.NewWalletFromMnemonic(mnemonic)
	if err != nil {
		t.Errorf("New Wallet From Mnemonic Failed:%s", err.Error())
		return
	}
}

func TestWallet(t *testing.T) {
	b := base.NewClient(cnf)
	c := NewClient(b)

	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := c.NewWalletFromMnemonic(mnemonic)
	if err != nil {
		t.Errorf("New Wallet From Mnemonic Failed:%s", err.Error())
		return
	}

	if len(wallet.Accounts()) != 0 {
		t.Error("expected 0")
		return
	}
	account, err := c.NewAccountFromPath(wallet, "m/44'/60'/0'/0/0", true)
	if err != nil {
		t.Errorf("New Account From Path Failed:%s", err.Error())
		return
	}
	if len(wallet.Accounts()) != 1 {
		t.Error("expected 1")
		return
	}
	if !wallet.Contains(account) {
		t.Error("expected to contain account")
		return
	}

	//Private Key
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		t.Errorf("parse Private Key Failed:%s", err.Error())
		return
	}
	fmt.Printf("Private Key:%s\n", hex.EncodeToString(privateKey))

	//publicKey Key
	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Public Key:%s\n", hex.EncodeToString(publicKey))

	//IBAX Address
	fmt.Println("IBAX Address:", account.Address)
	fmt.Println("IBAX KeyId:", crypto.Address(publicKey))
	cnf := b.GetConfig()
	if cnf.Cryptoer == "ECC_Secp256k1" {
		hash256 := crypto.Hash(crypto.CutPub(publicKey))
		fmt.Println("ETH Address:", "0x"+hex.EncodeToString(hash256[len(hash256)-20:]))
	}

	keccak := &hashalgo.Keccak256{}
	//Generate sub-accounts based on path
	for i := 1; i <= 10; i++ {
		account1, err := c.NewAccountFromPath(wallet, fmt.Sprintf("m/44'/60'/0'/0/%d", i), false)
		if err != nil {
			t.Errorf("New Account1 From Path Failed:%s", err.Error())
			return
		}

		publicKey, err := c.GetPublicKey(wallet, account1)
		if err != nil {
			t.Error(err)
			return
		}

		privateKey, err := c.GetPrivateKey(wallet, account1)
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println("PrivateKey:", hex.EncodeToString(privateKey))
		fmt.Println("Public Key:", hex.EncodeToString(publicKey))
		fmt.Println("IBAX Address:", account1.Address)
		fmt.Println("IBAX KeyId:", crypto.Address(publicKey))

		hash256 := keccak.GetHash(crypto.CutPub(publicKey))
		fmt.Println("ETH Address:", "0x"+hex.EncodeToString(hash256[len(hash256)-20:]))

	}
	if len(wallet.Accounts()) != 1 {
		t.Error("expected 1")
		return
	}
}

func Test_GetAddress(t *testing.T) {
	b := base.NewClient(cnf)
	c := NewClient(b)
	cnf := b.GetConfig()
	publicKey := cnf.PublicKey

	ibaxAddress := c.GetAddress(publicKey)
	fmt.Println("IBAX Address:", ibaxAddress)

	ibaxKeyId := c.GetKeyId(publicKey)
	fmt.Println("IBAX KeyId:", ibaxKeyId)

	ethAddress := c.GetETHAddress(publicKey)
	fmt.Println("ETH Address:", ethAddress)
}
