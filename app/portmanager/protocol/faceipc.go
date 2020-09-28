/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/11/21
 * Despcription: face ipc client define
 *
 */

package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"clc.hmu/app/portmanager/portnet"
	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolFaceIPC, generalFaceIPCDriverProtocol)
}

// Implement DriverProtocol
type faceIPCDriverProtocol struct {
	req *public.FaceIPCBindingPayload
	uri string
}

func (dp *faceIPCDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *faceIPCDriverProtocol) ClientID() string {
	return FaceIPCClientID
}

func (dp *faceIPCDriverProtocol) NewInstance() (PortClient, error) {
	return NewFaceIPCClient(*dp.req)
}

func generalFaceIPCDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeFaceIPCBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &faceIPCDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// FaceIPCClientID id
var FaceIPCClientID = "face-ipc-client-id"

// code define
const (
	faceCodeIPCRegister    = 0x00
	faceCodeHeartbeat      = 0x01
	faceCodeUpgrade        = 0x02
	faceCodeRecognise      = 0x03
	faceCodeTransmission   = 0x07
	faceCodeReboot         = 0x0a
	faceCodePersonRegister = 0x0b
	faceCodeGetConfig      = 0xc0
	faceCodeSetConfig      = 0xc1
)

// register result
const (
	faceRegisterFail      = 0
	faceRegisterSuccess   = 1
	faceUnregisterFail    = 2
	faceUnregisterSuccess = 3
)

// operate
const (
	faceOpRegister     = "register"
	faceOpUnregister   = "unregister"
	faceOpReboot       = "reboot"
	faceOpGetConfig    = "getconfig"
	faceOpSetConfig    = "setconfig"
	faceOpUpgrade      = "upgrade"
	faceOpTransmission = "transmission"
)

// server address flag
const (
	faceFlagIP     = 1
	faceFlagDomain = 2
)

// transmission mode
const (
	transmissionModeUDP = 0
	transmissionModeTCP = 1
)

type faceIPCFrame struct {
	Header          [2]byte // default 0x6868
	ProtocolVersion byte    // default 0x02
	FunctionID      byte    // communication code
	Parameter       []byte  // decision by function id
	Length          int32   // data's length
	Data            []byte  // data
	End             [2]byte // default 0x1616
}

// FaceIPCClient define
type FaceIPCClient struct {
	ClientID string

	conns map[string]net.Conn

	cresponse    chan bool
	responseinfo resonpseInfo

	muid string
	suid string

	appclient portnet.AppClient

	// upload
	uploadServer string
	author       string
	project      string
	token        string
	user         string
}

// resonpseInfo
type resonpseInfo struct {
	// recognise
	fid    string
	status byte

	// device config
	flag          byte
	serverAddress string
	serverPort    int
	deviceID      int
}

// NewFaceIPCClient new client
func NewFaceIPCClient(req public.FaceIPCBindingPayload) (PortClient, error) {
	var client FaceIPCClient

	client.ClientID = FaceIPCClientID
	client.cresponse = make(chan bool)
	client.conns = make(map[string]net.Conn)

	client.muid = req.MUID
	client.suid = req.SUID

	client.uploadServer = req.UploadServer
	client.author = req.Author
	client.project = req.Project
	client.token = req.Token
	client.user = req.User

	client.appclient = portnet.NewAppClient()

	addr := req.Host + ":" + req.Port
	client.startListen(addr)

	return &client, nil
}

// DecodeFaceIPCBindingPayload decode binding
func DecodeFaceIPCBindingPayload(payload string) (public.FaceIPCBindingPayload, error) {
	var p public.FaceIPCBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeFaceIPCOperatePayload decode binding
func DecodeFaceIPCOperatePayload(payload string) (public.FaceIPCOperationPayload, error) {
	var p public.FaceIPCOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// start listen
func (fc *FaceIPCClient) startListen(addr string) error {
	s, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return fmt.Errorf("resolve tcp address failed: %v", err)
	}

	l, err := net.ListenTCP("tcp", s)
	if err != nil {
		return fmt.Errorf("listen tcp failed: %v", err)
	}

	go func(l *net.TCPListener) {
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Printf("accept failed: %v\n", err)
				continue
			}

			fmt.Printf("client connect: %v\n", conn.RemoteAddr())

			go fc.handleConn(conn)
		}
	}(l)

	return nil
}

