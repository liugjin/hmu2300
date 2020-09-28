/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2019/03/21
 * Despcription: electrical fire client define
 *
 */

package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolElecFire, generalElecFireDriverProtocol)
}

// Implement DriverProtocol
type elecFireDriverProtocol struct {
	req *public.ElecFireBindingPayload
	uri string
}

func (dp *elecFireDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *elecFireDriverProtocol) ClientID() string {
	return ElecFireClientID + dp.req.Host
}

func (dp *elecFireDriverProtocol) NewInstance() (PortClient, error) {
	return NewElecFireClient(*dp.req)
}

func generalElecFireDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeElecFireBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &elecFireDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

const (
	cmdSetParam     byte = 0x82
	cmdRegister          = 0x84
	cmdAlarm             = 0x89
	cmdTransmission      = 0x90
	cmdUpload            = 0x91
	cmdTiming            = 0x93
)

const (
	addressUpload1       = 0x1000
	addressUpload1Length = 42

	addressUpload2       = 0x1204
	addressUpload2Length = 26

	addressUpload3       = 0x1300
	addressUpload3Length = 2

	addressAlarm1       = 0x1231
	addressAlarm1Length = 13

	addressAlarm2       = 0x102F
	addressAlarm2Length = 36
)

var (
	frameHeader = []byte{0x7b, 0x7b}
	frameEnd    = []byte{0x7d, 0x7d}
)

// ElecFireClientID id
var ElecFireClientID = "electrical-fire-client-id"

// ElecFireClient client
type ElecFireClient struct {
	ClientID string

	conns   map[net.Conn]string
	devices map[string]efDeviceData
}

type efDeviceData struct {
	SerialNumber   string
	CardNumber     string
	RSSI           byte
	UploadInterval byte

	AlarmData  []byte
	UploadData []byte
}

// NewElecFireClient new client
func NewElecFireClient(req public.ElecFireBindingPayload) (PortClient, error) {
	var client ElecFireClient

	client.ClientID = ElecFireClientID
	client.conns = make(map[net.Conn]string)
	client.devices = make(map[string]efDeviceData)

	addr := req.Host + ":" + req.Port
	go client.startListen(addr)

	return &client, nil
}

// DecodeElecFireBindingPayload decode binding payload
func DecodeElecFireBindingPayload(payload string) (public.ElecFireBindingPayload, error) {
	var p public.ElecFireBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeElecFireOperationPayload decode binding payload
func DecodeElecFireOperationPayload(payload string) (public.ElecFireOperationPayload, error) {
	var p public.ElecFireOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// start listen
func (ec *ElecFireClient) startListen(addr string) error {
	s, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return fmt.Errorf("resolve tcp address failed: %v", err)
	}

	l, err := net.ListenTCP("tcp", s)
	if err != nil {
		return fmt.Errorf("listen tcp failed: %v", err)
	}

	log.Printf("start listen at %s", addr)

	go func(l *net.TCPListener) {
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Printf("accept failed: %v\n", err)
				continue
			}

			fmt.Printf("client connect: %v\n", conn.RemoteAddr())

			go ec.handleConn(conn)
		}
	}(l)

	return nil
}

func (ec *ElecFireClient) handleConn(conn net.Conn) {
HANDLECONN:
	for {
		var frame []byte
		for {
			var buf = make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("read failed: %v\n", err)
				break HANDLECONN
			}

			frame = append(frame, buf[:n]...)
			l := len(frame)

			// check end
			if l > 2 && frame[l-2] == 0x7d && frame[l-1] == 0x7d {
				break
			}
		}

		fmt.Printf("data: %x, length: %v\n", frame, len(frame))

		l := len(frame)
		if l >= 7 {
			switch frame[2] {
			case cmdTiming:
				ec.handleTiming(conn)
			case cmdSetParam:
				ec.handleSetParam(conn)
			case cmdRegister:
				data := frame[3 : l-4]
				ec.handleRegister(conn, data)
			case cmdAlarm:
				data := frame[3 : l-4]
				ec.handleAlarm(conn, data)
			case cmdUpload:
				data := frame[3 : l-4]
				ec.handleUpload(conn, data)
			case cmdTransmission:
				ec.handleTransmission(conn)
			default:
				fmt.Printf("receive unknown frame: %x", frame)
			}
		}
	}
}

func (ec *ElecFireClient) send(conn net.Conn, msg []byte) (int, error) {
	if conn == nil {
		return 0, fmt.Errorf("invalid connection")
	}

	fmt.Printf("send: %x\n", msg)
	return conn.Write(msg)
}

func (ec *ElecFireClient) packFrame(cmd byte, d []byte) []byte {
	data := []byte{cmd}
	if d != nil {
		data = append(data, d...)
	}

	var crc public.CRC
	crc.Reset()
	crc.PushBytes(data)

	var bc = make([]byte, 2)
	binary.BigEndian.PutUint16(bc, crc.Value())

	data = append(frameHeader, data...)
	data = append(data, bc...)
	data = append(data, frameEnd...)

	return data
}

