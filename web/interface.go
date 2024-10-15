package web

// ApiStrategy 定义api接口策略
type ApiStrategy interface {
	Api(key string) (interface{}, error)
}

// PushStrategy 定义推送接口策略
type PushStrategy interface {
	Push() (interface{}, error)
}

type Agent struct {
	apiStrategy  ApiStrategy
	pushStrategy PushStrategy
}

func NewAgent() Agent {
	return Agent{}
}

func (a *Agent) SetApiStrategy(strategy ApiStrategy) {
	a.apiStrategy = strategy
}

func (a *Agent) SetPushStrategy(strategy PushStrategy) {
	a.pushStrategy = strategy
}

func (a *Agent) ExecuteApi(key string) (interface{}, error) {
	return a.apiStrategy.Api(key)
}

func (a *Agent) ExecutePush() (interface{}, error) {
	return a.pushStrategy.Push()
}
