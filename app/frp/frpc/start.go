package frpc

import (
	"clc.hmu/app/frp/frpc/sub"
	"github.com/gwaylib/errors"
)

var (
	frpcSrv *sub.Frpc
)

func StartSSH(f *FrpcSSH) error {
	// 因前期实放不太规范，暂写死sk， 以便可以运维
	frpcSrv, err := sub.NewFrpc("frp", f.String())
	if err != nil {
		return errors.As(err)
	}
	return errors.As(frpcSrv.Run())
}

func Start(iniPath string) error {
	tpl, err := OpenFrpcIni(iniPath)
	if err != nil {
		return errors.As(err)
	}
	frpcSrv, err = sub.NewFrpc("frp", tpl.String())
	if err != nil {
		return errors.As(err)
	}

	return errors.As(frpcSrv.Run())

}

func Reload(tpl *FrpcIni) error {
	return frpcSrv.Reload(tpl.String())
}
