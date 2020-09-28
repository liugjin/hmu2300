/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: app client
 *
 */

package module

import (
	"context"
	"time"

	pb "clc.hmu/app/appmanager/apppb"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/store/etc"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gwaylib/errors"
	"google.golang.org/grpc"
)

var appconn *grpc.ClientConn
var appconnAvailable = false

// AppClient app client
type AppClient struct {
	Client pb.AppClient
}

// ConnectAppServer connect app server
func ConnectAppServer() {
	address := etc.Etc.String("appmanager", "rpc_client")

	var err error
	appconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		buslog.LOG.Warningf("connect app server failed, errmsg {%v}", err)

		go reconnectAppServer(address)
		return
	}

	buslog.LOG.Infof("connect app server success")

	appconnAvailable = true
}

// DisconnectAppServer disconnect
func DisconnectAppServer() {
	appconnAvailable = false
	appconn.Close()
}

// NewAppClient new app client
func NewAppClient() AppClient {
	if !appconnAvailable {
		return AppClient{Client: nil}
	}

	return AppClient{Client: pb.NewAppClient(appconn)}
}

func reconnectAppServer(address string) {
	for {
		var err error
		appconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
		if err != nil {
			buslog.LOG.Warningf("reconnect app server failed, errmsg {%v}", err)
			continue
		} else {
			buslog.LOG.Infof("connect app server success")
			break
		}
	}

	appconnAvailable = true
}

// Notify notify
func (c *AppClient) Notify(topic, payload string) error {
	if c.Client == nil {
		return errors.New("app client unavaliable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	r, err := c.Client.Notify(ctx, &pb.NotifyRequest{
		Topic:   topic,
		Payload: payload,
		Caller:  public.CallerBusServer,
	})
	if err != nil {
		return errors.As(err, topic, payload)
	}

	if r.Status != public.StatusOK {
		return errors.New(r.Message)
	}

	return nil
}

var appclient AppClient

// SubMessageHandler handle subscribe message
func SubMessageHandler(client MQTT.Client, msg MQTT.Message) {
	// log.Debugf("RECEIVED TOPIC: %s MESSAGE: %s\n", msg.Topic(), string(msg.Payload()))

	// check appclient available
	if appclient.Client == nil {
		// buslog.LOG.Warningf("app client unavaliable, new client")

		appclient = NewAppClient()
		if appclient.Client == nil {
			// buslog.LOG.Warningf("new app client failed")
			return
		}

		buslog.LOG.Infof("new app client success")
	}

	topic := msg.Topic()
	payload := string(msg.Payload())
	// send message to app
	if err := appclient.Notify(topic, payload); err != nil {
		buslog.LOG.Warningf("notify app server failed, errmsg {%v}\n", errors.As(err))
	}
}
