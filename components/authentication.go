package components

import (
	"errors"

	"github.com/fatih/color"
	"github.com/samber/lo"
)

type IAuthentication interface {
	IAlchemyComponent

	Login() error
	Register() error
}

type Authentication struct{}

func (a *Authentication) Setup(component string) (func() error, error) {
	methods := map[string]func() error{
		"login":    a.Login,
		"register": a.Register,
	}

	if !lo.HasKey(methods, component) {
		return nil, errors.New("component is not available")
	}

	return methods[component], nil
}

func (a *Authentication) Login() (err error) {
	defer func() {
		if err == nil {
			color.Green("+ Authentication.Login")
		} else {
			color.Red("x Authentication.Login")
		}
	}()

	return err
}

func (a *Authentication) Register() (err error) {
	defer func() {
		if err == nil {
			color.Green("+ Authentication.Register")
		} else {
			color.Red("x Authentication.Register")
		}
	}()

	return err
}

func NewAuthentication() IAuthentication {
	return &Authentication{}
}
