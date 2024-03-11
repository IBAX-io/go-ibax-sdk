package query

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/modus"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/go-ibax-sdk/packages/response"
	"net/url"
	"strconv"
)

type query struct {
	modus.Base
}

func New(b modus.Base) modus.Query {
	return &query{Base: b}
}

func (q *query) GetKeyInfo(account string) (*response.KeyInfoResult, error) {
	var ret response.KeyInfoResult
	err := q.SendGet(`keyinfo/`+account, nil, &ret)
	if err != nil {
		return &ret, err
	}
	return &ret, nil
}

func (q *query) EcosystemName(ecosystem int64) (*response.EcosystemNameResult, error) {
	var result response.EcosystemNameResult
	reqUrl := fmt.Sprintf("ecosystemname?id=%d", ecosystem)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// SystemParams Find the specified name system params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (q *query) SystemParams(names string, params ...int) (*response.ParamsResult, error) {
	var result response.ParamsResult
	systemParamsUrl := fmt.Sprintf("systemparams?names=%s", names)
	err := q.SendGet(systemParamsUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) EcosystemParam(ecosystem int64, name string) (*response.ParamResult, error) {
	var result response.ParamResult
	reqUrl := fmt.Sprintf("ecosystemparam/%s?ecosystem=%d", name, ecosystem)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// EcosystemParams Find the specified name ecosystem params
// names Use‘,’ split multiple query parameters
// names if is zero value: get all
func (q *query) EcosystemParams(ecosystem int64, names string, params ...int) (*response.ParamsResult, error) {
	var result response.ParamsResult
	ecosystemParamsUrl := fmt.Sprintf("ecosystemparams?ecosystem=%d&names=%s", ecosystem, names)
	err := q.SendGet(ecosystemParamsUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetHistory(table string, id uint64) (*response.HistoryResult, error) {
	var result response.HistoryResult
	historyUrl := fmt.Sprintf("history/%s/%d", table, id)
	err := q.SendGet(historyUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetBalance(wallet string) (*response.TokenBalanceResult, error) {
	var result response.TokenBalanceResult
	reqUrl := fmt.Sprintf("balance/%s", wallet)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

type maxBlockResult struct {
	MaxBlockID int64 `json:"max_block_id"`
}

func (q *query) GetMaxBlockID() (int64, error) {
	var result maxBlockResult
	maxBlockIdUrl := fmt.Sprintf("maxblockid")
	err := q.SendGet(maxBlockIdUrl, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.MaxBlockID, nil
}

func (q *query) GetBlockInfo(id int64) (*response.BlockInfoResult, error) {
	var rlt response.BlockInfoHexResult
	var rets response.BlockInfoResult
	blockInfoUrl := fmt.Sprintf("block/%d", id)
	err := q.SendGet(blockInfoUrl, nil, &rlt)
	if err != nil {
		return &rets, err
	}
	rets.Hash = hex.EncodeToString(rlt.Hash)
	rets.RollbacksHash = hex.EncodeToString(rlt.RollbacksHash)
	rets.Time = rlt.Time
	rets.Tx = rlt.Tx
	rets.KeyID = rlt.KeyID
	rets.ConsensusMode = rlt.ConsensusMode
	rets.NodePosition = rlt.NodePosition
	return &rets, nil
}

func (q *query) Balance(wallet string, ecosystem int64) (*response.TokenBalanceResult, error) {
	var result response.TokenBalanceResult
	balanceUrl := fmt.Sprintf("balance/%s?ecosystem=%d", wallet, ecosystem)
	err := q.SendGet(balanceUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) AppParams(appid int64, names string, ecosystem int64, params ...int) (*response.AppParamsResult, error) {
	var rets response.AppParamsResult
	appParamsUrl := fmt.Sprintf("appparams/%d?ecosystem=%d&names=%s", appid, ecosystem, names)
	err := q.SendGet(appParamsUrl, nil, &rets)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}

func (q *query) AppParam(appid int64, name string, ecosystem int64) (*response.ParamResult, error) {
	var rets response.ParamResult
	appParamUrl := fmt.Sprintf("appparam/%d/%s?ecosystem=%d", appid, name, ecosystem)
	err := q.SendGet(appParamUrl, nil, &rets)
	if err != nil {
		return &rets, err
	}
	return &rets, nil
}

// GetRow return where table id = id AND Login ecosystem = ecosystem, rowsName optional
func (q *query) GetRow(tableName string, id int64, columns string, whereColumn string) (*response.RowResult, error) {
	var result response.RowResult
	getRowUrl := fmt.Sprintf("row/%s/%d", tableName, id)
	if columns != "" {
		getRowUrl += fmt.Sprintf("?columns=%s", columns)
	}
	err := q.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetRowExtend return where columns = value AND Login ecosystem = ecosystem, rowsName optional
func (q *query) GetRowExtend(tableName, columns string, value string, rowsName string) (*response.RowResult, error) {
	var result response.RowResult
	getRowUrl := fmt.Sprintf("row/%s/%s/%s", tableName, columns, value)
	if rowsName != "" {
		getRowUrl += fmt.Sprintf("?columns=%s", rowsName)
	}
	err := q.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetPageRow return page name = name and ecosystem id = login ecosystem id
func (q *query) GetPageRow(name string) (*response.PageResult, error) {
	var result response.PageResult
	getRowUrl := fmt.Sprintf("interface/page/%s", name)
	err := q.SendGet(getRowUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetMenuRow return menu name = name and ecosystem id = login ecosystem id
func (q *query) GetMenuRow(name string) (*response.MenuResult, error) {
	var result response.MenuResult
	getMenuUrl := fmt.Sprintf("interface/menu/%s", name)
	err := q.SendGet(getMenuUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetSnippetRow return snippet name = name and ecosystem id = login ecosystem id
func (q *query) GetSnippetRow(name string) (*response.SnippetResult, error) {
	var result response.SnippetResult
	getMenuUrl := fmt.Sprintf("interface/snippet/%s", name)
	err := q.SendGet(getMenuUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

// GetAppContent Get Obtain application-related information (including page, snippet, menu)
func (q *query) GetAppContent(appId int64) (*response.AppContentResult, error) {
	var result response.AppContentResult
	reqUrl := fmt.Sprintf("appcontent/%d", appId)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) BlocksTxInfo(blockId, count int64) (*map[int64][]response.BlockTxInfo, error) {
	var result map[int64][]response.BlockTxInfoHex
	var rlt map[int64][]response.BlockTxInfo
	reqUrl := fmt.Sprintf("blocks?count=%d&block_id=%d", count, blockId)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &rlt, err
	}
	rlt = make(map[int64][]response.BlockTxInfo, len(result))
	for block, v1 := range result {
		var list []response.BlockTxInfo
		for _, v2 := range v1 {
			var info response.BlockTxInfo
			info.Params = v2.Params
			info.Hash = hex.EncodeToString(v2.Hash)
			info.KeyID = v2.KeyID
			info.ContractName = v2.ContractName
			list = append(list, info)
		}
		rlt[block] = list
	}
	return &rlt, nil
}

func (q *query) DetailedBlocks(blockId, count int64) (*map[int64]response.BlockDetailedInfo, error) {
	var result map[int64]response.BlockDetailedInfo
	reqUrl := fmt.Sprintf("detailed_blocks?count=%d&block_id=%d", count, blockId)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

type countResult struct {
	Count int64 `json:"count"`
}

func (q *query) BlocksCount() (int64, error) {
	var result countResult
	reqUrl := fmt.Sprintf("metrics/blocks")
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (q *query) TransactionsCount() (int64, error) {
	var result countResult
	reqUrl := fmt.Sprintf("metrics/transactions")
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (q *query) EcosystemCount() (int64, error) {
	var result countResult
	reqUrl := fmt.Sprintf("metrics/ecosystems")
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (q *query) KeysCount() (int64, error) {
	var result countResult
	reqUrl := fmt.Sprintf("metrics/keys")
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (q *query) HonorNodesCount() (int64, error) {
	var result countResult
	reqUrl := fmt.Sprintf("metrics/honornodes")
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// DataVerify return request.FileType or string
func (q *query) DataVerify(tableName string, id int64, column, hash, fileName string) (result request.FileType, err error) {
	reqUrl := fmt.Sprintf("data/%s/%d/%s/%s", tableName, id, column, hash)
	if fileName != "" {
		result.Name = fileName
	}
	err = q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return request.FileType{}, err
	}
	return result, nil
}

// BinaryVerify return request.FileType or string
func (q *query) BinaryVerify(id int64, hash string, fileName string) (result request.FileType, err error) {
	reqUrl := fmt.Sprintf("data/%d/data/%s", id, hash)
	if fileName != "" {
		result.Name = fileName
	}
	err = q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return request.FileType{}, err
	}
	return result, nil
}

func (q *query) GetAvatar(account string, ecosystem int64, fileName string) (result request.FileType, err error) {
	reqUrl := fmt.Sprintf("avatar/%d/%s", ecosystem, account)
	var v *request.FileType
	if fileName != "" {
		v = request.NewFileType(fileName)
	} else {
		return result, errors.New("filename can't not be empty")
	}
	err = q.SendGet(reqUrl, nil, v)
	if err != nil {
		return result, err
	}
	result = *v
	return result, nil
}

func (q *query) GetTableCount(offset, limit int) (*response.TablesResult, error) {
	var v response.TablesResult
	reqUrl := fmt.Sprintf("tables?offset=%d&limit=%d", offset, limit)
	err := q.SendGet(reqUrl, nil, &v)
	if err != nil {
		return &v, err
	}
	return &v, nil
}

func (q *query) GetTable(tableName string) (*response.TableResult, error) {
	var result response.TableResult
	reqUrl := fmt.Sprintf("table/%s", tableName)
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetList(params request.GetList) (*response.ListResult, error) {
	var result response.ListResult
	getListUrl := fmt.Sprintf("list/%s?limit=%d&offset=%d&columns=%s", params.Name, params.Limit, params.Offset, params.Columns)
	err := q.SendGet(getListUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetSections(language string, offset, limit int) (*response.ListResult, error) {
	var result response.ListResult
	reqUrl := fmt.Sprintf("sections?offset=%d&limit=%d", offset, limit)
	if language != "" {
		reqUrl += fmt.Sprintf("&lang=%s", language)
	}
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetIBAXConfig(option string) (*string, error) {
	var result string
	reqUrl := "config/" + option
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetVersion() (*string, error) {
	var result string
	reqUrl := "version"
	err := q.SendGet(reqUrl, nil, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetListWhere(params request.GetList) (*response.ListResult, error) {
	var result response.ListResult
	reqUrl := fmt.Sprintf("listWhere/%s", params.Name)

	form := &url.Values{}
	if params.Order != nil {
		data, err := json.Marshal(params.Order)
		if err != nil {
			return &result, err
		}
		form.Set("order", string(data))
	}
	if params.Columns != "" {
		form.Set("columns", params.Columns)
	}
	if params.Offset < 0 || params.Limit <= 0 {
		return &result, errors.New("params invalid")
	}
	form.Set("offset", strconv.Itoa(params.Offset))
	if params.Where != nil {
		data, err := json.Marshal(params.Where)
		if err != nil {
			return &result, err
		}
		form.Set("where", string(data))
	}

	err := q.SendPost(reqUrl, form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) GetNodeListWhere(params request.GetList) (*response.ListResult, error) {
	var result response.ListResult
	reqUrl := fmt.Sprintf("nodelistWhere/%s", params.Name)

	form := &url.Values{}
	if params.Order != nil {
		data, err := json.Marshal(params.Order)
		if err != nil {
			return &result, err
		}
		form.Set("order", string(data))
	}
	if params.Columns != "" {
		form.Set("columns", params.Columns)
	}
	if params.Offset < 0 || params.Limit <= 0 {
		return &result, errors.New("params invalid")
	}
	form.Set("limit", strconv.Itoa(params.Limit))
	form.Set("offset", strconv.Itoa(params.Offset))
	if params.Where != nil {
		data, err := json.Marshal(params.Where)
		if err != nil {
			return &result, err
		}
		form.Set("where", string(data))
	}

	err := q.SendPost(reqUrl, form, &result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (q *query) BlockTxCount(blockIdOrHash any) (int64, error) {
	return 0, response.NotSupportError
}

func (q *query) DetailedBlock(blockIdOrBlockHash any) (*response.BlockDetailedInfo, error) {
	return nil, response.NotSupportError
}

func (q *query) EcosystemInfo(ecosystem int64) (*response.EcosystemInfo, error) {
	return nil, response.NotSupportError
}

func (q *query) GetMemberInfo(account string, ecosystem int64) (*response.MemberInfo, error) {
	return nil, response.NotSupportError
}
