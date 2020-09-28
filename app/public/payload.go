/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: modbus payload
 *
 */

package public

import (
	"encoding/json"

	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// 当发生采集时，响应此通用结构体用于是否向服务器发送数据
type SamplePayload struct {
	Send bool                   `json:"send"`   // 是否发送给服务器
	Data map[string]interface{} `json:"result"` // 需要发送的数据, key值是发送到mqtt的chanel通道
}

func ParseSamplePayload(data string) (*SamplePayload, error) {
	p := &SamplePayload{}
	if err := json.Unmarshal([]byte(data), p); err != nil {
		return nil, errors.As(err, data)
	}
	return p, nil
}

func NewSamplePayload(send bool) *SamplePayload {
	return &SamplePayload{
		Send: send,
		Data: map[string]interface{}{},
	}
}

func (p *SamplePayload) PutData(key string, val interface{}) {
	p.Data[key] = val
}

func (p *SamplePayload) Serial() string {
	data, err := json.Marshal(p)
	if err != nil {
		panic(*p)
	}
	return string(data)
}

// MessagePayload message payload
type MessagePayload struct {
	MonitoringUnitID string      `json:"monitoringUnitId"`
	SampleUnitID     string      `json:"sampleUnitId"`
	ChannelID        string      `json:"channelId"`
	Name             string      `json:"name"`
	Value            interface{} `json:"value"`
	Timestamp        string      `json:"timestamp,omitempty"`
	Cov              bool        `json:"cov"`
	State            int         `json:"state"`
}

// ModbusPayload is used for decoding request payload
type ModbusPayload struct {
	// shared, including bind, sample, command and
	Port    string `json:"port,omitempty"`
	Slaveid int32  `json:"slaveid,omitempty"`

	// binding port used
	BaudRate int32 `json:"baudrate,omitempty"`

	// Sensorflow use
	KeyNumber     int32  `json:"keyNumber,omitempty"`
	MUID          string `json:"muid,omitempty"`
	SUID          string `json:"suid,omitempty"`
	WANInterface  string `json:"WANInterface,omitempty"`
	WifiInterface string `json:"WifiInterface,omitempty"`

	// sample used
	Code     int32 `json:"code,omitempty"`
	Address  int32 `json:"address,omitempty"`
	Quantity int32 `json:"quantity,omitempty"`
	Timeout  int32 `json:"timeout,omitempty"`

	// command used
	Value interface{} `json:"value,omitempty"`

	// sensorflow used
	Mode       string `json:"mode,omitempty"`
	ColorTable string `json:"colorTable,omitempty"`
}

// SystemBindingPayload system binding request
type SystemBindingPayload struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Model string `json:"model"` // mu type, such as "hmu2000"
}

// SystemOperationPayload system operation request
type SystemOperationPayload struct {
	// sample request
	Model    string `json:"model"`
	Channel  string `json:"channel"`
	Quantity int    `json:"quantity"`

	// command request
	Request sys.HMUSystemReq
}

// PMBUSBindingPayload pmbus binding payload
type PMBUSBindingPayload struct {
	BaudRate int  `json:"baudrate"`
	Timeout  int  `json:"timeout"`
	SOI      byte `json:"soi"`
	VER      byte `json:"ver"`
	ADR      byte `json:"adr"`
	CID1     byte `json:"cid1"`
	EOI      byte `json:"eoi"`
}

// PMBUSOperationPayload pmbus operation payload
type PMBUSOperationPayload struct {
	ADR         byte   `json:"adr"`
	CID1        byte   `json:"cid1"`
	CID2        byte   `json:"cid2"`
	LENID       uint16 `json:"lenid"`
	COMMANDTYPE byte   `json:"commandtype"`
	COMMANDID   byte   `json:"commandid"`
}

// DeltaUPSBindingPayload delta ups binding payload
type DeltaUPSBindingPayload struct {
	BaudRate int  `json:"baudrate"`
	Timeout  int  `json:"timeout"`
	Header   byte `json:"header"`
	ID       int  `json:"id"`
}

// DeltaUPSOperationPayload delta ups operation payload
type DeltaUPSOperationPayload struct {
}

// OilMachineBindingPayload oil machine binding payload
type OilMachineBindingPayload struct {
	BaudRate int  `json:"baudrate"`
	Timeout  int  `json:"timeout"`
	SOI      byte `json:"soi"`
	EOI      byte `json:"eoi"`
}

// OilMachineOperationPayload oil machine operation payload
type OilMachineOperationPayload struct {
	ADR         byte   `json:"adr"`
	CID1        byte   `json:"cid1"`
	CID2        byte   `json:"cid2"`
	LENGTH      byte   `json:"length"`
	COMMANDINFO []byte `json:"info"`
}

// CommonSerialBindingPayload commond use serial binding payload
type CommonSerialBindingPayload struct {
	BaudRate int `json:"baudrate"`
	Timeout  int `json:"timeout"`
}

// LampWithOperationPayload lamp with operation payload
type LampWithOperationPayload struct {
	Value interface{} `json:"value"`
}

// CreditCardBindingPayload credit card binding payload
type CreditCardBindingPayload struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	SerialNumber string `json:"serialnum"`
}

// LuMiGatewayBindingPayload lumi gateway binding payload
type LuMiGatewayBindingPayload struct {
	SID          string `json:"sid"`
	Password     string `json:"password"`
	NetInterface string `json:"netinterface"`
}

