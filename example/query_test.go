package example

import (
	"encoding/hex"
	"github.com/IBAX-io/go-ibax-sdk/client"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/converter"
	"github.com/bitly/go-simplejson"
	"strconv"
	"testing"
)

func TestIBAX_WalletHistory(t *testing.T) {
	cnf := initFounderTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	//result, er := c.WalletHistory("all", 10, 10)
	//result, er := c.WalletHistory("income", 0, 20)
	_, er := c.WalletHistory("outcome", 0, 20)
	if er != nil {
		t.Errorf("wallet history failed :%s", er.Error())
		return
	}
}

func TestIBAX_EcosystemCount(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.EcosystemCount()
	if er != nil {
		t.Errorf("ecosystem count failed :%s", er.Error())
		return
	}

}
func TestIBAX_SystemParams(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	//result, er := c.SystemParams("")
	_, er := c.SystemParams("number_of_nodes,taxes_size")
	if er != nil {
		t.Errorf("system params failed :%s", er.Error())
		return
	}

}
func TestIBAX_EcosystemParam(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.EcosystemParam(1, "founder_account")
	if er != nil {
		t.Errorf("ecosystem param failed :%s", er.Error())
		return
	}

}
func TestIBAX_EcosystemParams(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	//_, er := c.EcosystemParams(1, "changing_page,changing_menu")
	_, er := c.EcosystemParams(1, "")
	if er != nil {
		t.Errorf("ecosystem params failed :%s", er.Error())
		return
	}

}
func TestIBAX_GetHistory(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetHistory("contracts", "15")
	if er != nil {
		t.Errorf("get history failed :%s", er.Error())
		return
	}

}