func (ec *ElecFireClient) handleTiming(conn net.Conn) (int, error) {
	// get current time
	now := time.Now()
	msg := fmt.Sprintf("%.2x%.2x%.2x%.2x%.2x%.2x%.2x", byte(now.Year()-2000), byte(now.Month()), byte(now.Day()), byte(now.Weekday()), byte(now.Hour()), byte(now.Minute()), byte(now.Second()))

	// pack
	tmp, _ := hex.DecodeString(msg)
	resp := ec.packFrame(cmdTiming, tmp)

	return ec.send(conn, resp)
}

func (ec *ElecFireClient) handleSetParam(conn net.Conn) (int, error) {
	return 0, nil
}

func (ec *ElecFireClient) handleRegister(conn net.Conn, data []byte) (int, error) {
	// serial number, 20 bytes
	sn := strings.Replace(string(data[:20]), "\x00", "", -1)

	// card number, 30 bytes
	cn := strings.Replace(string(data[20:50]), "\x00", "", -1)

	// RSSI, 1 byte
	rssi := data[50]

	// firmware version(1,2,3), 6 bytes, ignore

	// upload interval, 1 byte
	ui := data[57]

	var d = efDeviceData{
		SerialNumber:   sn,
		CardNumber:     cn,
		RSSI:           rssi,
		UploadInterval: ui,
	}

	ec.devices[sn] = d
	ec.conns[conn] = sn

	resp := ec.packFrame(cmdRegister, nil)
	return ec.send(conn, resp)
}

func (ec *ElecFireClient) handleAlarm(conn net.Conn, data []byte) (int, error) {
	sn, ok := ec.conns[conn]
	if !ok {
		return 0, fmt.Errorf("invalid connection")
	}

	d := ec.devices[sn]
	d.AlarmData = data
	ec.devices[sn] = d

	// response
	resp := ec.packFrame(cmdAlarm, nil)
	return ec.send(conn, resp)
}

func (ec *ElecFireClient) handleUpload(conn net.Conn, data []byte) (int, error) {
	sn, ok := ec.conns[conn]
	if !ok {
		return 0, fmt.Errorf("invalid connection")
	}

	d := ec.devices[sn]
	d.UploadData = data
	ec.devices[sn] = d

	// response
	resp := ec.packFrame(cmdUpload, nil)
	return ec.send(conn, resp)
}

func (ec *ElecFireClient) handleTransmission(conn net.Conn) (int, error) {
	return 0, nil
}

func (ec *ElecFireClient) parse(sn string, addr int, length int) ([]byte, error) {
	d, ok := ec.devices[sn]
	if !ok {
		return nil, fmt.Errorf("device not exist")
	}

	offset := addr - addressUpload1
	if offset >= 0 && (offset+length <= addressUpload1Length) {
		// upload 1 data start at 10th byte
		start := 10 + offset*2
		end := start + length*2

		if end > len(d.UploadData) {
			return nil, fmt.Errorf("out of range")
		}

		return d.UploadData[start:end], nil
	}

	offset = addr - addressUpload2
	if offset >= 0 && (offset+length <= addressUpload2Length) {
		// upload 2 data start at 110th byte
		start := 110 + offset*2
		end := start + length*2

		if end > len(d.UploadData) {
			return nil, fmt.Errorf("out of range")
		}

		return d.UploadData[start:end], nil
	}

	offset = addr - addressUpload3
	if offset >= 0 && (offset+length <= addressUpload3Length) {
		// upload 3 data start at 178th byte
		start := 178 + offset*2
		end := start + length*2

		if end > len(d.UploadData) {
			return nil, fmt.Errorf("out of range")
		}

		return d.UploadData[start:end], nil
	}

	offset = addr - addressAlarm1
	if offset >= 0 && (offset+length <= addressAlarm1Length) {
		// alarm 1 data start at 10th byte
		start := 10 + offset*2
		end := start + length*2

		if end > len(d.AlarmData) {
			return nil, fmt.Errorf("out of range")
		}

		return d.AlarmData[start:end], nil
	}

	offset = addr - addressAlarm2
	if offset >= 0 && (offset+length <= addressAlarm2Length) {
		// alarm 2 data start at 52th byte
		start := 52 + offset*2
		end := start + length*2

		if end > len(d.AlarmData) {
			return nil, fmt.Errorf("out of range")
		}

		return d.AlarmData[start:end], nil
	}

	return nil, fmt.Errorf("data do not exist")
}

// ID id
func (ec *ElecFireClient) ID() string {
	return ec.ClientID
}

// Sample get values
func (ec *ElecFireClient) Sample(payload string) (string, error) {
	p, err := DecodeElecFireOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode elec fire operation payload failed, errmsg {%v}", err)
	}

	// log.Println("sample:", p, "data:", ec.devices)

	result, err := ec.parse(p.SerialNumber, p.Address, p.Length)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", result), nil
}

// Command set values
func (ec *ElecFireClient) Command(payload string) (string, error) {

	return "ok", nil
}
