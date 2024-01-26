package goauth

import "github.com/curtisnewbie/miso/miso"

type LoadedGoAuthConfig struct {
	Config GoAuthConfig `mapstructure:"goauth"`
}

type GoAuthConfig struct {
	Resource []Resource `mapstructure:"resource"`
	Path     []Path     `mapstructure:"path"`
}

type Resource struct {
	Name string
	Code string
	Path []Path `mapstructure:"path"`
}

type Path struct {
	Url  string
	Type string
	Desc string
}

func LoadConfig() GoAuthConfig {
	var loaded LoadedGoAuthConfig
	miso.UnmarshalFromProp(&loaded)
	for i, p := range loaded.Config.Path {
		if p.Type == "" {
			p.Type = "PROTECTED"
			loaded.Config.Path[i] = p
		}
	}

	return loaded.Config
}
