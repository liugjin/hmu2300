/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: test file
 *
 */

package busmanager_test

import (
	"testing"

	"clc.hmu/app/public"
)

func TestSystem(t *testing.T) {
	address := "192.168.0.22:9988"

	var client public.SystemClient
	if err := client.ConnectSystemDaemon(address); err != nil {
		t.Fatal(err)
	}

	// resp, err := client.GPS()
	// resp, err := client.Time()
	resp, err := client.SystemInfo()
	// resp, err := client.Reboot()
	// resp, err := client.UUID()
	// resp, err := client.LAN()
	// resp, err := client.WAN()
	// resp, err := client.Wireless()
	// resp, err := client.Internet()
	if err != nil {
		t.Fatal(err)
	}

	// timeserver := "stdtime.gov.hk"
	// timeserver := "0.openwrt.pool.ntp.org 1.openwrt.pool.ntp.org 2.openwrt.pool.ntp.org 3.openwrt.pool.ntp.org"
	// resp, err := client.SetTimeServer(timeserver)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ip := "192.168.0.3"
	// mask := "255.255.255.0"
	// gateway := "192.168.0.1"
	// pdns := "114.114.114.114"
	// sdns := "8.8.8.8"
	// resp, err := client.SetEthStatic(ip, mask, gateway, pdns, sdns)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// resp, err := client.SetEthDHCP()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ssid := "clcdata"
	// key := "clc666888"
	// resp, err := client.SetWifi(ssid, key)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// resp, err := client.SetLTE()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// mode := "all"
	// resp, err := client.FactoryReset(mode)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	t.Logf("%v", resp)
}
