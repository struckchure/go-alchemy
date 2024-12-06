package components

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/struckchure/go-alchemy"
)

type IAuthentication interface {
	IAlchemyComponent

	Login() error
	Register() error
}

type Authentication struct{}

func (a *Authentication) Setup(component string) (func() error, error) {
	component = strings.ToLower(component)

	methods := map[string]func() error{
		"login":    a.Login,
		"register": a.Register,
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

	cfg, err := alchemy.ReadYaml[Config]("alchemy.yaml")
	if err != nil {
		return err
	}

	prevComponentConfig, componentExists := lo.Find(
		cfg.Components,
		func(c Component) bool { return c.Id == "Authentication" },
	)

	moduleName, err := GetModuleName()
	if err != nil {
		return err
	}

	values := map[string]interface{}{
		"User":       true,
		"ModuleName": moduleName,
		"Login":      true,
	}

	if componentExists {
		for _, s := range prevComponentConfig.Services {
			values[s.Id] = true
		}
	}

	var tmpls []GenerateTmplArgs = []GenerateTmplArgs{
		{
			TmplPath:   "_templates/schema.prisma.tmpl",
			OutputPath: "prisma/schema.prisma",
			Values: map[string]interface{}{
				"User": true,
			},
		},
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

	componentConfig := Component{
		Id: "Authentication",
		Models: []Dependency{
			{
				Id:   "User",
				Path: "prisma/schema.prisma",
			},
			{
				Id:   "UserDao",
				Path: "dao/user_dao.go",
			},
		},
		Services: []Dependency{
			{
				Id:   "Login",
				Path: "services/authentication_service.go",
			},
		},
	}

	if componentExists {
		componentConfig.Models = lo.Uniq(append(componentConfig.Models, prevComponentConfig.Models...))
		componentConfig.Services = lo.Uniq(append(componentConfig.Services, prevComponentConfig.Services...))

		_, idx, ok := lo.FindIndexOf(
			cfg.Components,
			func(c Component) bool { return c.Id == "Authentication" },
		)
		if ok {
			cfg.Components[idx] = componentConfig
		}
	} else {
		cfg.Components = append(cfg.Components, componentConfig)
	}

	err = alchemy.WriteYaml("alchemy.yaml", cfg)
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

	cfg, err := alchemy.ReadYaml[Config]("alchemy.yaml")
	if err != nil {
		return err
	}

	prevComponentConfig, componentExists := lo.Find(
		cfg.Components,
		func(c Component) bool { return c.Id == "Authentication" },
	)

	moduleName, err := GetModuleName()
	if err != nil {
		return err
	}

	values := map[string]interface{}{
		"User":       true,
		"ModuleName": moduleName,
		"Register":   true,
	}

	if componentExists {
		for _, s := range prevComponentConfig.Services {
			values[s.Id] = true
		}
	}

	var tmpls []GenerateTmplArgs = []GenerateTmplArgs{
		{
			TmplPath:   "_templates/schema.prisma.tmpl",
			OutputPath: "prisma/schema.prisma",
			Values:     values,
		},
		{
			TmplPath:   "_templates/prisma_user_dao.go.tmpl",
			OutputPath: "dao/user_dao.go",
			GoFormat:   true,
			Values:     values,
		},
		{
			TmplPath:   "_templates/authentication_service.go.tmpl",
			OutputPath: "services/authentication_service.go",
			GoFormat:   true,
			Values:     values,
		},
	}

	for _, tmpl := range tmpls {
		err := GenerateTmpl(tmpl)
		if err != nil {
			return err
		}

		color.Green("  + %s", tmpl.OutputPath)
	}

	componentConfig := Component{
		Id: "Authentication",
		Models: []Dependency{
			{
				Id:   "User",
				Path: "prisma/schema.prisma",
			},
			{
				Id:   "UserDao",
				Path: "dao/user_dao.go",
			},
		},
		Services: []Dependency{
			{
				Id:   "Register",
				Path: "services/authentication_service.go",
			},
		},
	}

	if componentExists {
		componentConfig.Models = lo.Uniq(append(componentConfig.Models, prevComponentConfig.Models...))
		componentConfig.Services = lo.Uniq(append(componentConfig.Services, prevComponentConfig.Services...))

		_, idx, ok := lo.FindIndexOf(
			cfg.Components,
			func(c Component) bool { return c.Id == "Authentication" },
		)
		if ok {
			cfg.Components[idx] = componentConfig
		}
	} else {
		cfg.Components = append(cfg.Components, componentConfig)
	}

	err = alchemy.WriteYaml("alchemy.yaml", cfg)
	if err != nil {
		return err
	}

	return a.PostSetup()
}

func NewAuthentication() IAuthentication {
	return &Authentication{}
}
