/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/12/25
 * Despcription: weigeng entrance guard protocol
 *
 */

package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolWeiGengEntry, generalWeiGengEntryDriverProtocol)
}

// Implement DriverProtocol
type weigengEntryDriverProtocol struct {
	req *public.WeiGengEntryBindingPayload
	uri string
}

func (dp *weigengEntryDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *weigengEntryDriverProtocol) ClientID() string {
	return WeiGengEntryClientID
}

func (dp *weigengEntryDriverProtocol) NewInstance() (PortClient, error) {
	return NewWeiGengEntryClient(
		dp.req.LocalAddress, dp.req.LocalPort,
		dp.req.DoorAddress, dp.req.DoorPort,
		dp.req.SerialNumber,
	)
}

func generalWeiGengEntryDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeWeiGengEntryBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &weigengEntryDriverProtocol{
		req: &req,
	}, nil
}

// ============= register driver end ==========================

// WeiGengEntryClientID id
var WeiGengEntryClientID = "wei-geng-entrance-gurad-id"

const (
	frameBytes = 64

	resultFail    = 0
	resultSuccess = 1

	strFunctionIDStatus     = "0x20"
	strFunctionIDDoorTime   = "0x32"
	strFunctionIDSyncTime   = "0x30"
	strFunctionIDOpenDoor   = "0x40"
	strFunctionIDAddCard    = "0x50"
	strFunctionIDRemoveCard = "0x52"

	strSequenceNumberCardRecord = "0xB3"
	strSequenceNumberStatus     = "0xC1"

	byteFunctionIDStatus     = 0x20
	byteFunctionIDDoorTime   = 0x32
	byteFunctionIDSyncTime   = 0x30
	byteFunctionIDOpenDoor   = 0x40
	byteFunctionIDAddCard    = 0x50
	byteFunctionIDRemoveCard = 0x52
)

type weigengData struct {
	recordIndex   uint32
	doorTime      string
	swipeCardTime string
	cardNumber    uint32
	result        byte
	doorNumber    byte
	door1status   byte
	door2status   byte
	door3status   byte
	door4status   byte
}

type weigengSet struct {
	// for remote open door
	doorNumber int

	// for add card and remove card
	cardid    string
	starttime string
	endtime   string
	doorright string
}

type weigengResponse struct {
	result int
	reason string
}

// WeiGengEntryClient client
type WeiGengEntryClient struct {
	ClientID string

	src   string
	dest  string
	devsn string

	conn     *net.UDPConn
	destAddr *net.UDPAddr

	cSet         chan int
	cSetComplete chan int
	setData      weigengSet
	response     weigengResponse

	cache weigengData
}

// NewWeiGengEntryClient new client
func NewWeiGengEntryClient(src, srcport, dest, destport, devsn string) (PortClient, error) {
	var client WeiGengEntryClient

	client.ClientID = WeiGengEntryClientID
	client.src = src + ":" + srcport
	client.dest = dest + ":" + destport
	client.devsn = devsn

	client.cSet = make(chan int)
	client.cSetComplete = make(chan int)

	if err := client.Start(); err != nil {
		return nil, err
	}

	go client.poll()

	return &client, nil
}

// Start start
func (wc *WeiGengEntryClient) Start() error {
	addr, err := net.ResolveUDPAddr("udp", wc.src)
	if err != nil {
		return err
	}

	s, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	destaddr, err := net.ResolveUDPAddr("udp", wc.dest)
	if err != nil {
		return err
	}

	log.Printf("start listen at: %v\n", wc.src)

	wc.conn = s
	wc.destAddr = destaddr

	return nil
}

