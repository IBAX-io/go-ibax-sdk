package response

type roleInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type notifyInfo struct {
	RoleID string `json:"role_id"`
	Count  int64  `json:"count"`
}

type keyEcosystemInfo struct {
	Ecosystem     string       `json:"ecosystem"`
	Name          string       `json:"name"`
	Digits        int64        `json:"digits"`
	Roles         []roleInfo   `json:"roles,omitempty"`
	Notifications []notifyInfo `json:"notifications,omitempty"`
}

type KeyInfoResult struct {
	Account    string              `json:"account"`
	Ecosystems []*keyEcosystemInfo `json:"ecosystems"`
}

type TokenBalanceResult struct {
	Amount      string `json:"amount"`
	Digits      int64  `json:"digits"`
	Total       string `json:"total"`
	Utxo        string `json:"utxo"`
	TokenSymbol string `json:"token_symbol"`
}

type MemberInfo struct {
	ID         int64  `json:"id"`
	MemberName string `json:"member_name"`
	ImageID    *int64 `json:"image_id"`
	MemberInfo string `json:"member_info"`
}
