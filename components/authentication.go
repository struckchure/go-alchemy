package components

type IAuthentication interface {
	Login() error
	Register() error
}

type Authentication struct{}

func (a *Authentication) Login() error {
	panic("unimplemented")
}

func (a *Authentication) Register() error {
	panic("unimplemented")
}

func NewAuthentication() IAuthentication {
	return &Authentication{}
}
