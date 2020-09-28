package busmanager_test

import (
	"testing"

	"clc.hmu/app/busmanager/src/config"
	"clc.hmu/app/busmanager/src/module"
)

func TestStartRouter(t *testing.T) {
	filename := "busmanager.json"
	module.BusConfigFilePath = filename
	config.OpenBusConfigFile(filename)
	module.StartRouter()
}
