/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/11/07
 * Despcription: lumi gateway client define
 *
 */

package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolLuMiGateway, generalLuMiGatewayDriverProtocol)
}

// Implement DriverProtocol
type lumiGatewayDriverProtocol struct {
	req *public.LuMiGatewayBindingPayload
	uri string
}

func (dp *lumiGatewayDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *lumiGatewayDriverProtocol) ClientID() string {
	return LuMiGatewayClientID
}

func (dp *lumiGatewayDriverProtocol) NewInstance() (PortClient, error) {
	return NewLuMiGatewayClient(
		dp.req.SID, dp.req.Password,
		dp.req.NetInterface,
	)
}

func generalLuMiGatewayDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeLuMiGatewayBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &lumiGatewayDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

/**

message example:

whois: 			{"cmd":"iam","port":"9898","sid":"7811dcb78b2f","model":"gateway","proto_version":"1.1.2","ip":"192.168.1.106"}
get_id_list: 	{"cmd":"get_id_list_ack","sid":"7811dcb78b2f","token":"XRwvMrwu16ZZWxPe","data":"[\"158d0002322bf7\",\"158d0002325b7e\",\"158d0002371943\",\"158d0002325aa7\",\"158d0001d8f1b6\",\"158d00022ca896\"]"}
heartbeat: 		{"cmd":"heartbeat","model":"gateway","sid":"7811dcb78b2f","short_id":"0","token":"oKlTd40GxrUa85wt","data":"{\"ip\":\"192.168.1.106\"}"}
report: 		{"cmd":"report","model":"motion","sid":"158d0002325b7e","short_id":43719,"data":"{\"no_motion\":\"1200\"}"}
read: 			{"cmd":"read_ack","model":"magnet","sid":"158d0002322bf7","short_id":30033,"data":"{\"voltage\":3015,\"status\":\"open\"}"}
read:			{"cmd":"read_ack","model":"motion","sid":"158d0002325b7e","short_id":43719,"data":"{\"voltage\":2995}"}
read:			{"cmd":"read_ack","model":"switch","sid":"158d0002371943","short_id":14851,"data":"{\"voltage\":3062}"}
read:			{"cmd":"read_ack","model":"plug","sid":"158d0002325aa7","short_id":7071,"data":"{\"voltage\":3600,\"status\":\"off\",\"inuse\":\"0\",\"power_consumed\":\"117\",\"load_power\":\"0.00\"}"}
read:			{"cmd":"read_ack","model":"smoke","sid":"158d0001d8f1b6","short_id":4847,"data":"{\"voltage\":3235,\"alarm\":\"0\"}"}
read:			{"cmd":"read_ack","model":"weather.v1","sid":"158d00022ca896","short_id":64336,"data":"{\"voltage\":2985,\"temperature\":\"2651\",\"humidity\":\"6565\",\"pressure\":\"100790\"}"}

**/

// LuMiGatewayClientID id
var LuMiGatewayClientID = "lumi-gateway-client-id"
var aesKeyIV = []byte{0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58, 0x56, 0x2e}

const (
	lumiCommandTypeWhois     = "whois"
	lumiCommandTypeGetIDList = "get_id_list"
	lumiCommandTypeRead      = "read"
	lumiCommandTypeWrite     = "write"

	lumiCommandResponseIam          = "iam"
	lumiCommandResponseGetIDListAck = "get_id_list_ack"
	lumiCommandResponseReadAck      = "read_ack"
	lumiCommandResponseWriteAck     = "write_ack"

	lumiCommandResponseHeartbeat = "heartbeat"
	lumiCommandResponseReport    = "report"
)