// deal with connection
func (fc *FaceIPCClient) handleConn(conn net.Conn) {
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

			// check end
			if buf[n-2] == 0x16 && buf[n-1] == 0x16 {
				break
			}
		}

		// fmt.Printf("data: %x, length: %v\n", frame, len(frame))
		if err := fc.handleReceive(conn, frame); err != nil {
			fmt.Printf("handle receive data failed: %v", err)
		}
	}

	// push ipc status to upper
	fc.pushMessage("_state", "设备状态", -1)
}

// handleReceive parse data, get code, param, data
func (fc *FaceIPCClient) handleReceive(conn net.Conn, data []byte) error {
	dl := len(data)
	if dl < 10 {
		return fmt.Errorf("data length too short")
	}

	// check header
	if data[0] != 0x68 || data[1] != 0x68 {
		return fmt.Errorf("data header error")
	}

	// check end
	if data[dl-1] != 0x16 || data[dl-2] != 0x16 {
		return fmt.Errorf("data end error")
	}

	// get code
	code := data[3]

	fmt.Printf("receive data %v bytes, code: %v\n", dl, code)

	switch code {
	case faceCodeIPCRegister:
		// param total 6 bytes: device id (2 bytes), hardware version (2 bytes), software version (2 bytes)
		id := data[4:6]
		hwver := data[6:8]
		swver := data[8:10]

		return fc.handleIPCRegister(conn, id, hwver, swver)
	case faceCodeHeartbeat:
		// no param
		return fc.handleHeartbeat(conn)
	case faceCodeRecognise:
		// param total 88 bytes:
		// session id (10 bytes), flag (2 bytes), device id (2 bytes)
		// record timestamp (4 bytes), face id (20 bytes), recognise info (50 bytes)
		sessionid := data[4:14]
		flag := data[14:16]
		did := data[16:18]
		rt := data[18:22]
		fid := data[22:42]
		rinfo := data[42:92]
		bl := data[92:96]
		d := data[96 : dl-2]

		return fc.handleRecognise(conn, sessionid, flag, did, rt, fid, rinfo, bl, d)
	case faceCodePersonRegister:
		// param 21 bytes: face id (20 bytes), status (1 byte)
		fid := data[4:24]
		status := data[24]

		return fc.handlePersonRegister(conn, fid, status)
	case faceCodeGetConfig:
		// param 59 bytes:
		// flag (1 byte), server address length (2 bytes), server address (50 bytes),
		// server port (4 bytes), device's id (2 bytes)
		flag := data[4]
		sal := data[5:7]
		sa := data[7:57]
		sp := data[57:61]
		did := data[61:63]

		return fc.handleQueryConfig(conn, flag, sal, sa, sp, did)
	}

	return nil
}

// handleIPCRegister save id and conn first, then response
func (fc *FaceIPCClient) handleIPCRegister(conn net.Conn, id, hwver, swver []byte) error {
	did := binary.LittleEndian.Uint16(id)
	sid := strconv.FormatUint(uint64(did), 10)

	// save connection
	fc.conns[sid] = conn
	fmt.Println(fc.conns)

	// param 10 bytes
	param := []byte{}

	// heartbeat interval, 2 bytes, default 0
	param = append(param, 0x3b, 0x00)

	// short connect port, 4 bytes, reserve
	param = append(param, 0x00, 0x00, 0x00, 0x00)

	// server current time, 4 bytes
	var tmp = make([]byte, 4)
	now := time.Now().Unix()
	binary.LittleEndian.PutUint32(tmp, uint32(now))

	param = append(param, tmp...)

	// response
	frame := fc.packFrame(faceCodeIPCRegister, param, nil)
	if err := fc.send(conn, frame); err != nil {
		return fmt.Errorf("send ipc register response failed: %v", err)
	}

	// push ipc status to upper
	return fc.pushMessage("_state", "采集单元状态", 0)
}

