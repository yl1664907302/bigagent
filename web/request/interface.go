package request

type RequestInterface interface {
	do() (interface{}, error)
}
