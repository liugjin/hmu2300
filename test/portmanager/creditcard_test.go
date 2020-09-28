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
	"testing"

	"clc.hmu/app/portmanager/src/protocol"
)

func TestCreditCard(t *testing.T) {
	username := "admin"
	password := "12345"
	serialnum := "220793197"

	c, err := protocol.NewCreditCardClient(username, password, serialnum)
	if err != nil {
		t.Error(err)
	}

	t.Log(c.Sample(""))
}