// handleHeartbeat response heartbeat only
func (fc *FaceIPCClient) handleHeartbeat(conn net.Conn) error {
	// response
	frame := fc.packFrame(faceCodeHeartbeat, nil, nil)
	if err := fc.send(conn, frame); err != nil {
		return fmt.Errorf("send heartbeat response failed: %v", err)
	}

	return nil
}

// handleRecognise save data, then response
func (fc *FaceIPCClient) handleRecognise(conn net.Conn, sessionid, flag, did, rt, fid, rinfo, bl, d []byte) error {
	// check length
	il := binary.LittleEndian.Uint32(bl)
	if int(il) != len(d) {
		return fmt.Errorf("data length error")
	}

	// save data to file
	sfid := string(bytes.TrimRight(fid, "\x00"))
	url, err := fc.saveFaceData(sfid, d)
	if err != nil {
		return fmt.Errorf("save face data failed: %v", err)
	}

	fmt.Printf("save data successful, url: %v\n", url)

	// param 11 bytes
	param := []byte{}

	// session id, 10 bytes, default set 0
	sid := make([]byte, 10)
	param = append(param, sid...)

	// status, 1 byte: 0 failed, 1 success, default set 1
	param = append(param, 0x01)

	// response
	frame := fc.packFrame(faceCodeRecognise, param, nil)
	if err := fc.send(conn, frame); err != nil {
		return fmt.Errorf("send recognise response failed: %v", err)
	}

	// push recognise info to upper
	chid := "face"
	chname := "识别结果"
	val := sfid + "," + url
	return fc.pushMessage(chid, chname, val)
}

// handlePersonRegister
func (fc *FaceIPCClient) handlePersonRegister(conn net.Conn, fid []byte, status byte) error {
	fc.responseinfo.fid = string(bytes.TrimRight(fid, "\x00"))
	fc.responseinfo.status = status

	fc.cresponse <- true

	return nil
}

// handleQueryConfig
func (fc *FaceIPCClient) handleQueryConfig(conn net.Conn, flag byte, sal, sa, sp, did []byte) error {
	length := int(binary.LittleEndian.Uint16(sal))
	if length > len(sa) {
		return fmt.Errorf("server address too large: %v", length)
	}

	address := string(sa[:length])
	port := binary.LittleEndian.Uint32(sp)
	deviceid := binary.LittleEndian.Uint16(did)

	fc.responseinfo.flag = flag
	fc.responseinfo.serverAddress = address
	fc.responseinfo.serverPort = int(port)
	fc.responseinfo.deviceID = int(deviceid)

	fc.cresponse <- true

	return nil
}

// pack frame
func (fc *FaceIPCClient) packFrame(code byte, param []byte, data []byte) []byte {
	var frame []byte

	// header, version, code
	frame = append(frame, 0x68, 0x68, 0x02, code)

	// param
	if param != nil {
		frame = append(frame, param...)
	}

	// data length
	length := 0
	if data != nil {
		length = len(data)
	}

	var tmp = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, uint32(length))

	frame = append(frame, tmp...)

	// data
	if data != nil {
		frame = append(frame, data...)
	}

	// end
	frame = append(frame, 0x16, 0x16)

	return frame
}

// response
func (fc *FaceIPCClient) send(conn net.Conn, data []byte) error {
	n, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("send data failed: %v", err)
	}

	// fmt.Printf("send data: %x\n", data)
	fmt.Printf("send data bytes: %v\n", n)

	return nil
}

// register
func (fc *FaceIPCClient) register(conn net.Conn, fid string, data []byte) error {
	if data == nil {
		return fmt.Errorf("data empty")
	}

	// set param, total 30 bytes
	var param []byte

	// flag, 2 bytes
	param = append(param, 0x01, 0x00)

	// validity, start time, 4 bytes, default 0
	param = append(param, 0x00, 0x00, 0x00, 0x00)

	// end time 4 bytes, default 0
	param = append(param, 0x00, 0x00, 0x00, 0x00)

	// face id, 20 bytes
	bfid := make([]byte, 20)
	copy(bfid, fid)
	param = append(param, bfid...)

	d := fc.packFrame(faceCodePersonRegister, param, data)
	return fc.send(conn, d)
}

