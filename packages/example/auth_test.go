package example

import (
	"github.com/IBAX-io/go-ibax-sdk/packages/client"
	"testing"
	"time"
)

func TestIBAX_Refresh(t *testing.T) {
	c := client.NewClient(cnf)
	err := c.AutoLogin()
	if err != nil {
		t.Errorf("auto login failed: %s", err.Error())
		return
	}

	cfg := c.GetConfig()
	//fmt.Println("TokenExpireTime:", cfg.TokenExpireTime)
	//cfg.TokenExpireTime = time.Now().Unix() //set token expire
	//time.Sleep(5 * time.Second)
	//fmt.Println("TokenExpireTime:", time.Unix(cfg.TokenExpireTime, 0).Sub(time.Now()))
	if time.Unix(cfg.TokenExpireTime, 0).Sub(time.Now()) < time.Minute*10 {
		cfg.Token = ""

		err := c.AutoLogin()
		if err != nil {
			t.Errorf("auto login failed: %s", err.Error())
			return
		}
		//fmt.Println("TokenExpireTime:", cfg.TokenExpireTime.String())
	}
}
