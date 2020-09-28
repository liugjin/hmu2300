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
	"encoding/json"
	"time"

	"clc.hmu/app/public/log"
	"github.com/gwaylib/errors"
)

// system command type
var (
	SystemCommandTypeGPS      = "gps"
	SystemCommandTypeTime     = "time"
	SystemCommandTypeSystem   = "system"
	SystemCommandTypeReboot   = "reboot"
	SystemCommandTypeUUID     = "uuid"
	SystemCommandTypeLAN      = "lan"
	SystemCommandTypeWAN      = "wan"
	SystemCommandTypeWireless = "wireless"
	SystemCommandTypeInternet = "internet"
	SystemCommandTypeFactory  = "factory"
	SystemCommandTypeAppLED   = "appled"
	SystemCommandTypeUpgrade  = "upgrade"
	SystemCommandTypeDI       = "di"
	SystemCommandTypeDO       = "do"
)

// net mode
var (
	NetModeETH  = "eth"
	NetModeWIFI = "wifi"
	NetModeLTE  = "lte"

	WANModeStatic = "static"
	WANModeDHCP   = "dhcp"

	WifiModeRepeater = "repeater"
	WifiModeAP       = "ap"
	WifiModeSTA      = "sta"
)

// HMUSystemReq system request
type HMUSystemReq struct {
	// gps/time/system/reboot/uuid/lan/wan/wireless/internet/factory
	Type  string `json:"type,omitempty"`
	IsSet int    `json:"isset"`

	// uuid
	UUID string `json:"uuid,omitempty"`

	// time
	TimeServer string `json:"timeserver,omitempty"`

	// reboot
	IsReboot int `json:"isreboot,omitempty"`

	// lan
	LANIP   string `json:"lanip,omitempty"`
	LANMask string `json:"lanmask,omitempty"`
	LANMAC  string `json:"lanmac,omitempty"`

	// wan
	WANMode    string `json:"wanmode,omitempty"` // static, dhcp
	WANIP      string `json:"wanip,omitempty"`
	WANMask    string `json:"wanmask,omitempty"`
	WANGateway string `json:"wangateway,omitempty"`
	WANPDNS    string `json:"wanpdns,omitempty"`
	WANSDNS    string `json:"wansdns,omitempty"`

	// wireless
	WIFIMode     string `json:"wifimode,omitempty"`
	RepeaterSSID string `json:"repeaterssid,omitempty"`
	RepeaterKey  string `json:"repeaterkey,omitempty"`
	SSID         string `json:"ssid,omitempty"`
	Encryption   string `json:"encryption,omitempty"`
	Key          string `json:"key,omitempty"`
	Channel      string `json:"channel,omitempty"`
	Hide         string `json:"hide,omitempty"`
	STASSID      string `json:"stassid,omitempty"`
	STAKey       string `json:"stakey,omitempty"`

	// internet
	InternetMode string `json:"internetmode,omitempty"`
	WifiSSID     string `json:"wifissid,omitempty"`
	WifiKey      string `json:"wifikey,omitempty"`

	// factory
	Mode string `json:"mode,omitempty"` // all, wifi, network

	// appled
	IsON int `json:"ison,omitempty"`

	// upgrade
	URL string `json:"url,omitempty"`

	// dido
	ID    string `json:"id,omitempty"`
	Value int    `json:"value"`
}

