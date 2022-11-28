package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/response"
	"net/url"
	"strconv"
)

func (chain *queryClient) GetKeyInfo(account string) (*response.KeyInfoResult, error) {
	var ret response.KeyInfoResult
	err := chain.SendGet(`keyinfo/`+account, nil, &ret)
	if err != nil {
		return &ret, err
	}
	return &ret, nil
}

// WalletHistory Find account transaction history
// search type (outcome:transfer, out income:transfer in, all:get all)
func (chain *queryClient) WalletHistory(searchType string, page, limit int64) (*response.WalletHistoryResult, error) {
	var result response.WalletHistoryResult
	form := url.Values{}
	reqUrl := fmt.Sprintf("walletHistory?searchType=%s&page=%d&limit=%d", searchType, page, limit)
	err := chain.SendGet(reqUrl, &form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemName(ecosystem int64) (*response.EcosystemNameResult, error) {
	var result response.EcosystemNameResult
	reqUrl := fmt.Sprintf("ecosystemname?id=%d", ecosystem)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemCount() (*response.CountResult, error) {
	var result response.CountResult
	form := url.Values{}
	ecosystemsUrl := fmt.Sprintf("/metrics/ecosystems")
	err := chain.SendGet(ecosystemsUrl, &form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// SystemParams Find the specified name system params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (chain *queryClient) SystemParams(names string) (*response.ParamsResult, error) {
	var result response.ParamsResult
	systemParamsUrl := fmt.Sprintf("systemparams?names=%s", names)
	err := chain.SendGet(systemParamsUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemParam(ecosystem int64, name string) (*response.ParamResult, error) {
	var result response.ParamResult
	reqUrl := fmt.Sprintf("ecosystemparam/%s?ecosystem=%d", name, ecosystem)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// EcosystemParams Find the specified name ecosystem params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (chain *queryClient) EcosystemParams(ecosystem int64, names string) (*response.ParamsResult, error) {
	var result response.ParamsResult
	ecosystemParamsUrl := fmt.Sprintf("ecosystemparams?ecosystem=%d&names=%s", ecosystem, names)
	err := chain.SendGet(ecosystemParamsUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func (chain *queryClient) GetHistory(table string, id string) (*response.HistoryResult, error) {
	var result response.HistoryResult
	historyUrl := fmt.Sprintf("history/%s/%s", table, id)
	err := chain.SendGet(historyUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func (chain *queryClient) GetBalance(wallet string) (*response.TokenBalanceResult, error) {
	var result response.TokenBalanceResult
	reqUrl := fmt.Sprintf("balance/%s", wallet)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func (chain *queryClient) GetMaxBlockID() (*response.MaxBlockResult, error) {
	var result response.MaxBlockResult
	maxBlockIdUrl := fmt.Sprintf("maxblockid")
	err := chain.SendGet(maxBlockIdUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetBlockInfo(id int64) (*response.BlockInfoResult, error) {
	var rets response.BlockInfoResult
	blockInfoUrl := fmt.Sprintf("block/%d", id)
	err := chain.SendGet(blockInfoUrl, nil, &rets)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
func (chain *queryClient) Balance(wallet string, ecosystem int64) (*response.TokenBalanceResult, error) {
	var result response.TokenBalanceResult
	balanceUrl := fmt.Sprintf("balance/%s?ecosystem=%d", wallet, ecosystem)
	err := chain.SendGet(balanceUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func (chain *queryClient) AppParams(appid int64, names string, ecosystem int64) (*response.AppParamsResult, error) {
	var rets response.AppParamsResult
	appParamsUrl := fmt.Sprintf("appparams/%d?ecosystem=%d&names=%s", appid, ecosystem, names)
	err := chain.SendGet(appParamsUrl, nil, &rets)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}
func (chain *queryClient) AppParam(appid int64, name string, ecosystem int64) (*response.ParamResult, error) {
	var rets response.ParamResult
	appParamUrl := fmt.Sprintf("appparam/%d/%s?ecosystem=%d", appid, name, ecosystem)
	err := chain.SendGet(appParamUrl, nil, &rets)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}

// GetRow return where table id = id AND Login ecosystem = ecosystem, rowsName optional
func (chain *queryClient) GetRow(tableName string, id int64, rowsName string) (*response.RowResult, error) {
	var result response.RowResult
	getRowUrl := fmt.Sprintf("row/%s/%d", tableName, id)
	if rowsName != "" {
		getRowUrl += fmt.Sprintf("?columns=%s", rowsName)
	}
	err := chain.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetRowExtend return where columns = value AND Login ecosystem = ecosystem, rowsName optional
func (chain *queryClient) GetRowExtend(tableName, columns string, value string, rowsName string) (*response.RowResult, error) {
	var result response.RowResult
	getRowUrl := fmt.Sprintf("row/%s/%s/%s", tableName, columns, value)
	if rowsName != "" {
		getRowUrl += fmt.Sprintf("?columns=%s", rowsName)
	}
	err := chain.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetPageRow return page name = name and ecosystem id = login ecosystem id
func (chain *queryClient) GetPageRow(name string) (*response.PageResult, error) {
	var result response.PageResult
	getRowUrl := fmt.Sprintf("interface/page/%s", name)
	err := chain.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetMenuRow return menu name = name and ecosystem id = login ecosystem id
func (chain *queryClient) GetMenuRow(name string) (*response.MenuResult, error) {
	var result response.MenuResult
	getMenuUrl := fmt.Sprintf("interface/menu/%s", name)
	err := chain.SendGet(getMenuUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetSnippetRow return snippet name = name and ecosystem id = login ecosystem id
func (chain *queryClient) GetSnippetRow(name string) (*response.SnippetResult, error) {
	var result response.SnippetResult
	getMenuUrl := fmt.Sprintf("interface/snippet/%s", name)
	err := chain.SendGet(getMenuUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) BlockTxInfo(count, blockId int64) (*map[int64][]response.TxInfo, error) {
	var result map[int64][]response.TxInfo
	reqUrl := fmt.Sprintf("blocks?count=%d&block_id=%d", count, blockId)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) DetailedBlocks(count, blockId int64) (*map[int64]response.BlockDetailedInfo, error) {
	var result map[int64]response.BlockDetailedInfo
	reqUrl := fmt.Sprintf("detailed_blocks?count=%d&block_id=%d", count, blockId)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) BlocksCount() (*response.CountResult, error) {
	var result response.CountResult
	reqUrl := fmt.Sprintf("metrics/blocks")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) TransactionsCount() (*response.CountResult, error) {
	var result response.CountResult
	reqUrl := fmt.Sprintf("metrics/transactions")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) EcosystemsCount() (*response.CountResult, error) {
	var result response.CountResult
	reqUrl := fmt.Sprintf("metrics/ecosystems")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) KeysCount() (*response.CountResult, error) {
	var result response.CountResult
	reqUrl := fmt.Sprintf("metrics/keys")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) HonorNodesCount() (*response.CountResult, error) {
	var result response.CountResult
	reqUrl := fmt.Sprintf("metrics/honornodes")
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) DataVerify(tableName string, id int64, column, hash string) (*string, error) {
	var result string
	reqUrl := fmt.Sprintf("data/%s/%d/%s/%s", tableName, id, column, hash)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) BinaryVerify(id int64, hash string) (*any, error) {
	reqUrl := fmt.Sprintf("data/%d/data/%s", id, hash)
	var v any
	err := chain.SendGet(reqUrl, nil, &v)
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (chain *queryClient) GetTables(limit, offset int) (*response.TablesResult, error) {
	var v response.TablesResult
	reqUrl := fmt.Sprintf("tables?offset=%d&limit=%d", offset, limit)
	err := chain.SendGet(reqUrl, nil, &v)
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (chain *queryClient) GetTable(tableName string) (*response.TableResult, error) {
	var result response.TableResult
	reqUrl := fmt.Sprintf("table/%s", tableName)
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetList(tableName, columns string, limit, offset int64) (*response.ListResult, error) {
	var result response.ListResult
	getListUrl := fmt.Sprintf("list/%s?limit=%d&offset=%d&columns=%s", tableName, limit, offset, columns)
	err := chain.SendGet(getListUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetSections(language string, offset, limit int) (*response.ListResult, error) {
	var result response.ListResult
	reqUrl := fmt.Sprintf("sections?offset=%d&limit=%d", offset, limit)
	if language != "" {
		reqUrl += fmt.Sprintf("&lang=%s", language)
	}
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetCentrifugo() (*string, error) {
	var result string
	reqUrl := "config/centrifugo"
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetVersion() (*string, error) {
	var result string
	reqUrl := "version"
	err := chain.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error) {
	var result response.ListResult
	reqUrl := fmt.Sprintf("listWhere/%s", tableName)

	form := &url.Values{}
	if order != "" {
		form.Set("order", order)
	}
	if columns != "" {
		form.Set("columns", columns)
	}
	if page < 0 || limit <= 0 {
		return &result, errors.New("params invalid")
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	form.Set("offset", strconv.Itoa(offset))
	if where != nil {
		data, err := json.Marshal(where)
		if err != nil {
			return &result, err
		}
		form.Set("where", string(data))
	}

	err := chain.SendPost(reqUrl, form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (chain *queryClient) GetNodeListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error) {
	var result response.ListResult
	reqUrl := fmt.Sprintf("nodelistWhere/%s", tableName)

	form := &url.Values{}
	if order != "" {
		form.Set("order", order)
	}
	if columns != "" {
		form.Set("columns", columns)
	}
	if page < 0 || limit <= 0 {
		return &result, errors.New("params invalid")
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	form.Set("offset", strconv.Itoa(offset))
	if where != nil {
		data, err := json.Marshal(where)
		if err != nil {
			return &result, err
		}
		form.Set("where", string(data))
	}

	err := chain.SendPost(reqUrl, form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
