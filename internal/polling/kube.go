package polling

import (
	"bigagent/internal/check/result"
)

type KubePolling struct {
}

func NewKubePolling() *KubePolling {
	return &KubePolling{}
}

func (k *KubePolling) RunPolling(fn PollFunc) result.Result {
	return fn()
}
