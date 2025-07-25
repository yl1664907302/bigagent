package request

import (
	"bigagent/internal/service"
	"net/http"
)

type PostVeops struct {
	k        bool
	h        string
	KeyFirst bool
	c        *http.Client
}

func NewPostVeops(host string, keyFirst bool) *PostVeops {
	return &PostVeops{h: host, c: &http.Client{}, KeyFirst: keyFirst}
}

func (p *PostVeops) Do() (interface{}, error) {
	service.CheckClient(p.h)
	if p.KeyFirst {
		service.CreateMachineUUID()
		p.KeyFirst = false
	}
	service.UpdateMachineData()
	return nil, nil
}
