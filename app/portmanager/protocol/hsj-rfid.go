/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/02/28
 * Despcription: hongshunjie rfid implement
 *
 */

package protocol

import (
	"fmt"
	"strings"
	"time"

	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
	"github.com/tarm/serial"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolHSJRFID, generalHSJRFIDDriverProtocol)
}

// Implement DriverProtocol
type hsjRFIDDriverProtocol struct {
	req *public.ModbusPayload
	uri string
}

func (dp *hsjRFIDDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *hsjRFIDDriverProtocol) ClientID() string {
	return HSJRFIDClientID
}

func (dp *hsjRFIDDriverProtocol) NewInstance() (PortClient, error) {
	return NewHSJRFIDClient(
		dp.uri,
		int(dp.req.BaudRate), int(dp.req.Timeout),
	)
}

func generalHSJRFIDDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &hsjRFIDDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// HSJRFIDClientID id
var HSJRFIDClientID = "hsj-rfid-client-id"

// HSJRFIDClient client
type HSJRFIDClient struct {
	ClientID string
	Port     *serial.Port
}

// NewHSJRFIDClient new client
func NewHSJRFIDClient(port string, baudrate, timeout int) (PortClient, error) {
	cfg := &serial.Config{Name: port, Baud: baudrate, ReadTimeout: time.Millisecond * time.Duration(timeout)}
	sp, err := serial.OpenPort(cfg)
	if err != nil {
		return &HSJRFIDClient{}, err
	}

	var client HSJRFIDClient
	client.ClientID = HSJRFIDClientID
	client.Port = sp

	return &client, nil
}

func (hc *HSJRFIDClient) collect() (string, error) {
	chunks := []byte{}
	endflag := byte('\n')

	for {
		data := make([]byte, 256)
		n, err := hc.Port.Read(data)
		if err != nil {
			return "", fmt.Errorf("read failed, errmsg {%v}", err)
		}

		chunks = append(chunks, data[:n]...)

		if (n == 0) || (data[n-1] == endflag) {
			break
		}
	}

	return string(chunks), nil
}

// ID id
func (hc *HSJRFIDClient) ID() string {
	return hc.ClientID
}

// Sample sample, get values
func (hc *HSJRFIDClient) Sample(payload string) (string, error) {
	seqs, err := hc.collect()
	if err != nil {
		return "", err
	}

	rs := []string{}
	list := make(map[string]bool)

	ss := strings.Split(seqs, "\n")
	for _, s := range ss {
		s = strings.Replace(s, "\r", "", -1)
		if s == "" {
			continue
		}

		list[s] = true
	}

	for k := range list {
		rs = append(rs, k)
	}

	result := ""
	if len(rs) != 0 {
		result = strings.Join(rs, ",")
	}

	return fmt.Sprintf("%x", result), nil
}

// Command command, set values
func (hc *HSJRFIDClient) Command(payload string) (string, error) {
	return "", nil
}