const (
	lumiGateway            = "gateway"
	lumiSensorMagnet       = "magnet"
	lumiSensorMotion       = "motion"
	lumiSensorSwitch       = "switch"
	lumiSensorPlug         = "plug"
	lumiSensorCtrlNeutral1 = "ctrl_neutral1"
	lumiSensorCtrlNeutral2 = "ctrl_neutral2"
	lumiSensor86SW1        = "86sw1"
	lumiSensor86SW2        = "86sw2"
	lumiSensorHT           = "sensor_ht"
	lumiSensorCube         = "cube"
	lumiSensorCurtain      = "curtain"
	lumiSensorCtrlLn1      = "ctrl_ln1"
	lumiSensorCtrlLn2      = "ctrl_ln2"
	lumiSensor86Plug       = "86plug"
	lumiSensorNatgas       = "natgas"
	lumiSensorSmoke        = "smoke"
	lumiSensorMagnetAq2    = "sensor_magnet.aq2"
	lumiSensorMotionAq2    = "sensor_motion.aq2"
	lumiSensorSwitchAq2    = "sensor_switch.aq2"
	lumiSensorWeatherV1    = "weather.v1"
	lumiSensorWleakAq1     = "sensor_wleak.aq1"
	lumiSensorLockAq1      = "lock.aq1"
)

type lumiRequest struct {
	Command string `json:"cmd,omitempty"`
	Model   string `json:"model,omitempty"`
	SID     string `json:"sid,omitempty"`
	ShortID int    `json:"short_id,omitempty"`
	Data    string `json:"data,omitempty"`
}

type lumiResponse struct {
	Command      string `json:"cmd,omitempty"`
	Port         string `json:"port,omitempty"`
	SID          string `json:"sid,omitempty"`
	Model        string `json:"model,omitempty"`
	ProtoVersion string `json:"proto_version,omitempty"`
	IP           string `json:"ip,omitempty"`

	// get id list ack
	Token string `json:"token,omitempty"`
	Data  string `json:"data,omitempty"`

	// heartbeat
	ShortID interface{} `json:"short_id,omitempty"`
}

type lumiData struct {
	RGB            uint   `json:"rgb,omitempty"`
	Illumination   int    `json:"illumination,omitempty"`
	ProtoVersion   string `json:"proto_version,omitempty"`
	MID            int    `json:"mid,omitempty"`
	JoinPermission string `json:"join_permission,omitempty"`
	RemoveDevice   string `json:"remove_device,omitempty"`
	Voltage        int    `json:"voltage,omitempty"`
	Status         string `json:"status,omitempty"`
	LoadVoltage    string `json:"load_voltage,omitempty"`
	LoadPower      string `json:"load_power,omitempty"`
	PowerConsumed  string `json:"power_consumed,omitempty"`
	InUse          string `json:"inuse,omitempty"`
	Alarm          string `json:"alarm,omitempty"`
	Temperature    string `json:"temperature,omitempty"`
	Humidity       string `json:"humidity,omitempty"`
	Pressure       string `json:"pressure,omitempty"`
	Channel0       string `json:"channel0,omitempty"`
	Channel1       string `json:"channel1,omitempty"`
	Rotate         string `json:"rotate,omitempty"`
	CurtainLevel   string `json:"curtain_level,omitempty"`
	Lux            string `json:"lux,omitempty"`
	VerifiedWrong  string `json:"verified_wrong,omitempty"`
	FingVerified   string `json:"fing_verified,omitempty"`

	// for write command
	ShortID int    `json:"short_id,omitempty"`
	Key     string `json:"key,omitempty"`
}

type lumiSampleValue struct {
	SID  string `json:"sid"`
	Data string `json:"data"`
}

type lumiControl struct {
	SID     string `json:"sid"`
	Model   string `json:"model"`
	ShortID int    `json:"short_id"`
	Data    string `json:"data"`
}

// LuMiGatewayClient client
type LuMiGatewayClient struct {
	ClientID string

	conn *net.UDPConn
	mtx  sync.Mutex

	gatewayid    string
	password     string
	netinterface string

	token string

	bgetdest chan bool
	dest     string

	bgetidlist chan bool
	idlist     []string

	bcontrol    chan bool
	controldata lumiControl

	devices map[string][]string // cache devices' sid by model
	data    map[string]lumiData // cache devices' data by sid
}

