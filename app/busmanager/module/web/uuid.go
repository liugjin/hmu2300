package web

import "clc.hmu/app/public/sys"

// get uuid
func getUUID() string {
	// read id from config file
	return sys.GetMonitoringUnitCfg().ID
}