// HMUSystemResp system response
type HMUSystemResp struct {
	// gps
	GPSLatitude  string `json:"gpslat,omitempty"`
	GPSLongitude string `json:"gpslon,omitempty"`
	GPSTime      string `json:"gpstime,omitempty"`

	// time
	UpTime      string `json:"uptime,omitempty"`
	CurrentTime string `json:"currtime,omitempty"`
	TimeServer  string `json:"timeserver,omitempty"`

	// system
	CPU           string `json:"cpu,omitempty"`
	Model         string `json:"model,omitempty"`
	SWVersion     string `json:"swversion,omitempty"`
	HWVersion     string `json:"hwversion,omitempty"`
	RAMTotal      string `json:"ramtotal,omitempty"`
	RAMFree       string `json:"ramfree,omitempty"`
	FlashTotal    string `json:"flashtotal,omitempty"`
	FlashFree     string `json:"flashfree,omitempty"`
	SDTotal       string `json:"sdtotal,omitempty"`
	SDFree        string `json:"sdfree,omitempty"`
	KernelVersion string `json:"kernelver,omitempty"`
	PCIEDevice    string `json:"pciedev,omitempty"`
	LTEIMSI       string `json:"lteimsi,omitempty"`
	LTECCID       string `json:"lteccid,omitempty"`
	LTECSQ        string `json:"ltecsq,omitempty"` // 信号强度量

	// reboot, factory
	ISOK int `json:"isOK,omitempty"`

	// uuid
	UUID string `json:"uuid,omitempty"`

	// lan
	LANIP   string `json:"lanip,omitempty"`
	LANMask string `json:"lanmask,omitempty"`
	LANMAC  string `json:"lanmac,omitempty"`

	// wan
	WANMode    string `json:"wanmode,omitempty"`
	WANIP      string `json:"wanip,omitempty"`
	WANMask    string `json:"wanmask,omitempty"`
	WANMAC     string `json:"wanmac,omitempty"`
	WANGateway string `json:"wangateway,omitempty"`
	WANPDNS    string `json:"wanpdns,omitempty"`
	WANSDNS    string `json:"wansdns,omitempty"`

	// wireless
	WIFIMode     string `json:"wifimode,omitempty"`
	RepeaterIP   string `json:"repip,omitempty"`
	RepeaterMask string `json:"repmask,omitempty"`
	RepeaterSSID string `json:"repssid,omitempty"`
	RepeaterKey  string `json:"reprkey,omitempty"`
	SSID         string `json:"ssid,omitempty"`
	Encryption   string `json:"encryption,omitempty"`
	Key          string `json:"key,omitempty"`
	MAC          string `json:"mac,omitempty"`
	Hidden       string `json:"hidden,omitempty"`
	Channel      string `json:"channel,omitempty"`
	Hide         string `json:"hide,omitempty"`
	STAIP        string `json:"staip,omitempty"`
	STAMask      string `json:"stamask,omitempty"`
	STASSID      string `json:"stassid,omitempty"`
	STAKey       string `json:"stakey,omitempty"`
	STAMAC       string `json:"stamac,omitempty"`

	// internet
	InternetMode string `json:"internetmode,omitempty"`
	// WANMode      string `json:"wanmode,omitempty"`
	StaticIP      string `json:"staticip,omitempty"`
	StaticMask    string `json:"staticmask,omitempty"`
	StaticGateway string `json:"staticgateway,omitempty"`
	StaticPDNS    string `json:"staticpdns,omitempty"`
	StaticSDNS    string `json:"staticsdns,omitempty"`
	WIFISSID      string `json:"wifissid,omitempty"`
	WIFIPassword  string `json:"wifipass,omitempty"`
	LTEIP         string `json:"lteip,omitempty"`
	LTEMask       string `json:"ltemask,omitempty"`
	LTEMAC        string `json:"ltemac,omitempty"`
	// WANIP         string `json:"wanip.omitempty"`
	// WANMask       string `json:"wanmask,omitempty"`
	// WANGateway    string `json:"wangateway,omitempty"`
	// WANPDNS       string `json:"wanpdns,omitempty"`
	// WANSDNS       string `json:"wansdns,omitempty"`
	// WIFIMode      string `json:"wifimode,omitempty"`
	// RepeaterIP    string `json:"repip,omitempty"`
	// RepeaterMask  string `json:"repmask,omitempty"`
	// RepeaterSSID  string `json:"repssid,omitempty"`
	// RepeaterKey   string `json:"reprkey,omitempty"`
	// SSID          string `json:"ssid,omitempty"`
	// Encryption    string `json:"encryption,omitempty"`
	// Key           string `json:"key,omitempty"`
	// MAC           string `json:"mac,omitempty"`
	// Hidden        string `json:"hidden,omitempty"`
	// Channel       string `json:"channel,omitempty"`

	// dido
	Value int `json:"value,omitempty"`
}

func (r HMUSystemReq) String() string {
	br, err := json.Marshal(r)
	if err != nil {
		return ""
	}

	return string(br)
}

func byteToResponse(b []byte) *HMUSystemResp {
	var resp = &HMUSystemResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		log.Warning(errors.As(err, string(b)))
		return resp
	}

	return resp
}

type SystemClient interface {
	New(opts *SystemServerOption) SystemClient
	ModelName() string
	Disconnect() error
	Request(req *HMUSystemReq) (*HMUSystemResp, error)
	// 自动检查网络，检查网络过程中会依次切换wan口的各种上网模式。
	// 需要注意，此接口会阻塞线程
	// 若成功，返回error等于空，并返回当前用于wan中信息，详情见Internet接口返回。
	AutoCheckNetworking(urls []string, timeout time.Duration) (*HMUSystemResp, error)
	GPS() (*HMUSystemResp, error)
	Time() (*HMUSystemResp, error)
	SetTimeServer(timeserver string) (*HMUSystemResp, error)
	SystemInfo() (*HMUSystemResp, error)
	Reboot() (*HMUSystemResp, error)
	UUID() (*HMUSystemResp, error)
	SetUUID(id string) (*HMUSystemResp, error)
	LAN() (*HMUSystemResp, error)
	WAN() (*HMUSystemResp, error)
	Wireless() (*HMUSystemResp, error)
	SetAP(ssid, encryption, key, channel, hide string) (*HMUSystemResp, error)
	SetLAN(ip, mask string) (*HMUSystemResp, error)
	Internet() (*HMUSystemResp, error)
	SetEthStatic(ip, mask, gateway, pdns, sdns string) (*HMUSystemResp, error)
	SetEthDHCP() (*HMUSystemResp, error)
	SetWifi(ssid, key string) (*HMUSystemResp, error)
	SetLTE() (*HMUSystemResp, error)
	LTECSQ() (*HMUSystemResp, error)
	FactoryReset(mode string) (*HMUSystemResp, error)
	SetAppLED(status int) (*HMUSystemResp, error)
	DIDO(id string) (*HMUSystemResp, error)
	SetDO(id string, value int) (*HMUSystemResp, error)
}

var sysDrvs = map[string]SystemClient{}
var sysClients = map[string]SystemClient{}

// Register client model
// model -- see ./common.go#model
func RegSysClientModel(model string, client SystemClient) {
	sysDrvs[model] = client
}

// ConnectSystemDaemon connect system daemon
//
// model -- see ./common.go#model
// uri -- connect uri
//
func ConnectSystemDaemon(model string, opts *SystemServerOption) SystemClient {
	drv, ok := sysDrvs[model]
	if !ok {
		// Get default driver
		// log.Debugf("Not Found model:%s, use default instead.", model)
		drv, _ = sysDrvs[MODEL_DEFAULT]
	}

	// cache uri client
	cacheKey := opts.Uri + opts.Vals
	client, ok := sysClients[cacheKey]
	if !ok {
		client = drv.New(opts)
		sysClients[cacheKey] = client
	}
	return client
}
