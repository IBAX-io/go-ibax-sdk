package config

// Config
// If you want to modify the configuration, you need to pay attention to the multi-threaded call problem, it is recommended to use the SetConfig() method
type Config struct {
	PrivateKey      string `json:"private_key" yaml:"private_key"` // private key.Do not use clear text. You can set the environment variable. The key controls access to your funds!
	PublicKey       []byte `json:"public_key" yaml:"-"`            // public key
	Ecosystem       int64  `json:"ecosystem" yaml:"ecosystem"`     // Login ecosystem Id.
	ApiAddress      string `json:"api_address" yaml:"api_address"` // api address. Restful Api address or RPC api address. depends on enable_rpc
	ApiPath         string `json:"api_path" yaml:"api_path"`       // api path. Restful Api Path
	JwtPrefix       string `json:"jwt_prefix" yaml:"jwt_prefix"`   // jwt prefix
	Token           string `json:"token" yaml:"-"`                 // token
	TokenExpireTime int64  `json:"token_expire_time" yaml:"-"`     // token expire time
	KeyId           int64  `json:"key_id" yaml:"-"`                // key id
	Account         string `json:"account" yaml:"-"`               // account address
	UID             string `json:"uid" yaml:"-"`                   // uid
	NetworkId       int64  `json:"network_id" yaml:"-"`            // network id
	Cryptoer        string `json:"cryptoer" yaml:"cryptoer"`       // cryptoer
	Hasher          string `json:"hasher" yaml:"hasher"`           // hasher crypto
	EnableRpc       bool   `json:"enable_rpc" yaml:"enable_rpc"`   // enable rpc
}

// Version SDK Version
const Version = "1.1.0"