// DecodeWeiGengEntryBindingPayload decode binding payload
func DecodeWeiGengEntryBindingPayload(payload string) (public.WeiGengEntryBindingPayload, error) {
	var p public.WeiGengEntryBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DecodeWeiGengEntryOperationPayload decode operation payload
func DecodeWeiGengEntryOperationPayload(payload string) (public.WeiGengEntryOperationPayload, error) {
	var p public.WeiGengEntryOperationPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

type cardRecord struct {
	Time       string `json:"time"`
	CardNumber uint32 `json:"cardNo"`
	Result     byte   `json:"result"`
	DoorNumber byte   `json:"door"`
}

// poll
func (wc *WeiGengEntryClient) poll() {
	for {
		select {
		case code := <-wc.cSet:
			// send set data
			var err error

			switch code {
			case byteFunctionIDSyncTime:
				err = wc.setDoorTime()
			case byteFunctionIDOpenDoor:
				err = wc.remoteOpenDoor(wc.setData.doorNumber)
			case byteFunctionIDAddCard:
				{
					cardid := wc.setData.cardid
					starttime := wc.setData.starttime
					endtime := wc.setData.endtime
					right := wc.setData.doorright
					err = wc.addCard(cardid, starttime, endtime, right)
				}
			case byteFunctionIDRemoveCard:
				err = wc.removeCard(wc.setData.cardid)
			}

			if err != nil {
				wc.response.result = resultFail
				wc.response.reason = err.Error()
				wc.cSetComplete <- 1
				break
			}

			wc.response.result = resultSuccess
			wc.cSetComplete <- 1

		default:
			// read for status
			if err := wc.requestForStatus(); err != nil {
				fmt.Println(err)
			}

			// read for door time
			if err := wc.requestForDoorTime(); err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(time.Second)
	}
}

// request for status
func (wc *WeiGengEntryClient) requestForStatus() error {
	t := byte(0x17)
	d := make([]byte, 32)
	f, err := wc.packFrame(t, byteFunctionIDStatus, d)
	if err != nil {
		return fmt.Errorf("pack status frame fail: %v", err)
	}

	resp, err := wc.request(f)
	if err != nil {
		return fmt.Errorf("request for status failed: %v", err)
	}

	data, err := wc.parse(resp)
	if err != nil {
		return fmt.Errorf("parse status response failed: %v", err)
	}

	wc.cache.recordIndex = data.recordIndex
	wc.cache.result = data.result
	wc.cache.doorNumber = data.doorNumber
	wc.cache.cardNumber = data.cardNumber
	wc.cache.swipeCardTime = data.swipeCardTime
	wc.cache.door1status = data.door1status
	wc.cache.door2status = data.door2status
	wc.cache.door3status = data.door3status
	wc.cache.door4status = data.door4status

	return nil
}

// request for door time
func (wc *WeiGengEntryClient) requestForDoorTime() error {
	t := byte(0x17)
	d := make([]byte, 32)
	f, err := wc.packFrame(t, byteFunctionIDDoorTime, d)
	if err != nil {
		return fmt.Errorf("pack door time frame fail: %v", err)
	}

	resp, err := wc.request(f)
	if err != nil {
		return fmt.Errorf("request for door time failed: %v", err)
	}

	data, err := wc.parse(resp)
	if err != nil {
		return fmt.Errorf("parse door time response failed: %v", err)
	}

	wc.cache.doorTime = data.doorTime

	return nil
}

// set door time
func (wc *WeiGengEntryClient) setDoorTime() error {
	t := byte(0x17)
	d := make([]byte, 32)

	now := time.Now()
	ny := now.Year()
	nmo := int(now.Month())
	nd := now.Day()
	nh := now.Hour()
	nmi := now.Minute()
	ns := now.Second()

	sy := strconv.Itoa(ny)
	iy1, _ := strconv.ParseInt(sy[:2], 16, 32)
	iy2, _ := strconv.ParseInt(sy[2:], 16, 32)
	imo, _ := strconv.ParseInt(strconv.Itoa(nmo), 16, 32)
	id, _ := strconv.ParseInt(strconv.Itoa(nd), 16, 32)
	ih, _ := strconv.ParseInt(strconv.Itoa(nh), 16, 32)
	imi, _ := strconv.ParseInt(strconv.Itoa(nmi), 16, 32)
	is, _ := strconv.ParseInt(strconv.Itoa(ns), 16, 32)

	tmp := append([]byte{}, byte(iy1), byte(iy2), byte(imo), byte(id), byte(ih), byte(imi), byte(is))
	copy(d, tmp)

	f, err := wc.packFrame(t, byteFunctionIDSyncTime, d)
	if err != nil {
		return fmt.Errorf("pack door time frame fail: %v", err)
	}

	resp, err := wc.request(f)
	if err != nil {
		return fmt.Errorf("set door time failed: %v", err)
	}

	data, err := wc.parse(resp)
	if err != nil {
		return fmt.Errorf("parse set door time response failed: %v", err)
	}

	ot := fmt.Sprintf("%x%x-%x-%x %x:%x:%x", iy1, iy2, imo, id, ih, imi, is)
	if ot != data.doorTime {
		fmt.Printf("set: %v; response: %v", ot, data.doorTime)
		return fmt.Errorf("set door time failed")
	}

	return nil
}

// remote open door
func (wc *WeiGengEntryClient) remoteOpenDoor(doornumber int) error {
	t := byte(0x17)
	d := make([]byte, 32)

	d[0] = byte(doornumber)

	f, err := wc.packFrame(t, byteFunctionIDOpenDoor, d)
	if err != nil {
		return fmt.Errorf("pack door time frame fail: %v", err)
	}

	resp, err := wc.request(f)
	if err != nil {
		return fmt.Errorf("request for door time failed: %v", err)
	}

	data, err := wc.parse(resp)
	if err != nil {
		return fmt.Errorf("parse door time response failed: %v", err)
	}

	if data.result != resultSuccess {
		return fmt.Errorf("remote open door [%v] failed", doornumber)
	}

	return nil
}

type doorRight struct {
	Door1 int `json:"1"`
	Door2 int `json:"2"`
	Door3 int `json:"3"`
	Door4 int `json:"4"`
}

// add card, start time and end time like '20181228', door right like '{"1":1,"2":1,"3":0,"4":0}'
func (wc *WeiGengEntryClient) addCard(cardid, starttime, endtime, doorright string) error {
	t := byte(0x17)
	d := make([]byte, 32)

	// card id
	icid, err := strconv.Atoi(cardid)
	if err != nil {
		return fmt.Errorf("card id illegal: %v", err)
	}

	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, uint32(icid))

	copy(d, tmp)

	// start time
	st, err := time.Parse("20060102", starttime)
	if err != nil {
		return fmt.Errorf("start time illegal: %v", err)
	}

	ny := st.Year()
	nm := int(st.Month())
	nd := st.Day()

	sy := strconv.Itoa(ny)
	iy1, _ := strconv.ParseInt(sy[:2], 16, 32)
	iy2, _ := strconv.ParseInt(sy[2:], 16, 32)
	im, _ := strconv.ParseInt(strconv.Itoa(nm), 16, 32)
	id, _ := strconv.ParseInt(strconv.Itoa(nd), 16, 32)

	tmp = append([]byte{}, byte(iy1), byte(iy2), byte(im), byte(id))
	copy(d[4:], tmp)

	// end time
	et, err := time.Parse("20060102", endtime)
	if err != nil {
		return fmt.Errorf("end time illegal: %v", err)
	}

	ny = et.Year()
	nm = int(et.Month())
	nd = et.Day()

	sy = strconv.Itoa(ny)
	iy1, _ = strconv.ParseInt(sy[:2], 16, 32)
	iy2, _ = strconv.ParseInt(sy[2:], 16, 32)
	im, _ = strconv.ParseInt(strconv.Itoa(nm), 16, 32)
	id, _ = strconv.ParseInt(strconv.Itoa(nd), 16, 32)

	tmp = append([]byte{}, byte(iy1), byte(iy2), byte(im), byte(id))
	copy(d[8:], tmp)

	// door right
	var right doorRight
	if err := json.Unmarshal([]byte(doorright), &right); err != nil {
		return fmt.Errorf("door right illegal: %v", doorright)
	}

	d[12] = byte(right.Door1)
	d[13] = byte(right.Door2)
	d[14] = byte(right.Door3)
	d[15] = byte(right.Door4)

	f, err := wc.packFrame(t, byteFunctionIDAddCard, d)
	if err != nil {
		return fmt.Errorf("pack door time frame fail: %v", err)
	}

	resp, err := wc.request(f)
	if err != nil {
		return fmt.Errorf("request for door time failed: %v", err)
	}

	data, err := wc.parse(resp)
	if err != nil {
		return fmt.Errorf("parse door time response failed: %v", err)
	}

	if data.result != resultSuccess {
		return fmt.Errorf("add card [%v] failed", cardid)
	}

	return nil
}

// remove card
func (wc *WeiGengEntryClient) removeCard(cardid string) error {
	t := byte(0x17)
	d := make([]byte, 32)

	// card id
	icid, err := strconv.Atoi(cardid)
	if err != nil {
		return fmt.Errorf("card id illegal: %v", err)
	}

	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, uint32(icid))

	copy(d, tmp)

	f, err := wc.packFrame(t, byteFunctionIDRemoveCard, d)
	if err != nil {
		return fmt.Errorf("pack door time frame fail: %v", err)
	}

	resp, err := wc.request(f)
	if err != nil {
		return fmt.Errorf("request for door time failed: %v", err)
	}

	data, err := wc.parse(resp)
	if err != nil {
		return fmt.Errorf("parse door time response failed: %v", err)
	}

	if data.result != resultSuccess {
		return fmt.Errorf("add card [%v] failed", cardid)
	}

	return nil
}

// weigeng frame has constant length which is 64 bytes
type weigengFrame struct {
	Type       byte
	FunctionID byte
	Reserved   int16
	DeviceSN   int32
	Data       [32]byte
	SequenceID int32
	ExternData [20]byte
}

func (wc *WeiGengEntryClient) packFrame(t, fid byte, data []byte) ([]byte, error) {
	var frame []byte

	frame = append(frame, t, fid, 0x00, 0x00)

	// serail number to bytes
	sn, err := strconv.Atoi(wc.devsn)
	if err != nil {
		return nil, err
	}

	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, uint32(sn))

	frame = append(frame, tmp...)
	frame = append(frame, data...)

	// sequence id use zero default
	frame = append(frame, 0x00, 0x00, 0x00, 0x00)

	// extern data use zero default
	ed := make([]byte, 20)
	frame = append(frame, ed...)

	fmt.Printf("%x\n", frame)

	return frame, nil
}

