/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/03/21
 * Despcription: dido client define
 *
 */

package protocol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"clc.hmu/app/public"
	"clc.hmu/app/public/sys"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolDIDO, generalDIDODriverProtocol)
}

// Implement DriverProtocol
type didoDriverProtocol struct {
	req string
	uri string
}

func (dp *didoDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *didoDriverProtocol) ClientID() string {
	return DIDOClientID + dp.uri
}

func (dp *didoDriverProtocol) NewInstance() (PortClient, error) {
	return NewDIDOClient(dp.uri)
}

func generalDIDODriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	return &didoDriverProtocol{
		req: payload,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// DIDOClientID id
var DIDOClientID = "dido-client-id"

// DIDOClient client
type DIDOClient struct {
	ClientID string
	Model    string

	filename string

	sysclient sys.SystemClient
}

// NewDIDOClient new client
func NewDIDOClient(port string) (PortClient, error) {
	var client DIDOClient

	cfg := sys.GetBusManagerCfg()

	client.Model = cfg.Model
	client.ClientID = DIDOClientID + port

	switch cfg.Model {
	case sys.MODEL_HMU2500, sys.MODEL_HMU2400:
		client.filename = port + "/value"
	case sys.MODEL_HMU2300:
		client.filename = filepath.Base(port)
		client.sysclient = sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	}

	return &client, nil
}

// DecodeDIDOOperationPayload decode binding payload
func DecodeDIDOOperationPayload(payload string) (public.DIDOOperationPayload, error) {
	var p public.DIDOOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// ID id
func (dc *DIDOClient) ID() string {
	return dc.ClientID
}

// Sample get values
func (dc *DIDOClient) Sample(payload string) (string, error) {
	var val string

	switch dc.Model {
	case sys.ModelHMU2500, sys.MODEL_HMU2400:
		d, err := ioutil.ReadFile(dc.filename)
		if err != nil {
			return "", fmt.Errorf("read data failed, %v", err)
		}

		v := string(d)
		val = strings.Replace(v, "\n", "", -1)
	case sys.MODEL_HMU2300:
		resp, err := dc.sysclient.DIDO(dc.filename)
		if err != nil {
			return "", fmt.Errorf("get value failed, %v", err)
		}

		val = strconv.Itoa(resp.Value)
	}

	return val, nil
}

// Command set values
func (dc *DIDOClient) Command(payload string) (string, error) {
	p, err := DecodeDIDOOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode lamp with operation payload failed, errmsg {%v}", err)
	}

	val := ""

	switch v := p.Value.(type) {
	case int:
		val = strconv.Itoa(v)
	case float64:
		val = strconv.Itoa(int(v))
	case string:
		val = v
	}

	switch dc.Model {
	case sys.ModelHMU2500:
		if err := ioutil.WriteFile(dc.filename, []byte(val), 0644); err != nil {
			return "", fmt.Errorf("write data failed, %v", err)
		}
	case sys.MODEL_HMU2300:
		v, _ := strconv.Atoi(val)
		_, err := dc.sysclient.SetDO(dc.filename, v)
		if err != nil {
			return "", fmt.Errorf("set do failed, %v", err)
		}
	case sys.MODEL_HMU2400:
		enfilename := "/ual/do_en/value"

		// first set enable to false
		if err := ioutil.WriteFile(enfilename, []byte("0"), 0644); err != nil {
			return "", fmt.Errorf("write data to [%s] failed, %v", enfilename, err)
		}

		if err := ioutil.WriteFile(dc.filename, []byte(val), 0644); err != nil {
			return "", fmt.Errorf("write data to [%s] failed, %v", enfilename, err)
		}

		// then set enable
		if err := ioutil.WriteFile(enfilename, []byte("1"), 0644); err != nil {
			return "", fmt.Errorf("write data to [%s] failed, %v", enfilename, err)
		}
	}

	return "ok", nil
}
