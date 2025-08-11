package request

import (
	"bigagent/internal/service"
	"net/http"
)

type PostK8s struct {
	k        bool
	h        string
	KeyFirst bool
	c        *http.Client
}

func NewPostK8s(host string, keyFirst bool) *PostK8s {
	return &PostK8s{h: host, c: &http.Client{}, KeyFirst: keyFirst}
}

func (p *PostK8s) Do() (interface{}, error) {
	// 这里的占位逻辑与 veops 保持一致：进行一次本地数据刷新或校验
	service.CheckClient(p.h)
	if p.KeyFirst {
		service.CreateMachineUUID()
		p.KeyFirst = false
	}
	service.UpdateMachineData()
	return nil, nil
}
