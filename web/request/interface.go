package request

type RequestInterface interface {
	do() (interface{}, error)
}

type Post struct {
	do RequestInterface
}
