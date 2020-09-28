package extend

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// AppLED app led control
type AppLED struct {
	model string
	sf    *os.File // sensorflow app led device file
}

// Prepare prepare
func (al *AppLED) Prepare(model string) error {
	al.model = model

	switch model {
	case sys.ModelSensorflow:
	}

	return nil
}

// CleanUp clean up
func (al *AppLED) CleanUp() error {
	switch al.model {
	case sys.ModelSensorflow:
	}

	return nil
}

// SetLEDStatus set status
func (al *AppLED) SetLEDStatus(status int) error {
	switch al.model {
	case sys.ModelHMU2000:
		return al.setHMU2000APPLED(status)
	case sys.ModelSensorflow:
		return al.setSensorflowAPPLED(status)
	}

	return fmt.Errorf("unknown device")
}

// SetHMU2000APPLED set hmu2000 app led
func (al *AppLED) setHMU2000APPLED(status int) error {
	cfg := sys.GetBusManagerCfg()
	// address := host + ":" + port
	// address := cfg.SystemServer.Uri

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	if _, err := client.SetAppLED(status); err != nil {
		return errors.As(err, status)
	}

	return client.Disconnect()
}

// SetSensorflowAPPLED set sensorflow app led
func (al *AppLED) setSensorflowAPPLED(status int) error {
	cmd := exec.Command("/usr/bin/aggregation/appled.sh", strconv.Itoa(status))
	return cmd.Run()
}
