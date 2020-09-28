/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: bus server net
 *
 */

package appnet

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/applog"
	"github.com/gwaylib/errors"

	pb "clc.hmu/app/busmanager/buspb"
	"google.golang.org/grpc"
)

var bmconn *grpc.ClientConn
var bmcli pb.BusClient

const bmtimeout = 3

// ConnectBusManager connect bus server
func ConnectBusManager(address string) {
	var err error

	bmconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*bmtimeout))
	if err != nil {
		go reconnectBusManager(address)
		return
	}

	applog.LOG.Info("connect bus manager success")

	bmcli = pb.NewBusClient(bmconn)
}

func reconnectBusManager(address string) {
	for {
		applog.LOG.Warning("reconnect bus manager...")

		var err error
		bmconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*bmtimeout))
		if err == nil {
			break
		}
	}

	applog.LOG.Info("connect bus manager success")

	bmcli = pb.NewBusClient(bmconn)
}

// DisconnectBusManager disconnect
func DisconnectBusManager() error {
	return bmconn.Close()
}

// SetHeartbeat set heartbead
func SetHeartbeat(topic, payload string) {
	go func() {
		for {
			if bmcli != nil {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*bmtimeout)
				defer cancel()

				r, err := bmcli.Publish(ctx, &pb.PublishRequest{
					Topic:   topic,
					Payload: string(payload),
				})

				if err != nil {
					log.Printf("send heartbeat failed, errmsg {%v}", err)
				} else if r.Status != public.StatusOK {
					log.Printf("send heartbeat failed, result [%s]", r.Message)
				} else {
					log.Printf("send heartbeat success")
				}
			}

			time.Sleep(time.Second * 2)
		}
	}()
}

// PublishSampleValues public messages
func PublishSampleValues(p public.MessagePayload) error {
	if bmcli == nil {
		return fmt.Errorf("bus client unavailable")
	}

	topic := "sample-values/" + p.MonitoringUnitID + "/" + p.SampleUnitID + "/" + p.ChannelID

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal payload failed")
	}
	// log.Printf("topic:%s\npayload:%+v", topic, p)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*bmtimeout)
	defer cancel()

	r, err := bmcli.Publish(ctx, &pb.PublishRequest{
		Topic:   topic,
		Payload: string(payload),
	})

	if err != nil {
		return fmt.Errorf("publish failed, errmsg [%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("publish failed, result [%s]", r.Message)
	}

	return nil
}

// Subscribe subscribe
func Subscribe(topic string) error {
	if bmcli == nil {
		return fmt.Errorf("bus client unavailable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*bmtimeout)
	defer cancel()

	r, err := bmcli.Subscribe(ctx, &pb.SubscribeRequest{
		Topic: topic,
	})

	if err != nil {
		return fmt.Errorf("subscribe failed, errmsg [%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("subscribe failed, result [%s]", r.Message)
	}

	return nil
}

// PublishDiscovery public messages
func PublishDiscovery(muid string, data string) error {
	if bmcli == nil {
		return fmt.Errorf("bus client unavailable")
	}

	topic := "discovery/" + muid + "/info"

	// applog.LOG.Infof("topic:%s\npayload:%+v", topic, data)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*bmtimeout)
	defer cancel()

	r, err := bmcli.Publish(ctx, &pb.PublishRequest{
		Topic:   topic,
		Payload: data,
	})

	if err != nil {
		return fmt.Errorf("publish failed, errmsg [%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("publish failed, result [%s]", r.Message)
	}

	return nil
}

// ReplyCommand reply command
func ReplyCommand(topic string, data interface{}) error {
	if bmcli == nil {
		return errors.New("bus client unavailable")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return errors.As(err, data)
	}
	// applog.LOG.Infof("topic:%s\npayload:%+v", topic, string(payload))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*bmtimeout)
	defer cancel()

	r, err := bmcli.Publish(ctx, &pb.PublishRequest{
		Topic:   topic,
		Payload: string(payload),
	})
	if err != nil {
		return errors.As(err, topic, string(payload))
	}

	if r.Status != public.StatusOK {
		return errors.New("publish failed, result").As(r.Message)
	}

	return nil
}