func lumiTransformChannelname(name string) string {
	switch name {
	case "rgb":
		return "RGB"
	case "illumination":
		return "Illumination"
	case "proto_version":
		return "ProtoVersion"
	case "mid":
		return "MID"
	case "join_permission":
		return "JoinPermission"
	case "remove_device":
		return "RemoveDevice"
	case "voltage":
		return "Voltage"
	case "status":
		return "Status"
	case "load_voltage":
		return "LoadVoltage"
	case "load_power":
		return "LoadPower"
	case "power_consumed":
		return "PowerConsumed"
	case "inuse":
		return "InUse"
	case "alarm":
		return "Alarm"
	case "temperature":
		return "Temperature"
	case "humidity":
		return "Humidity"
	case "pressure":
		return "Pressure"
	case "channel0":
		return "Channel0"
	case "channel1":
		return "Channel1"
	case "rotate":
		return "Rotate"
	case "curtain_level":
		return "CurtainLevel"
	case "lux":
		return "Lux"
	case "verified_wrong":
		return "VerifiedWrong"
	case "fing_verified":
		return "FingVerified"
	}

	return ""
}

// NewLuMiGatewayClient new lumi gateway client
func NewLuMiGatewayClient(sid, password, netinterface string) (PortClient, error) {
	var client = LuMiGatewayClient{
		ClientID:     LuMiGatewayClientID,
		gatewayid:    sid,
		password:     password,
		netinterface: netinterface,
	}

	client.bgetdest = make(chan bool)
	client.bgetidlist = make(chan bool)
	client.bcontrol = make(chan bool)
	client.devices = make(map[string][]string)
	client.data = make(map[string]lumiData)

	go client.Start()

	return &client, nil
}

