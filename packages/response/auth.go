package response

type GetUIDResult struct {
	UID         string `json:"uid,omitempty"`
	Token       string `json:"token,omitempty"`
	Expire      string `json:"expire,omitempty"`
	EcosystemID string `json:"ecosystem_id,omitempty"`
	KeyID       string `json:"key_id,omitempty"`
	Address     string `json:"address,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
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
