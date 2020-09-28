/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/10/30
 * Despcription: test file
 *
 */

package portmanager_test

import (
	"fmt"
	"testing"
	"time"

	"clc.hmu/app/portmanager/src/protocol"
)

func TestLuMiGateway(t *testing.T) {
	sid := "7811dcb78b2f"
	password := "3EA66EEB96434CBB"
	ift := "ens33"
	c, err := protocol.NewLuMiGatewayClient(sid, password, ift)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 5)

	fmt.Println(c.Sample(""))
	fmt.Println(c.Sample("plug"))
	fmt.Println(c.Sample("magnet"))
	fmt.Println(c.Sample("motion"))
	fmt.Println(c.Sample("switch"))
	fmt.Println(c.Sample("smoke"))
	fmt.Println(c.Sample("weather.v1"))

	fmt.Println(c.Command(`{"model":"plug","sid":"158d0002325aa7","value":"{\"status\":\"on\"}"}`))
	time.Sleep(time.Second * 20)
}
