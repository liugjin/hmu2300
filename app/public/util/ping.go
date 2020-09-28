package util

import (
	"log"
	"net"
	"time"

	"github.com/gwaylib/errors"
)

// 需要root权限
func ICMPPing(ip string) bool {
	t := 3 * time.Second
	conn, err := net.DialTimeout("ip4:icmp", ip, t)
	if err != nil {
		log.Println(errors.As(err))
		return false
	}

	payload := []byte{0x08, 0x00, 0x4d, 0x4b, 0x00, 0x01, 0x00, 0x10, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69}
	_, err = conn.Write(payload)
	if err != nil {
		log.Println(errors.As(err))
		return false
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	buf := make([]byte, 2048)
	num, err := conn.Read(buf[0:])
	if err != nil {
		log.Println(errors.As(err))
		return false
	}

	if string(buf[0:num]) != "" {
		return true
	}

	return false
}

func TCPPing(uri string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", uri, timeout)
	if err != nil {
		log.Println(errors.As(err, uri))
		return false
	}
	defer conn.Close()
	return true
}
