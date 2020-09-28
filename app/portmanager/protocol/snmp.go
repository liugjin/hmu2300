/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/12/07
 * Despcription: snmp client define
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
	g "github.com/soniah/gosnmp"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolSNMP, generalSNMPDriverProtocol)
}

// Implement DriverProtocol
type snmpDriverProtocol struct {
	req *public.SNMPBindingPayload
	uri string
}

func (dp *snmpDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *snmpDriverProtocol) ClientID() string {
	return SNMPClientID + dp.req.Target
}

func (dp *snmpDriverProtocol) NewInstance() (PortClient, error) {
	return NewSNMPClient(
		dp.req.Version, dp.req.Target,
		dp.req.Port, dp.req.ReadCommunity,
		dp.req.WriteCommunity,
	)
}

func generalSNMPDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeSNMPBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &snmpDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// version
const (
	SNMPVersion1  = "v1"
	SNMPVersion2c = "v2c"
	SNMPVersion3  = "v3"
)

// SNMPClientID id
var SNMPClientID = "snmp-client-id"

// SNMPClient snmp client
type SNMPClient struct {
	ClientID string

	readConn  *g.GoSNMP
	writeConn *g.GoSNMP

	Version        string
	Target         string
	Port           int
	ReadCommunity  string
	WriteCommunity string
}

// NewSNMPClient new snmp client
func NewSNMPClient(version, target, port, readcommunity, writecommunity string) (PortClient, error) {
	var client SNMPClient
	client.ClientID = SNMPClientID
	client.Target = target
	client.ReadCommunity = readcommunity
	client.WriteCommunity = writecommunity

	if err := client.Start(version, target, port, readcommunity, writecommunity); err != nil {
		log.Printf("start snmp failed, errmsg: %v\n", err)
		go func() {
			for {
				time.Sleep(time.Second * 3)
				if err := client.Start(version, target, port, readcommunity, writecommunity); err != nil {
					log.Printf("restart snmp failed, errmsg: %v\n", err)
					continue
				}

				break
			}
		}()
	}

	return &client, nil
}

// DecodeSNMPBindingPayload decode binding
func DecodeSNMPBindingPayload(payload string) (public.SNMPBindingPayload, error) {
	var p public.SNMPBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeSNMPOperatePayload decode binding
func DecodeSNMPOperatePayload(payload string) (public.SNMPOperationPayload, error) {
	var p public.SNMPOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// ID id
func (sc *SNMPClient) ID() string {
	return sc.ClientID + sc.Target
}

// Start start
func (sc *SNMPClient) Start(version, target, port, readcommunity, writecommunity string) error {
	ver := g.Version2c
	switch version {
	case SNMPVersion1:
		ver = g.Version1
	case SNMPVersion2c:
		ver = g.Version2c
	case SNMPVersion3:
		ver = g.Version3
	}

	p, _ := strconv.Atoi(port)
	if p == 0 {
		p = 161
	}

	rc := &g.GoSNMP{
		Target:    target,
		Port:      uint16(p),
		Community: readcommunity,
		Version:   ver,
		Timeout:   time.Duration(2) * time.Second,
		// Logger:    log.New(os.Stdout, "", 0),
	}

	if err := rc.Connect(); err != nil {
		return fmt.Errorf("connect target [%v] failed: %v", sc.Target, err)
	}

	sc.readConn = rc

	wc := &g.GoSNMP{
		Target:    target,
		Port:      uint16(p),
		Community: writecommunity,
		Version:   ver,
		Timeout:   time.Duration(2) * time.Second,
		// Logger:    log.New(os.Stdout, "", 0),
	}

	if err := wc.Connect(); err != nil {
		return fmt.Errorf("connect target [%v] failed: %v", sc.Target, err)
	}

	sc.writeConn = wc

	return nil
}

// Sample sample, get values
func (sc *SNMPClient) Sample(payload string) (string, error) {
	p, err := DecodeSNMPOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	if sc.readConn == nil {
		return "", fmt.Errorf("read conn unavaliable")
	}

	r, err := sc.readConn.Get(p.OIDS)
	if err != nil {
		return "", fmt.Errorf("get info failed: %v", err)
	}

	var result string
	for _, variable := range r.Variables {
		result += variable.Name + ","

		switch variable.Type {
		case g.OctetString:
			result += string(variable.Value.([]byte))
		default:
			result += g.ToBigInt(variable.Value).String()
		}

		result += ","
	}

	// trim the last ","
	result = result[:len(result)-1]

	return result, nil
}

// Command command, set values
func (sc *SNMPClient) Command(payload string) (string, error) {
	p, err := DecodeSNMPOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	if sc.writeConn == nil {
		return "", fmt.Errorf("write conn unavaliable")
	}

	var pdus []g.SnmpPDU

	var pdu g.SnmpPDU

	pdu.Name = p.OID

	switch p.Value.(type) {
	case bool:
		pdu.Type = g.Boolean
		pdu.Value = p.Value
	case float64:
		pdu.Type = g.Integer
		pdu.Value = int(p.Value.(float64))
	case string:
		pdu.Type = g.OctetString
		pdu.Value = p.Value
	}

	pdus = append(pdus, pdu)

	_, err = sc.writeConn.Set(pdus)
	if err != nil {
		return "", fmt.Errorf("set value failed: %v", err)
	}

	return "ok", nil
}
