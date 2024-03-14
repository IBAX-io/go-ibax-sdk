package example

import (
	"encoding/hex"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"sync"
	"testing"
	"time"
)

func TestNewMnemonic(t *testing.T) {
	c := client.NewClient(cnf)

	le := request.WordsLenTwelve
	//le = request.WordsLenFifteen
	//le = request.WordsLenEighteen
	//le = request.WordsTwentyOne
	//le = request.WordsLenTwentyTwo
	mnemonic, err := c.NewMnemonic(le)
	if err != nil {
		t.Errorf("new mnemonic failed :%s", err.Error())
		return
	}

	fmt.Println("mnemonic:", mnemonic)
}

func TestNewWallet(t *testing.T) {
	c := client.NewClient(cnf)

	le := request.WordsLenTwelve
	_, err := c.NewWallet(le)
	if err != nil {
		t.Errorf("New Wallet Failed:%s", err.Error())
		return
	}
}

func TestNewWalletFromMnemonic(t *testing.T) {
	c := client.NewClient(cnf)

	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	_, err := c.NewWalletFromMnemonic(mnemonic)
	if err != nil {
		t.Errorf("New Wallet From Mnemonic Failed:%s", err.Error())
		return
	}
}

func TestWallet(t *testing.T) {
	c := client.NewClient(cnf)

	//mnemonic := "tag volcano eight thank ti danger cost health above argue embrace heavy"
	//wallet, err := c.NewWalletFromMnemonic(mnemonic)
	//if err != nil {
	//t.Errorf("New Wallet From Mnemonic Failed:%s", err.Error())
	//return
	//}
	wallet, err := c.NewWallet(request.WordsLenTwentyTwo)
	if err != nil {
		t.Errorf("New wallet Failed:%s", err.Error())
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
	if cnf.Cryptoer == "ECC_Secp256k1" {
		hash256 := crypto.Hash(crypto.CutPub(publicKey))
		fmt.Println("ETH Address:", "0x"+hex.EncodeToString(hash256[len(hash256)-20:]))
	}

	//keccak := &hashalgo.Keccak256{}
	wg := sync.WaitGroup{}
	number := 100
	//Generate sub-accounts based on path
	for i := 1; i < number; i++ {
		account1, err := c.NewAccountFromPath(wallet, fmt.Sprintf("m/44'/60'/0'/0/%d", i), true)
		if err != nil {
			t.Errorf("New Account1 From Path Failed:%s", err.Error())
			return
		}

		//publicKey, err := c.GetPublicKey(wallet, account1)
		//if err != nil {
		//	t.Error(err)
		//	return
		//}

		privateKey, err := c.GetPrivateKey(wallet, account1)
		if err != nil {
			t.Error(err)
			return
		}

		//fmt.Println("PrivateKey:", hex.EncodeToString(privateKey))
		//fmt.Println("Public Key:", hex.EncodeToString(publicKey))
		//fmt.Println("IBAX Address:", account1.Address)
		//fmt.Println("IBAX KeyId:", crypto.Address(publicKey))

		//hash256 := keccak.GetHash(crypto.CutPub(publicKey))
		//fmt.Println("ETH Address:", "0x"+hex.EncodeToString(hash256[len(hash256)-20:]))

		c1 := c.GetConfig()
		c1.Ecosystem = 31
		c1.PrivateKey = hex.EncodeToString(privateKey)
		time.Sleep(1 * time.Second)
		go func(cg config.Config) {
			wg.Add(1)
			defer wg.Done()
			cli := client.NewClient(cg)
			err = cli.AutoLogin()
			if err != nil {
				fmt.Printf("auto login:%s\n", err.Error())
				return
			}
		}(c1)
	}
	wg.Wait()
	if len(wallet.Accounts()) != number {
		t.Error(fmt.Sprintf("expected %d,got:%d", number, len(wallet.Accounts())))
		return
	}
}

func Test_GetAddress(t *testing.T) {
	c := client.NewClient(cnf)
	publicKey := cnf.PublicKey

	ibaxAddress := c.GetAddress(publicKey)
	fmt.Println("IBAX Address:", ibaxAddress)

	ibaxKeyId := c.GetKeyId(publicKey)
	fmt.Println("IBAX KeyId:", ibaxKeyId)

	ethAddress := c.GetETHAddress(publicKey)
	fmt.Println("ETH Address:", ethAddress)
}
