package decision

type decision interface {
	Recovery() error
	Judge() (interface{}, error)
}