// unregister
func (fc *FaceIPCClient) unregister(conn net.Conn, fid string) error {
	// set param, total 30 bytes
	var param []byte

	// flag, 2 bytes
	param = append(param, 0x03, 0x00)

	// validity, start time, 4 bytes, default 0
	param = append(param, 0x00, 0x00, 0x00, 0x00)

	// end time 4 bytes, default 0
	param = append(param, 0x00, 0x00, 0x00, 0x00)

	// face id, 20 bytes
	bfid := make([]byte, 20)
	copy(bfid, fid)
	param = append(param, bfid...)

	d := fc.packFrame(faceCodePersonRegister, param, nil)
	return fc.send(conn, d)
}

// reboot
func (fc *FaceIPCClient) reboot(conn net.Conn) error {
	// without param
	d := fc.packFrame(faceCodeReboot, nil, nil)
	return fc.send(conn, d)
}

// getconfig
func (fc *FaceIPCClient) getconfig(conn net.Conn) error {
	// without param
	d := fc.packFrame(faceCodeGetConfig, nil, nil)
	return fc.send(conn, d)
}

// setconfig
func (fc *FaceIPCClient) setconfig(conn net.Conn, flag, length, port, did int, address string) error {
	// check param false
	if flag != faceFlagIP && flag != faceFlagDomain {
		return fmt.Errorf("invalid flag: %v", flag)
	}

	if length > len(address) && length > 50 {
		return fmt.Errorf("invalid server address length: %v", length)
	}

	// set param, total 59 bytes
	var param []byte

	// flag, 1 byte
	param = append(param, byte(flag))

	// length, 2 bytes
	var bl = make([]byte, 2)
	binary.LittleEndian.PutUint16(bl, uint16(length))
	param = append(param, bl...)

	// address, 50 bytes
	var baddr = make([]byte, 50)
	copy(baddr, address)
	param = append(param, baddr...)

	// port, 4 bytes
	var bp = make([]byte, 4)
	binary.LittleEndian.PutUint32(bp, uint32(port))
	param = append(param, bp...)

	// device id, 2 bytes
	var bdid = make([]byte, 2)
	binary.LittleEndian.PutUint16(bdid, uint16(did))
	param = append(param, bdid...)

	d := fc.packFrame(faceCodeSetConfig, param, nil)
	return fc.send(conn, d)
}

// upgrade
func (fc *FaceIPCClient) upgrade(conn net.Conn, data []byte) error {
	// TODO check validity of data here

	// without param
	d := fc.packFrame(faceCodeUpgrade, nil, data)
	return fc.send(conn, d)
}

// transmission mode param (0: udp; 1: tcp)
func (fc *FaceIPCClient) transmission(conn net.Conn, dest string, port, mode int, data []byte) error {
	// check dest, must be ip v4
	ip := net.ParseIP(dest)
	if ip == nil {
		return fmt.Errorf("invalid dest ip address")
	}

	// check mode, only support 0 and 1
	if mode != transmissionModeUDP && mode != transmissionModeTCP {
		return fmt.Errorf("invalid mode [%v], only support udp(0) and tcp(1)", mode)
	}

	// set param, total 8 bytes
	var param []byte

	// dest, 4 bytes
	ipv4 := ip.To4()
	param = append(param, ipv4[0], ipv4[1], ipv4[2], ipv4[3])

	// port, 2 bytes
	var bp = make([]byte, 2)
	binary.LittleEndian.PutUint16(bp, uint16(port))
	param = append(param, bp...)

	// transfer mode
	var bm = make([]byte, 2)
	binary.LittleEndian.PutUint16(bm, uint16(mode))
	param = append(param, bm...)

	d := fc.packFrame(faceCodeTransmission, param, data)
	return fc.send(conn, d)
}

