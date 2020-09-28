/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: hmu2000 definition
 *
 */

package sys

import (
	"time"

	"clc.hmu/app/public/store/muid"
	"github.com/gwaylib/errors"
)

// SystemClient system client
type SystemNone struct {
}

func init() {
	RegSysClientModel(MODEL_DEFAULT, &SystemNone{})
}

var DefaultHMUSystemResp = &HMUSystemResp{
	// TODO:make default value
	UUID: muid.GetMuID(),
}

func (c *SystemNone) New(opts *SystemServerOption) SystemClient {
	return &SystemNone{}
}
func (c *SystemNone) ModelName() string {
	return MODEL_DEFAULT
}

// Disconnect disconnect
func (c *SystemNone) Disconnect() error {
	return nil
}

// Request requst to system daemon
func (c *SystemNone) Request(req *HMUSystemReq) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

func (c *SystemNone) AutoCheckNetworking(urls []string, timeout time.Duration) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// GPS query gps info
func (c *SystemNone) GPS() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// Time query time info
func (c *SystemNone) Time() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetTimeServer set time server
func (c *SystemNone) SetTimeServer(timeserver string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SystemInfo query system info, including memory, cpu and so on
func (c *SystemNone) SystemInfo() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// Reboot reboot device
func (c *SystemNone) Reboot() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// UUID query device's uuid
func (c *SystemNone) UUID() (*HMUSystemResp, error) {
	// 默认总是返回空id
	return DefaultHMUSystemResp, nil
}

// SetUUID set device's uuid
func (c *SystemNone) SetUUID(id string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// LAN query LAN info
func (c *SystemNone) LAN() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// WAN query WAN info
func (c *SystemNone) WAN() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// Wireless query wireless info
func (c *SystemNone) Wireless() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetAP set ap mode
func (c *SystemNone) SetAP(ssid, encryption, key, channel, hide string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetLAN set lan
func (c *SystemNone) SetLAN(ip, mask string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// Internet query internet info
func (c *SystemNone) Internet() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetEthStatic set eth static mode
func (c *SystemNone) SetEthStatic(ip, mask, gateway, pdns, sdns string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetEthDHCP set eth dhcp mode
func (c *SystemNone) SetEthDHCP() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetWifi set wifi mode
func (c *SystemNone) SetWifi(ssid, key string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetLTE set lte mode
func (c *SystemNone) SetLTE() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetLTE set lte mode
func (c *SystemNone) LTECSQ() (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// FactoryReset factory reset
func (c *SystemNone) FactoryReset(mode string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetAppLED set appled
func (c *SystemNone) SetAppLED(status int) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// DIDO get dido value
func (c *SystemNone) DIDO(id string) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}

// SetDO set do
func (c *SystemNone) SetDO(id string, value int) (*HMUSystemResp, error) {
	// 默认未实现，总是返回不通
	return nil, errors.New("UnImplements")
}
