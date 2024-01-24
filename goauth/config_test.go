package goauth

import (
	"testing"

	"github.com/curtisnewbie/miso/miso"
)

func TestLoadConfig(t *testing.T) {
	rail := miso.EmptyRail()
	miso.LoadConfigFromFile("../conf.yml", rail)
	config := LoadConfig()
	t.Logf("%+v", config)
}
