package core

import (
	"context"
	"fmt"
	"time"

	"clc.hmu/app/appmanager/appnet/pmnet"
	pb "clc.hmu/app/portmanager/portpb"
	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// 将端口绑定到硬件通讯
// 直接将配置文件透传过去自行解析
// TODO:从配置文件池中进行统一更新与读取
//
// 参数
// uri -- 资源地址、设备端口等, 必须。

func binding(uri, protocol, suid string) error {
	if len(uri) == 0 || len(protocol) == 0 || len(suid) == 0 {
		return errors.New("Invalid argument").As(uri, protocol)
	}
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     uri,
		Protocol: protocol,
		Suid:     suid,
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}
