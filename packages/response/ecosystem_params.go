package response

type ParamResult struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Conditions string `json:"conditions"`
}

type ParamsResult struct {
	List []ParamResult `json:"list"`
}

type EcosystemInfo struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Digits       int64  `json:"digits"`
	TokenSymbol  string `json:"token_symbol"`
	TokenName    string `json:"token_name"`
	TotalAmount  string `json:"total_amount"`
	IsWithdraw   bool   `json:"is_withdraw"`
	Withdraw     string `json:"withdraw"`
	IsEmission   bool   `json:"is_emission"`
	Emission     string `json:"emission"`
	Introduction string `json:"introduction"`
	Logo         int64  `json:"logo"`
	Creator      string `json:"creator"`
}

type EcosystemNameResult struct {
	EcosystemName string `json:"ecosystem_name"`
}
