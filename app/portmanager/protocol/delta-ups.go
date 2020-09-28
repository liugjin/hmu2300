/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/30
 * Despcription: delta ups implement
 *
 */

package protocol

import (
	"encoding/binary"
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
	RegDriverProtocol(public.ProtocolDeltaUPS, generalDeltaUPSDriverProtocol)
}

// Implement DriverProtocol
type deltaUPSDriverProtocol struct {
	req *public.DeltaUPSBindingPayload
	uri string
}

func (dp *deltaUPSDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *deltaUPSDriverProtocol) ClientID() string {
	return DeltaUPSClientID
}

func (dp *deltaUPSDriverProtocol) NewInstance() (PortClient, error) {
	return NewDeltaUPSClient(
		dp.uri,
		dp.req.BaudRate, dp.req.Timeout,
		dp.req.Header, uint16(dp.req.ID),
	)
}

func generalDeltaUPSDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeDeltaUPSBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &deltaUPSDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// transfer type
const (
	TransferTypeReject   = 'R' // UPS -> Computer, command rejected due to not support
	TransferTypeAccepted = 'A' // UPS -> Computer, command accepted
	TransferTypePolling  = 'P' // Computer -> UPS, polling command
	TransferTypeSet      = 'S' // Computer -> UPS, set command
	TransferTypeData     = 'D' // UPS -> Computer, data returned
)

// DeltaUPSClientID delta ups clientid
var DeltaUPSClientID = "delataupsclient"

// DeltaUPSFrame frame
type DeltaUPSFrame struct {
	Header   byte
	ID       [2]byte
	Type     byte
	Length   [3]byte
	Data     []byte // 128 bytes max
	Checksum [2]byte
}

// DeltaUPSClient delta ups client
type DeltaUPSClient struct {
	ClientID string

	Port  *serial.Port
	Frame DeltaUPSFrame
}

// NewDeltaUPSClient new client
func NewDeltaUPSClient(port string, baudrate, timeout int, header byte, id uint16) (*DeltaUPSClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}

	// transfer id to bytes
	bid := make([]byte, 2)
	binary.BigEndian.PutUint16(bid, id)

	var client DeltaUPSClient
	client.Frame.Header = header
	client.Frame.ID[0] = bid[0]
	client.Frame.ID[1] = bid[1]

	client.Port = sp

	return &client, nil
}

// DecodeDeltaUPSBindingPayload decode binding payload
func DecodeDeltaUPSBindingPayload(payload string) (public.DeltaUPSBindingPayload, error) {
	var p public.DeltaUPSBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeDeltaUPSOperationPayload decode operation payload
func DecodeDeltaUPSOperationPayload(payload string) (public.DeltaUPSOperationPayload, error) {
	var p public.DeltaUPSOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// EncodeRequest encode request
func (dc *DeltaUPSClient) EncodeRequest(req public.DeltaUPSOperationPayload) ([]byte, error) {
	return nil, nil
}

// Request send request
func (dc *DeltaUPSClient) Request(req []byte) ([]byte, error) {
	n, err := dc.Port.Write(req)
	if err != nil {
		return nil, err
	}

	if n != len(req) {
		return nil, fmt.Errorf("write incomplete, send %v bytes", n)
	}

	chunks := []byte{}

	data := make([]byte, 256)
	n, err = dc.Port.Read(data)
	if err != nil {
		return nil, fmt.Errorf("read response failed, errmsg {%v}", err)
	}

	chunks = append(chunks, data[:n]...)

	return chunks, nil
}

// ID client's id
func (dc *DeltaUPSClient) ID() string {
	return dc.ClientID
}

// Sample delta ups sample implement
func (dc *DeltaUPSClient) Sample(payload string) (string, error) {
	return "", nil
}

// Command delta ups command implement
func (dc *DeltaUPSClient) Command(payload string) (string, error) {
	return "", nil
}
