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
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	pb "clc.hmu/app/portmanager/portmanager"
	"clc.hmu/app/portmanager/src/protocol"
	"clc.hmu/app/public"
	"google.golang.org/grpc"
)

func TestEncodeRequest(t *testing.T) {
	port := "/dev/ttyS0"
	baudrate := 9600
	timeout := 3000
	soi := byte(0x7E)
	ver := byte(0x21)
	adr := byte(0x01)
	cid1 := byte(0x2A)
	eoi := byte(0x0D)
	client, err := protocol.NewPMBusClient(port, baudrate, timeout, soi, ver, adr, cid1, eoi)
	if err != nil {
		t.Error(err)
	}

	req := public.PMBUSOperationPayload{
		CID2:  byte(0x42),
		LENID: 0,
	}

	frame, err := client.EncodeRequest(req)
	if err != nil {
		t.Error(err)
	}

	log.Println(string(frame), len(frame))
}

func TestPMBusBinding(t *testing.T) {
	// address := "127.0.0.1:50051"
	address := "192.168.1.1:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pb.NewPortClient(conn)

	// port := "/dev/ttyS0"
	port := "/dev/COM2"

	req := public.PMBUSBindingPayload{
		BaudRate: 9600,
		Timeout:  3000,
		SOI:      byte(0x7E),
		VER:      byte(0x21),
		ADR:      byte(0x01),
		CID1:     byte(0x2A),
		EOI:      byte(0x0D),
	}

	payload, _ := json.Marshal(req)

	r, err := client.Binding(ctx, &pb.BindingRequest{
		Protocol: public.ProtocolPMBUS,
		Port:     port,
		Payload:  string(payload),
	})

	if err != nil {
		t.Errorf("could not bind: %v", err)
	}

	t.Log(r)
}

func TestPMBusSample(t *testing.T) {
	// address := "127.0.0.1:50051"
	address := "192.168.1.1:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	client := pb.NewPortClient(conn)

	// port := "/dev/ttyS0"
	port := "/dev/COM2"

	req := public.PMBUSOperationPayload{
		CID2:  byte(0x43),
		LENID: 0,
	}

	payload, _ := json.Marshal(req)

	r, err := client.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		t.Errorf("could not sample: %v", err)
	}

	t.Log(r)
}
