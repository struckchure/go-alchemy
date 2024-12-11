package components

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/struckchure/go-alchemy/internals"
)

type IAuthentication interface {
	IAlchemyComponent

	Login() error
	Register() error
}

type Authentication struct{}

func (a *Authentication) Setup(component string) (func() error, error) {
	methods := map[string]func() error{
		"Login":    a.Login,
		"Register": a.Register,
	}

	if !lo.HasKey(methods, component) {
		return nil, fmt.Errorf("component `%s` is not available", component)
	}

	err := a.PreSetup()
	if err != nil {
		return nil, err
	}

	return methods[component], nil
}

func (a *Authentication) PreSetup() error {
	if !internals.FileExists("alchemy.yaml") {
		return internals.ErrAlchemyConfigNotFound
	}

	cmd := exec.Command("go", "get", "github.com/steebchen/prisma-client-go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(err, errors.New(string(out)))
	}

	return nil
}

func (a *Authentication) PostSetup() error {
	cmd := exec.Command("go", "run", "github.com/steebchen/prisma-client-go", "db", "push")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(err, errors.New(string(out)))
	}

	cmd = exec.Command("go", "mod", "tidy")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(err, errors.New(string(out)))
	}

	return nil
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

	err = GenerateMultipleTmpls(GenerateMultipleTmplsArgs{
		ComponentId: strings.Split(componentId, ".")[0],
		Tmpls:       LoginTmpls,
		Values: map[string]interface{}{
			"Login":      true,
			"User":       true,
			"ModuleName": moduleName,
		},
	})
	if err != nil {
		return err
	}

	return a.PostSetup()
}

func (a *Authentication) Register() (err error) {
	componentId := "Authentication.Register"

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

	err = GenerateMultipleTmpls(GenerateMultipleTmplsArgs{
		ComponentId: strings.Split(componentId, ".")[0],
		Tmpls:       RegisterTmpls,
		Values: map[string]interface{}{
			"Register":   true,
			"User":       true,
			"ModuleName": moduleName,
		},
	})
	if err != nil {
		return err
	}

	return a.PostSetup()
}

func NewAuthentication() IAuthentication {
	return &Authentication{}
}
