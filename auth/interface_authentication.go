package auth

import (
	"github.com/IBAX-io/go-ibax-sdk/base"
	"github.com/IBAX-io/go-ibax-sdk/config"
	"github.com/IBAX-io/go-ibax-sdk/response"
)

// Authentication
// Auth IBAX Functions
type Authentication interface {
	// GetUid
	// @method get
	// Get the temporary token for login and save prepare config
	GetUid() error
	// GetUID
	// @method get
	// @return uid string "signature line"
	// @return token string "temporary token to pass in login. Currently, the lifetime of a temporary token is 5 seconds"
	GetUID() (*response.GetUIDResult, error)
	// Login
	// @method post
	// Log in to the specified role and save the token
	Login(roleId string) (err error)
	// LOGIN
	// @method post
	// Log in to the specified role and return login information
	LOGIN(roleId string) (*response.LoginResult, error)
	// AutoLogin
	// no role auto login
	AutoLogin() error
}

type authClient struct {
	Config     *config.IbaxConfig `yaml:"chain_sdk"`
	baseClient base.Base
}

func NewClient(config *config.IbaxConfig, b base.Base) Authentication {
	return &authClient{Config: config, baseClient: b}
}
