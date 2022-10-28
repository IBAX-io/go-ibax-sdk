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

type CountResult struct {
	Count int64 `json:"count"`
}

type EcosystemNameResult struct {
	EcosystemName string `json:"ecosystem_name"`
}
