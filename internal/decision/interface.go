package decision

type decision interface {
	Recovery(bool2 bool) error
	Judge() func()
}
