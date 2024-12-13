package components

import "github.com/struckchure/go-alchemy/internals"

var prismaTmpls []GenerateSingleTmplArgs = []GenerateSingleTmplArgs{
	{
		Id:         "Models.User",
		TmplPath:   "prisma/schema.prisma",
		OutputPath: "prisma/schema.prisma",
	},
	{
		Id:         "Models.UserDao",
		TmplPath:   "orms/prisma/user.go",
		OutputPath: "dao/user.go",
		GoFormat:   true,
	},
}

var gormTmpls []GenerateSingleTmplArgs = []GenerateSingleTmplArgs{
	{
		Id:         "Models.UserDao",
		TmplPath:   "orms/gorm/user.go",
		OutputPath: "dao/user.go",
		GoFormat:   true,
	},
	{
		Id:         "Models.Utils",
		TmplPath:   "orms/gorm/utils.go",
		OutputPath: "dao/utils.go",
		GoFormat:   true,
	},
}

var sharedTmpls []GenerateSingleTmplArgs = []GenerateSingleTmplArgs{
	{
		Id:         "Services.Utils",
		TmplPath:   "services/utils.go",
		OutputPath: "services/utils.go",
		GoFormat:   true,
	},
	{
		Id:         "Services.Jwt",
		TmplPath:   "services/jwt.go",
		OutputPath: "services/jwt.go",
		GoFormat:   true,
	},
}

var loginTmpls []GenerateSingleTmplArgs = append([]GenerateSingleTmplArgs{
	{
		Id:         "Services.Login",
		TmplPath:   "services/authentication.go",
		OutputPath: "services/authentication.go",
		GoFormat:   true,
	},
}, sharedTmpls...)

var registerTmpls []GenerateSingleTmplArgs = append([]GenerateSingleTmplArgs{
	{
		Id:         "Services.Register",
		TmplPath:   "services/authentication.go",
		OutputPath: "services/authentication.go",
		GoFormat:   true,
	},
}, sharedTmpls...)

func GetLoginTemplates() ([]GenerateSingleTmplArgs, error) {
	cfg, err := internals.ReadYaml[Config]("alchemy.yaml")
	if err != nil {
		return nil, err
	}

	switch cfg.Orm.Name {
	case "Prisma":
		loginTmpls = append(loginTmpls, prismaTmpls...)
	case "Gorm":
		loginTmpls = append(loginTmpls, gormTmpls...)
	}

	return loginTmpls, nil
}

func GetRegisterTemplates() ([]GenerateSingleTmplArgs, error) {
	cfg, err := internals.ReadYaml[Config]("alchemy.yaml")
	if err != nil {
		return nil, err
	}

	switch cfg.Orm.Name {
	case "Prisma":
		registerTmpls = append(registerTmpls, prismaTmpls...)
	case "Gorm":
		registerTmpls = append(registerTmpls, gormTmpls...)
	}

	return loginTmpls, nil
}
