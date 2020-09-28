/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/11/21
 * Despcription: test file
 *
 */

package portmanager_test

import (
	"encoding/json"
	"testing"
	"time"

	"clc.hmu/app/public"

	"clc.hmu/app/portmanager/src/protocol"
)

func TestFaceIPCClient(t *testing.T) {
	host := "192.168.10.22"
	port := "20020"
	c, err := protocol.NewFaceIPCClient(host, port)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 15)

	var p public.FaceIPCOperationPayload
	p.CameraID = "0"
	p.FaceID = "ewven"
	p.FaceURL = "http://192.168.10.1/ewven_cj.jpg"
	p.Perpose = "register"

	d, _ := json.Marshal(p)

	c.Command(string(d))
}
