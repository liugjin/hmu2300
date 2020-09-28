/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/04/03
 * Despcription: es5200 entry implement
 *
 */

package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
	"github.com/tarm/serial"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolES5200, generalES5200DriverProtocol)
}

// Implement DriverProtocol
type es5200DriverProtocol struct {
	req *public.CommonSerialBindingPayload
	uri string
}

func (dp *es5200DriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *es5200DriverProtocol) ClientID() string {
	return ES5200ClientID
}

func (dp *es5200DriverProtocol) NewInstance() (PortClient, error) {
	return NewES5200Client(
		dp.uri,
		int(dp.req.BaudRate), int(dp.req.Timeout),
	)
}

func generalES5200DriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeES5200BindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &es5200DriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// ES5200ClientID clientid
var ES5200ClientID = "es5200-client-id"

// ES5200Frame pmbus frame
type ES5200Frame struct {
	SOI    byte
	VER    byte
	ADR    byte
	CID1   byte
	CID2   byte
	LENGTH int16
	INFO   []byte
	CHKSUM int16
	EOI    byte
}

// ES5200Client client
type ES5200Client struct {
	ClientID string

	Port  *serial.Port
	Frame ES5200Frame
}

// NewES5200Client new es5200 client
func NewES5200Client(port string, baudrate, timeout int) (PortClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return &ES5200Client{}, err
	}

	frame := ES5200Frame{
		SOI: 0x7E,
		VER: 0x10,
		EOI: 0x0D,
	}

	return &ES5200Client{ClientID: ES5200ClientID, Frame: frame, Port: sp}, nil
}

// DecodeES5200BindingPayload decode binding payload
func DecodeES5200BindingPayload(payload string) (public.CommonSerialBindingPayload, error) {
	var p public.CommonSerialBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeES5200OperationPayload decode operation payload
func DecodeES5200OperationPayload(payload string) (public.ES5200OperationPayload, error) {
	var p public.ES5200OperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// EncodeRequest encode request
func (ec *ES5200Client) EncodeRequest(req public.ES5200OperationPayload) ([]byte, error) {
	// compute length
	length, err := computeLength(uint16(req.LENID))
	if err != nil {
		return nil, fmt.Errorf("compute length failed, errmsg {%v}", err)
	}

	// pack frame
	frame := []byte{ec.Frame.VER, req.ADR, req.CID1, req.CID2}
	frame = append(frame, length...)

	// plus command info when require
	if req.LENID > 0 {
		frame = append(frame, req.COMMANDGROUP, req.COMMANDTYPE)
		frame = append(frame, req.INFO...)
	}

	// compute chksum
	chksum, err := computeCheckSum(frame)
	if err != nil {
		return nil, fmt.Errorf("compute check sum failed, errmsg {%v}", err)
	}

	frame = append(frame, chksum...)

	// transfer main info to HEX-ASCII bytes
	frame = []byte(strings.ToUpper(hex.EncodeToString(frame)))

	// pack entire frame
	frame = append([]byte{ec.Frame.SOI}, frame...)
	frame = append(frame, ec.Frame.EOI)

	return frame, nil
}

// DecodeResponse decode response
func (ec *ES5200Client) DecodeResponse(resp []byte) ([]byte, error) {
	resplength := len(resp)

	// number of response at least larger than 18
	if resplength < 18 {
		return nil, fmt.Errorf("response incompelete")
	}

	// RTN in 8th to 9th
	rtn := resp[7:9]
	if string(rtn) != "00" {
		return nil, fmt.Errorf("error ocurred, errcode {%v}", string(rtn))
	}

	// LENGTH in 10th to 13th, lenid in 11th to 13th, which is number of data bytes
	length := resp[9:13]

	// verify length
	blength, err := hex.DecodeString(string(length))
	if err != nil {
		return nil, fmt.Errorf("decode response length to string failed, errmsg {%v}", err)
	}

	ilength := binary.BigEndian.Uint16(blength)

	clength, err := computeLength(ilength)
	if err != nil {
		return nil, fmt.Errorf("verify length compute failed, errmsg {%v}", err)
	}

	if string(clength) != string(blength) {
		return nil, fmt.Errorf("lchksum inconsistent, response length {%v}, compute length {%v}", length, clength)
	}

	// CHKSUM trim last charater in response which is EOI, chksum in last four bytes
	chksum := resp[resplength-5 : resplength-1]

	bchksum, err := hex.DecodeString(string(chksum))
	if err != nil {
		return nil, fmt.Errorf("decode response chksum to string failed, errmsg {%v}", err)
	}

	// verify chksum, trim SOI(1), EOI(1), and CHKSUM(4)
	info := resp[1 : resplength-5]
	cchksum, err := computeResponseCheckSum(info)
	if err != nil {
		return nil, fmt.Errorf("verify checksum compute failed, errmsg {%v}", err)
	}

	if string(cchksum) != string(bchksum) {
		return nil, fmt.Errorf("chksum inconsistent, response chksum {%v}, compute chksum {%v}", chksum, cchksum)
	}

	// data in 13th to last character front CHKSUM
	data := resp[13 : resplength-5]

	return data, nil
}

// Request request
func (ec *ES5200Client) Request(req []byte) ([]byte, error) {
	n, err := ec.Port.Write(req)
	if err != nil {
		return nil, err
	}

	if n != len(req) {
		return nil, fmt.Errorf("write incomplete, send %v bytes", n)
	}

	chunks := []byte{}
	endflag := ec.Frame.EOI

	for {
		data := make([]byte, 256)
		n, err := ec.Port.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read response failed, errmsg {%v}", err)
		}

		chunks = append(chunks, data[:n]...)

		if data[n-1] == endflag {
			break
		}
	}

	fmt.Println("resp:", string(chunks))

	return chunks, nil
}

type es5200Record struct {
	CardID    int64  `json:"cardId"`
	TimeStamp string `json:"timestamp"`
	Status    int64  `json:"status"`
	Remark    int64  `json:"remark"`
}

func parseResponse(data []byte, cgroup, ctype byte) (string, error) {
	switch ctype {
	case 0xE0:
		if len(data) < 16 {
			return "", fmt.Errorf("response length illegal")
		}

		// realtime
		return fmt.Sprintf("%s-%s-%s %s:%s:%s", data[:4], data[4:6], data[6:8], data[10:12], data[12:14], data[14:16]), nil
	case 0xE2:
		if len(data) < 28 {
			return "", fmt.Errorf("response length illegal")
		}

		var r es5200Record

		// card id, 5 bytes
		strid := string(data[:10])
		r.CardID, _ = strconv.ParseInt(strid, 16, 32)
		r.TimeStamp = fmt.Sprintf("%s-%s-%s %s:%s:%s", data[10:14], data[14:16], data[16:18], data[18:20], data[20:22], data[22:24])
		r.Status, _ = strconv.ParseInt(string(data[24:26]), 16, 32)
		r.Remark, _ = strconv.ParseInt(string(data[26:28]), 16, 32)

		br, _ := json.Marshal(r)
		return string(br), nil
	}

	return "", nil
}

// ID client's id
func (ec *ES5200Client) ID() string {
	return ec.ClientID
}

// Sample sample implement
func (ec *ES5200Client) Sample(payload string) (string, error) {
	req, err := DecodeES5200OperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode es5200 operation payload failed, errmsg {%v}", err)
	}

	if req.COMMANDGROUP == 0xF2 {
		switch req.COMMANDTYPE {
		case 0xE0:
			req.INFO = []byte{0x00}
		case 0xE2:
		}
	}

	frame, err := ec.EncodeRequest(req)
	if err != nil {
		return "", fmt.Errorf("encode request failed, errmsg {%v}", err)
	}

	fmt.Printf("req: %v\n", string(frame))
	resp, err := ec.Request(frame)
	if err != nil {
		return "", fmt.Errorf("request for data failed, errmsg {%v}", err)
	}

	data, err := ec.DecodeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("decode response failed, errmsg {%v}", err)
	}

	// parse response
	return parseResponse(data, req.COMMANDGROUP, req.COMMANDTYPE)
}

