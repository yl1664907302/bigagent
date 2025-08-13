package polling

type KubePolling struct {
	P polling
}

func NewKubePolling(p polling) *KubePolling {
	return &KubePolling{P: p}
}
