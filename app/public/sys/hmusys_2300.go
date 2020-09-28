/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/08/15
 * Despcription: hmu2300 definition
 *
 */

package sys

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/muid"
	"clc.hmu/app/public/util"
	"github.com/gwaylib/errors"
)

// SystemHMU2300 system client
type SystemHMU2300 struct {
	uri string

	mutex sync.Mutex // 资源锁
	Conn  *net.UDPConn

	uuid string // 缓存uuid
}

func init() {
	RegSysClientModel(MODEL_HMU2300, &SystemHMU2300{})
}

// New new
func (c *SystemHMU2300) New(opts *SystemServerOption) SystemClient {
	return &SystemHMU2300{
		uri: opts.Uri,
	}
}

// ModelName model name
func (c *SystemHMU2300) ModelName() string {
	return MODEL_HMU2300
}

// Disconnect disconnect
func (c *SystemHMU2300) Disconnect() error {
	// not need to implements
	return nil
}

func (c *SystemHMU2300) connect() error {
	addr := c.uri
	// connect to hmc system server
	var err error
	raddr, err := net.ResolveUDPAddr("udp", addr)
	c.Conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Debug(errors.As(err, addr))

		return errors.As(err, addr)
	}

	return nil
}

func (c *SystemHMU2300) disconnect() error {
	if c.Conn != nil {
		err := c.Conn.Close()
		c.Conn = nil
		return errors.As(err)
	}
	return nil
}

// Request requst to system daemon
func (c *SystemHMU2300) Request(req *HMUSystemReq) (*HMUSystemResp, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 这是一个短连接操作
	if err := c.connect(); err != nil {
		return nil, errors.As(err)
	}
	defer c.disconnect()

	// set deadline
	c.Conn.SetDeadline(time.Now().Add(time.Second * 60))

	// write to daemon
	content := fmt.Sprintf("%s", req)
	n, err := c.Conn.Write([]byte(content))
	if err != nil {
		return nil, errors.As(err)
	}

	// read response
	var buf = make([]byte, 512)
	n, err = c.Conn.Read(buf)
	if err != nil {
		// io中断，1，接口程序挂了，2，系统重启了
		// if err == io.EOF {
		// }
		return nil, errors.As(err)
	}
	// buslog.LOG.Infof("response: %s, bytes: %d", string(buf[:n]), n)

	return byteToResponse(buf[:n]), nil
}

func (c *SystemHMU2300) checkNetworking(urls []string, timeout time.Duration) (string, error) {
	// 检查当前的网络是否可用
	for _, val := range urls {
		if !util.TCPPing(val, timeout) {
			continue
		}
		return val, nil
	}
	return "", errors.New("All invalid")
}

// AutoCheckNetworking auto check network
func (c *SystemHMU2300) AutoCheckNetworking(urls []string, timeout time.Duration) (*HMUSystemResp, error) {
	return &HMUSystemResp{}, nil
}

// GPS query gps info
func (c *SystemHMU2300) GPS() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeGPS

	return c.Request(req)
}

// Time query time info
func (c *SystemHMU2300) Time() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeTime
	req.IsSet = 0

	return c.Request(req)
}

// SetTimeServer set time server
func (c *SystemHMU2300) SetTimeServer(timeserver string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeTime
	req.IsSet = 1
	req.TimeServer = timeserver

	return c.Request(req)
}

// SystemInfo query system info, including memory, cpu and so on
func (c *SystemHMU2300) SystemInfo() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeSystem

	return c.Request(req)
}

// Reboot reboot device
func (c *SystemHMU2300) Reboot() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeReboot
	req.IsReboot = 1

	return c.Request(req)
}

// UUID query device's uuid
func (c *SystemHMU2300) UUID() (*HMUSystemResp, error) {
	return &HMUSystemResp{
		UUID: muid.GetMuID(),
	}, nil
}

// SetUUID set device's uuid
func (c *SystemHMU2300) SetUUID(id string) (*HMUSystemResp, error) {
	// TODO:确认此处逻辑
	// 此处应检查id是否设置过，若设置过，不应再设置，否则可能会被改写为不合法的ID

	muid.SetMemID(id) // 更新缓存

	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeUUID
	req.IsSet = 1
	req.UUID = id

	return c.Request(req)
}

// LAN query LAN info
func (c *SystemHMU2300) LAN() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeLAN

	return c.Request(req)
}

// WAN query WAN info
func (c *SystemHMU2300) WAN() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeWAN

	return c.Request(req)
}

// Wireless query wireless info
func (c *SystemHMU2300) Wireless() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeWireless

	return c.Request(req)
}

// SetAP set ap mode
func (c *SystemHMU2300) SetAP(ssid, encryption, key, channel, hide string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeWireless
	req.IsSet = 1
	req.WIFIMode = WifiModeAP
	req.SSID = ssid
	req.Encryption = encryption
	req.Key = key
	req.Channel = channel
	req.Hide = hide

	return c.Request(req)
}

// SetLAN set lan
func (c *SystemHMU2300) SetLAN(ip, mask string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeLAN
	req.IsSet = 1
	req.LANIP = ip
	req.LANMask = mask

	return c.Request(req)
}

// Internet query internet info
func (c *SystemHMU2300) Internet() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet

	return c.Request(req)
}

// SetEthStatic set eth static mode
func (c *SystemHMU2300) SetEthStatic(ip, mask, gateway, pdns, sdns string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeETH
	req.WANMode = WANModeStatic
	req.WANIP = ip
	req.WANMask = mask
	req.WANGateway = gateway
	req.WANPDNS = pdns
	req.WANSDNS = sdns

	return c.Request(req)
}

// SetEthDHCP set eth dhcp mode
func (c *SystemHMU2300) SetEthDHCP() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeETH
	req.WANMode = WANModeDHCP

	return c.Request(req)
}

// SetWifi set wifi mode
func (c *SystemHMU2300) SetWifi(ssid, key string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeWIFI
	req.WifiSSID = ssid
	req.WifiKey = key

	return c.Request(req)
}

// SetLTE set lte mode
func (c *SystemHMU2300) SetLTE() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeLTE

	return c.Request(req)
}

// LTECSQ 获取LTE信号强度
func (c *SystemHMU2300) LTECSQ() (*HMUSystemResp, error) {

	return &HMUSystemResp{}, nil
}

// FactoryReset factory reset
func (c *SystemHMU2300) FactoryReset(mode string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeFactory
	req.Mode = mode

	return c.Request(req)
}

// SetAppLED set appled
func (c *SystemHMU2300) SetAppLED(status int) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeAppLED
	req.IsON = status

	return c.Request(req)
}

// DIDO get dido value
func (c *SystemHMU2300) DIDO(id string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}

	if strings.Contains(id, "di") {
		req.Type = SystemCommandTypeDI
	} else {
		req.Type = SystemCommandTypeDO
	}

	req.ID = id

	return c.Request(req)
}

// SetDO set do
func (c *SystemHMU2300) SetDO(id string, value int) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeDO
	req.IsSet = 1
	req.ID = id
	req.Value = value

	return c.Request(req)
}
