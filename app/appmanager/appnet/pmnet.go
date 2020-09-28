/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: port server net
 *
 */

package appnet

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"clc.hmu/app/appmanager/appnet/pmnet"
	"clc.hmu/app/public"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"

	pb "clc.hmu/app/portmanager/portpb"
)

//
// 调用portmanager的Operate操作
//
// 参数
// timeout -- 请求的超时间，若填写0，使用默认值
// kind -- 操作的类别
// port -- 操作的端口，比如：/dev/com1, http://...等等
// value -- 操作的数据, 自具体协议实现
//
// 返回
// string -- 操作的结果值
// error -- 是否操作成功
func Operate(timeout time.Duration, kind, port, suid, value string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := timeout
	if pmtimeout == 0 {
		pmtimeout = pmnet.GetTimeout() * time.Second
	}

	if pmcli == nil {
		return "", errors.New("port client unavailable").As(pmtimeout, kind, port, suid, value)
	}

	ctx, cancel := context.WithTimeout(pmparentctx, pmtimeout)
	defer cancel()

	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    kind,
		Suid:    suid,
		Payload: value,
	})
	if err != nil {
		return "", errors.As(err, pmtimeout, kind, port, suid, value)
	}

	return r.Data, nil
}

// 默认的采集透传
func CommonSample(timeout time.Duration, port, suid string) (string, error) {
	return Operate(timeout, public.OperateSample, port, suid, "")
}

// 默认服务器下行指令透传
func CommonCommand(timeout time.Duration, port, suid string, value string) (string, error) {
	return Operate(timeout, public.OperateCommand, port, suid, value)
}

