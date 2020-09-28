/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/06/21
 * Despcription: system invoke
 *
 */

package protocol

import (
	"encoding/json"
	"fmt"
	"sync"

	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolHYIOTMU, generalSystemDriverProtocol)
}

// Implement DriverProtocol
type systemDriverProtocol struct {
	req *public.SystemBindingPayload
	uri string
}

func (dp *systemDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *systemDriverProtocol) ClientID() string {
	return dp.req.Model
}

func (dp *systemDriverProtocol) NewInstance() (PortClient, error) {
	return NewSystemClient(*dp.req)
}

func generalSystemDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeSystemBindingRequest(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &systemDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

var (
	channelID              = "id"
	channelNetMode         = "netmode"
	channelLongitude       = "longitude"
	channelLatitude        = "latitude"
	channelGPS             = "gps" // including longitude and latitude
	channelUTCTime         = "utctime"
	channelUpTime          = "uptime"
	channelCurrentTime     = "currtime"
	channelTimeServer      = "timeserver"
	channelCPU             = "cpu"
	channelModel           = "model"
	channelSoftwareVersion = "swversion"
	channelHardwareVersion = "hwversion"
	channelRAMTotal        = "ramtotal"
	channelRAMFree         = "ramfree"
	channelFlashTotal      = "flashtotal"
	channelFlashFree       = "flashfree"
	channelSDTotal         = "sdtotal"
	channelSDFree          = "sdfree"
	channelKernelVersion   = "kernelver"
	channelPCIEDevice      = "pciedev"
	channelLTEIMSI         = "lteimsi"
	channelLTECCID         = "lteccid"
	channelLTECSQ          = "ltecsq"
	channelLANIP           = "lanip"
	channelLANMask         = "lanmask"
	channelLANMAC          = "lanmac"
	channelWANIP           = "wanip"
	channelWANMask         = "wanmask"
	channelWANGateway      = "wangateway"
	channelWANMAC          = "wanmac"
	channelWANPDNS         = "wanpdns"
	channelWANSDNS         = "wansdns"
	channelAPSSID          = "apssid"
	channelAPKey           = "apkey"
	channelWifiSSID        = "wifissid"
	channelWifiPassword    = "wifipass"
)

// SystemChannel system channels, cache channel values
type SystemChannel struct {
	ID              string
	NetMode         string
	Longitude       string
	Latitude        string
	UTCTime         string
	UpTime          string
	CurrentTime     string
	TimeServer      string
	CPU             string
	Model           string
	SoftwareVersion string
	HardwareVersion string
	RAMTotal        string
	RAMFree         string
	FlashTotal      string
	FlashFree       string
	SDTotal         string
	SDFree          string
	KernelVersion   string
	PCIEDevice      string
	LTEIMSI         string
	LTECCID         string
	LTECSQ          string
	LANIP           string
	LANMask         string
	LANMAC          string
	WANIP           string
	WANMask         string
	WANGateway      string
	WANMAC          string
	WANPDNS         string
	WANSDNS         string
	APSSID          string
	APKey           string
	WifiSSID        string
	WifiPassword    string
}

// SystemClient system client
type SystemClient struct {
	Model string `json:"model"` // mu model

	mtx sync.Mutex

	// cache
	Channels    SystemChannel `json:"channels"`
	SampleCount int           `json:"samplecount"` // set reacquire interval
}

// DecodeSystemPayload decode system payload
func DecodeSystemPayload(payload string) (*sys.HMUSystemReq, error) {
	var p sys.HMUSystemReq
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// EncodeSystemPayload encode system response
func EncodeSystemPayload(data *sys.HMUSystemResp) (string, error) {
	bytedata, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytedata), nil
}

// DecodeSystemBindingRequest decode system binding request
func DecodeSystemBindingRequest(payload string) (public.SystemBindingPayload, error) {
	var p public.SystemBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeSystemOperationRequest decode system operation request
func DecodeSystemOperationRequest(payload string) (public.SystemOperationPayload, error) {
	var p public.SystemOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// NewSystemClient new system client
func NewSystemClient(req public.SystemBindingPayload) (*SystemClient, error) {
	return &SystemClient{Model: req.Model}, nil
}

// SystemSample system sample
func (sc *SystemClient) SystemSample() error {
	// TODO: 确认是否需要多次重连
	// new system client, and then request for data
	cfg := sys.GetBusManagerCfg()
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	// id
	resp, err := client.UUID()
	if err != nil {
		return fmt.Errorf("query uuid info failed, errmsg {%v}", err)
	}

	sc.Channels.ID = resp.UUID

	// gps

	resp, err = client.GPS()
	if err != nil {
		return fmt.Errorf("query gps info failed, errmsg {%v}", err)
	}

	sc.Channels.Latitude = resp.GPSLatitude
	sc.Channels.Longitude = resp.GPSLongitude
	sc.Channels.UTCTime = resp.GPSTime

	// time
	resp, err = client.Time()
	if err != nil {
		return fmt.Errorf("query time info failed, errmsg {%v}", err)
	}

	sc.Channels.UpTime = resp.UpTime
	sc.Channels.CurrentTime = resp.CurrentTime
	sc.Channels.TimeServer = resp.TimeServer

	// ltecsq
	resp, err = client.LTECSQ()
	if err != nil {
		// 这里的4G可以没有开启
		log.Debug(errors.As(err))
		resp = sys.DefaultHMUSystemResp
	}
	sc.Channels.LTECSQ = resp.LTECSQ

	// system
	resp, err = client.SystemInfo()
	if err != nil {
		return fmt.Errorf("query system info failed, errmsg {%v}", err)
	}

	sc.Channels.CPU = resp.CPU
	sc.Channels.Model = resp.Model
	sc.Channels.SoftwareVersion = resp.SWVersion
	sc.Channels.HardwareVersion = resp.HWVersion
	sc.Channels.RAMTotal = resp.RAMTotal
	sc.Channels.RAMFree = resp.RAMFree
	sc.Channels.FlashTotal = resp.FlashTotal
	sc.Channels.FlashFree = resp.FlashFree
	sc.Channels.SDTotal = resp.SDTotal
	sc.Channels.SDFree = resp.SDFree
	sc.Channels.KernelVersion = resp.KernelVersion
	sc.Channels.PCIEDevice = resp.PCIEDevice
	sc.Channels.LTEIMSI = resp.LTEIMSI
	sc.Channels.LTECCID = resp.LTECCID

	// lan
	resp, err = client.LAN()
	if err != nil {
		return fmt.Errorf("query lan info failed, errmsg {%v}", err)
	}

	sc.Channels.LANIP = resp.LANIP
	sc.Channels.LANMask = resp.LANMask
	sc.Channels.LANMAC = resp.LANMAC

	// wan
	resp, err = client.WAN()
	if err != nil {
		return fmt.Errorf("query lan info failed, errmsg {%v}", err)
	}

	sc.Channels.WANIP = resp.WANIP
	sc.Channels.WANMask = resp.WANMask
	sc.Channels.WANMAC = resp.WANMAC

	// wireless
	resp, err = client.Wireless()
	if err != nil {
		return fmt.Errorf("query wireless info failed, errmsg {%v}", err)
	}

	sc.Channels.APSSID = resp.SSID
	sc.Channels.APKey = resp.Key

	// internet
	resp, err = client.Internet()
	if err != nil {
		return fmt.Errorf("query internet info failed, errmsg {%v}", err)
	}

	if resp.InternetMode == sys.NetModeETH {
		sc.Channels.NetMode = resp.WANMode
	} else {
		sc.Channels.NetMode = resp.InternetMode
	}

	sc.Channels.WANGateway = resp.WANGateway
	sc.Channels.WANPDNS = resp.WANPDNS
	sc.Channels.WANSDNS = resp.WANSDNS
	sc.Channels.WifiSSID = resp.WIFISSID
	sc.Channels.WifiPassword = resp.WIFIPassword

	return nil
}

// ID client's id
func (sc *SystemClient) ID() string {
	return sc.Model
}

// Sample system sample
func (sc *SystemClient) Sample(payload string) (string, error) {
	// decode payload
	req, err := DecodeSystemOperationRequest(payload)
	if err != nil {
		return "", fmt.Errorf("decode system payload failed, errmsg {%v}", err)
	}

	{
		sc.mtx.Lock()
		defer sc.mtx.Unlock()

		if sc.SampleCount%req.Quantity == 0 {
			sc.SystemSample()
			log.Println(sc.Channels, sc.SampleCount)
		}

		sc.SampleCount++
	}

	// query data
	result := ""
	switch req.Channel {
	case channelID:
		result = sc.Channels.ID
	case channelNetMode:
		result = sc.Channels.NetMode
	case channelLongitude:
		result = sc.Channels.Longitude
	case channelLatitude:
		result = sc.Channels.Latitude
	case channelGPS:
		result = sc.Channels.Longitude + "," + sc.Channels.Latitude
	case channelUTCTime:
		result = sc.Channels.UTCTime
	case channelUpTime:
		result = sc.Channels.UpTime
	case channelCurrentTime:
		result = sc.Channels.CurrentTime
	case channelTimeServer:
		result = sc.Channels.TimeServer
	case channelCPU:
		result = sc.Channels.CPU
	case channelModel:
		result = sc.Channels.Model
	case channelSoftwareVersion:
		result = sc.Channels.SoftwareVersion
	case channelHardwareVersion:
		result = sc.Channels.HardwareVersion
	case channelRAMTotal:
		result = sc.Channels.RAMTotal
	case channelRAMFree:
		result = sc.Channels.RAMFree
	case channelFlashTotal:
		result = sc.Channels.FlashTotal
	case channelFlashFree:
		result = sc.Channels.FlashFree
	case channelSDTotal:
		result = sc.Channels.SDTotal
	case channelSDFree:
		result = sc.Channels.SDFree
	case channelKernelVersion:
		result = sc.Channels.KernelVersion
	case channelPCIEDevice:
		result = sc.Channels.PCIEDevice
	case channelLTEIMSI:
		result = sc.Channels.LTEIMSI
	case channelLTECCID:
		result = sc.Channels.LTECCID
	case channelLTECSQ:
		result = sc.Channels.LTECSQ
	case channelLANIP:
		result = sc.Channels.LANIP
	case channelLANMask:
		result = sc.Channels.LANMask
	case channelLANMAC:
		result = sc.Channels.LANMAC
	case channelWANIP:
		result = sc.Channels.WANIP
	case channelWANMask:
		result = sc.Channels.WANMask
	case channelWANGateway:
		result = sc.Channels.WANGateway
	case channelWANMAC:
		result = sc.Channels.WANMAC
	case channelWANPDNS:
		result = sc.Channels.WANPDNS
	case channelWANSDNS:
		result = sc.Channels.WANSDNS
	case channelAPSSID:
		result = sc.Channels.APSSID
	case channelAPKey:
		result = sc.Channels.APKey
	case channelWifiSSID:
		result = sc.Channels.WifiSSID
	case channelWifiPassword:
		result = sc.Channels.WifiPassword
	}

	return result, nil
}

// Command system command
func (sc *SystemClient) Command(payload string) (string, error) {
	// decode payload
	data, err := DecodeSystemOperationRequest(payload)
	if err != nil {
		return "", fmt.Errorf("decode system payload failed, errmsg {%v}", err)
	}

	// new system client, and then request for data
	cfg := sys.GetBusManagerCfg()
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)

	// set
	req := data.Request
	var resp *sys.HMUSystemResp
	switch req.Type {
	case sys.SystemCommandTypeTime:
		resp, err = client.SetTimeServer(req.TimeServer)
		if err != nil {
			return "", fmt.Errorf("set time info failed, errmsg {%v}", err)
		}
	case sys.SystemCommandTypeReboot:
		resp, err = client.Reboot()
		if err != nil {
			return "", fmt.Errorf("reboot failed, errmsg {%v}", err)
		}
	case sys.SystemCommandTypeUUID:
		resp, err = client.SetUUID(req.UUID)
		if err != nil {
			return "", fmt.Errorf("set uuid info failed, errmsg {%v}", err)
		}
	case sys.SystemCommandTypeAppLED:
		resp, err = client.SetAppLED(req.IsON)
		if err != nil {
			return "", fmt.Errorf("set app led failed, errmsg {%v}", err)
		}
	case sys.SystemCommandTypeInternet:
		if req.InternetMode == sys.NetModeETH && req.WANMode == sys.WANModeStatic {
			resp, err = client.SetEthStatic(req.WANIP, req.WANMask, req.WANGateway, req.WANPDNS, req.WANSDNS)
			if err != nil {
				return "", fmt.Errorf("set eth static mode failed, errmsg {%v}", err)
			}
		} else if req.InternetMode == sys.NetModeETH && req.WANMode == sys.WANModeDHCP {
			resp, err = client.SetEthDHCP()
			if err != nil {
				return "", fmt.Errorf("set eth dhcp mode failed, errmsg {%v}", err)
			}
		} else if req.InternetMode == sys.NetModeLTE {
			resp, err = client.SetLTE()
			if err != nil {
				return "", fmt.Errorf("set lte mode failed, errmsg {%v}", err)
			}
		} else if req.InternetMode == sys.NetModeWIFI {
			resp, err = client.SetWifi(req.WifiSSID, req.WifiKey)
			if err != nil {
				return "", fmt.Errorf("set wifi mode failed, errmsg {%v}", err)
			}
		}

		// reboot
		// client.Reboot()

		// reboot by tty
		// cmd := exec.Command("reboot")
		// if err := cmd.Run(); err != nil {
		// 	return "", err
		// }
	}

	return EncodeSystemPayload(resp)
}