func TestIBAX_GetBalance(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.GetBalance("0666-0819-7161-7879-5186")
	if er != nil {
		t.Errorf("get balance failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetMaxBlockID(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.GetMaxBlockID()
	if er != nil {
		t.Errorf("get max block id failed :%s", er.Error())
		return
	}

}
func TestIBAX_GetBlockInfo(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.GetBlockInfo(2)
	if er != nil {
		t.Errorf("get block info failed :%s", er.Error())
		return
	}
}
func TestIBAX_Balance(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.Balance("0666-0819-7161-7879-5186", 1)
	if er != nil {
		t.Errorf("get ibax balance failed :%s", er.Error())
		return
	}

}
func TestIBAX_AppParams(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.AppParams(1, "", 1)
	if er != nil {
		t.Errorf("app params failed :%s", er.Error())
		return
	}
}

func TestIBAX_AppParam(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.AppParam(1, "role_governancer", 1)
	if er != nil {
		t.Errorf("app param failed :%s", er.Error())
		return
	}

}

func TestIBAX_EcosystemName(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.EcosystemName(2)
	if er != nil {
		t.Errorf("get ecosystem name failed :%s", er.Error())
		return
	}

}

func TestIBAX_GetList(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetList("keys", "id,amount", 10, 0)
	if er != nil {
		t.Errorf("get list failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetRow(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	//_, er := c.GetRow("keys", 5555, "")
	_, er := c.GetRow("pages", 87, "value")
	if er != nil {
		t.Errorf("get row failed :%s", er.Error())
		return
	}
}

func TestIbax_DataVerify(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	table := "pages"
	id := int64(87)
	rowsName := "value"
	rets, er := c.GetRow(table, id, rowsName)
	if er != nil {
		t.Errorf("get row failed :%s", er.Error())
		return
	}
	cfg := c.GetConfig()

	if v, ok := rets.Value[rowsName]; ok {
		hash := crypto.Hash([]byte(v))

		tableName := converter.ParseTable(table, cfg.Ecosystem)
		_, err := c.DataVerify(tableName, id, rowsName, hex.EncodeToString(hash))
		if err != nil {
			t.Errorf("data verify failed :%s", err.Error())
			return
		}
	}

}

func TestIbax_BinaryVerify(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	table := "binaries"
	id := int64(3)
	rowsName := "data"
	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}
	jst.Set("id", map[string]any{"$eq": strconv.FormatInt(id, 10)})
	rets, er := c.GetNodeListWhere(table, jst, rowsName, "", 1, 1)
	if er != nil {
		t.Errorf("get row failed:%s", er.Error())
		return
	}

	if rets.Count > 0 {
		for _, list := range rets.List {
			if v, ok := list[rowsName]; ok {
				data, err := hex.DecodeString(v)
				if err != nil {
					t.Errorf("data decode failed:%s", er.Error())
					return
				}
				hash := crypto.Hash(data)
				_, err = c.BinaryVerify(id, hex.EncodeToString(hash))
				if err != nil {
					t.Errorf("binary verify failed :%s", err.Error())
					return
				}
			}
		}
	}

}

func TestIBAX_GetRowExtend(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetRowExtend("app_params", "name", "role_governancer", "app_id,name")
	//_, er := c.GetRowExtend("app_params", "name", "role_governancer", "")
	if er != nil {
		t.Errorf("get row extend failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetPageRow(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetPageRow("default_page")
	if er != nil {
		t.Errorf("get page row failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetMenuRow(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetMenuRow("default_menu")
	if er != nil {
		t.Errorf("get menu row failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetSnippetRow(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetSnippetRow("export_info")
	if er != nil {
		t.Errorf("get snippet row failed :%s", er.Error())
		return
	}
}

func TestIBAX_BlockTxInfo(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.BlockTxInfo(2, 87)
	if er != nil {
		t.Errorf("block tx info failed :%s", er.Error())
		return
	}
}

func TestIBAX_DetailBlocks(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.DetailedBlocks(1, 87)
	if er != nil {
		t.Errorf("detailed blocks failed :%s", er.Error())
		return
	}
}

func TestIBAX_BlocksCount(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.BlocksCount()
	if er != nil {
		t.Errorf("blocks count failed :%s", er.Error())
		return
	}
}

func TestIBAX_TransactionsCount(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.TransactionsCount()
	if er != nil {
		t.Errorf("transactions count failed :%s", er.Error())
		return
	}
}

func TestIBAX_EcosystemsCount(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.EcosystemsCount()
	if er != nil {
		t.Errorf("ecosystems count failed :%s", er.Error())
		return
	}
}

func TestIBAX_KeysCount(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.KeysCount()
	if er != nil {
		t.Errorf("keys count failed :%s", er.Error())
		return
	}
}

func TestIBAX_HonorNodesCount(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.HonorNodesCount()
	if er != nil {
		t.Errorf("honor nodes count failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetTables(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetTables(10, 0)
	if er != nil {
		t.Errorf("get tables failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetTable(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetTable("app_params")
	if er != nil {
		t.Errorf("get table failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetSections(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()
	_, er := c.GetSections("zh", 0, 10)
	if er != nil {
		t.Errorf("get Sections failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetCentrifugo(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.GetCentrifugo()
	if er != nil {
		t.Errorf("get Centrifugo failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetVersion(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	_, er := c.GetVersion()
	if er != nil {
		t.Errorf("get version failed :%s", er.Error())
		return
	}
}

func TestIBAX_GetListWhere1(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()

	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}

	j1, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json j1 failed :%s", err.Error())
		return
	}
	j2, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json j2 failed :%s", err.Error())
		return
	}

	j1.Set("key", map[string]any{"$eq": "avatar"})
	j2.Set("account", map[string]any{"$eq": "0666-0819-7161-7879-5186"})

	jst.Set("$or", []any{j1, j2})

	_, er := c.GetListWhere("buffer_data", jst, "ecosystem,key,value,account", "id desc", 1, 10)
	if er != nil {
		t.Errorf("get list where failed :%s", er.Error())
		return
	}

}

func TestIBAX_GetListWhere2(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()

	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}

	jst.Set("key", map[string]any{"$eq": "avatar"})
	jst.Set("account", map[string]any{"$eq": "0666-0819-7161-7879-5186"})

	_, er := c.GetListWhere("buffer_data", jst, "ecosystem,key,value,account", "id desc", 1, 10)
	if er != nil {
		t.Errorf("get list where failed :%s", er.Error())
		return
	}

}

func TestIBAX_GetNodeListWhere(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)
	c.AutoLogin()

	jst, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		t.Errorf("new json failed :%s", err.Error())
		return
	}

	jst.Set("key", map[string]any{"$eq": "avatar"})
	jst.Set("account", map[string]any{"$eq": "0666-0819-7161-7879-5186"})

	_, er := c.GetNodeListWhere("buffer_data", jst, "ecosystem,key,value,account", "id desc", 1, 10)
	if er != nil {
		t.Errorf("get node list where failed :%s", er.Error())
		return
	}

}

func TestChain_GetKeyInfo(t *testing.T) {
	cnf := initNormalTest()
	c := client.NewClient(cnf)

	_, err := c.GetKeyInfo(cnf.Account)
	if err != nil {
		t.Errorf("get key info failed:%s", err.Error())
		return
	}
}
