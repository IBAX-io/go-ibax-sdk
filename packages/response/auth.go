package response

type GetUIDResult struct {
	UID         string `json:"uid"`
	Token       string `json:"token"`
	Expire      string `json:"expire"`
	EcosystemID string `json:"ecosystem_id"`
	KeyID       string `json:"key_id"`
	Address     string `json:"address"`
	NetworkID   string `json:"network_id"`
	Cryptoer    string `json:"cryptoer"`
	Hasher      string `json:"hasher"`
}

type LoginResult struct {
	Token       string        `json:"token"`
	EcosystemID string        `json:"ecosystem_id"`
	KeyID       string        `json:"key_id"`
	Account     string        `json:"account"`
	NotifyKey   string        `json:"notify_key"`
	IsNode      bool          `json:"isnode"`
	IsOwner     bool          `json:"isowner"`
	IsCLB       bool          `json:"clb"`
	Timestamp   string        `json:"timestamp"`
	Roles       []rolesResult `json:"roles"`
}

type rolesResult struct {
	RoleId   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
}

type AuthStatusResponse struct {
	IsActive  bool  `json:"active"`
	ExpiresAt int64 `json:"exp,omitempty"`
}
