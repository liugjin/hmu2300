/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/10/17
 * Despcription: port client define
 *
 */

package protocol

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"github.com/gwaylib/errors"
	"github.com/tarm/serial"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolLampWith, generalLampWithDriverProtocol)
}

// Implement DriverProtocol
type lampWithDriverProtocol struct {
	req *public.CommonSerialBindingPayload
	uri string
}

func (dp *lampWithDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *lampWithDriverProtocol) ClientID() string {
	return LampWithClientID
}

func (dp *lampWithDriverProtocol) NewInstance() (PortClient, error) {
	return NewLampWithClient(
		dp.uri,
		int(dp.req.BaudRate), dp.req.Timeout,
	)
}

func generalLampWithDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeLampWithBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &lampWithDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// LampWithClientID id
var LampWithClientID = "lamp-with-client-id"

// const
const (
	LampWithLightOff = iota
	LampWithLightBlinkRed
	LampWithLightBlinkGreen
	LampWithLightBlinkBlue
	LampWithLightBreatheRed
	LampWithLightBreatheGreen
	LampWithLightBreatheBlue
	LampWithLightAlwaysRed
	LampWithLightAlwaysGreen
	LampWithLightAlwaysBlue
)

// LampWithClient client
type LampWithClient struct {
	ClientID string

	Port *serial.Port
}

// NewLampWithClient new lamp with client
func NewLampWithClient(port string, baudrate, timeout int) (*LampWithClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return &LampWithClient{}, err
	}

	return &LampWithClient{ClientID: LampWithClientID, Port: sp}, nil
}

// DecodeLampWithBindingPayload decode binding payload
func DecodeLampWithBindingPayload(payload string) (public.CommonSerialBindingPayload, error) {
	var p public.CommonSerialBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeLampWithOperationPayload decode binding payload
func DecodeLampWithOperationPayload(payload string) (public.LampWithOperationPayload, error) {
	var p public.LampWithOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

func computeLWCheckSum(frame []byte) byte {
	sum := uint16(0)
	for _, v := range frame {
		sum += uint16(v)
	}

	sum = sum % 256

	return byte(sum)
}

func packLampWithFrame(mode int) []byte {
	switch mode {
	case LampWithLightOff:
		{
			frame := []byte{0x06, 0xff, 0xa0, 0x00, 0x00, 0x00}
			checksum := computeLWCheckSum(frame)

			frame = append([]byte{0xa5}, frame...)
			return append(frame, checksum, 0x5a)
		}
	case LampWithLightBlinkRed:
		return []byte{0xa5, 0x0a, 0xff, 0xa1, 0x02, 0x02, 0x10, 0x10, 0xff, 0x00, 0x00, 0xee, 0x5a}
	case LampWithLightBlinkGreen:
		return []byte{0xa5, 0x0a, 0xff, 0xa1, 0x02, 0x03, 0x10, 0x10, 0x00, 0xff, 0x00, 0xee, 0x5a}
	case LampWithLightBlinkBlue:
		return []byte{0xa5, 0x0a, 0xff, 0xa1, 0x02, 0x04, 0x10, 0x10, 0x00, 0x00, 0xff, 0xee, 0x5a}
	case LampWithLightBreatheRed:
		return []byte{0xa5, 0x06, 0xff, 0xa1, 0x01, 0x01, 0x10, 0xee, 0x5a}
	case LampWithLightBreatheGreen:
		return []byte{0xa5, 0x06, 0xff, 0xa1, 0x01, 0x02, 0x10, 0xee, 0x5a}
	case LampWithLightBreatheBlue:
		return []byte{0xa5, 0x06, 0xff, 0xa1, 0x01, 0x03, 0x10, 0xee, 0x5a}
	case LampWithLightAlwaysRed:
		{
			frame := []byte{0x06, 0xff, 0xa0, 0xff, 0x00, 0x00}
			checksum := computeLWCheckSum(frame)

			frame = append([]byte{0xa5}, frame...)
			return append(frame, checksum, 0x5a)
		}
	case LampWithLightAlwaysGreen:
		{
			frame := []byte{0x06, 0xff, 0xa0, 0x00, 0xff, 0x00}
			checksum := computeLWCheckSum(frame)

			frame = append([]byte{0xa5}, frame...)
			return append(frame, checksum, 0x5a)
		}
	case LampWithLightAlwaysBlue:
		{
			frame := []byte{0x06, 0xff, 0xa0, 0x00, 0x00, 0xff}
			checksum := computeLWCheckSum(frame)

			frame = append([]byte{0xa5}, frame...)
			return append(frame, checksum, 0x5a)
		}
	}

	return nil
}

// Request request
func (lc *LampWithClient) Request(req []byte) ([]byte, error) {
	_, err := lc.Port.Write(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ID client's id
func (lc *LampWithClient) ID() string {
	return lc.ClientID
}

// Sample lamp with sample implement
func (lc *LampWithClient) Sample(payload string) (string, error) {
	return "", nil
}

// Command lamp with command implement
func (lc *LampWithClient) Command(payload string) (string, error) {
	p, err := DecodeLampWithOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode lamp with operation payload failed, errmsg {%v}", err)
	}

	val := 0

	switch v := p.Value.(type) {
	case int:
		val = v
	case float64:
		val = int(v)
	case string:
		tmp, err := strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("unknown value {%v}", v)
		}

		val = tmp
	}

	frame := packLampWithFrame(val)
	log.Printf("%x", frame)

	if frame == nil {
		return "", fmt.Errorf("not support value {%v}", val)
	}

	_, err = lc.Request(frame)
	if err != nil {
		return "", fmt.Errorf("contorl lampwith failed, errmsg {%v}", err)
	}

	return "ok", nil
}
