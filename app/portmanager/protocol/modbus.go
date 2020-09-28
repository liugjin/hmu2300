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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"clc.hmu/app/public"

	"clc.hmu/app/public/log/portlog"
	"github.com/goburrow/modbus"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolModbusSerial, generalModbusSerialDriverProtocol)
	RegDriverProtocol(public.ProtocolModbusTCP, generalModbusTCPDriverProtocol)
}

// Implement DriverProtocol
type modbusSerialDriverProtocol struct {
	req *public.ModbusPayload
	uri string
}

func (dp *modbusSerialDriverProtocol) Payload() interface{} {
	return dp.req
}
func (dp *modbusSerialDriverProtocol) ClientID() string {
	return strconv.Itoa(int(dp.req.Slaveid)) + strconv.Itoa(int(dp.req.BaudRate))
}

func (dp *modbusSerialDriverProtocol) NewInstance() (PortClient, error) {
	return NewModbusSerialClient(dp.uri, dp.req.BaudRate, dp.req.Timeout, byte(dp.req.Slaveid))
}

func generalModbusSerialDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &modbusSerialDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// Implement DriverProtocol
type modbusTCPDriverProtocol struct {
	req *public.ModbusPayload
	uri string
}

func (dp *modbusTCPDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *modbusTCPDriverProtocol) ClientID() string {
	return strconv.Itoa(int(dp.req.Slaveid)) + strconv.Itoa(int(dp.req.BaudRate))
}

func (dp *modbusTCPDriverProtocol) NewInstance() (PortClient, error) {
	return NewModbusTCPClient(dp.uri, dp.req.BaudRate, dp.req.Timeout, byte(dp.req.Slaveid))
}

func generalModbusTCPDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}
	return &modbusTCPDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// ModbusClient modbus client
type ModbusClient struct {
	SlaveID  byte
	BaudRate int32
	Handler  modbus.ClientHandler
	Client   modbus.Client

	mtx sync.Mutex
}

// DecodeModbusPayload decode modbus payload
func DecodeModbusPayload(payload string) (public.ModbusPayload, error) {
	var p public.ModbusPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// NewModbusClient new modbus client
func NewModbusClient(protocol, port, payload string) (*ModbusClient, error) {
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return &ModbusClient{}, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	if protocol == public.ProtocolModbusSerial {
		return NewModbusSerialClient(port, req.BaudRate, req.Timeout, byte(req.Slaveid))
	} else if protocol == public.ProtocolModbusTCP {
		return NewModbusTCPClient(port, req.BaudRate, req.Timeout, byte(req.Slaveid))
	}

	return &ModbusClient{}, fmt.Errorf("unknown protocol")
}

// NewModbusSerialClient port is local physical port, like /dev/ttyS0
func NewModbusSerialClient(port string, baudrate, timeout int32, slaveid byte) (*ModbusClient, error) {
	// new handler
	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = int(baudrate)
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = slaveid
	handler.Timeout = time.Millisecond * time.Duration(timeout)
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)

	if err := handler.Connect(); err != nil {
		return &ModbusClient{}, fmt.Errorf("connect port[%s] failed, baudrate[%v], slaveid[%v], errmsg[%v]", port, baudrate, slaveid, err)
	}

	portlog.LOG.Infof("connect port[%s] success,  baudrate[%v], slaveid[%v]", port, baudrate, slaveid)

	// defer handler.Close()
	var client = ModbusClient{
		SlaveID:  slaveid,
		BaudRate: baudrate,
		Handler:  handler,
		Client:   modbus.NewClient(handler),
	}

	return &client, nil
}

