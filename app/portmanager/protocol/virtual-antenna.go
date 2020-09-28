package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolVirtualAntenna, generalVirtualAntennaDriverProtocol)
}

// Implement DriverProtocol
type virtualAntennaDriverProtocol struct {
	req *public.VirtualAntennaBindingPayload
	uri string
}

func (dp *virtualAntennaDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *virtualAntennaDriverProtocol) ClientID() string {
	return VirtualAntennaClientID
}

func (dp *virtualAntennaDriverProtocol) NewInstance() (PortClient, error) {
	return NewVirtualAntennaClient(*dp.req)
}

func generalVirtualAntennaDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	req, err := DecodeVirtualAntennaBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}

	return &virtualAntennaDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// VirtualAntennaClientID id
var VirtualAntennaClientID = "virtual-antenna-client-id"

type VirtualAntennaClient struct {
	ClientID string

	muid string
	suid string

	conn net.Conn

	B9Data
}

//总帧格式
type Frame2 struct {
	Stx  uint16  // 2字节，0x00FFFFH
	Ver  uint8   // 1字节, 0x00H
	Seq  uint8   // 1字节
	Len  uint32  // 4字节
	Data []uint8 // 变长节字
	Crc  uint16  // 2字节
}

//C0帧格式
type Data struct {
	FrameType    uint8   // 1字节, 0xC0
	Datetime     []uint8 // 7字节
	LaneMode     uint8
	BSTInterval  uint8
	TxPower      uint8
	PLLChannelID uint8
	Worktype     uint8
	RecordNo     uint8
	HeartBeat    uint8
	ProvinceID   uint16
}

//C4帧格式
type C4Data struct {
	FrameType uint8 // 1字节, 0xC0
	Control   uint8
}

//B9帧格式
type B9Data struct {
	ControlCount     uint8
	ControlStatusN   []uint8 // len为ControlCount
	AntennaCount     uint8
	AntennaInfoNList []AntennaInfoN //len为AntennaCount * 4
}

type AntennaInfoN struct {
	Status  uint8
	Power   uint8
	Channel uint8
	Control uint8
}

func (f *Frame2) Serial() []byte {
	//enc := make([]byte, 10+len(f.Data))
	var enc []byte
	//var result []byte
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, uint16(f.Stx))
	enc = append(enc, result[:2]...)

	binary.BigEndian.PutUint16(result, uint16(f.Ver))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(f.Seq))
	enc = append(enc, result[1:2]...)

	lenResult := make([]byte, 4)
	binary.BigEndian.PutUint32(lenResult, uint32(f.Len))
	enc = append(enc, lenResult...)

	var Data []byte
	for _, v := range f.Data {
		binary.BigEndian.PutUint16(result, uint16(v))
		Data = append(Data, result[1:2]...)
	}
	enc = append(enc, Data...)

	binary.BigEndian.PutUint16(result, uint16(f.Crc))
	enc = append(enc, result...)
	return enc
}

func (d *Data) Serial() []byte {
	var enc []byte
	result := make([]byte, 8)
	binary.BigEndian.PutUint16(result, uint16(d.FrameType))
	enc = append(enc, result[1:2]...)

	var Data []byte
	for _, v := range d.Datetime {
		binary.BigEndian.PutUint16(result, uint16(v))
		Data = append(Data, result[1:2]...)
	}
	enc = append(enc, Data...)

	binary.BigEndian.PutUint16(result, uint16(d.LaneMode))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.BSTInterval))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.TxPower))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.PLLChannelID))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.Worktype))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.RecordNo))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.HeartBeat))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(d.ProvinceID))
	enc = append(enc, result[:2]...)
	return enc
}

func (c4 *C4Data) Serial() []byte {
	//enc := make([]byte, 10+len(f.Data))
	var enc []byte
	//var result []byte
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, uint16(c4.FrameType))
	enc = append(enc, result[1:2]...)

	binary.BigEndian.PutUint16(result, uint16(c4.Control))
	enc = append(enc, result...)
	return enc
}

func NewVirtualAntennaClient(req public.VirtualAntennaBindingPayload) (PortClient, error) {
	var client VirtualAntennaClient
	addr := req.Host + ":" + req.Port
	client.ClientID = VirtualAntennaClientID
	client.startListen(addr)
	return &client, nil
}

// DecodeVirtualAntennaBindingPayload decode binding
func DecodeVirtualAntennaBindingPayload(payload string) (public.VirtualAntennaBindingPayload, error) {
	var p public.VirtualAntennaBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}
	return p, nil
}

