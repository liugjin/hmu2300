/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: hmu2000 definition
 *
 */

package sys

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"net"
	"os/exec"
	"time"

	"clc.hmu/app/public/at"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/muid"
	"clc.hmu/app/public/util"
	"github.com/gwaylib/errors"
)

// SystemClient system client
type SystemHmu2000 struct {
	uri   string
	atCmd at.AT

	mutex sync.Mutex // 资源锁
	Conn  net.Conn

	uuid string // 缓存uuid
}

func init() {
	RegSysClientModel(MODEL_HMU2000, &SystemHmu2000{})
}

// RestartSystemDaemon restart system daemon
func restartSystemHmu2000Daemon() error {
	// execute read public key
	cmd := exec.Command("/etc/init.d/system_daemon", "restart")
	if err := cmd.Run(); err != nil {
		return errors.As(err)
	}

	return nil
}

func (c *SystemHmu2000) New(opts *SystemServerOption) SystemClient {
	vals, err := url.ParseQuery(opts.Vals)
	if err != nil {
		panic(err)
	}
	atFile := vals.Get("at_file")
	atTimeout, err := strconv.Atoi(vals.Get("at_timeout"))
	if err != nil {
		atTimeout = 60 * 1000
	}
	return &SystemHmu2000{
		uri:   opts.Uri,
		atCmd: at.NewFileATCmd(atFile, time.Duration(atTimeout)*1e6),
	}
}

func (c *SystemHmu2000) ModelName() string {
	return MODEL_HMU2000
}

func (c *SystemHmu2000) Disconnect() error {
	// not need to implements
	return nil
}

func (c *SystemHmu2000) connect() error {
	addr := c.uri
	// connect to hmc system server
	var err error
	c.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Debug(errors.As(err, addr))
		// try restart system daemon once
		if err := restartSystemHmu2000Daemon(); err != nil {
			return errors.As(err, addr)
		}

		time.Sleep(time.Second * 3)

		// reconnect system daemon
		c.Conn, err = net.Dial("tcp", addr)
		if err != nil {
			return errors.As(err, addr)
		}
	}

	return nil
}

func (c *SystemHmu2000) disconnect() error {
	if c.Conn != nil {
		err := c.Conn.Close()
		c.Conn = nil
		return errors.As(err)
	}
	return nil
}

// Request requst to system daemon
func (c *SystemHmu2000) Request(req *HMUSystemReq) (*HMUSystemResp, error) {
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

func (c *SystemHmu2000) checkNetworking(urls []string, timeout time.Duration) (string, error) {
	// 检查当前的网络是否可用
	for _, val := range urls {
		if !util.TCPPing(val, timeout) {
			continue
		}
		return val, nil
	}
	return "", errors.New("All invalid")
}

func (c *SystemHmu2000) makeNetReqFromNetResp(resp *HMUSystemResp, kind string) *HMUSystemReq {
	var req = &HMUSystemReq{}
	switch kind {
	case "dhcp":
		req.Type = SystemCommandTypeInternet
		req.IsSet = 1
		req.InternetMode = NetModeETH
		req.WANMode = WANModeDHCP

	case "static":
		req.Type = SystemCommandTypeInternet
		req.IsSet = 1
		req.InternetMode = NetModeETH
		req.WANMode = WANModeStatic
		req.WANIP = resp.StaticIP
		req.WANMask = resp.StaticMask
		req.WANGateway = resp.StaticGateway
		req.WANPDNS = resp.StaticPDNS
		req.WANSDNS = resp.StaticSDNS
	case "wifi":
		req.Type = SystemCommandTypeInternet
		req.IsSet = 1
		req.InternetMode = NetModeWIFI
		req.WifiSSID = resp.WIFISSID
		req.WifiKey = resp.WIFIPassword

	case "lte":
		req.Type = SystemCommandTypeInternet
		req.IsSet = 1
		req.InternetMode = NetModeLTE
	}
	return req

}

func (c *SystemHmu2000) AutoCheckNetworking(urls []string, timeout time.Duration) (*HMUSystemResp, error) {
	// 备份当前网络
	bakResp, err := c.Internet()
	if err != nil {
		return nil, errors.As(err)
	}
	_ = bakResp

	// 检查当前网络是否正确
	if _, err := c.checkNetworking(urls, timeout); err == nil {
		return bakResp, nil
	}

	internetMode := []string{
		"dhcp",   // 设置动态获取ip
		"lte",    // 设置使用4G
		"wifi",   // 设置wifi参数
		"static", // 设置静态ip
	}
	// 依次设置网络进行检查
	for _, val := range internetMode {
		log.Debug("Auto Connect :", val)
		req := c.makeNetReqFromNetResp(bakResp, val)
		if _, err := c.Request(req); err != nil {
			// io中断，底层切换网络会造成此异常
			if !errors.Equal(io.EOF, err) {
				log.Debug(errors.As(err))
			}
		}
		_, err := c.checkNetworking(urls, timeout)
		if err == nil {
			log.Debug("Auto Connect success")
			return c.Internet()
		}
	}

	// 恢复原设置
	kind := "dhcp"
	switch bakResp.InternetMode {
	case "eth":
		if bakResp.WANMode == "dhcp" {
			kind = "dhcp"
		} else {
			kind = "static"
		}
	default:
		kind = bakResp.InternetMode
	}
	req := c.makeNetReqFromNetResp(bakResp, kind)
	if _, err := c.Request(req); err != nil {
		log.Debug(errors.As(err))
	}

	// 返回所有设置不可用的操作
	return nil, errors.New("NotFound")
}

// GPS query gps info
func (c *SystemHmu2000) GPS() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeGPS

	return c.Request(req)
}

