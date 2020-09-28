package web
import (
	"context"
	"time"

	pb "clc.hmu/app/appmanager/apppb"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/store/etc"
	"github.com/gwaylib/errors"
	"google.golang.org/grpc"
)

var appWebconn *grpc.ClientConn
var appWebconnAvailable = false

// AppWebClient app client
type AppWebClient struct {
	Client pb.AppClient
}

// ConnectAppServer connect app server
func ConnectAppServer() {
	address := etc.Etc.String("appmanager", "rpc_client")

	var err error
	appWebconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		buslog.LOG.Warningf("connect app server failed, errmsg {%v}", err)

		go reconnectAppServer(address)
		return
	}

	buslog.LOG.Infof("connect app server success")

	appWebconnAvailable = true
}

// DisconnectAppServer disconnect
func DisconnectAppServer() {
	appWebconnAvailable = false
	appWebconn.Close()
}

// NewAppWebClient new app client
func NewAppWebClient() AppWebClient {
	if !appWebconnAvailable {
		return AppWebClient{Client: nil}
	}

	return AppWebClient{Client: pb.NewAppClient(appWebconn)}
}

func reconnectAppServer(address string) {
	for {
		var err error
		appWebconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
		if err != nil {
			buslog.LOG.Warningf("reconnect app server failed, errmsg {%v}", err)
			continue
		} else {
			buslog.LOG.Infof("connect app server success")
			break
		}
	}

	appWebconnAvailable = true
}

// Notify notify
func (c *AppWebClient) Notify(topic, payload string) error {
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

