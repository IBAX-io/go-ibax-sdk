package example

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/client"
	"github.com/IBAX-io/go-ibax-sdk/tx/utxo"
	"net/url"
	"testing"
)

// UTXO Account Token Send
func TestIBAX_UtxoTokenSend(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1000000000000"}, "comment": {""}}
	result, err := c.AutoCallUtxo(utxo.TypeTransfer, form, "")
	if err != nil {
		t.Errorf("utxo token send failed: %s", err.Error())
		return
	}

	fmt.Printf("utxo token send result:%+v\n", result)
}

// UTXO Account Transfer Self Contract Account
func TestIBAX_UtxoTransferSelf(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"amount": {"1000000000000"}}
	result, err := c.AutoCallUtxo(utxo.TypeUTXOToContract, form, "")
	if err != nil {
		t.Errorf("utxo transfer self failed: %s", err.Error())
		return
	}

	fmt.Printf("utxo transfer self result:%+v\n", result)
}

// Contract Account Transfer Self UTXO Account
func TestIBAX_ContractTransferSelf(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"amount": {"1000000000000"}}
	result, err := c.AutoCallUtxo(utxo.TypeContractToUTXO, form, "")
	if err != nil {
		t.Errorf("contract transfer self failed: %s", err.Error())
		return
	}

	fmt.Printf("contract transfer self result:%+v\n", result)
}
