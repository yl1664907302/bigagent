package model

import (
	"fmt"
)

// 认证策略

type TokenAuth struct {
	Token string `json:"auth_token"`
}

func (t *TokenAuth) ApplyAuth(args ...interface{}) error {
	// 确保参数正确
	if len(args) != 1 {
		return fmt.Errorf("token认证参数不存在")
	}
	t.Token = args[0].(string)
	return nil
}
