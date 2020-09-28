/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: test file
 *
 */

package portmanager_test

import (
	"testing"

	"clc.hmu/app/portmanager/src/protocol"
	"clc.hmu/app/public"
)

func TestSystemSample(t *testing.T) {
	req := public.SystemBindingPayload{Host: "192.168.10.1", Port: "9988"}
	client, err := protocol.NewSystemClient(req)
	if err != nil {
		t.Errorf("new system client failed, errmsg {%v}", err)
	}

	var pc protocol.PortClient
	pc = client

	payload := []string{
		`{"type": "gps"}`,
		`{"type": "time"}`,
		`{"type": "system"}`,
		`{"type": "uuid"}`,
		`{"type": "internet"}`,
	}

	for _, p := range payload {
		result, err := pc.Sample(p)
		if err != nil {
			t.Errorf("sample failed, errmsg {%v}", err)
		}

		t.Log(result)
	}
}

func TestSystemCommand(t *testing.T) {
	req := public.SystemBindingPayload{Host: "192.168.10.1", Port: "9988"}
	client, err := protocol.NewSystemClient(req)
	if err != nil {
		t.Errorf("new system client failed, errmsg {%v}", err)
	}

	var pc protocol.PortClient
	pc = client

	payload := []string{
		`{"type": "time"}`,
		`{"type": "system"}`,
		`{"type": "uuid"}`,
		`{"type": "internet"}`,
	}

	for _, p := range payload {
		result, err := pc.Command(p)
		if err != nil {
			t.Errorf("sample failed, errmsg {%v}", err)
		}

		t.Log(result)
	}
}
