package internals

var OrmOptions []string = []string{
	"Prisma",
	"Gorm",
}

var OrmMappings map[string][]string = map[string][]string{
	"Prisma": PrismaOptions,
	"Gorm":   GormOptions,
}

var PrismaOptions []string = []string{
	"PostgreSQL",
	"MySQl",
	"SQLite",
	"SQLServer",
	"MongoDB",
	"CockroachDB",
}

var GormOptions []string = []string{
	"PostgreSQL",
	"MySQl",
	"SQLite",
	"SQLServer",
	"Clickhouse",
}
