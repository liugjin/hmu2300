/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/25
 * Despcription: test file
 *
 */

package portmanager_test

import (
	"encoding/json"
	"testing"

	"clc.hmu/app/portmanager/src/protocol"
	"clc.hmu/app/public"
)

func TestOilMachineClient(t *testing.T) {
	port := "/dev/ttyS0"
	baudrate := 9600
	timeout := 3000
	soi := byte(0x7E)
	eoi := byte(0x0D)
	client, err := protocol.NewOilMachineClient(port, baudrate, timeout, soi, eoi)
	if err != nil {
		t.Error(err)
	}

	req := public.OilMachineOperationPayload{
		ADR:    0x60,
		CID1:   0x42,
		CID2:   0x45,
		LENGTH: 0,
	}

	bytereq, _ := json.Marshal(req)

	resp, err := client.Sample(string(bytereq))
	if err != nil {
		t.Error(err)
	}

	br := []byte(resp)
	t.Logf("%X, %d", br, len(br))
}
