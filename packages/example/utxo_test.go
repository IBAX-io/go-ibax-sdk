package example

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/shopspring/decimal"
	"net/url"
	"strconv"
	"testing"
	"time"
)

// UTXO Account Token Send
func TestIBAX_UtxoTokenSend(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	var comment string
	testMessage := "testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;" +
		"testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;" +
		"testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;" +
		"testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;" +
		"testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;" +
		"testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;" +
		"testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;testMessage;"
	for i := 0; i < 500; i++ {
		comment += testMessage
	}
	form := url.Values{"recipient": {"0000-0000-0000-0000-5555"}, "amount": {"1"}, "comment": {""}}
	form.Set("comment", comment)
	result, err := c.AutoCallUtxo(request.TypeTransfer, form, "")
	if err != nil {
		t.Errorf("utxo token send failed: %s", err.Error())
		return
	}

	fmt.Printf("utxo token send result:%+v\n", result)
}

// UTXO Account Transfer Self Contract Account
func TestIBAX_UtxoTransferSelf(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	amount := decimal.NewFromInt(1000000).Mul(decimal.NewFromInt(1000000000000))
	form := url.Values{"amount": {amount.String()}}
	result, err := c.AutoCallUtxo(request.TypeUTXOToContract, form, "")
	if err != nil {
		t.Errorf("utxo transfer self failed: %s", err.Error())
		return
	}

	fmt.Printf("utxo transfer self result:%+v\n", result)
}

// Contract Account Transfer Self UTXO Account
func TestIBAX_ContractTransferSelf(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"amount": {"100000000000000"}}
	result, err := c.AutoCallUtxo(request.TypeContractToUTXO, form, "")
	if err != nil {
		t.Errorf("contract transfer self failed: %s", err.Error())
		return
	}

	fmt.Printf("contract transfer self result:%+v\n", result)
}

func TestChain_CallUtxo(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"10000"}, "comment": {""}}
	smartTx, err := c.NewUtxoSmartTransaction(request.TypeContractToUTXO, form, "")
	if err != nil {
		t.Errorf("New Utxo Smart Transaction form1 failed:%s", err.Error())
		return
	}

	data, hash, err := c.NewUtxoTransaction(*smartTx)
	if err != nil {
		t.Errorf("New Utxo Transaction form1 failed: %s", err.Error())
		return
	}
	arrData := make(map[string][]byte)
	arrData[fmt.Sprintf("%x", hash)] = data
	hashMap, err := c.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var txHash string
		for _, v := range *hashMap {
			txHash = v
		}
		result, err := c.TxStatus(txHash, 10, time.Millisecond*500)
		if err != nil {
			t.Errorf("Tx Status failed:%s", err.Error())
			return
		}
		fmt.Println("result:", result)
	}
}

func TestChain_CallMoreUtxo(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form1 := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1000"}, "comment": {""}}
	form2 := url.Values{"recipient": {"1196-2490-5275-7101-3496"}, "amount": {"500"}}
	smartTx, err := c.NewUtxoSmartTransaction(request.TypeTransfer, form1, "")
	if err != nil {
		t.Errorf("New Utxo Smart Transaction form1 failed:%s", err.Error())
		return
	}

	data1, hash1, err := c.NewUtxoTransaction(*smartTx)
	if err != nil {
		t.Errorf("New Utxo Transaction form1 failed: %s", err.Error())
		return
	}

	smartTx, err = c.NewUtxoSmartTransaction(request.TypeTransfer, form2, "1")
	if err != nil {
		t.Errorf("New Utxo Smart Transaction form2 failed :%s", err.Error())
		return
	}

	data2, hash2, err := c.NewUtxoTransaction(*smartTx)
	if err != nil {
		t.Errorf("New Utxo Transaction form2 failed :%s", err.Error())
		return
	}
	arrData := make(map[string][]byte)
	arrData[fmt.Sprintf("%x", hash1)] = data1
	arrData[fmt.Sprintf("%x", hash2)] = data2

	hashMap, err := c.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var hashes []string
		for _, v := range *hashMap {
			hashes = append(hashes, v)
		}
		results, err := c.TxsStatus(hashes, time.Millisecond*500)
		if err != nil {
			t.Errorf("Txs Status failed:%s", err.Error())
			return
		}
		fmt.Println("results:", results)
	}
}

func TestChain_CallMoreUtxo2(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	arrData := make(map[string][]byte)

	form := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1"}, "comment": {strconv.FormatInt(time.Now().UnixMilli(), 10)}}
	smartTx, err := c.NewUtxoSmartTransaction(request.TypeTransfer, form, "")
	if err != nil {
		t.Errorf("New Utxo Smart Transaction failed:%s", err.Error())
		return
	}

	data, hash, err := c.NewUtxoTransaction(*smartTx)
	if err != nil {
		t.Errorf("New Utxo Transaction failed: %s", err.Error())
		return
	}

	arrData[fmt.Sprintf("%x", hash)] = data

	hashMap, err := c.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var hashes []string
		for _, v := range *hashMap {
			hashes = append(hashes, v)
		}
		results, err := c.TxsStatus(hashes, time.Millisecond*500)
		if err != nil {
			t.Errorf("Txs Status failed:%s", err.Error())
			return
		}
		fmt.Println("results:", results)
	}
}
