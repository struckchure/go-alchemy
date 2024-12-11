package components

var sharedTmpls []GenerateSingleTmplArgs = []GenerateSingleTmplArgs{
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

var LoginTmpls []GenerateSingleTmplArgs = append([]GenerateSingleTmplArgs{
	{
		Id:         "Services.Login",
		TmplPath:   "services/authentication.go",
		OutputPath: "services/authentication.go",
		GoFormat:   true,
	},
}, sharedTmpls...)

var RegisterTmpls []GenerateSingleTmplArgs = append([]GenerateSingleTmplArgs{
	{
		Id:         "Services.Register",
		TmplPath:   "services/authentication.go",
		OutputPath: "services/authentication.go",
		GoFormat:   true,
	},
}, sharedTmpls...)
