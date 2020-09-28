package sys

import (
	"fmt"
	"testing"
)

func TestHmu2000Net(t *testing.T) {
	conn, err := ConnectSystemDaemon(MODEL_HMU2000, &SystemServerOption{
		Uri:  "192.168.20.1:9988",
		Vals: "at_file=/dev/ttyUSB5&at_timeout=60000",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Disconnect()

	resp, err := conn.AutoCheckNetworking([]string{
		"baidu.com:80",
		"lab.huayuan-iot.com:1883",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", *resp)

}

func TestHmu2000Func(t *testing.T) {
	conn, err := ConnectSystemDaemon(MODEL_HMU2000, &SystemServerOption{
		Uri:  "192.168.20.1:9988",
		Vals: "at_file=/dev/ttyUSB5&at_timeout=60000",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(conn.ModelName())

	gpsInfo, err := conn.GPS()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("gps:%+v\n", gpsInfo)

	timeInfo, err := conn.Time()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("time:%+v\n", timeInfo)

	sysInfo, err := conn.Time()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("sys:%+v\n", sysInfo)

	uuidInfo, err := conn.UUID()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("uuid:%+v\n", uuidInfo)

	internetInfo, err := conn.Internet()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("internet:%+v\n", internetInfo)

	dhcpInfo, err := conn.SetEthDHCP()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("internet:%+v\n", dhcpInfo)

}
