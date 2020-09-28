/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/09/19
 * Despcription: video implement
 *
 */

package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"clc.hmu/app/public"
	"clc.hmu/app/public/util"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolCamera, generalCameraDriverProtocol)
}

// Implement DriverProtocol
type cameraDriverProtocol struct {
	req *public.CameraBindingPayload
	uri string
}

func (dp *cameraDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *cameraDriverProtocol) ClientID() string {
	return VideoClientID + dp.req.Host
}

func (dp *cameraDriverProtocol) NewInstance() (PortClient, error) {
	return NewVideoClient(dp.req.Host, dp.req.User, dp.req.Password)
}

func generalCameraDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeCameraBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &cameraDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// VideoClientID id
var VideoClientID = "video-client"

// VideoClient video client
type VideoClient struct {
	ClientID string

	host string

	filename string
	bufsize  int
	handler  func(string) string
	conn     *net.UnixConn

	username string
	password string
}

// NewVideoClient new video client
func NewVideoClient(host, user, password string) (PortClient, error) {
	size := 10480
	fn := "rtspclient.sock"

	vc := VideoClient{
		ClientID: VideoClientID + host,
		filename: fn,
		bufsize:  size,
		host:     host,
		username: user,
		password: password,
	}

	go vc.Start()

	return &vc, nil
}

// DecodeCameraBindingPayload decode binding
func DecodeCameraBindingPayload(payload string) (public.CameraBindingPayload, error) {
	var p public.CameraBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeCameraOperationPayload decode binding
func DecodeCameraOperationPayload(payload string) (public.CameraOperationPayload, error) {
	var p public.CameraOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

func isPingOK(ip string) bool {
	return util.ICMPPing(ip)
}

// Start start
func (vc *VideoClient) Start() {

}

func (vc *VideoClient) connectServer() error {
	addr, err := net.ResolveUnixAddr("unix", vc.filename)
	if err != nil {
		return fmt.Errorf("Cannot resolve unix addr: %v", err)
	}

	vc.conn, err = net.DialUnix("unix", nil, addr)
	if err != nil {
		return fmt.Errorf("DialUnix failed. %v", err)
	}

	return nil
}

func (vc *VideoClient) write(context string) error {
	if vc.conn == nil {
		return fmt.Errorf("conn unavaliable")
	}

	vc.conn.SetDeadline(time.Now().Add(time.Second * 3))
	_, err := vc.conn.Write([]byte(context))
	if err != nil {
		return fmt.Errorf("Writes failed. %v", err)
	}

	return nil
}

func (vc *VideoClient) read() (string, error) {
	if vc.conn == nil {
		return "", fmt.Errorf("conn unavaliable")
	}

	buf := make([]byte, vc.bufsize)
	vc.conn.SetDeadline(time.Now().Add(time.Second))
	nr, err := vc.conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("Read: " + err.Error())
	}

	return string(buf[0:nr]), nil
}

// ID client id
func (vc *VideoClient) ID() string {
	return vc.ClientID
}

// Sample get values
func (vc *VideoClient) Sample(payload string) (string, error) {
	p, err := DecodeCameraOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	switch p.Channel {
	case "state":
		if isPingOK(vc.host) {
			return "ok", nil
		}
		return "", fmt.Errorf("ping failed")
	}

	return "", nil
}

// Command set values
func (vc *VideoClient) Command(payload string) (string, error) {
	p, err := DecodeCameraOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	switch p.Channel {
	case "capture":
		data, err := public.OnvifHTTPCapture(vc.host, vc.username, vc.password)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}

	return "", nil
}
