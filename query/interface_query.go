package query

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/response"
)

type Query interface {
	// GetKeyInfo
	// @method get
	// get key info
	GetKeyInfo(account string) (*response.KeyInfoResult, error)
	// GetBalance
	// @method get
	GetBalance(wallet string) (*response.TokenBalanceResult, error)
	// Balance
	// @method get
	Balance(wallet string, ecosystem int64) (*response.TokenBalanceResult, error)
	// WalletHistory
	// @method get
	WalletHistory(searchType string, page, limit int64) (*response.WalletHistoryResult, error)
	EcosystemCount() (*response.CountResult, error)
	GetVersion() (*string, error)
	GetCentrifugo() (*string, error)

	//block
	GetMaxBlockID() (*response.MaxBlockResult, error)
	DetailedBlocks(count, blockId int64) (*map[int64]response.BlockDetailedInfo, error)
	GetBlockInfo(id int64) (*response.BlockInfoResult, error)
	BlockTxInfo(count, blockId int64) (*map[int64][]response.TxInfo, error)

	//metrics
	BlocksCount() (*response.CountResult, error)
	TransactionsCount() (*response.CountResult, error)
	EcosystemsCount() (*response.CountResult, error)
	KeysCount() (*response.CountResult, error)
	HonorNodesCount() (*response.CountResult, error)

	//data verify
	DataVerify(tableName string, id int64, column, hash string) (*string, error)
	BinaryVerify(id int64, hash string) (*any, error)

	//ecosystem
	EcosystemName(ecosystem int64) (*response.EcosystemNameResult, error)
	AppParams(appid int64, names string, ecosystem int64) (*response.AppParamsResult, error)
	AppParam(appid int64, name string, ecosystem int64) (*response.ParamResult, error)
	EcosystemParam(ecosystem int64, name string) (*response.ParamResult, error)
	EcosystemParams(ecosystem int64, names string) (*response.ParamsResult, error)
	GetTables(limit, offset int) (*response.TablesResult, error)
	GetTable(tableName string) (*response.TableResult, error)
	GetList(tableName, columns string, limit, offset int64) (*response.ListResult, error)
	GetSections(language string, offset, limit int) (*response.ListResult, error)
	GetRow(tableName string, id int64, rowsName string) (*response.RowResult, error)
	GetRowExtend(tableName, columns string, value string, rowsName string) (*response.RowResult, error)
	SystemParams(names string) (*response.ParamsResult, error)
	GetHistory(table string, id string) (*response.HistoryResult, error)
	GetPageRow(name string) (*response.PageResult, error)
	GetMenuRow(name string) (*response.MenuResult, error)
	GetSnippetRow(name string) (*response.SnippetResult, error)
	GetListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error)
	GetNodeListWhere(tableName string, where any, columns, order string, page, limit int) (*response.ListResult, error)
}

type queryClient struct {
	base.Base
}

func NewClient(b base.Base) Query {
	return &queryClient{Base: b}
}