// DecodeFaceIPCOperatePayload decode binding
func DecodeVirtualAntennaOperatePayload(payload string) (public.VirtualAntennaOperationPayload, error) {
	var p public.VirtualAntennaOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}
	return p, nil
}

// start listen
func (vc *VirtualAntennaClient) startListen(addr string) error {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("accept failed: %v\n", err)
			time.Sleep(time.Second * 2)
			continue
		}

		fmt.Printf("client connect: %v\n", conn.RemoteAddr())

		break
	}
	vc.conn = conn
	go vc.handleConn(conn) // 有判断conn是否失效 go
	return nil
}

// deal with connection
func (vc *VirtualAntennaClient) handleConn(conn net.Conn) {
HANDLECONN:
	for {
		var frame []byte
		for {
			var buf = make([]byte, 1024)
			n, err := vc.conn.Read(buf)

			if err != nil {
				//vc.conn = nil
				fmt.Printf("read failed: %v\n", err)
				break HANDLECONN
			}

			frame = append(frame, buf[:n]...)

			fmt.Println(frame)

			//判断b0帧
			if frame[8] == 0xb0 {
				if len(frame) == 20 {
					break
				}
			}

			//判断b9帧
			if frame[8] == 0xb9 {
				if len(frame) == 27 {
					break
				}
			}

		}
		if err := vc.handleReceive(conn, frame); err != nil {
			fmt.Printf("handle receive data failed: %v", err)
		}
	}

	// push ipc status to upper
	//vc.pushMessage("_state", "设备状态", -1)
}

func (vc *VirtualAntennaClient) handleReceive(conn net.Conn, data []byte) error {
	dl := len(data)
	if dl < 10 {
		return fmt.Errorf("data length too short")
	}

	// check header
	if data[0] != 0xFF || data[1] != 0xFF {
		return fmt.Errorf("data header error")
	}

	// get code
	code := data[8]

	fmt.Printf("receive data %v bytes, code: %v\n", dl, code)

	switch code {
	case 0xb0:
		if data[9] != 0x00 {
			return fmt.Errorf("response B0 Heartbeat response  error code: %v", data[9])
		}
	case 0xb9:
		// no param
		status := data[8 : len(data)-2]
		return vc.handleHeartbeat(conn, status)

	}

	return nil
}

func (vc *VirtualAntennaClient) handleHeartbeat(conn net.Conn, data []byte) error {
	fmt.Printf("%x\n", data)
	if data[1] != 0x00 {
		return fmt.Errorf("response B9 Heartbeat response  error code: %v", data[1])
	}

	controlCount := data[9]
	controlStatusN := data[10 : 10+data[9]]
	antennaCount := data[10+data[9]]
	var antennaInfoNList []AntennaInfoN
	var a AntennaInfoN
	for i := 0; i < int(antennaCount); i++ {
		for j := 0; j < 4; j++ {
			a.Status = data[10+data[9]+4*(antennaCount-1)+1]
			a.Power = data[10+data[9]+4*(antennaCount-1)+2]
			a.Channel = data[10+data[9]+4*(antennaCount-1)+3]
			a.Control = data[10+data[9]+4*(antennaCount-1)+4]
		}
		antennaInfoNList = append(antennaInfoNList, a)
	}

	b9 := B9Data{
		ControlCount:     controlCount,
		ControlStatusN:   controlStatusN,
		AntennaCount:     antennaCount,
		AntennaInfoNList: antennaInfoNList,
	}

	//怎么样mqtt发送状态
	fmt.Println(b9)
	vc.B9Data = b9

	return nil
}

//返回c0帧
func C0Frame() []byte {
	t, _ := hex.DecodeString(DateTimeChange())

	//以下为组帧
	d := Data{
		FrameType:    0xC0,
		Datetime:     t,
		LaneMode:     0x9,
		BSTInterval:  0x0a,
		TxPower:      0x01,
		PLLChannelID: 0x00,
		Worktype:     0x05,
		RecordNo:     0x00,
		HeartBeat:    0x0a,
		ProvinceID:   0x0005,
	}
	fmt.Printf("C0帧： %x\n", d.Serial())

	data := d.Serial()
	l := len(data)

	f := Frame2{
		Stx:  0xffff,
		Ver:  0x00,
		Seq:  0x10,
		Len:  uint32(l),
		Data: data,
		Crc:  0x0000,
	}
	fmt.Printf("%v\n", f.Serial())

	FrameData := f.Serial()
	f.Crc = Crc_16_x25(FrameData[2 : len(FrameData)-2])
	fmt.Printf("C0总帧：%x\n", f.Serial())
	return f.Serial()
}

