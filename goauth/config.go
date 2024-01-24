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
}

type Path struct {
	Url    string
	Method string
	Type   string
	Desc   string
	Code   string
}

func LoadConfig() GoAuthConfig {
	var loaded LoadedGoAuthConfig
	miso.UnmarshalFromProp(&loaded)
	return loaded.Config
}
