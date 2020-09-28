/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/25
 * Despcription: pmbus implement
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
	RegDriverProtocol(public.ProtocolPMBUS, generalPMBUSDriverProtocol)
	RegDriverProtocol(public.ProtocolYDN23, generalPMBUSDriverProtocol)
}

// Implement DriverProtocol
type pmbusDriverProtocol struct {
	req *public.PMBUSBindingPayload
	uri string
}

func (dp *pmbusDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *pmbusDriverProtocol) ClientID() string {
	return PMBUSClientID + strconv.Itoa(int(dp.req.ADR))
}

func (dp *pmbusDriverProtocol) NewInstance() (PortClient, error) {
	return NewPMBusClient(
		dp.uri,
		dp.req.BaudRate, dp.req.Timeout,
		dp.req.SOI, dp.req.VER,
		dp.req.ADR, dp.req.CID1,
		dp.req.EOI,
	)
}

func generalPMBUSDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodePMBUSBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &pmbusDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// PMBUSClientID pmbus clientid
var PMBUSClientID = "pmbusclient"

// PMBusFrame pmbus frame
type PMBusFrame struct {
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

// PMBusClient pmbus client
type PMBusClient struct {
	ClientID string

	Port    *serial.Port
	Frame   PMBusFrame
	Address byte
}

// NewPMBusClient new pmbus client
func NewPMBusClient(port string, baudrate, timeout int, soi, ver, adr, cid1, eoi byte) (*PMBusClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return &PMBusClient{}, err
	}

	frame := PMBusFrame{
		SOI:  soi,
		VER:  ver,
		ADR:  adr,
		CID1: cid1,
		EOI:  eoi,
	}

	return &PMBusClient{ClientID: PMBUSClientID, Frame: frame, Port: sp, Address: adr}, nil
}

// DecodePMBUSBindingPayload decode binding payload
func DecodePMBUSBindingPayload(payload string) (public.PMBUSBindingPayload, error) {
	var p public.PMBUSBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodePMBUSOperationPayload decode operation payload
func DecodePMBUSOperationPayload(payload string) (public.PMBUSOperationPayload, error) {
	var p public.PMBUSOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// Release release client
func (pc *PMBusClient) Release() error {
	return pc.Port.Close()
}

func computeLength(lenid uint16) ([]byte, error) {
	// transfer lenid to bytes
	byteid := make([]byte, 2)
	binary.BigEndian.PutUint16(byteid, lenid)

	// transfer to hex mode, get four ASCII code
	hexid := strings.ToUpper(hex.EncodeToString(byteid))

	// transfer the second, third, fourth ASCII code to int then plus, mod 16, not it and plus 1 for chksum
	second, _ := hex.DecodeString("0" + string(hexid[1]))
	third, _ := hex.DecodeString("0" + string(hexid[2]))
	fourth, _ := hex.DecodeString("0" + string(hexid[3]))

	sum := (second[0] + third[0] + fourth[0]) % 16
	chksum := ^sum + 1

	// transfer chksum to hex string
	bytechksum := make([]byte, 2)
	binary.BigEndian.PutUint16(bytechksum, uint16(chksum))
	hexchksum := hex.EncodeToString(bytechksum)

	// set up string length
	sumlen := len(hexchksum)
	strlength := strings.ToUpper(hexchksum[sumlen-1:] + hexid[1:])

	// decode to hex
	return hex.DecodeString(strlength)
}

func computeCheckSum(frame []byte) ([]byte, error) {
	// transfer bytes to hex string and ensure all characters capital
	hexstrframe := strings.ToUpper(hex.EncodeToString(frame))

	// sum, mod 65535, not it and plus 1 for chksum
	hexbyteframe := []byte(hexstrframe)

	sum := 0
	for _, b := range hexbyteframe {
		sum += int(b)
	}

	chksum := ^(sum % 65535) + 1

	// transfer checksum to bytes
	bytechksum := make([]byte, 2)
	binary.BigEndian.PutUint16(bytechksum, uint16(chksum))

	return bytechksum, nil
}

func computeResponseCheckSum(info []byte) ([]byte, error) {
	// sum, mod 65535, not it and plus 1 for chksum
	sum := 0
	for _, b := range info {
		sum += int(b)
	}

	chksum := ^(sum % 65535) + 1

	// transfer checksum to bytes
	bytechksum := make([]byte, 2)
	binary.BigEndian.PutUint16(bytechksum, uint16(chksum))

	return bytechksum, nil
}

// EncodeRequest encode request
func (pc *PMBusClient) EncodeRequest(req public.PMBUSOperationPayload) ([]byte, error) {
	// compute length
	length, err := computeLength(req.LENID)
	if err != nil {
		return nil, fmt.Errorf("compute length failed, errmsg {%v}", err)
	}

	// pack frame
	frame := []byte{pc.Frame.VER, req.ADR, req.CID1, req.CID2}
	frame = append(frame, length...)

	// plus command info when require
	if req.LENID > 0 {
		frame = append(frame, req.COMMANDTYPE, req.COMMANDID)
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
	frame = append([]byte{pc.Frame.SOI}, frame...)
	frame = append(frame, pc.Frame.EOI)

	return frame, nil
}

// DecodeResponse decode response
func (pc *PMBusClient) DecodeResponse(resp []byte) ([]byte, error) {
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
func (pc *PMBusClient) Request(req []byte) ([]byte, error) {
	n, err := pc.Port.Write(req)
	if err != nil {
		return nil, err
	}

	if n != len(req) {
		return nil, fmt.Errorf("write incomplete, send %v bytes", n)
	}

	chunks := []byte{}
	endflag := pc.Frame.EOI

	for {
		data := make([]byte, 256)
		n, err := pc.Port.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read response failed, errmsg {%v}", err)
		}

		chunks = append(chunks, data[:n]...)

		if data[n-1] == endflag {
			break
		}
	}

	return chunks, nil
}

// ID client's id
func (pc *PMBusClient) ID() string {
	return pc.ClientID + strconv.Itoa(int(pc.Address))
}

// Sample pmbus sample implement
func (pc *PMBusClient) Sample(payload string) (string, error) {
	req, err := DecodePMBUSOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode pmbus operation payload failed, errmsg {%v}", err)
	}

	frame, err := pc.EncodeRequest(req)
	if err != nil {
		return "", fmt.Errorf("encode request failed, errmsg {%v}", err)
	}

	fmt.Printf("req: %v\n", string(frame))
	resp, err := pc.Request(frame)
	if err != nil {
		return "", fmt.Errorf("request for data failed, errmsg {%v}", err)
	}

	data, err := pc.DecodeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("decode response failed, errmsg {%v}", err)
	}
	fmt.Printf("resp: %v\n", string(data))

	return string(data), nil
}

// Command pmbus command implement
func (pc *PMBusClient) Command(payload string) (string, error) {
	return "", nil
}
