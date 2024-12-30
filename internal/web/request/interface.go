package request

type RequestInterface interface {
	Do() (interface{}, error)
}