// Time query time info
func (c *SystemHmu2000) Time() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeTime
	req.IsSet = 0

	return c.Request(req)
}

// SetTimeServer set time server
func (c *SystemHmu2000) SetTimeServer(timeserver string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeTime
	req.IsSet = 1
	req.TimeServer = timeserver

	return c.Request(req)
}

// SystemInfo query system info, including memory, cpu and so on
func (c *SystemHmu2000) SystemInfo() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeSystem

	return c.Request(req)
}

// Reboot reboot device
func (c *SystemHmu2000) Reboot() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeReboot
	req.IsReboot = 1

	return c.Request(req)
}

// UUID query device's uuid
func (c *SystemHmu2000) UUID() (*HMUSystemResp, error) {
	return &HMUSystemResp{
		UUID: muid.GetMuID(),
	}, nil
}

// SetUUID set device's uuid
func (c *SystemHmu2000) SetUUID(id string) (*HMUSystemResp, error) {
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
func (c *SystemHmu2000) LAN() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeLAN

	return c.Request(req)
}

// WAN query WAN info
func (c *SystemHmu2000) WAN() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeWAN

	return c.Request(req)
}

// Wireless query wireless info
func (c *SystemHmu2000) Wireless() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeWireless

	return c.Request(req)
}

// SetAP set ap mode
func (c *SystemHmu2000) SetAP(ssid, encryption, key, channel, hide string) (*HMUSystemResp, error) {
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
func (c *SystemHmu2000) SetLAN(ip, mask string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeLAN
	req.IsSet = 1
	req.LANIP = ip
	req.LANMask = mask

	return c.Request(req)
}

// Internet query internet info
func (c *SystemHmu2000) Internet() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet

	return c.Request(req)
}

// SetEthStatic set eth static mode
func (c *SystemHmu2000) SetEthStatic(ip, mask, gateway, pdns, sdns string) (*HMUSystemResp, error) {
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
func (c *SystemHmu2000) SetEthDHCP() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeETH
	req.WANMode = WANModeDHCP

	return c.Request(req)
}

// SetWifi set wifi mode
func (c *SystemHmu2000) SetWifi(ssid, key string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeWIFI
	req.WifiSSID = ssid
	req.WifiKey = key

	return c.Request(req)
}

// SetLTE set lte mode
func (c *SystemHmu2000) SetLTE() (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeInternet
	req.IsSet = 1
	req.InternetMode = NetModeLTE

	return c.Request(req)
}

// 获取LTE信号强度
func (c *SystemHmu2000) LTECSQ() (*HMUSystemResp, error) {
	val, err := c.atCmd.CSQ()
	if err != nil {
		return nil, errors.As(err)
	}
	vals := strings.Split(val, ":")
	if len(vals) != 2 {
		return nil, errors.New("error protocol").As(val)
	}

	return &HMUSystemResp{
		LTECSQ: vals[1],
	}, nil
}

// FactoryReset factory reset
func (c *SystemHmu2000) FactoryReset(mode string) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeFactory
	req.Mode = mode

	return c.Request(req)
}

// SetAppLED set appled
func (c *SystemHmu2000) SetAppLED(status int) (*HMUSystemResp, error) {
	var req = &HMUSystemReq{}
	req.Type = SystemCommandTypeAppLED
	req.IsON = status

	return c.Request(req)
}

// DIDO get dido value
func (c *SystemHmu2000) DIDO(id string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetDO set do
func (c *SystemHmu2000) SetDO(id string, value int) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}
