/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: test file
 *
 */
package appnet_test

import (
	"log"
	"testing"
	"time"

	"clc.hmu/app/appmanager/appnet"
)

func TestConcurrentSample(t *testing.T) {
	appnet.ConnectPortManager("192.168.0.2:50051")
	port := "/run/com1"
	code := 3

	quantity := 1
	timeout := 3000

	// go func() {
	// 	slaveid := 1
	// 	address := 0
	// 	for {
	// 		result, err := appnet.Sample(client, port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
	// 		if err != nil {
	// 			log.Printf("sample failed, errmsg [%s]", err)
	// 		}

	// 		log.Println(slaveid, address, result)
	// 		time.Sleep(time.Second * 3)
	// 	}
	// }()

	// go func() {
	// 	slaveid := 1
	// 	address := 1
	// 	for {
	// 		result, err := appnet.Sample(client, port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
	// 		if err != nil {
	// 			log.Printf("sample failed, errmsg [%s]", err)
	// 		}

	// 		log.Println(slaveid, address, result)
	// 		time.Sleep(time.Second * 3)
	// 	}
	// }()

	// go func() {
	// 	slaveid := 2
	// 	address := 0
	// 	for {
	// 		result, err := appnet.Sample(client, port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
	// 		if err != nil {
	// 			log.Printf("sample failed, errmsg [%s]", err)
	// 		}

	// 		log.Println(slaveid, address, result)
	// 		time.Sleep(time.Second * 3)
	// 	}
	// }()

	// go func() {
	// 	slaveid := 2
	// 	address := 1
	// 	for {
	// 		result, err := appnet.Sample(client, port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
	// 		if err != nil {
	// 			log.Printf("sample failed, errmsg [%s]", err)
	// 		}

	// 		log.Println(slaveid, address, result)
	// 		time.Sleep(time.Second * 3)
	// 	}
	// }()

	for {
		time.Sleep(time.Second * 3)

		slaveid := 1
		address := 0
		result, err := appnet.Sample(port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
		if err != nil {
			log.Printf("sample failed, errmsg [%s]", err)
		}

		log.Println(slaveid, address, result)

		slaveid = 1
		address = 1
		result, err = appnet.Sample(port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
		if err != nil {
			log.Printf("sample failed, errmsg [%s]", err)
		}

		log.Println(slaveid, address, result)

		slaveid = 2
		address = 0
		result, err = appnet.Sample(port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
		if err != nil {
			log.Printf("sample failed, errmsg [%s]", err)
		}

		log.Println(slaveid, address, result)

		slaveid = 2
		address = 1
		result, err = appnet.Sample(port, int32(code), int32(slaveid), int32(address), int32(quantity), int32(timeout))
		if err != nil {
			log.Printf("sample failed, errmsg [%s]", err)
		}
		log.Println(slaveid, address, result)
	}
}
