/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/08/10
 * Despcription: dc oil machine implement
 *
 */

package protocol

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"clc.hmu/app/public"

	"github.com/gwaylib/errors"
	"github.com/tarm/serial"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolOilMachine, generalOilMachineDriverProtocol)
}

// Implement DriverProtocol
type oilMachineDriverProtocol struct {
	req *public.OilMachineBindingPayload
	uri string
}

func (dp *oilMachineDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *oilMachineDriverProtocol) ClientID() string {
	return OilMachineClientID
}

func (dp *oilMachineDriverProtocol) NewInstance() (PortClient, error) {
	return NewOilMachineClient(
		dp.uri,
		dp.req.BaudRate, dp.req.Timeout,
		dp.req.SOI, dp.req.EOI,
	)
}

func generalOilMachineDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeOilMachineBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &oilMachineDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// OilMachineClientID client id
var OilMachineClientID = "oilmachineclient"

// OilMachineFrame oil machine client
type OilMachineFrame struct {
	SOI    byte
	ADR    byte
	CID1   byte
	CID2   byte
	LENGTH byte
	INFO   []byte
	CHKSUM byte
	EOI    byte
}

// OilMachineClient client
type OilMachineClient struct {
	ClientID string

	Port  *serial.Port
	Frame OilMachineFrame
}

// NewOilMachineClient new client
func NewOilMachineClient(port string, baudrate, timeout int, soi, eoi byte) (*OilMachineClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return &OilMachineClient{}, err
	}

	frame := OilMachineFrame{
		SOI: soi,
		EOI: eoi,
	}

	return &OilMachineClient{ClientID: OilMachineClientID, Frame: frame, Port: sp}, nil
}

// DecodeOilMachineBindingPayload decode binding payload
func DecodeOilMachineBindingPayload(payload string) (public.OilMachineBindingPayload, error) {
	var p public.OilMachineBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeOilMachineOperationPayload decode operation payload
func DecodeOilMachineOperationPayload(payload string) (public.OilMachineOperationPayload, error) {
	var p public.OilMachineOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// Release release client
func (oc *OilMachineClient) Release() error {
	return oc.Port.Close()
}

func normalHexToTransferHex(b int) []byte {
	return append([]byte{}, byte((((b&0xf0)>>4)&0x0f)+0x30), byte((b&0x0f)+0x30))
}

func transferHexToNormalHex(data []byte) byte {
	if len(data) != 2 {
		return 0
	}

	tmp := []byte{data[0] - 0x30, data[1] - 0x30}
	return (tmp[0]<<4)&0xf0 + tmp[1]&0x0f
}

// EncodeRequest encode request
func (oc *OilMachineClient) EncodeRequest(req public.OilMachineOperationPayload) ([]byte, error) {
	var info = []byte{req.ADR, req.CID1, req.CID2, req.LENGTH}
	info = append(info, req.COMMANDINFO...)

	// normal hex data to transfer data format
	tran := []byte{}
	for _, b := range info {
		tran = append(tran, normalHexToTransferHex(int(b))...)
	}

	sum := 0
	for _, b := range tran {
		sum += int(b)
	}

	// check sum
	sum = sum % 256

	tran = append(tran, normalHexToTransferHex(sum)...)
	tran = append([]byte{oc.Frame.SOI}, tran...)
	tran = append(tran, oc.Frame.EOI)

	return tran, nil
}

// Request request
func (oc *OilMachineClient) Request(req []byte) ([]byte, error) {
	n, err := oc.Port.Write(req)
	if err != nil {
		return nil, err
	}

	if n != len(req) {
		return nil, fmt.Errorf("write incomplete, send %v bytes", n)
	}

	chunks := []byte{}
	endflag := oc.Frame.EOI

	for {
		data := make([]byte, 256)
		n, err := oc.Port.Read(data)
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

// DecodeResponse decode response
func (oc *OilMachineClient) DecodeResponse(resp []byte) ([]byte, error) {
	rlen := len(resp)
	if rlen < 12 {
		return nil, fmt.Errorf("response length no enough")
	}

	// verify length
	dlen := resp[7:9]
	ilen := int(transferHexToNormalHex(dlen))

	data := resp[9 : rlen-3]
	if ilen != len(data) {
		return nil, fmt.Errorf("length no match, rlen(%v), ilen(%v", rlen, ilen)
	}

	// transfer to normal hex
	info := []byte{}
	for i := 0; i < ilen; i = i + 2 {
		d := data[i : i+2]
		b := transferHexToNormalHex(d)
		info = append(info, b)
	}

	return info, nil
}

// ID client's id
func (oc *OilMachineClient) ID() string {
	return oc.ClientID
}

// Sample oil machine sample implement
func (oc *OilMachineClient) Sample(payload string) (string, error) {
	req, err := DecodeOilMachineOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode oil machine operation payload failed, errmsg {%v}", err)
	}

	frame, err := oc.EncodeRequest(req)
	if err != nil {
		return "", fmt.Errorf("encode request failed, errmsg {%v}", err)
	}

	resp, err := oc.Request(frame)
	if err != nil {
		return "", fmt.Errorf("request for data failed, errmsg {%v}", err)
	}

	data, err := oc.DecodeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("decode response failed, errmsg {%v}", err)
	}

	return hex.EncodeToString(data), nil
}

// Command oil machine command implement
func (oc *OilMachineClient) Command(payload string) (string, error) {
	return "", nil
}