// Command pmbus command implement
func (ec *ES5200Client) Command(payload string) (string, error) {
	req, err := DecodeES5200OperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode es5200 operation payload failed, errmsg {%v}", err)
	}

	// get permission
	var p public.ES5200OperationPayload
	p.ADR = req.ADR
	p.CID1 = req.CID1
	p.CID2 = 0x48
	p.LENID = 14
	p.COMMANDGROUP = 0xF0
	p.COMMANDTYPE = 0xE0
	p.INFO = []byte{0x00, 0x00, 0x00, 0x00, 0x00}

	frame, err := ec.EncodeRequest(p)
	if err != nil {
		return "", fmt.Errorf("encode request failed, errmsg {%v}", err)
	}

	fmt.Printf("req: %v\n", string(frame))
	resp, err := ec.Request(frame)
	if err != nil {
		return "", fmt.Errorf("request for permission failed, errmsg {%v}", err)
	}

	_, err = ec.DecodeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("decode response failed, errmsg {%v}", err)
	}

	switch req.COMMANDTYPE {
	case 0xE0:
		now := time.Now().UTC()
		strnow := fmt.Sprintf("%.4d%.2d%.2d%.2d%.2d%.2d%.2d", now.Year(), now.Month(), now.Day(), now.Weekday(), now.Hour()+8, now.Minute(), now.Second())
		req.INFO, _ = hex.DecodeString(strnow)
	case 0xE3:
		info := fmt.Sprintf("%.10X%.8d%.4s%.8s%.2X", req.CardID, req.UserID, req.Password, req.ExpireDate, req.Permission)
		req.INFO, _ = hex.DecodeString(info)
	case 0xE4:
		info := fmt.Sprintf("%.2X%.10X", 0, req.CardID)
		req.INFO, _ = hex.DecodeString(info)
	case 0xED:
		req.INFO = []byte{0x01}
	}

	// set
	frame, err = ec.EncodeRequest(req)
	if err != nil {
		return "", fmt.Errorf("encode request failed, errmsg {%v}", err)
	}

	fmt.Printf("req: %v\n", string(frame))
	resp, err = ec.Request(frame)
	if err != nil {
		return "", fmt.Errorf("request for data failed, errmsg {%v}", err)
	}

	_, err = ec.DecodeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("decode response failed, errmsg {%v}", err)
	}

	return "ok", nil
}