// NewModbusTCPClient port is local or remote software port, like "127.0.0.1:502", timeout unit use millisecond
func NewModbusTCPClient(port string, baudrate, timeout int32, slaveid byte) (*ModbusClient, error) {
	handler := modbus.NewTCPClientHandler(port)
	handler.Timeout = time.Duration(timeout) * time.Millisecond
	handler.SlaveId = slaveid

	if err := handler.Connect(); err != nil {
		return &ModbusClient{}, fmt.Errorf("connect address[%s] failed, slaveid[%v], errmsg[%v]", port, slaveid, err)
	}

	portlog.LOG.Infof("connect address[%s] success, slaveid[%v]", port, slaveid)

	// defer handler.Close()

	var client = ModbusClient{
		SlaveID:  slaveid,
		BaudRate: baudrate,
		Handler:  handler,
		Client:   modbus.NewClient(handler),
	}

	return &client, nil
}

// Release release implement
func (mc *ModbusClient) Release() {
	// mc.Handler.Close()
}

// ID client's id
func (mc *ModbusClient) ID() string {
	return strconv.Itoa(int(mc.SlaveID)) + strconv.Itoa(int(mc.BaudRate))
}

// Sample modbus sample
func (mc *ModbusClient) Sample(payload string) (string, error) {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()

	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}
	result, err := mc.sample(&req)
	if err != nil {
		return "", errors.As(err)
	}

	return fmt.Sprintf("%x", result), nil
}
func (mc *ModbusClient) sample(req *public.ModbusPayload) ([]byte, error) {
	var result []byte

	address := uint16(req.Address)
	quantity := uint16(req.Quantity)
	client := mc.Client

	switch req.Code {
	case 0x01:
		r, err := client.ReadCoils(address, quantity)
		if err != nil {
			return nil, fmt.Errorf("read failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	case 0x02:
		r, err := client.ReadDiscreteInputs(address, quantity)
		if err != nil {
			return nil, fmt.Errorf("read failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	case 0x03:
		r, err := client.ReadHoldingRegisters(address, quantity)
		if err != nil {
			return nil, fmt.Errorf("read failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	case 0x04:
		r, err := client.ReadInputRegisters(address, quantity)
		if err != nil {
			return nil, fmt.Errorf("read failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	}

	// fmt.Printf("code %v, %x\n", req.Code, result)
	return result, nil
}

// Command modbus command
func (mc *ModbusClient) Command(payload string) (string, error) {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()

	// decode payload
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	result, err := mc.command(&req)
	if err != nil {
		return "", errors.As(err)
	}
	return fmt.Sprintf("%x", result), nil
}

func (mc *ModbusClient) command(req *public.ModbusPayload) ([]byte, error) {
	// portlog.LOG.Infof("payload {%+v}", req)

	address := uint16(req.Address)
	quantity := uint16(req.Quantity)

	ival := uint16(0)
	bval := []byte{}

	switch v := req.Value.(type) {
	case int:
		ival = uint16(v)
	case float64:
		ival = uint16(v)
	case string:
		sep := ","
		vals := strings.Split(v, sep)

		if len(vals) != int(quantity) {
			return nil, fmt.Errorf("number of value not match, want[%v], input[%v], split by ','", quantity, len(vals))
		}

		// tranfer to bytes
		for _, v := range vals {
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("invalid param [%v]", v)
			}

			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(i))
			bval = append(bval, b...)
		}
	default:
		return nil, fmt.Errorf("value type not support")
	}

	client := mc.Client

	var result []byte

	switch req.Code {
	case 0x05:
		r, err := client.WriteSingleCoil(address, ival)
		if err != nil {
			return nil, fmt.Errorf("write failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	case 0x06:
		r, err := client.WriteSingleRegister(address, ival)
		if err != nil {
			return nil, fmt.Errorf("write failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	case 0x0F:
		r, err := client.WriteMultipleCoils(address, quantity, bval)
		if err != nil {
			return nil, fmt.Errorf("write failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	case 0x10:
		r, err := client.WriteMultipleRegisters(address, quantity, bval)
		if err != nil {
			return nil, fmt.Errorf("write failed, port[%s], errmsg[%v]", req.Port, err)
		}

		result = r
	}

	return result, nil
}
