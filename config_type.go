package alchemy

type Component struct {
	Id           string   `yaml:"Id"`
	Path         string   `yaml:"Path"`
	Requirements []string `yaml:"Requirements"`
}

type Orm struct {
	Name             string `yaml:"Name"`
	DatabaseProvider string `yaml:"DatabaseProvider"`
}

type Config struct {
	ProjectName string      `yaml:"ProjectName"`
	Root        string      `yaml:"Root"`
	Orm         Orm         `yaml:"Orm"`
	Components  []Component `yaml:"Components"`
}
