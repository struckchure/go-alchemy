package alchemy

type Dependency struct {
	Id   string `yaml:"Id"`
	Path string `yaml:"Path"`
}

type Component struct {
	Id       string       `yaml:"Id"`
	Path     string       `yaml:"Path"`
	Models   []Dependency `yaml:"Models"`
	Services []Dependency `yaml:"Services"`
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
