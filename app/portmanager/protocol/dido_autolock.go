/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/03/21
 * Despcription: autoLockDIDO client define
 *
 */

package protocol

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/portlog"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolAutoLockDIDO, generalAutoLockDIDODriverProtocol)
}

// Implement DriverProtocol
type autoLockDIDODriverProtocol struct {
	uri     string
	suid    string // 采集单元的配置文件ID
	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit
}

func (dp *autoLockDIDODriverProtocol) Payload() interface{} {
	return dp.suid
}

func (dp *autoLockDIDODriverProtocol) ClientID() string {
	return dp.suid
}

func (dp *autoLockDIDODriverProtocol) NewInstance() (PortClient, error) {
	return NewAutoLockDIDOClient(dp.uri, dp.muCfg, dp.portCfg, dp.unitCfg)
}

func generalAutoLockDIDODriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	cfg := sys.GetMonitoringUnitCfg()
	portCfg := cfg.GetSamplePort(uri)
	unitCfg := portCfg.GetSampleUnit(suid)
	return &autoLockDIDODriverProtocol{
		uri:     uri,
		suid:    suid,
		muCfg:   cfg,
		portCfg: portCfg,
		unitCfg: unitCfg,
	}, nil
}

// ============= register driver end ==========================

// AutoLockDIDOClient client
type AutoLockDIDOClient struct {
	clientID string
	uri      string

	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit

	lockOpen bool
}

// NewAutoLockDIDOClient new client
func NewAutoLockDIDOClient(uri string, muCfg *sys.MonitoringUnit, portCfg *sys.SamplePort, unitCfg *sys.SampleUnit) (PortClient, error) {
	return &AutoLockDIDOClient{
		clientID: unitCfg.ID,
		uri:      uri,
		muCfg:    muCfg,
		portCfg:  portCfg,
		unitCfg:  unitCfg,
	}, nil

}

// Start start
func (vc *AutoLockDIDOClient) Start() {
	log.Info("Set DIDO AutoLock to default.")
	// 掉电时重启时默认将设备设置为锁关闭的状态
	if err := vc.lock("0"); err != nil {
		portlog.LOG.Warning(errors.As(err))
		return
	}

}

// ID id
func (dc *AutoLockDIDOClient) ID() string {
	return dc.clientID
}

func (dc *AutoLockDIDOClient) lock(c string) error {
	if err := ioutil.WriteFile(dc.unitCfg.Setting.AutoLockDoValue, []byte(c), 0644); err != nil {
		return errors.As(err)
	}
	return nil
}

func (dc *AutoLockDIDOClient) diStatus() (string, error) {
	data, err := ioutil.ReadFile(dc.unitCfg.Setting.AutoLockDiValue)
	if err != nil {
		return "", errors.As(err, dc.unitCfg.Setting.AutoLockDiValue)
	}
	val := strings.TrimSpace(string(data))
	switch val {
	case "0", "1":
		// pass
	default:
		return "", errors.New("Unknow value").As(val)
	}
	return val, nil
}

func (dc *AutoLockDIDOClient) doStatus() (string, error) {
	data, err := ioutil.ReadFile(dc.unitCfg.Setting.AutoLockDoValue)
	if err != nil {
		return "", errors.As(err, dc.unitCfg.Setting.AutoLockDoValue)
	}
	val := strings.TrimSpace(string(data))
	switch val {
	case "0", "1":
		// pass
	default:
		return "", errors.New("Unknow value").As(val)
	}
	return val, nil
}

// Sample get values
func (dc *AutoLockDIDOClient) Sample(payload string) (string, error) {
	diStatus, err := dc.diStatus()
	if err != nil {
		return "", errors.As(err, payload)
	}

	result := public.NewSamplePayload(true)
	// 如果di的手柄是开着的，优先使用手柄的状态
	if diStatus == "1" {
		if !dc.lockOpen {
			dc.lockOpen = true
			if err := dc.lock("1"); err != nil {
				return "", errors.As(err)
			}
		}
	} else {
		dc.lockOpen = false
	}
	result.PutData("val", diStatus)
	return result.Serial(), nil
}

// 网络下行的指令, 用于mqtt push
type DidoLockOperationParam struct {
	Value string `json:"value"` // 支持两种:支持0与1，0关，1开
}

// Command set values
func (dc *AutoLockDIDOClient) Command(payload string) (string, error) {
	op := &DidoLockOperationParam{}
	if err := json.Unmarshal([]byte(payload), op); err != nil {
		return "", errors.As(err, payload)
	}

	switch op.Value {
	case "1":
		if err := dc.lock(op.Value); err != nil {
			return "", errors.As(err)
		}
	case "0":
		if err := dc.lock(op.Value); err != nil {
			return "", errors.As(err)
		}
	default:
		return "", errors.New("Unknow command param:").As(op.Value)
	}
	return "ok", nil
}
