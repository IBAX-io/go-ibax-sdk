package query

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	response2 "github.com/IBAX-io/go-ibax-sdk/response"
)

type Query interface {
	// GetKeyInfo
	// @method get
	// get key info
	GetKeyInfo(account string) (*response2.KeyInfoResult, error)
	// GetBalance
	// @method get
	GetBalance(wallet string) (*response2.TokenBalanceResult, error)
	// Balance
	// @method get
	Balance(wallet string, ecosystem int64) (*response2.TokenBalanceResult, error)
	// WalletHistory
	// @method get
	WalletHistory(searchType string, page, limit int64) (*response2.WalletHistoryResult, error)
	EcosystemCount() (*response2.CountResult, error)
	GetVersion() (*string, error)
	GetCentrifugo() (*string, error)

	//block
	GetMaxBlockID() (*response2.MaxBlockResult, error)
	DetailedBlocks(count, blockId int64) (*map[int64]response2.BlockDetailedInfo, error)
	GetBlockInfo(id int64) (*response2.BlockInfoResult, error)
	BlockTxInfo(count, blockId int64) (*map[int64][]response2.TxInfo, error)

	//metrics
	BlocksCount() (*response2.CountResult, error)
	TransactionsCount() (*response2.CountResult, error)
	EcosystemsCount() (*response2.CountResult, error)
	KeysCount() (*response2.CountResult, error)
	HonorNodesCount() (*response2.CountResult, error)

	//data verify
	DataVerify(tableName string, id int64, column, hash string) (*string, error)
	BinaryVerify(id int64, hash string) (*any, error)

	//ecosystem
	EcosystemName(ecosystem int64) (*response2.EcosystemNameResult, error)
	AppParams(appid int64, names string, ecosystem int64) (*response2.AppParamsResult, error)
	AppParam(appid int64, name string, ecosystem int64) (*response2.ParamResult, error)
	EcosystemParam(ecosystem int64, name string) (*response2.ParamResult, error)
	EcosystemParams(ecosystem int64, names string) (*response2.ParamsResult, error)
	GetTables(limit, offset int) (*response2.TablesResult, error)
	GetTable(tableName string) (*response2.TableResult, error)
	GetList(tableName, columns string, limit, offset int64) (*response2.ListResult, error)
	GetSections(language string, offset, limit int) (*response2.ListResult, error)
	GetRow(tableName string, id int64, rowsName string) (*response2.RowResult, error)
	GetRowExtend(tableName, columns string, value string, rowsName string) (*response2.RowResult, error)
	SystemParams(names string) (*response2.ParamsResult, error)
	GetHistory(table string, id string) (*response2.HistoryResult, error)
	GetPageRow(name string) (*response2.PageResult, error)
	GetMenuRow(name string) (*response2.MenuResult, error)
	GetSnippetRow(name string) (*response2.SnippetResult, error)
	GetListWhere(tableName string, where any, columns, order string, page, limit int) (*response2.ListResult, error)
	GetNodeListWhere(tableName string, where any, columns, order string, page, limit int) (*response2.ListResult, error)
}

type queryClient struct {
	base.Base
}

func NewClient(b base.Base) Query {
	return &queryClient{Base: b}
}
