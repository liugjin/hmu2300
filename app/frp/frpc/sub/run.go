package sub

import (
	"context"
	"net"
	"strings"

	"clc.hmu/app/frp/frpc/client"
	"clc.hmu/app/frp/frpc/g"
	"clc.hmu/app/frp/frpc/models/config"
	"clc.hmu/app/frp/frpc/utils/log"
	"github.com/fatedier/golib/crypto"
	"github.com/gwaylib/errors"
)

type Frpc struct {
	*client.Service
}

func (f *Frpc) Run() error {
	// Capture the exit signal if we use kcp.
	if g.GlbClientCfg.Protocol == "kcp" {
		go handleSignal(f.Service)
	}

	if err := f.Service.Run(); err != nil {
		return errors.As(err)
	}
	if g.GlbClientCfg.Protocol == "kcp" {
		<-kcpDoneCh
	}
	return nil

}

func (f *Frpc) Reload(content string) error {
	if err := parseClientCommonCfg(CfgFileTypeIni, content); err != nil {
		return errors.As(err)
	}

	pxyCfgs, visitorCfgs, err := config.LoadAllConfFromIni(g.GlbClientCfg.User, content, g.GlbClientCfg.Start)
	if err != nil {
		return errors.As(err)
	}

	return f.Service.ReloadConf(pxyCfgs, visitorCfgs)
}

func NewFrpc(salt, content string) (*Frpc, error) {
	crypto.DefaultSalt = salt
	if err := parseClientCommonCfg(CfgFileTypeIni, content); err != nil {
		return nil, errors.As(err)
	}

	pxyCfgs, visitorCfgs, err := config.LoadAllConfFromIni(g.GlbClientCfg.User, content, g.GlbClientCfg.Start)
	if err != nil {
		return nil, errors.As(err)
	}

	log.InitLog(g.GlbClientCfg.LogWay, g.GlbClientCfg.LogFile, g.GlbClientCfg.LogLevel, g.GlbClientCfg.LogMaxDays)
	if g.GlbClientCfg.DnsServer != "" {
		s := g.GlbClientCfg.DnsServer
		if !strings.Contains(s, ":") {
			s += ":53"
		}
		// Change default dns server for frpc
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", s)
			},
		}
	}

	srv, err := client.NewService(pxyCfgs, visitorCfgs)
	if err != nil {
		return nil, errors.As(err)
	}

	f := &Frpc{Service: srv}
	return f, nil
}