func (wc *WeiGengEntryClient) request(data []byte) ([]byte, error) {
	wc.conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := wc.conn.WriteToUDP(data, wc.destAddr); err != nil {
		return nil, err
	}

	var resp = make([]byte, 1024)
	wc.conn.SetDeadline(time.Now().Add(time.Second))
	n, _, err := wc.conn.ReadFromUDP(resp)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%x\n", resp[:n])

	return resp[:n], nil
}

func (wc *WeiGengEntryClient) parse(data []byte) (weigengData, error) {
	var d weigengData

	if len(data) != frameBytes {
		return d, fmt.Errorf("frame bytes abnormal")
	}

	// get function id
	fid := data[1]
	switch fid {
	case byteFunctionIDStatus:
		{
			// index
			bi := data[8:12]
			d.recordIndex = binary.LittleEndian.Uint32(bi)

			// result
			d.result = data[13]

			// door number
			d.doorNumber = data[14]

			// card number
			bcn := data[16:20]
			d.cardNumber = binary.LittleEndian.Uint32(bcn)

			// swipe card time
			d.swipeCardTime = fmt.Sprintf("%x%x-%x-%x %x:%x:%x", data[20], data[21], data[22], data[23], data[24], data[25], data[26])

			// status
			d.door1status = data[28]
			d.door2status = data[29]
			d.door3status = data[30]
			d.door4status = data[31]
		}
	case byteFunctionIDDoorTime, byteFunctionIDSyncTime:
		{
			// door time
			d.doorTime = fmt.Sprintf("%x%x-%x-%x %x:%x:%x", data[8], data[9], data[10], data[11], data[12], data[13], data[14])
		}
	case byteFunctionIDOpenDoor, byteFunctionIDAddCard, byteFunctionIDRemoveCard:
		{
			// result
			d.result = data[8]
		}
	}

	return d, nil
}

