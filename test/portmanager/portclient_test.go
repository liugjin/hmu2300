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
	"context"
	"testing"
	"time"

	"clc.hmu/app/public"

	pb "clc.hmu/app/portmanager/portmanager"
	"google.golang.org/grpc"
)

func TestBindingPort(t *testing.T) {
	address := "127.0.0.1:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pb.NewPortClient(conn)

	port := "/dev/ttyS0"
	payload := `{"slaveid": 1}`
	r, err := client.Binding(ctx, &pb.BindingRequest{
		Protocol: public.ProtocolModbusSerial,
		Port:     port,
		Payload:  payload,
	})

	if err != nil {
		t.Errorf("could not bind: %v", err)
	}

	t.Log(r)
}

func TestBindingSystem(t *testing.T) {
	address := "127.0.0.1:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pb.NewPortClient(conn)

	port := "/dev/self"
	payload := `{"host": "192.168.10.1", "port": "9988", "model": "hmu2000"}`
	r, err := client.Binding(ctx, &pb.BindingRequest{
		Protocol: public.ProtocolHYIOTMU,
		Port:     port,
		Payload:  payload,
	})

	if err != nil {
		t.Errorf("could not bind: %v", err)
	}

	t.Log(r)
}

func TestSample(t *testing.T) {
	address := "127.0.0.1:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := pb.NewPortClient(conn)

	port := "/dev/ttyS0"
	payload := `{"slaveid": 1}`
	r, err := client.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: payload,
	})

	if err != nil {
		t.Errorf("could not sample: %v", err)
	}

	t.Log(r)
}
