package frpc

import "fmt"

const (
	// 需要写死在代码中, 以防止查看
	DefaultSSHSK       = "SHkXhM833093hNHvmoRMYmjVoH9kbQqE"
	DefaultServerAddr  = "frp.huayuan-iot.com"
	DefaultServerPort  = "7001"
	DefaultServerToken = "elh529bQ3ei5t2TW9VQbZmya2uS60M0C93AY^Q0CAc&VrgeWv6NiK1aYm@eAW^tc"
)
const sshData = `
[common]
server_addr = %s
server_port = %s
token = %s

[%s_ssh]
type = stcp
sk = %s
local_ip = 127.0.0.1
local_port = %s
`

type FrpcSSH struct {
	ServerHost  string
	ServerPort  string
	ServerToken string
	MuID        string
	MuSK        string
	SSHPort     string
}

func (f *FrpcSSH) String() string {
	return fmt.Sprintf(
		sshData,
		f.ServerHost,
		f.ServerPort,
		f.ServerToken,
		f.MuID,
		f.MuSK,
		f.SSHPort,
	)
}
