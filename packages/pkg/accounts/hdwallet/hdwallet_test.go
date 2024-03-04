/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package hdwallet

import (
	"encoding/hex"
	"fmt"
	accounts2 "github.com/IBAX-io/go-ibax-sdk/packages/pkg/accounts"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/common/crypto"
	"net/url"
	"reflect"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

func init() {
	crypto.InitAsymAlgo("ECC_Secp256k1")
	crypto.InitHashAlgo("KECCAK256")
}

func TestWallet(t *testing.T) {
	tests := []struct {
		name    string
		want    *Wallet
		wantErr bool
	}{
		{name: "test0", want: nil, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}

	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err0 := NewWalletFromMnemonic(mnemonic, "")
	if err0 != nil {
		t.Error(err0)
	}
	//var _ accounts2.Wallet = wallet

	path0, err0 := ParseDerivationPath("m/44'/60'/0'/0/0")

	// path0, err0 := ParseDerivationPath("m44'60'0'00") 	//! "m44'60'0'00"

	if err0 != nil {
		t.Error(err0)
	}
	if len(wallet.Accounts()) != 0 {
		t.Error("expected 0")
	}
	account0, err0 := wallet.Derive(path0, true)
	if err0 != nil {
		t.Error(err0)
	}

	path0_, err := ParseDerivationPath(account0.URL.Path)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("account0: %v\n", account0)
	fmt.Printf("account0's path0:  %s\n", account0.URL.Path)
	fmt.Printf("account0's path0_:  %s\n", path0_)
	priv0_0, err0 := wallet.derivePrivateKey(path0)
	if err0 != nil {
		t.Error(err0)
	}

	priv0_1, err1 := wallet.PrivateKey(account0)
	if err1 != nil {
		t.Error(err1)
	}
	fmt.Printf("priv0_0: %v,priv0_0:%s\n", priv0_0, hex.EncodeToString(priv0_0))
	fmt.Printf("priv0_1: %v,priv0_1:%s\n", priv0_1, hex.EncodeToString(priv0_1))

	publicKey, err := wallet.derivePublicKey(path0)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("publickey:%s,key_id:%d\n", hex.EncodeToString(publicKey), crypto.Address(publicKey))

	if len(wallet.Accounts()) != 1 {
		t.Error("expected 1")
	}
	if !wallet.Contains(account0) {
		t.Error("expected to contain account")
	}

	for i := 0; i <= 10; i++ {
		path1, err1 := ParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		if err1 != nil {
			t.Error(err1)
		}
		account1, err1 := wallet.Derive(path1, true)
		if err1 != nil {
			t.Error(err1)
		}
		fmt.Println("path", path1)
		publicKey, err := wallet.derivePublicKey(path1)
		if err != nil {
			t.Error(err)
		}

		privateKey, err := wallet.derivePrivateKey(path1)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("PrivateKey", hex.EncodeToString(privateKey))
		hash256 := crypto.Hash(crypto.CutPub(publicKey))
		fmt.Println("IBAX Address", account1.Address)
		fmt.Println("ETH Address", "0x"+hex.EncodeToString(hash256[len(hash256)-20:]))

		fmt.Println("IBAX KeyId", crypto.Address(publicKey))
		fmt.Println("Public Key", hex.EncodeToString(publicKey))
	}

}

func TestWallet_SignHash(t *testing.T) {
	type fields struct {
		mnemonic  string
		masterKey *hdkeychain.ExtendedKey
		seed      []byte
		url       accounts2.URL
		paths     map[accounts2.Address]accounts2.DerivationPath
		accounts  []accounts2.Account
		stateLock sync.RWMutex
	}
	type args struct {
		account accounts2.Account
		hash    []byte
	}

	// @todo [ ] Why is Init() needed?

	mnemonic, _ := NewMnemonic()

	wallet, _ := NewWalletFromMnemonic(mnemonic, "123456")

	//var _ accounts2.Wallet = wallet

	path0, err := ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		t.Error(err)
	}

	account0, err := wallet.Derive(path0, true)
	if err != nil {
		t.Error(err)
	}

	if len(wallet.Accounts()) != 1 {
		t.Error("expected 1")
	}
	if !wallet.Contains(account0) {
		t.Error("expected to contain account")
	}

	// privateKey, err := wallet.PrivateKey(account0)
	privateKey, err := wallet.derivePrivateKey(path0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("privateKey", privateKey)

	// @todo [ ]
	publicKey_, err := crypto.PrivateToPublic(privateKey)
	if err != nil {
		t.Error(err)
	}
	// publicKey, err := wallet.PublicKey(account0)
	publicKey, err := wallet.derivePublicKey(path0)
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(publicKey_) != hex.EncodeToString(publicKey) {
		t.Error("public key not match")
	}

	testHash := []byte("test")
	// signature_, err := crypto.Sign(privateKey, testHash)
	// if err != nil {
	// 	t.Error(err)
	// }
	signature, err := wallet.SignHash(account0, testHash)
	if err != nil {
		t.Error(err)
	}

	// if hex.EncodeToString(signature_) != hex.EncodeToString(signature) {
	// 	fmt.Println("signature_", signature_)
	// 	fmt.Println("signature", signature)
	// 	t.Error("signature not match")
	// }

	verified, err := crypto.Verify(publicKey, testHash, signature)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("verified", verified)

	// t.Run(tt.name, func(t *testing.T) {
	// 	w := &Wallet{
	// 		mnemonic:  tt.fields.mnemonic,
	// 		masterKey: tt.fields.masterKey,
	// 		seed:      tt.fields.seed,
	// 		url:       tt.fields.url,
	// 		paths:     tt.fields.paths,
	// 		accounts:  tt.fields.accounts,
	// 		stateLock: tt.fields.stateLock,
	// 	}
	// 	got, err := w.SignHash(tt.args.account, tt.args.hash)
	// 	if (err != nil) != tt.wantErr {
	// 		t.Errorf("Wallet.SignHash() error = %v, wantErr %v", err, tt.wantErr)
	// 		return
	// 	}
	// 	if !reflect.DeepEqual(got, tt.want) {
	// 		t.Errorf("Wallet.SignHash() = %v, want %v", got, tt.want)
	// 	}
	// }
	// )
}

func TestWallet_MultiSignForm(t *testing.T) {
	type fields struct {
		mnemonic  string
		masterKey *hdkeychain.ExtendedKey
		seed      []byte
		url       accounts2.URL
		paths     map[accounts2.Address]accounts2.DerivationPath
		accounts  []accounts2.Account
		stateLock sync.RWMutex
	}
	type args struct {
		Data      string
		Voter     int
		Threshold int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    url.Values
		wantErr bool
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet := &Wallet{
				mnemonic:  tt.fields.mnemonic,
				masterKey: tt.fields.masterKey,
				seed:      tt.fields.seed,
				url:       tt.fields.url,
				paths:     tt.fields.paths,
				accounts:  tt.fields.accounts,
				stateLock: tt.fields.stateLock,
			}
			got, err := wallet.MultiSignForm(tt.args.Data, tt.args.Voter, tt.args.Threshold)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wallet.MultiSignForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wallet.MultiSignForm() = %v, want %v", got, tt.want)
			}
		})
	}
}
