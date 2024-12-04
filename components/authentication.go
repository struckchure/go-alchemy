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
	componentId := "Authentication.Login"
	defer func() {
		if err == nil {
			color.Green("+ %s", componentId)
		} else {
			color.Red("x %s", componentId)
		}
	}()

	color.Green("Creating %s component", componentId)

	moduleName, err := GetModuleName()
	if err != nil {
		return err
	}

	var tmpls []GenerateTmplArgs = []GenerateTmplArgs{
		{
			TmplPath:   "_templates/prisma_user_dao.go.tmpl",
			OutputPath: "dao/user_dao.go",
			GoFormat:   true,
			Values: map[string]interface{}{
				"ModuleName": moduleName,
				"Login":      true,
			},
		},
		{
			TmplPath:   "_templates/authentication_service.go.tmpl",
			OutputPath: "services/authentication_service.go",
			GoFormat:   true,
			Values: map[string]interface{}{
				"ModuleName": moduleName,
				"Login":      true,
			},
		},
	}

	for _, tmpl := range tmpls {
		err := GenerateTmpl(tmpl)
		if err != nil {
			return err
		}

		color.Green("  + %s", tmpl.OutputPath)
	}

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