// BindingSerialPort binding serial port
func BindingSerialPort(port, protocol string, baudrate, slaveid, timeout int32) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ModbusPayload{
		BaudRate: baudrate,
		Slaveid:  slaveid,
		Timeout:  timeout,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// BindingSensorflowPort binding sensorflow port
func BindingSensorflowPort(port, protocol string, baudrate, timeout, slaveid, keynumber int32, muid, suid, wanift, wifiift string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ModbusPayload{
		BaudRate:      baudrate,
		Timeout:       timeout,
		Slaveid:       slaveid,
		KeyNumber:     keynumber,
		MUID:          muid,
		SUID:          suid,
		WANInterface:  wanift,
		WifiInterface: wifiift,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// Sample sample
func Sample(port string, baudrate, code, slaveid, address, quantity, timeout int32, muid, suid string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ModbusPayload{
		BaudRate: baudrate,
		Code:     code,
		Slaveid:  slaveid,
		Address:  address,
		Quantity: quantity,
		Timeout:  timeout,
		MUID:     muid,
		SUID:     suid,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// Command command
func Command(port string, baudRate, code, slaveid, address, value, timeout int32) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ModbusPayload{
		Port:     port,
		BaudRate: baudRate,
		Code:     code,
		Slaveid:  slaveid,
		Address:  address,
		Value:    value,
		Timeout:  timeout,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateCommand,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingSystemPort binding system port
func BindingSystemPort(port, protocol, host, tcpport, model string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.SystemBindingPayload{
		Host:  host,
		Port:  tcpport,
		Model: model,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// SelfSample sample self info
func SelfSample(port, model, channel string, quantity int) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.SystemOperationPayload{
		Model:    model,
		Channel:  channel,
		Quantity: quantity,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// SelfCommand self command
func SelfCommand(port, model, value, channel string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout*pmtimeout)
	defer cancel()

	// get param
	var ri sys.HMUSystemReq
	if err := json.Unmarshal([]byte(value), &ri); err != nil {
		return "", err
	}

	var p = public.SystemOperationPayload{
		Model:   model,
		Request: ri,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateCommand,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingPMBusPort binding pmbus port
func BindingPMBusPort(port, protocol string, baudrate, timeout int, soi, ver, adr, eoi byte) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.PMBUSBindingPayload{
		BaudRate: baudrate,
		Timeout:  timeout,
		SOI:      soi,
		VER:      ver,
		ADR:      adr,
		EOI:      eoi,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// PMBusSample pmbus sample
func PMBusSample(port string, cid1, cid2, adr byte, lenid uint16) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.PMBUSOperationPayload{
		ADR:   adr,
		CID1:  cid1,
		CID2:  cid2,
		LENID: lenid,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingOilMachine binding oil machine
func BindingOilMachine(port, protocol string, baudrate, timeout int, soi, eoi byte) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.OilMachineBindingPayload{
		BaudRate: baudrate,
		Timeout:  timeout,
		SOI:      soi,
		EOI:      eoi,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// OilMachineSample oil machine sample
func OilMachineSample(port string, cid1, cid2, adr byte) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.OilMachineOperationPayload{
		ADR:    adr,
		CID1:   cid1,
		CID2:   cid2,
		LENGTH: 0,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// CommonSerialBinding common serial binding
func CommonSerialBinding(port, protocol string, baudrate, timeout int) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.CommonSerialBindingPayload{
		BaudRate: baudrate,
		Timeout:  timeout,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// BindingCreditCard binding oil machine
func BindingCreditCard(port, protocol string, username, password, serialnum string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.CreditCardBindingPayload{
		Username:     username,
		Password:     password,
		SerialNumber: serialnum,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// BindingLuMiGateway binding lumi gateway
func BindingLuMiGateway(port, protocol string, sid, password, netinterface string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.LuMiGatewayBindingPayload{
		SID:          sid,
		Password:     password,
		NetInterface: netinterface,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// LuMiGatewaySample lumi gateway sample
func LuMiGatewaySample(port string, model, sid, value string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.LuMiGatewayOperationPayload{
		SID:   sid,
		Model: model,
		Value: value,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingFaceIPC binding face ipc
func BindingFaceIPC(port, protocol string, host, listenport, muid, suid, server, author, project, token, user string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.FaceIPCBindingPayload{
		Host:         host,
		Port:         listenport,
		MUID:         muid,
		SUID:         suid,
		UploadServer: server,
		Author:       author,
		Project:      project,
		Token:        token,
		User:         user,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// BindingSNMP binding snmp
func BindingSNMP(port, protocol, version, target, targetport, readcommunity, writecommunity string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.SNMPBindingPayload{
		Version:        version,
		Target:         target,
		Port:           targetport,
		ReadCommunity:  readcommunity,
		WriteCommunity: writecommunity,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// SNMPSample snmp sample
func SNMPSample(port string, target string, oids []string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.SNMPOperationPayload{
		Target: target,
		OIDS:   oids,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingWeiGengEntry binding weigeng entry
func BindingWeiGengEntry(port, protocol, localhostAddress, localhostPort, doorAddress, doorPort, serialNo string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.WeiGengEntryBindingPayload{
		LocalAddress: localhostAddress,
		LocalPort:    localhostPort,
		DoorAddress:  doorAddress,
		DoorPort:     doorPort,
		SerialNumber: serialNo,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// EntrySample entry sample
func EntrySample(port string, seqno, code string, group int) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.WeiGengEntryOperationPayload{
		SequenceNumber: seqno,
		FunctionID:     code,
		Group:          group,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingDIDO binding dido
func BindingDIDO(port, protocol string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  "",
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// DIDOSample dido sample
func DIDOSample(port string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: "",
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// ES5200Sample es5200 sample
func ES5200Sample(port string, cid1, cid2, adr, cgroup, ctype byte, lenid int) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ES5200OperationPayload{
		ADR:          adr,
		CID1:         cid1,
		CID2:         cid2,
		LENID:        lenid,
		COMMANDGROUP: cgroup,
		COMMANDTYPE:  ctype,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingCamera binding camera
func BindingCamera(port, protocol, host string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.CameraBindingPayload{
		Host: host,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// CameraSample camera sample
func CameraSample(port, host, channel string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.CameraOperationPayload{
		Host:    host,
		Channel: channel,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// BindingElecFire binding elec fire
func BindingElecFire(port, protocol, host, hostport string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ElecFireBindingPayload{
		Host: host,
		Port: hostport,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// ElecFireSample elec fire sample
func ElecFireSample(port, serialnum string, addr, length int) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.ElecFireOperationPayload{
		SerialNumber: serialnum,
		Address:      addr,
		Length:       length,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}

// Binding Virtual antenna binding Virtual antenna
func BindingVirtualAntenna(port, protocol, host, hostport string) error {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.VirtualAntennaBindingPayload{
		Host: host,
		Port: hostport,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("encode payload failed")
	}

	r, err := pmcli.Binding(ctx, &pb.BindingRequest{
		Port:     port,
		Protocol: protocol,
		Payload:  string(payload),
	})

	if err != nil {
		return fmt.Errorf("could not binding, errmsg[%v]", err)
	}

	if r.Status != public.StatusOK {
		return fmt.Errorf("binding failed, result[%v]", r.Message)
	}

	return nil
}

// VirtualAntennaSample elec fire sample
func VirtualAntennaSample(port, channel string) (string, error) {
	pmcli := pmnet.GetClient()
	pmparentctx := pmnet.GetRootContext()
	pmtimeout := pmnet.GetTimeout()

	if pmcli == nil {
		return "", fmt.Errorf("port client unavailable")
	}

	ctx, cancel := context.WithTimeout(pmparentctx, time.Second*pmtimeout)
	defer cancel()

	var p = public.VirtualAntennaOperationPayload{
		Channel: channel,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("encode payload failed")
	}

	// pmcli := pb.NewPortClient(pmconn)
	r, err := pmcli.Operate(ctx, &pb.OperateRequest{
		Port:    port,
		Type:    public.OperateSample,
		Payload: string(payload),
	})

	if err != nil {
		return "", fmt.Errorf("could not sample: %v", err)
	}

	return r.Data, nil
}
