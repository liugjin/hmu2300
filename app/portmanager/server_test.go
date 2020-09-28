/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: test file
 *
 */

package portmanager

import (
	"log"
	"testing"

	"clc.hmu/app/portmanager/protocol"
)

func TestNewPortClient(t *testing.T) {
	var ps PortServer
	var mc protocol.ModbusClient
	port := "port"
	mc.SlaveID = 1
	ps.Clients = make(map[string][]protocol.PortClient)
	ps.Clients[port] = append(ps.Clients[port], &mc)

	var mc2 protocol.ModbusClient
	mc2.SlaveID = 2
	ps.Clients[port] = append(ps.Clients[port], &mc2)
	log.Println(ps)
	log.Println(ps.Clients[port][0].ID())
	log.Println(ps.Clients[port][1].ID())
}
