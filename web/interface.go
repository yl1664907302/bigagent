package web

// ApiStrategy 定义api接口策略
type ApiStrategy interface {
	Api() (interface{}, error)
}

// PushStrategy 定义推送接口策略
type PushStrategy interface {
	Push() (interface{}, error)
}

type Agent struct {
	apiStrategy  ApiStrategy
	pushStrategy PushStrategy
}

func (a *Agent) SetApiStrategy(strategy ApiStrategy) {
	a.apiStrategy = strategy
}

func (a *Agent) SetPushStrategy(strategy PushStrategy) {
	a.pushStrategy = strategy
}

func (a *Agent) ExecuteApi() (interface{}, error) {
	data, err := a.apiStrategy.Api()
	if err != nil {
		return nil, err
	}
	return data, err
}

func (a *Agent) ExecutePush() (interface{}, error) {
	data, err := a.apiStrategy.Api()
	if err != nil {
		return nil, err
	}
	return data, err
}
