/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: sample port
 *
 */

package core

import (
	"strings"
	"time"

	"clc.hmu/app/appmanager/appnet"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log/applog"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// SamplePort sample port
type SamplePort struct {
	sys.SamplePort
}

// Start start
func (p *SamplePort) Start(muid string) ([]SampleUnitEx, error) {
	var sus []SampleUnitEx
	for _, su := range p.SampleUnits {
		if !su.Enable {
			applog.LOG.Infof("sample unit [%s] has been disable", su.ID)
			continue
		}

		switch p.Protocol {
		case public.ProtocolModbusSerial, public.ProtocolModbusTCP, public.ProtocolHSJRFID:
			if err := appnet.BindingSerialPort(p.Setting.Port, p.Protocol, p.Setting.BaudRate, su.Setting.Address, su.Timeout); err != nil && !strings.Contains(err.Error(), "client has exist") {
				applog.LOG.Warningf("binding serial port failed, port[%s], slaveid[%d], errmsg[%v]", p.Setting.Port, su.Setting.Address, err)
				continue
			}
		case public.ProtocolSensorflow:
			if err := appnet.BindingSensorflowPort(p.Setting.Port, p.Protocol, p.Setting.BaudRate, su.Timeout, su.Setting.Address, p.Setting.KeyNumber, muid, su.ID, p.Setting.WANInterface, p.Setting.WifiInterface); err != nil && !strings.Contains(err.Error(), "client has exist") {
				applog.LOG.Warningf("binding sensorflow port failed, port[%s], slaveid[%d], keynumber[%d], errmsg[%v]", p.Setting.Port, su.Setting.Address, p.Setting.KeyNumber, err)
				continue
			}
		case public.ProtocolHYIOTMU:
			if err := appnet.BindingSystemPort(p.Setting.Port, p.Protocol, su.Setting.Host, su.Setting.Port, su.Setting.Model); err != nil {
				applog.LOG.Warningf("binding system port failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolPMBUS, public.ProtocolYDN23:
			soi := byte(0x7E)
			ver := byte(0x21)
			adr := byte(su.Setting.Address)
			eoi := byte(0x0D)
			if err := appnet.BindingPMBusPort(p.Setting.Port, p.Protocol, int(p.Setting.BaudRate), int(su.Timeout), soi, ver, adr, eoi); err != nil {
				applog.LOG.Warningf("binding pmbus port failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolOilMachine:
			soi := byte(0x7E)
			eoi := byte(0x0D)
			if err := appnet.BindingOilMachine(p.Setting.Port, p.Protocol, int(p.Setting.BaudRate), int(su.Timeout), soi, eoi); err != nil {
				applog.LOG.Warningf("binding oil machine failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolLampWith:
			if err := appnet.CommonSerialBinding(p.Setting.Port, p.Protocol, int(p.Setting.BaudRate), int(su.Timeout)); err != nil {
				applog.LOG.Warningf("binding lamp with failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolCreditCard:
			if err := appnet.BindingCreditCard(p.Setting.Port, p.Protocol, su.Setting.Username, su.Setting.Password, su.Setting.SerialNumber); err != nil {
				applog.LOG.Warningf("binding credit card failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolLuMiGateway:
			if err := appnet.BindingLuMiGateway(p.Setting.Port, p.Protocol, p.Setting.SID, p.Setting.Password, p.Setting.NetInterface); err != nil {
				applog.LOG.Warningf("binding lumi gateway failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolFaceIPC:
			if err := appnet.BindingFaceIPC(p.Setting.Port, p.Protocol, p.Setting.Host, p.Setting.Port, muid, su.ID, p.Setting.UploadServer, p.Setting.Author, p.Setting.Project, p.Setting.Token, p.Setting.User); err != nil {
				applog.LOG.Warningf("binding face ipc failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolSNMP:
			if err := appnet.BindingSNMP(p.Setting.Port, p.Protocol, su.Setting.Version, su.Setting.Target, su.Setting.Port, su.Setting.ReadCommunity, su.Setting.WriteCommunity); err != nil {
				applog.LOG.Warningf("binding snmp failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolWeiGengEntry:
			if err := appnet.BindingWeiGengEntry(p.Setting.Port, p.Protocol, p.Setting.LocalhostAddress, p.Setting.LocalhostPort, su.Setting.DoorAddress, su.Setting.DoorPort, su.Setting.SerialNo); err != nil {
				applog.LOG.Warningf("binding weigeng entry failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolDIDO:
			if err := appnet.BindingDIDO(p.Setting.Port, p.Protocol); err != nil {
				applog.LOG.Warningf("binding dido failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolES5200:
			if err := appnet.CommonSerialBinding(p.Setting.Port, p.Protocol, int(p.Setting.BaudRate), int(su.Timeout)); err != nil && !strings.Contains(err.Error(), "client has exist") {
				applog.LOG.Warningf("binding es5200 failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolCamera:
			if err := appnet.BindingCamera(p.Setting.Port, p.Protocol, su.Setting.Host); err != nil {
				applog.LOG.Warningf("binding camera failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolElecFire:
			if err := appnet.BindingElecFire(p.Setting.Port, p.Protocol, p.Setting.Host, p.Setting.Port); err != nil {
				applog.LOG.Warningf("binding camera failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		case public.ProtocolVirtualAntenna:
			if err := appnet.BindingVirtualAntenna(p.Setting.Port, p.Protocol, p.Setting.Host, p.Setting.Port); err != nil {
				applog.LOG.Warningf("binding camera failed, port[%s], errmsg[%v]", p.Setting.Port, err)
				continue
			}
		default:
			if err := binding(p.Setting.Port, p.Protocol, su.ID); err != nil {
				applog.LOG.Warningf("binding [%s] failed, port[%s], errmsg[%v]", p.Protocol, p.Setting.Port, err)
				continue
			}
		}

		var suex SampleUnitEx
		suex.SU = &SampleUnit{su}

		if err := suex.Start(); err != nil {
			applog.LOG.Warningf("start su failed, errmsg {%v}", err)
			continue
		}

		sus = append(sus, suex)
	}

	if len(sus) == 0 {
		applog.LOG.Warningf("no sample unit works!!!")
		return nil, nil
	}

	go func(port, protocol string, baudrate int32) {
		applog.LOG.Infof("sample go coroutine, port{%v}, protocol{%v}, baudrate{%v}", port, protocol, baudrate)
		for {
			// TODO:考虑出现panic时线程挂掉的情况(by shu).
			sulen := len(sus)
			for i := 0; i < sulen; i++ {
				if err := sus[i].Sample(muid, port, protocol, baudrate); err != nil {
					applog.LOG.Warning(errors.As(err))
				}
			}

			// use first one's period
			period := time.Duration(sus[0].SU.Period)
			if period <= 0 {
				// 需要注意，此值为0时，只会调用一遍，不会启动循环采集功能。
				return
			}
			time.Sleep(time.Millisecond * period)
		}
	}(p.Setting.Port, p.Protocol, p.Setting.BaudRate)

	return sus, nil
}
