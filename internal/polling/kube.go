package polling

import (
	"bigagent/internal/result"
)

type KubePolling struct {
}

func NewKubePolling() *KubePolling {
	return &KubePolling{}
}

func (k *KubePolling) RunPolling(fn PollFunc) result.Result {
	return fn()
}