// DecodeLuMiGatewayBindingPayload decode lumi gateway payload
func DecodeLuMiGatewayBindingPayload(payload string) (public.LuMiGatewayBindingPayload, error) {
	var p public.LuMiGatewayBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeLuMiGatewayOperatePayload decode lumi gateway payload
func DecodeLuMiGatewayOperatePayload(payload string) (public.LuMiGatewayOperationPayload, error) {
	var p public.LuMiGatewayOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// AesEncrypt encrypt
func (lc *LuMiGatewayClient) AesEncrypt(encodeStr string, key []byte) (string, error) {
	encodeBytes := []byte(encodeStr)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCEncrypter(block, aesKeyIV)
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return fmt.Sprintf("%x", crypted), nil
}

// PackCommand pack command
func (lc *LuMiGatewayClient) PackCommand(cmd, model, sid string, shortid int, data string) ([]byte, error) {
	var req = lumiRequest{
		Command: cmd,
		Model:   model,
		SID:     sid,
		ShortID: shortid,
		Data:    data,
	}

	d, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Read read devices' data
func (lc *LuMiGatewayClient) Read(idlist []string, s *net.UDPConn, addr *net.UDPAddr) {
	for {
		for _, id := range idlist {
			d, err := lc.PackCommand(lumiCommandTypeRead, "", id, 0, "")
			if err != nil {
				fmt.Printf("pack read command failed: %v", err)
				continue
			}

			s.WriteToUDP(d, addr)
		}

		time.Sleep(time.Second)
	}
}

// Start start
func (lc *LuMiGatewayClient) Start() error {
	// get specified gateway address first
	dest, err := lc.getGatewayIP()
	if err != nil {
		fmt.Printf("get gateway ip failed: %v", err)
		return err
	}

	// start listen
	s, err := lc.listen()
	if err != nil {
		fmt.Printf("listen failed: %v", err)
		return err
	}

	lc.conn = s

	addr, err := net.ResolveUDPAddr("udp4", dest)
	if err != nil {
		return err
	}

	// start get id list
	go func(c *LuMiGatewayClient, s *net.UDPConn, addr *net.UDPAddr) {
	IDLIST:
		for {
			select {
			case <-c.bgetidlist:
				// start read devices' data when gain device list
				go c.Read(lc.idlist, s, addr)

				break IDLIST
			default:
				// get id list
				d, err := c.PackCommand(lumiCommandTypeGetIDList, "", "", 0, "")
				if err != nil {
					fmt.Printf("pack read command failed: %v", err)
					continue
				}

				s.WriteToUDP(d, addr)

				time.Sleep(time.Second)
			}
		}
	}(lc, s, addr)

	return nil
}

// whois command discovery gateways, gain gateway's ip and so on
func (lc *LuMiGatewayClient) getGatewayIP() (string, error) {
	src := "224.0.0.50:9999"
	dest := "224.0.0.50:4321"

	laddr, err := net.ResolveUDPAddr("udp4", src)
	if err != nil {
		return "", err
	}

	// get specified net interface
	var ift net.Interface
	ifts, _ := net.Interfaces()
	for _, i := range ifts {
		if i.Name == lc.netinterface {
			ift = i
		}
	}

	// multicast listen
	socket, err := net.ListenMulticastUDP("udp4", &ift, laddr)
	if err != nil {
		return "", err
	}

	go func(c *LuMiGatewayClient, s *net.UDPConn) {
		defer s.Close()
		for {
			time.Sleep(time.Second)

			// send whois multicast
			addr, err := net.ResolveUDPAddr("udp4", dest)
			if err != nil {
				fmt.Println(err)
				continue
			}

			data, err := c.PackCommand(lumiCommandTypeWhois, "", "", 0, "")
			if err != nil {
				fmt.Println(err)
				continue
			}

			_, err = s.WriteToUDP(data, addr)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// wait for response
			d := make([]byte, 1024)
			s.SetDeadline(time.Now().Add(time.Second * 3))
			read, _, err := s.ReadFromUDP(d)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// parse
			var resp lumiResponse
			if err := json.Unmarshal(d[:read], &resp); err != nil {
				fmt.Println(err)
				continue
			}

			if resp.SID != c.gatewayid {
				fmt.Printf("get sid [%s], want [%s], inconformity\n", resp.SID, c.gatewayid)
				continue
			}

			c.dest = (resp.IP + ":" + resp.Port)
			lc.bgetdest <- true
			break
		}
	}(lc, socket)

	// wait
	<-lc.bgetdest

	return lc.dest, nil
}

func (lc *LuMiGatewayClient) listen() (*net.UDPConn, error) {
	laddr, err := net.ResolveUDPAddr("udp4", "224.0.0.50:9898")
	if err != nil {
		return nil, err
	}

	var ift net.Interface
	ifts, _ := net.Interfaces()
	for _, i := range ifts {
		if i.Name == lc.netinterface {
			ift = i
		}
	}

	socket, err := net.ListenMulticastUDP("udp4", &ift, laddr)
	if err != nil {
		return nil, err
	}

	go lc.gainAndParse(socket)
	go lc.control(socket)

	return socket, nil
}

// gain and parse data
func (lc *LuMiGatewayClient) gainAndParse(s *net.UDPConn) {
	for {
		d := make([]byte, 1024)
		read, _, err := s.ReadFromUDP(d)
		if err != nil {
			fmt.Printf("read from udp failed: %v", err)
			continue
		}

		var resp lumiResponse
		if err := json.Unmarshal(d[:read], &resp); err != nil {
			fmt.Printf("unmarshal response failed: %v", err)
			continue
		}

		switch resp.Command {
		case lumiCommandResponseGetIDListAck:
			{
				var sids []string
				if err := json.Unmarshal([]byte(resp.Data), &sids); err != nil {
					fmt.Printf("unmarshal get id list faield: %v", err)
					continue
				}

				lc.token = resp.Token
				lc.idlist = sids

				lc.bgetidlist <- true
			}
		case lumiCommandResponseHeartbeat:
			{
				lc.token = resp.Token
			}
		case lumiCommandResponseReport:
			{
				// TODO push data to upper server
				// fmt.Printf("report: %s\n", d[:read])
			}
		case lumiCommandResponseReadAck:
			{
				var d lumiData
				if err := json.Unmarshal([]byte(resp.Data), &d); err != nil {
					fmt.Println(err)
					continue
				}

				// cache sid
				ids, ok := lc.devices[resp.Model]
				if !ok {
					lc.devices[resp.Model] = append(lc.devices[resp.Model], resp.SID)
				} else {
					exist := false
					for _, id := range ids {
						if id == resp.SID {
							exist = true
							break
						}
					}

					if !exist {
						lc.devices[resp.Model] = append(lc.devices[resp.Model], resp.SID)
					}
				}

				// cache data
				d.ShortID = int(resp.ShortID.(float64))
				lc.data[resp.SID] = d
			}
		default:
			fmt.Printf("response: %s\n", d[:read])
		}

	}
}

// control control devices
func (lc *LuMiGatewayClient) control(s *net.UDPConn) {
	for {
		select {
		case <-lc.bcontrol:
			addr, err := net.ResolveUDPAddr("udp4", lc.dest)
			if err != nil {
				fmt.Printf("resolve udp address failed: %v", err)
				continue
			}

			req, err := lc.PackCommand(lumiCommandTypeWrite, lc.controldata.Model, lc.controldata.SID, lc.controldata.ShortID, lc.controldata.Data)
			if err != nil {
				fmt.Printf("pack command failed: %v", err)
				continue
			}

			fmt.Printf("%v", string(req))

			s.WriteToUDP(req, addr)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// ID id
func (lc *LuMiGatewayClient) ID() string {
	return lc.ClientID
}

// Sample get values
func (lc *LuMiGatewayClient) Sample(payload string) (string, error) {
	p, err := DecodeLuMiGatewayOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("unmarshal command payload failed: %v", err)
	}

	// find model mapping data
	ids, ok := lc.devices[p.Model]
	if !ok {
		return "", nil
	}

	name := lumiTransformChannelname(p.Value)
	if name == "" {
		return "", fmt.Errorf("field [%v] not found", p.Value)
	}

	// get data from cache
	for _, id := range ids {
		if id == p.SID {
			d := lc.data[id]

			object := reflect.ValueOf(&d)
			myref := object.Elem()
			t := myref.Type()

			nf := myref.NumField()
			for i := 0; i < nf; i++ {
				field := myref.Field(i)

				// value must be data's filed
				if t.Field(i).Name == name {
					r := ""
					v := field.Interface()

					switch field.Type().String() {
					case "string":
						r = v.(string)
					case "int":
						tmp := v.(int)
						r = strconv.Itoa(tmp)
					case "uint":
						tmp := v.(uint)
						r = strconv.FormatUint(uint64(tmp), 10)
					}

					return r, nil
				}
			}

			return "", fmt.Errorf("field [%v] not found", p.Value)
		}
	}

	return "", fmt.Errorf("slave device [%v] do not exist", p.SID)
}

// Command set values
func (lc *LuMiGatewayClient) Command(payload string) (string, error) {
	lc.mtx.Lock()
	defer lc.mtx.Unlock()

	p, err := DecodeLuMiGatewayOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("unmarshal command payload failed: %v", err)
	}

	d, ok := lc.data[p.SID]
	if !ok {
		return "", fmt.Errorf("sid [%s] do not exist", p.SID)
	}

	var val lumiData
	if err := json.Unmarshal([]byte(p.Value), &val); err != nil {
		return "", fmt.Errorf("unmarshal payload value failed: %v", err)
	}

	k, err := lc.AesEncrypt(lc.token, []byte(lc.password))
	if err != nil {
		return "", fmt.Errorf("encrypt key failed: %v", err)
	}

	val.Key = k

	bval, err := json.Marshal(val)
	if err != nil {
		return "", fmt.Errorf("marshal control value failed: %v", err)
	}

	lc.controldata.SID = p.SID
	lc.controldata.Model = p.Model
	lc.controldata.ShortID = d.ShortID
	lc.controldata.Data = string(bval)

	lc.bcontrol <- true

	return "ok", nil
}
