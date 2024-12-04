package components

var ComponentCategoryOptions []string = []string{
	"Authentication",
	"Authorization",
	"Products",
	"Orders",
	"Media",
}

var ComponentMapping map[string][]string = map[string][]string{
	"Authentication": AuthenticationOptions,
	"Authorization":  AuthorizationOptions,
	"Products":       ProductsOptions,
	"Orders":         OrdersOptions,
	"Media":          MediaOptions,
}

var AuthenticationOptions []string = []string{
	"Login",
	"Register",
}

var AuthorizationOptions []string = []string{
	"RoleBaseAccessControl",
	"AttributeBaseAccessControl",
}

var ProductsOptions []string = []string{}

var OrdersOptions []string = []string{}

var MediaOptions []string = []string{}
