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
	"time"

	"clc.hmu/app/busmanager/src/module"
)

func TestPublish(t *testing.T) {
	var mqttclient module.MQTTClient
	if err := mqttclient.PublishSampleValues("sample-values/test", `{"muid": "testid"}`); err == nil {
		t.Errorf("mqtt client should be unavaliable to publish message before be new")
	}

	mqttclient = module.NewMQTTClient(nil, "", "")

	if err := mqttclient.ConnectServer(); err != nil {
		t.Errorf("connect mqtt server failed, errmsg {%v}", err)
	}
	defer mqttclient.DisconnectServer()

	if err := mqttclient.PublishSampleValues("sample-values/test", `{"muid": "testid"}`); err != nil {
		t.Errorf("publish sample values failed, errmsg {%v}", err)
	}
}

func TestSubscribe(t *testing.T) {
	var mqttclient module.MQTTClient
	if err := mqttclient.Subscribe("command/#"); err == nil {
		t.Errorf("mqtt client should be unavaliable to subsribe message before be new")
	}

	mqttclient = module.NewMQTTClient(module.SubMessageHandler, "", "")

	if err := mqttclient.ConnectServer(); err != nil {
		t.Errorf("connect mqtt server failed, errmsg {%v}", err)
	}
	defer mqttclient.DisconnectServer()

	// go func() {
	// 	if err := mqttclient.PublishSampleValues("sample-values/test", `{"muid": "testid"}`); err != nil {
	// 		t.Errorf("publish sample values failed, errmsg {%v}", err)
	// 	}
	// }()

	// module.ConnectAppServer()

	if err := mqttclient.Subscribe("command/#"); err != nil {
		t.Errorf("subscribe failed, errmsg {%v}", err)
	}

	if err := mqttclient.Subscribe("test/#"); err != nil {
		t.Errorf("subscribe failed, errmsg {%v}", err)
	}

	time.Sleep(time.Second * 5)
}
