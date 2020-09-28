/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/03/21
 * Despcription: autoLockIOBlock client define
 *
 */

package protocol

import (
	"fmt"
	"sync"

	"clc.hmu/app/public"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
// 使用的是C2000-A1-PDD4040-BB1的设备
func init() {
	RegDriverProtocol(public.ProtocolAutoLockIOBlock, generalAutoLockIOBlockDriverProtocol)
}

// Implement DriverProtocol
type autoLockIOBlockDriverProtocol struct {
	uri     string
	suid    string // 采集单元的配置文件ID
	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit
}

func (dp *autoLockIOBlockDriverProtocol) Payload() interface{} {
	return dp.suid
}

func (dp *autoLockIOBlockDriverProtocol) ClientID() string {
	return dp.suid
}

func (dp *autoLockIOBlockDriverProtocol) NewInstance() (PortClient, error) {
	return NewAutoLockIOBlockClient(dp.uri, dp.muCfg, dp.portCfg, dp.unitCfg)
}

func generalAutoLockIOBlockDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	cfg := sys.GetMonitoringUnitCfg()
	portCfg := cfg.GetSamplePort(uri)
	unitCfg := portCfg.GetSampleUnit(suid)
	return &autoLockIOBlockDriverProtocol{
		uri:     uri,
		suid:    suid,
		muCfg:   cfg,
		portCfg: portCfg,
		unitCfg: unitCfg,
	}, nil
}

// ============= register driver end ==========================

// AutoLockIOBlockClient client
type AutoLockIOBlockClient struct {
	clientID string
	uri      string

	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit

	lockOpen []bool

	ioSync sync.Mutex

	client *ModbusClient
}

// NewAutoLockIOBlockClient new client
func NewAutoLockIOBlockClient(uri string, muCfg *sys.MonitoringUnit, portCfg *sys.SamplePort, unitCfg *sys.SampleUnit) (PortClient, error) {
	client, err := NewModbusSerialClient(
		portCfg.Setting.Port,
		portCfg.Setting.BaudRate,
		unitCfg.Timeout,
		byte(unitCfg.Setting.Address),
	)
	if err != nil {
		return nil, errors.As(err)
	}
	return &AutoLockIOBlockClient{
		clientID: unitCfg.ID,
		uri:      uri,
		muCfg:    muCfg,
		portCfg:  portCfg,
		unitCfg:  unitCfg,
		client:   client,
		lockOpen: make([]bool, unitCfg.Setting.AutoLockNum),
	}, nil

}

// Start start
func (vc *AutoLockIOBlockClient) Start() {
}

// ID id
func (dc *AutoLockIOBlockClient) ID() string {
	return dc.clientID
}

// pos从0开始，依次为0,1,2,3
func (dc *AutoLockIOBlockClient) lock(pos, c int) ([]byte, error) {
	result, err := dc.client.command(&public.ModbusPayload{
		Code:     5,
		Slaveid:  dc.unitCfg.Setting.Address,
		Address:  100 + int32(pos),
		Quantity: 1,
		Value:    c,
	})
	if err != nil {
		return nil, errors.As(err)
	}
	return result, nil
}

// Sample get values
func (dc *AutoLockIOBlockClient) Sample(payload string) (string, error) {
	dc.ioSync.Lock()
	defer dc.ioSync.Unlock()

	doRead, err := dc.client.sample(&public.ModbusPayload{
		Code:     1,
		Slaveid:  dc.unitCfg.Setting.Address,
		Address:  100,
		Quantity: dc.unitCfg.Setting.AutoLockNum,
	})
	if err != nil {
		return "", errors.As(err)
	}
	doStatus := []byte{}
	for i := int32(0); i < dc.unitCfg.Setting.AutoLockNum; i++ {
		doStatus = append(doStatus, (doRead[0]>>uint32(i))&0x1)
	}

	diRead, err := dc.client.sample(&public.ModbusPayload{
		Code:     2,
		Slaveid:  dc.unitCfg.Setting.Address,
		Address:  200,
		Quantity: dc.unitCfg.Setting.AutoLockNum,
	})
	if err != nil {
		return "", errors.As(err)
	}
	diStatus := []byte{}
	for i := int32(0); i < dc.unitCfg.Setting.AutoLockNum; i++ {
		diStatus = append(diStatus, (diRead[0]>>uint32(i))&0x1)
	}

	result := public.NewSamplePayload(true)
	for i, di := range diStatus {
		result.PutData(fmt.Sprintf("di", i), di)
		result.PutData(fmt.Sprintf("do", i), doStatus[i])
		if di == 1 {
			if !dc.lockOpen[i] {
				// 若非此值0xff00，需要调线
				if _, err := dc.lock(int(i), 65280); err != nil {
					return "", errors.As(err)
				}
				dc.lockOpen[i] = true
			}
		} else {
			dc.lockOpen[i] = false
		}
	}

	return result.Serial(), nil
}

// Command set values
func (dc *AutoLockIOBlockClient) Command(payload string) (string, error) {
	dc.ioSync.Lock()
	defer dc.ioSync.Unlock()
	// TODO: 透传channel值过来

	result, err := dc.lock(0, 0)
	if err != nil {
		return "", errors.As(err)
	}
	return fmt.Sprintf("%x", result), nil
}