// ID id
func (wc *WeiGengEntryClient) ID() string {
	return wc.ClientID
}

// Sample sample
func (wc *WeiGengEntryClient) Sample(payload string) (string, error) {
	p, err := DecodeWeiGengEntryOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	switch p.FunctionID {
	case strFunctionIDStatus:
		{
			if p.SequenceNumber == strSequenceNumberStatus {
				switch p.Group {
				case 1:
					return strconv.Itoa(int(wc.cache.door1status)), nil
				case 2:
					return strconv.Itoa(int(wc.cache.door2status)), nil
				case 3:
					return strconv.Itoa(int(wc.cache.door3status)), nil
				case 4:
					return strconv.Itoa(int(wc.cache.door4status)), nil
				default:
					return "", fmt.Errorf("unknown group [%v]", p.Group)
				}
			} else if p.SequenceNumber == strSequenceNumberCardRecord {
				var cr cardRecord
				cr.Time = wc.cache.swipeCardTime
				cr.CardNumber = wc.cache.cardNumber
				cr.Result = wc.cache.result
				cr.DoorNumber = wc.cache.doorNumber

				bcr, err := json.Marshal(cr)
				if err != nil {
					return "", fmt.Errorf("marshal card record failed: %v", err)
				}

				return string(bcr), nil
			} else {
				return "", fmt.Errorf("unknown sequence number [%v]", p.SequenceNumber)
			}
		}
	case strFunctionIDDoorTime:
		{
			return wc.cache.doorTime, nil
		}
	default:
		return "", fmt.Errorf("unknown function id [%v]", p.FunctionID)
	}
}

// Command command
func (wc *WeiGengEntryClient) Command(payload string) (string, error) {
	p, err := DecodeWeiGengEntryOperationPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed: %v", err)
	}

	switch p.FunctionID {
	case strFunctionIDAddCard:
		{
			wc.setData.cardid = p.CardID
			wc.setData.starttime = time.Now().Format("20060102")
			wc.setData.endtime = p.ExpireDate
			wc.setData.doorright = p.DoorRight

			wc.cSet <- byteFunctionIDAddCard
		}
	case strFunctionIDRemoveCard:
		{
			wc.setData.cardid = p.CardID

			wc.cSet <- byteFunctionIDRemoveCard
		}
	case strFunctionIDOpenDoor:
		{
			wc.setData.doorNumber = p.Door

			wc.cSet <- byteFunctionIDOpenDoor
		}
	case strFunctionIDSyncTime:
		{
			wc.cSet <- byteFunctionIDSyncTime
		}
	default:
		return "", fmt.Errorf("unknown function id [%v]", p.FunctionID)
	}

	<-wc.cSetComplete

	if wc.response.result != resultSuccess {
		return "", fmt.Errorf(wc.response.reason)
	}

	return "ok", nil
}
