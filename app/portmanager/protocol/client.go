/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/06/21
 * Despcription: port client define
 *
 */

package protocol

import (
	"clc.hmu/app/public/log"
	"github.com/gwaylib/errors"
)

// PortClient port client, one port client is related to one sample unit
type PortClient interface {
	// specified client's ID, use for searching
	ID() string

	// sample, get values
	Sample(payload string) (string, error)

	// command, set values
	Command(payload string) (string, error)
}

// 用于注册实例驱动的入口
type DriverProtocol interface {
	// 返回解析后的协议模板
	// 在老协议中，该值是原Operate的payload值
	// 在V1中，该值是采集单元配置文件的ID值
	Payload() interface{}
	// 通过解析payload返回驱动协议ClientID
	ClientID() string
	// 使用模板协议进行设备连接并生成连接的客户端
	NewInstance() (PortClient, error)
}

// 生成驱动协议的入口
//
// 参数
// uri -- 设备的端口或地址
// suid -- 采集单元的ID
// payload -- 兼容原老协议而设置此值
type DriverGeneral func(uri, suid, payload string) (DriverProtocol, error)

// 驱动协议注册
var drvs = map[string]DriverGeneral{}

func RegDriverProtocol(p string, g DriverGeneral) {
	_, ok := drvs[p]
	if ok {
		log.Fatalf("protocol has registered:%s", p)
	}
	drvs[p] = g
}

// 获取驱动接入驱动接口
// payload是一个需要驱动生成需要的协议模板
func GetDriverProtocol(drvName, uri, suid, payload string) (DriverProtocol, error) {
	fn, ok := drvs[drvName]
	if !ok {
		return nil, errors.ErrNoData.As(drvName)
	}
	return fn(uri, suid, payload)
}
