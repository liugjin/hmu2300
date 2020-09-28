/*
 *
 * Copyright 2019 huayuan-iot
 *
 * Author: saul
 * Date: 2019/07/01
 * Despcription: lte client define
 *
 */

package protocol

import (
	"clc.hmu/app/public"
	"clc.hmu/app/public/at"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolLTE, generalLTEDriverProtocol)
}

// Implement DriverProtocol
type lteDriverProtocol struct {
	uri     string
	suid    string // 采集单元的配置文件ID
	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit
}

func (dp *lteDriverProtocol) Payload() interface{} {
	return dp.suid
}

func (dp *lteDriverProtocol) ClientID() string {
	return dp.suid
}

func (dp *lteDriverProtocol) NewInstance() (PortClient, error) {
	return NewLTEClient(dp.uri, dp.muCfg, dp.portCfg, dp.unitCfg)
}

func generalLTEDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	cfg := sys.GetMonitoringUnitCfg()
	portCfg := cfg.GetSamplePort(uri)
	unitCfg := portCfg.GetSampleUnit(suid)
	return &lteDriverProtocol{
		uri:     uri,
		suid:    suid,
		muCfg:   cfg,
		portCfg: portCfg,
		unitCfg: unitCfg,
	}, nil
}

// ============= register driver end ==========================

// LTEClient client
type LTEClient struct {
	clientID string
	uri      string

	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit

	atCmd at.AT
}

// NewLTEClient new client
func NewLTEClient(uri string, muCfg *sys.MonitoringUnit, portCfg *sys.SamplePort, unitCfg *sys.SampleUnit) (PortClient, error) {
	switch portCfg.Protocol {
	case public.ProtocolLTE:
		return &LTEClient{
			clientID: unitCfg.ID,
			uri:      uri,
			muCfg:    muCfg,
			portCfg:  portCfg,
			unitCfg:  unitCfg,

			atCmd: at.NewFileATCmd(uri, 60*1e9),
		}, nil
	}

	return nil, errors.New("Not implement").As(portCfg.Protocol)
}

// Start start
func (vc *LTEClient) Start() {
	log.Info("Ignore LTEClient start")
}

// ID id
func (dc *LTEClient) ID() string {
	return dc.clientID
}

// Sample get values
func (dc *LTEClient) Sample(payload string) (string, error) {
	// payload参数已无意义
	csq, err := dc.atCmd.CSQ()
	if err != nil {
		return "", errors.As(err, payload)
	}
	result := public.NewSamplePayload(true)
	result.PutData("val", csq)
	return result.Serial(), nil
}

// Command set values
func (dc *LTEClient) Command(payload string) (string, error) {
	return "", errors.New("Not implement").As(payload)
}
