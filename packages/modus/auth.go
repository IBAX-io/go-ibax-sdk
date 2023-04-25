package modus

import "github.com/IBAX-io/go-ibax-sdk/packages/response"

type Authentication interface {
	// GetUid
	// @method get
	// @return uid string "signature line"
	// @return token string "temporary token to pass in login. Currently, the lifetime of a temporary token is 5 seconds"
	// Get the temporary token for login and save prepare config
	GetUid() error
	// Login
	// @method post
	// Log in to the specified role and save the token
	Login(roleId int64) (err error)
	// AutoLogin
	// no role auto login
	AutoLogin() error
	GetAuthStatus() (*response.AuthStatusResponse, error)
}