// LuMiGatewayOperationPayload lumi gateway operation payload
type LuMiGatewayOperationPayload struct {
	SID   string `json:"sid"`
	Model string `json:"model"`
	Value string `json:"value"`
}

// FaceIPCBindingPayload face ipc binding payload
type FaceIPCBindingPayload struct {
	Host         string `json:"host"`         // listen host
	Port         string `json:"port"`         // listen port
	MUID         string `json:"muid"`         // monitoring unit id
	SUID         string `json:"suid"`         // sample unit id
	UploadServer string `json:"uploadServer"` // upload param: server address
	Author       string `json:"author"`       // upload param: author
	Project      string `json:"project"`      // upload param: project
	Token        string `json:"token"`        // upload param: token
	User         string `json:"user"`         // upload param: user
}

// FaceIPCOperationPayload face ipc operation payload
type FaceIPCOperationPayload struct {
	CameraID string `json:"camera_id"` // on behalf of camera
	Perpose  string `json:"perpose"`   // register or unregister
	FaceID   string `json:"face_id"`   // on behalf of person
	FaceURL  string `json:"face_url"`  // url of face image

	Flag          int    `json:"flag"`           // server address flag [1: ip; 2: domain]
	AddressLength int    `json:"address_length"` // server address length, max 50 bytes
	Address       string `json:"address"`        // server address
	Port          int    `json:"port"`           // server port (or destination port for transmission)
	DeviceID      int    `json:"deviceid"`       // device id

	Data        string `json:"data"`        // data for upgrade or transmission
	Destination string `json:"destination"` // destination host (IP) for transmission
	Mode        int    `json:"mode"`        // mode for transmission, only support udp(0) and tcp(1)
}

// SNMPBindingPayload snmp binding payload
type SNMPBindingPayload struct {
	Version        string `json:"version"`        // version
	Target         string `json:"target"`         // target
	Port           string `json:"port"`           // target port
	ReadCommunity  string `json:"readCommunity"`  // read community
	WriteCommunity string `json:"writeCommunity"` // write community
}

// SNMPOperationPayload snmp operation payload
type SNMPOperationPayload struct {
	Target string   `json:"target"` // target
	OIDS   []string `json:"oids"`

	OID   string      `json:"oid"`
	Value interface{} `json:"value"`
}

// WeiGengEntryBindingPayload snmp binding payload
type WeiGengEntryBindingPayload struct {
	LocalAddress string `json:"localhostAddress"` // localhost address
	LocalPort    string `json:"localhostPort"`    // localhost port
	DoorAddress  string `json:"doorAddress"`      // door address
	DoorPort     string `json:"doorPort"`         // door port
	SerialNumber string `json:"serialNo"`         // device serial number
}

// WeiGengEntryOperationPayload snmp operation payload
type WeiGengEntryOperationPayload struct {
	SequenceNumber string `json:"seqno"`
	FunctionID     string `json:"code"`
	Group          int    `json:"group"`
	Type           string `json:"type"`

	Door int `json:"door"`

	CardID       string `json:"cardNo"`
	UserID       string `json:"userId"`
	UserPassword string `json:"userPassword"`
	ExpireDate   string `json:"expireDate"`
	CardType     int    `json:"cardType"`
	CardStatus   int    `json:"byCardValid"`
	DoorRight    string `json:"byDoorRight"`
}

// DIDOOperationPayload dido operation payload
type DIDOOperationPayload struct {
	Value interface{} `json:"value"`
}

// ES5200OperationPayload es5200 operation payload
type ES5200OperationPayload struct {
	ADR          byte   `json:"adr"`
	CID1         byte   `json:"cid1"`
	CID2         byte   `json:"cid2"`
	LENID        int    `json:"lenid"`
	COMMANDGROUP byte   `json:"commandgroup"`
	COMMANDTYPE  byte   `json:"commandtype"`
	INFO         []byte `json:"info"`

	CardID     int    `json:"cardId"`
	UserID     int    `json:"userId"`
	Password   string `json:"password"`
	ExpireDate string `json:"expireDate"`
	Permission int    `json:"permission"`
}

// CameraBindingPayload camera binding payload
// TODO: 编写协议文档
type CameraBindingPayload struct {
	Host     string `json:"host"`
	User     string `json:"username"`
	Password string `json:"password"`
}

// CameraOperationPayload camera operation payload
// TODO: 编写协议文档
type CameraOperationPayload struct {
	Host    string `json:"host"`
	Channel string `json:"channel"`

	// 当channel为sample时启用以下参数
	// sample_image, sample_vedio
	SampleCmd string `json:"sample_cmd"`
	SampleArg string `json:"sample_arg"`
}

// ElecFireBindingPayload elec fire binding payload
type ElecFireBindingPayload struct {
	Host string `json:"host"` // listen host
	Port string `json:"port"` // listen port
}

// ElecFireOperationPayload elec fire operation payload
type ElecFireOperationPayload struct {
	SerialNumber string `json:"serialnum"`
	Address      int    `json:"address"`
	Length       int    `json:"length"`
}

type VirtualAntennaBindingPayload struct {
	Host string `json:"host"` // listen host
	Port string `json:"port"` // listen port
}

type VirtualAntennaOperationPayload struct {
	Channel string `json:"channel"`
	Perpose string `json:"perpose"` // register or unregister
	Value   int    `json:"value"`
}