func DateTimeChange() string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	y1, _ := strconv.Atoi(year[:2])
	y2, _ := strconv.Atoi(year[2:])

	month := strconv.Itoa(int(now.Month()))
	m, _ := strconv.Atoi(month)

	day := strconv.Itoa(now.Day())
	d, _ := strconv.Atoi(day)

	hour := strconv.Itoa(now.Hour())
	h, _ := strconv.Atoi(hour)

	minute := strconv.Itoa(now.Minute())
	min, _ := strconv.Atoi(minute)

	second := strconv.Itoa(now.Second())
	s, _ := strconv.Atoi(second)

	//b = []byte{uint8(y1), uint8(y2), uint8(m), uint8(d), uint8(h), uint8(min), uint8(s)}
	ss := fmt.Sprintf("%.2d%.2d%.2d%.2d%.2d%.2d%.2d", y1, y2, m, d, h, min, s)

	return ss
}

//返回C4帧 0是关 1开
func C4Frame(control int) []byte {

	//以下为组帧
	d := C4Data{
		FrameType: 0xC4,
		Control:   uint8(control),
	}

	data := d.Serial()
	l := len(data)

	f := Frame2{
		Stx:  0xffff,
		Ver:  0x00,
		Seq:  0x90,
		Len:  uint32(l),
		Data: data,
		Crc:  0x0000,
	}
	fmt.Printf("%v\n", f.Serial())

	FrameData := f.Serial()
	f.Crc = Crc_16_x25(FrameData[2 : len(FrameData)-2])
	fmt.Printf("C4总帧：%x\n", f.Serial())
	return f.Serial()
}

func Crc_16_x25(data []byte) uint16 {
	var h uint16
	h = 0xffff
	for _, v := range data {
		h ^= uint16(v)
		for j := 0; j < 8; j++ {
			lsb := h & 0x0001
			h >>= 1
			if lsb == 1 {
				h ^= 0x8408
			}
		}
	}
	h ^= 0xffff
	return h
}

func (vc *VirtualAntennaClient) send(conn net.Conn, data []byte) error {
	n, err := conn.Write(data)
	if err != nil {
		vc.conn = nil
		return fmt.Errorf("send data failed: %v", err)
	}

	// fmt.Printf("send data: %x\n", data)
	fmt.Printf("send data bytes: %v\n", n)

	return nil
}

//control 0是关 1开
func (vc *VirtualAntennaClient) changeP30Status(conn net.Conn, control int) error {
	switch control {
	case 0:
		c4 := C4Frame(control)
		err := vc.send(conn, c4)
		return err
	case 1:
		c0 := C0Frame()
		err := vc.send(conn, c0)
		return err
	}
	return nil
}

// pushMessage push message to upper
/*func (vc *VirtualAntennaClient)pushMessage(chanid, name string, value interface{}) error {
	// pack payload
	payload := public.MessagePayload{
		MonitoringUnitID: vc.muid,
		SampleUnitID:     vc.suid,
		ChannelID:        chanid,
		Name:             name,
		Value:            value,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	topic := "sample-values/" + vc.muid + "/" + vc.suid + "/" + chanid
	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}*/

// ID client id
func (vc *VirtualAntennaClient) ID() string {
	return vc.ClientID
}

// Sample get values
func (vc *VirtualAntennaClient) Sample(payload string) (string, error) {
	p, err := DecodeVirtualAntennaOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	if vc.conn == nil {
		return "", fmt.Errorf("device disconnect")
	}

	switch p.Channel {
	case "status":
		if len(vc.B9Data.AntennaInfoNList) == 0 {
			return "", fmt.Errorf("vc.B9Data.AntennaInfoNList is nil")
		}
		s := strconv.Itoa(int(vc.B9Data.AntennaInfoNList[0].Control))
		return s, nil
	default:
		return "invalite channel id", nil
	}

	return "", nil
}

// Command set values
func (vc *VirtualAntennaClient) Command(payload string) (string, error) {
	req, err := DecodeVirtualAntennaOperatePayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode VirtualAntenna operation payload failed, errmsg {%v}", err)
	}
	switch req.Perpose {

	//根据value值 开关p30
	case "changeP30Status":
		if err := vc.changeP30Status(vc.conn, req.Value); err != nil {
			return "", fmt.Errorf("transmission failed: %v", err)
		}
		// return direct
		return "ok", nil
	}

	return "", nil
}
