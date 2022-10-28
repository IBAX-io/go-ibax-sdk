package response

type AppParamsResult struct {
	App  string        `json:"app_id"`
	List []ParamResult `json:"list"`
}
