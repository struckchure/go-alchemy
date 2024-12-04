package internals

var ComponentCategories []string = []string{
	"Authentication",
	"Authorization",
	"Products",
	"Orders",
	"Media",
}

var ComponentMapping map[string][]string = map[string][]string{
	"Authentication": Authentication,
	"Authorization":  Authorization,
	"Products":       Products,
	"Orders":         Orders,
	"Media":          Media,
}

var Authentication []string = []string{"Login", "Register"}

var Authorization []string = []string{"RoleBaseAccessControl", "AttributeBaseAccessControl"}

var Products []string = []string{}

var Orders []string = []string{}

var Media []string = []string{}
