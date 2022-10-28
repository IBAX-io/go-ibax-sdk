package utxo

import (
	"encoding/hex"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/auth"
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/tx"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

func initConfig() *config.IbaxConfig {
	cnf := &config.IbaxConfig{}
	cnf.Ecosystem = 1
	cnf.JwtPrefix = "Bearer "
	cnf.PrivateKey = os.Getenv("IBAX_Private_Key1")
	cnf.ApiPath = "/api/v2/"
	cnf.ApiAddress = "http://localhost:7179"
	cnf.Cryptoer = "ECC_Secp256k1"
	cnf.Hasher = "KECCAK256"

	return cnf
}

func TestChain_NewUtxoTransaction(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := &url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1000"}, "comment": {""}}
	//smartTx, err := c.NewUtxoSmartTransaction(TypeContractToUTXO, form, "")
	//smartTx, err := c.NewUtxoSmartTransaction(TypeUTXOToContract, form, "")
	//smartTx, err := c.NewUtxoSmartTransaction(TypeTransfer, form, "")
	smartTx, err := c.NewUtxoSmartTransaction(TypeUTXOToContract, form, "")
	if err != nil {
		t.Errorf("New Utxo Smart Transaction failed: %s", err.Error())
		return
	}
	_, hash, err := c.NewUtxoTransaction(*smartTx)
	if err != nil {
		t.Errorf("New Utxo Transaction failed: %s", err.Error())
		return
	}
	fmt.Println("hash:", hex.EncodeToString(hash))
}

func TestChain_CallUtxo(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1000"}, "comment": {""}}
	smartTx, err := c.NewUtxoSmartTransaction(TypeContractToUTXO, form, "")
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
	hashMap, err := x.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var txHash string
		for _, v := range *hashMap {
			txHash = v
		}
		result, err := x.TxStatus(txHash, 10, time.Millisecond*500)
		if err != nil {
			t.Errorf("Tx Status failed:%s", err.Error())
			return
		}
		fmt.Println("result:", result)
	}
}

func TestChain_CallMoreUtxo(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form1 := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1000"}, "comment": {""}}
	form2 := url.Values{"recipient": {"1196-2490-5275-7101-3496"}, "amount": {"500"}}
	smartTx, err := c.NewUtxoSmartTransaction(TypeTransfer, form1, "")
	if err != nil {
		t.Errorf("New Utxo Smart Transaction form1 failed:%s", err.Error())
		return
	}

	data1, hash1, err := c.NewUtxoTransaction(*smartTx)
	if err != nil {
		t.Errorf("New Utxo Transaction form1 failed: %s", err.Error())
		return
	}

	smartTx, err = c.NewUtxoSmartTransaction(TypeTransfer, form2, "1")
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

	hashMap, err := x.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var hashs []string
		for _, v := range *hashMap {
			hashs = append(hashs, v)
		}
		results, err := x.TxsStatus(hashs, time.Millisecond*500)
		if err != nil {
			t.Errorf("Txs Status failed:%s", err.Error())
			return
		}
		fmt.Println("results:", results)
	}
}

func TestChain_CallMoreUtxo2(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	arrData := make(map[string][]byte)

	form := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"1"}, "comment": {strconv.FormatInt(time.Now().UnixMilli(), 10)}}
	smartTx, err := c.NewUtxoSmartTransaction(TypeTransfer, form, "")
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

	hashMap, err := x.SendTx(arrData)
	if err != nil {
		t.Errorf("Send Tx failed :%s", err.Error())
		return
	}
	if hashMap != nil {
		var hashs []string
		for _, v := range *hashMap {
			hashs = append(hashs, v)
		}
		results, err := x.TxsStatus(hashs, time.Millisecond*500)
		if err != nil {
			t.Errorf("Txs Status failed:%s", err.Error())
			return
		}
		fmt.Println("results:", results)
	}
}

func TestChain_AutoCallUtxo(t *testing.T) {
	cnf := initConfig()
	b := base.NewClient(cnf)
	x := tx.NewClient(b)
	c := NewClient(cnf, b, x)
	a := auth.NewClient(cnf, b)
	err := a.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}
	form := url.Values{"recipient": {"0000-0000-0000-0000-0000"}, "amount": {"5"}, "comment": {""}}
	result, err := c.AutoCallUtxo(TypeTransfer, form, "")
	if err != nil {
		t.Errorf("auto call utxo failed: %s", err.Error())
		return
	}

	fmt.Printf("auto call utxo result:%+v\n", result)
}
