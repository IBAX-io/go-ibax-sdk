package response

type ListResult struct {
	Count int64               `json:"count"`
	List  []map[string]string `json:"list"`
}

type RowResult struct {
	Value map[string]string `json:"value"`
}

type PageResult struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Menu          string `json:"menu"`
	ValidateCount int64  `json:"nodesCount"`
	AppID         int64  `json:"app_id"`
	Conditions    string `json:"conditions"`
}

type MenuResult struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Value      string `json:"value"`
	Conditions string `json:"conditions"`
}

type SnippetResult struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Conditions string `json:"conditions"`
}

type NameInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AppContentResult struct {
	Snippets  []NameInfo `json:"snippets"`
	Pages     []NameInfo `json:"pages"`
	Contracts []NameInfo `json:"contracts"`
}

type HistoryResult struct {
	List []map[string]string `json:"list"`
}

type tableInfo struct {
	Name  string `json:"name"`
	Count string `json:"count"`
}

type TablesResult struct {
	Count int64       `json:"count"`
	List  []tableInfo `json:"list"`
}

type columnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Perm string `json:"perm"`
}

type TableResult struct {
	Name       string       `json:"name"`
	Insert     string       `json:"insert"`
	NewColumn  string       `json:"new_column"`
	Update     string       `json:"update"`
	Read       string       `json:"read,omitempty"`
	Filter     string       `json:"filter,omitempty"`
	Conditions string       `json:"conditions"`
	AppID      string       `json:"app_id"`
	Columns    []columnInfo `json:"columns"`
}
