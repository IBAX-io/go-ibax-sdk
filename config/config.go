package config

import (
	"time"
)

type IbaxConfig struct {
	PrivateKey      string        `json:"private_key" yaml:"private_key"` // private key
	PublicKey       []byte        `json:"public_key" yaml:"-"`            // public key
	Ecosystem       int64         `json:"ecosystem" yaml:"ecosystem"`     // ecosystem
	ApiAddress      string        `json:"api_address" yaml:"api_address"` // api address
	IsMobile        bool          `json:"is_mobile" yaml:"is_mobile"`     // is mobile
	ApiPath         string        `json:"api_path" yaml:"api_path"`       // api path
	JwtPrefix       string        `json:"jwt_prefix" yaml:"jwt_prefix"`   // jwt prefix
	Token           string        `json:"token"`                          // token
	TokenExpireTime time.Duration `json:"token_expire_time"`              // token expire time
	Policy          string        `json:"policy" yaml:"policy"`           // address select policy :random,all,better default random
	KeyId           int64         `json:"key_id"`                         // key id
	Account         string        `json:"account"`                        // account address
	UID             string        `json:"uid"`                            // uid
	NetworkId       int64         `json:"network_id"`                     // network id
	Enable          bool          `yaml:"enable"`                         // enable
	Cryptoer        string        `json:"cryptoer" yaml:"cryptoer"`       // cryptoer
	Hasher          string        `json:"hasher" yaml:"hasher"`           // hasher
}
