package frpc

import (
	"strings"

	"github.com/gwaylib/conf/ini"
	"github.com/gwaylib/errors"
)

const (
	customTime = "2006-01-02 15:04:05"
)

type FrpcIni struct {
	*ini.File

	path string
}

func OpenFrpcIni(path string) (*FrpcIni, error) {
	frpcEtc, err := ini.GetFile(path)
	if err != nil {
		return nil, errors.As(err)
	}

	// 若连接frp.huayuan-iot.com的服务器，自动替换token
	sAddr := frpcEtc.String("common", "server_addr")
	sPort := frpcEtc.String("common", "server_port")
	if sAddr == DefaultServerAddr && sPort == DefaultServerPort {
		section := frpcEtc.Section("common")
		section.Key("token").SetValue(DefaultServerToken)
	}

	return &FrpcIni{
		path: path,
		File: frpcEtc,
	}, nil
}

func (f *FrpcIni) String() string {
	buf := &strings.Builder{}
	f.WriteToIndent(buf, "") // ignore error
	return buf.String()
}