// getFaceData get from http server by default
func (fc *FaceIPCClient) getFaceData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// saveFaceData post data to http server, return url
func (fc *FaceIPCClient) saveFaceData(fid string, data []byte) (string, error) {
	filepre := os.TempDir()
	filename := fid + ".jpg"
	fp := filepath.Join(filepre, filename)
	if err := ioutil.WriteFile(fp, data, 0x0666); err != nil {
		return "", fmt.Errorf("write file [%s] failed: %v", filename, err)
	}

	return public.UploadFile(fp, filename, fc.uploadServer, fc.author, fc.project, fc.token, fc.user)
}

// pushMessage push message to upper
func (fc *FaceIPCClient) pushMessage(chanid, name string, value interface{}) error {
	// pack payload
	payload := public.MessagePayload{
		MonitoringUnitID: fc.muid,
		SampleUnitID:     fc.suid,
		ChannelID:        chanid,
		Name:             name,
		Value:            value,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	topic := "sample-values/" + fc.muid + "/" + fc.suid + "/" + chanid
	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}

// ID id
func (fc *FaceIPCClient) ID() string {
	return fc.ClientID
}

// Sample get values
func (fc *FaceIPCClient) Sample(payload string) (string, error) {
	return "", nil
}

// Command set values
func (fc *FaceIPCClient) Command(payload string) (string, error) {
	p, err := DecodeFaceIPCOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	// get device id
	did := p.CameraID

	// get connection
	conn, ok := fc.conns[did]
	if !ok {
		return "", fmt.Errorf("face ipc [%v] do not exist", did)
	}

	// register or unregister
	switch p.Perpose {
	case faceOpRegister:
		{
			// get face data
			data, err := fc.getFaceData(p.FaceURL)
			if err != nil {
				return "", fmt.Errorf("read face data failed: %v", err)
			}

			if err := fc.register(conn, p.FaceID, data); err != nil {
				return "", fmt.Errorf("regitser failed: %v", err)
			}
		}
	case faceOpUnregister:
		{
			if err := fc.unregister(conn, p.FaceID); err != nil {
				return "", fmt.Errorf("unregitser failed: %v", err)
			}
		}
	case faceOpReboot:
		{
			if err := fc.reboot(conn); err != nil {
				return "", fmt.Errorf("reboot failed: %v", err)
			}

			// return direct
			return "ok", nil
		}
	case faceOpGetConfig:
		{
			if err := fc.getconfig(conn); err != nil {
				return "", fmt.Errorf("getconfig failed: %v", err)
			}
		}
	case faceOpSetConfig:
		{
			if err := fc.setconfig(conn, p.Flag, p.AddressLength, p.Port, p.DeviceID, p.Address); err != nil {
				return "", fmt.Errorf("setconfig failed: %v", err)
			}

			// return direct
			return "ok", nil
		}
	case faceOpTransmission:
		{
			if err := fc.transmission(conn, p.Destination, p.Port, p.Mode, []byte(p.Data)); err != nil {
				return "", fmt.Errorf("transmission failed: %v", err)
			}

			// return direct
			return "ok", nil
		}
	case faceOpUpgrade:
		{
			if err := fc.upgrade(conn, []byte(p.Data)); err != nil {
				return "", fmt.Errorf("upgrade failed: %v", err)
			}

			// return direct
			return "ok", nil
		}
	}

	select {
	case <-fc.cresponse:
		switch p.Perpose {
		case faceOpRegister, faceOpUnregister:
			fmt.Printf("operate result: %v %v %v\n", p.Perpose, fc.responseinfo.fid, fc.responseinfo.status)
			if fc.responseinfo.status == faceRegisterFail || fc.responseinfo.status == faceUnregisterFail {
				return "", fmt.Errorf("operate fail: %v", p.Perpose)
			}
		case faceOpGetConfig:
			fmt.Printf("device config, flag[%v], server host[%v], server port[%v], device id[%v]\n", fc.responseinfo.flag, fc.responseinfo.serverAddress, fc.responseinfo.serverPort, fc.responseinfo.deviceID)
			return fmt.Sprintf("%s,%d,%d", fc.responseinfo.serverAddress, fc.responseinfo.serverPort, fc.responseinfo.deviceID), nil
		}
	case <-time.After(time.Second * 10):
		return "", fmt.Errorf("operate timeout: %v", p.Perpose)
	}

	return "ok", nil
}
