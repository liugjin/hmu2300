/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/11/23
 * Despcription: app client
 *
 */

package portnet

import (
	"context"
	"fmt"
	"time"

	pb "clc.hmu/app/appmanager/apppb"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log/portlog"
	"clc.hmu/app/public/store/etc"
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
		portlog.LOG.Warningf("port server connect app server failed, errmsg {%v}", err)

		go reconnectAppServer(address)
		return
	}

	portlog.LOG.Infof("port server connect app server success")

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
		address := etc.Etc.String("appmanager", "rpc_client")
		appconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
		if err != nil {
			portlog.LOG.Warningf("port server reconnect app server failed, errmsg {%v}", err)
			continue
		} else {
			portlog.LOG.Infof("port server connect app server success")
			break
		}
	}

	appconnAvailable = true
}

// Notify notify
func (c *AppClient) Notify(topic, payload string) error {
	if c.Client == nil {
		return fmt.Errorf("app client unavaliable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Client.Notify(ctx, &pb.NotifyRequest{
		Topic:   topic,
		Payload: payload,
		Caller:  public.CallerPortServer,
	})

	if err != nil {
		return fmt.Errorf("could not publish: %v", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("%s", r.Message)
	}

	return nil
}

// DefaultAppClient default app client
var DefaultAppClient AppClient

// DefaultNotify use default app client to notify
func DefaultNotify(topic, payload string) error {
	if err := DefaultAppClient.Notify(topic, payload); err != nil {
		fmt.Printf("notify failed: %v\n", err)

		// fail, renew client
		DefaultAppClient = NewAppClient()
		return DefaultAppClient.Notify(topic, payload)
	}

	return nil
}
