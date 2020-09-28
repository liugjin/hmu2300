package pmnet

import (
	"context"
	"sync"
	"time"

	pb "clc.hmu/app/portmanager/portpb"
	"clc.hmu/app/public/log/applog"
	"google.golang.org/grpc"
)

var (
	pmconn      *grpc.ClientConn
	pmcli       pb.PortClient
	pmparentctx context.Context
	pmMutex     sync.Mutex
)

const pmtimeout = 60

func GetClientConn() *grpc.ClientConn {
	pmMutex.Lock()
	defer pmMutex.Unlock()
	// if pmconn == nil {
	// 	panic("need init")
	// }
	return pmconn
}

func GetClient() pb.PortClient {
	pmMutex.Lock()
	defer pmMutex.Unlock()
	// if pmcli == nil {
	// 	panic("need init")
	// }
	return pmcli
}

func GetRootContext() context.Context {
	pmMutex.Lock()
	defer pmMutex.Unlock()
	// if pmparentctx == nil {
	// 	panic("need init")
	// }
	return pmparentctx
}

func GetTimeout() time.Duration {
	return pmtimeout
}

// ConnectPortManager connect port server
func ConnectPortManager(address string) {
	pmMutex.Lock()
	defer pmMutex.Unlock()

	var err error

	pmconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*pmtimeout))
	if err != nil {
		go reconnectPortManager(address)
		return
	}

	applog.LOG.Info("connect port manager success")

	pmparentctx = context.Background()

	pmcli = pb.NewPortClient(pmconn)
}

// NewPortClient new port client
func NewPortClient() pb.PortClient {
	return pb.NewPortClient(pmconn)
}

func reconnectPortManager(address string) {
	for {
		applog.LOG.Warning("reconnect port manager...")

		var err error
		pmconn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*pmtimeout))
		if err == nil {
			break
		}
	}

	applog.LOG.Info("connect port manager success")

	pmparentctx = context.Background()

	pmcli = pb.NewPortClient(pmconn)
}

// DisconnectPortManager disconnect
func DisconnectPortManager() error {
	pmMutex.Lock()
	defer pmMutex.Unlock()
	return pmconn.Close()
}
