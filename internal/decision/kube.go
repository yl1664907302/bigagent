package decision

type KubeDecision struct {
	C check
}

func NewKubeDecision(c check) *KubeDecision {
	return &KubeDecision{C: c}
}

type check interface {
	StartCheck() error
}
