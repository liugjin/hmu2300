/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/09/11
 * Despcription: sensorflow implement
 *
 */

package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"clc.hmu/app/public/log/portlog"
	"clc.hmu/app/portmanager/portnet"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log"

	"github.com/goburrow/modbus"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolSensorflow, generalSensorflowDriverProtocol)
}

// Implement DriverProtocol
type sensorflowDriverProtocol struct {
	req *public.ModbusPayload
	uri string
}

func (dp *sensorflowDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *sensorflowDriverProtocol) ClientID() string {
	return SensorflowClientID
}

func (dp *sensorflowDriverProtocol) NewInstance() (PortClient, error) {
	return NewSensorflowClient(
		dp.uri,
		int(dp.req.BaudRate), int(dp.req.Timeout),
		int(dp.req.KeyNumber), dp.req.MUID,
		dp.req.SUID, dp.req.WANInterface,
		dp.req.WifiInterface,
	)
}

func generalSensorflowDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &sensorflowDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

/**
	NOTE: 1 register including 2 bytes
**/

// const
const (
	MaxUTrakerNumber          = 10   // MaxUTrakerNumber max u wei number
	MaxMessageSize            = 1024 // default bytes of message
	RegisterAddressListNumber = 6
	MaxTransmisionSize        = 256
	IntervalPoll              = 50
	MaxDataSize               = 2560 // cache all data
	SampleTimeout             = 100
	CommandPrepareTimeConsume = 15
	MaxSampleCount            = 5
)

// Tag type
const (
	TagTypeRFID         = 0 // RFID
	TagTypeTH           = 1 // temperature and humudity
	TagTypeDoorMagnetic = 2 // door magnatic
	TagTypeAcoustoOptic = 3 // acousto-optic
	TagTypeSmoke        = 4 // smoke
	TagTypePT100        = 5 // PT100
	TagTypeElectricity  = 6 // 4-20mA
	TagTypeVoltage      = 7 // 0-10V
	TagTypeLampWith     = 8 // lamp with
)

// lamp with display mode
const (
	DisplayCustom  = "custom"
	DisplayMarquee = "marquee"
)

// value
var (
	SensorflowClientID   = "sensorflow-client"                           // SensorflowClientID id
	RegisterAddressIndex = [MaxUTrakerNumber]int{0, 32, 0, 68, 0, 104}   // start address
	RegisterNumberIndex  = [MaxUTrakerNumber]int{32, 36, 32, 36, 32, 36} // register number
)

type color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type colorEx struct {
	UWeiIndex int
	LEDIndex  int
	Red       int
	Green     int
	Blue      int
}

// SensorflowClient sensorflow client
type SensorflowClient struct {
	Clients []modbus.Client

	ClientID string
	mtx      sync.Mutex

	UKeyNumber int // number of keys in one u tracker

	OnlineModuleNumber int                   // online u wei number
	PollNumber         int                   // poll number
	Overtime           [MaxUTrakerNumber]int // over time

	Index int // current index

	TransmissionData   [MaxUTrakerNumber][MaxTransmisionSize]byte // set tag
	TransmissionNumber [MaxUTrakerNumber]int                      // number of set data

	LEDSet      [MaxUTrakerNumber]bool
	LEDPollTest [MaxUTrakerNumber]bool
	LEDIndex    [MaxUTrakerNumber]int // u position, max 6 led
	LEDRed      [MaxUTrakerNumber]int
	LEDGreen    [MaxUTrakerNumber]int
	LEDBlue     [MaxUTrakerNumber]int

	Data [MaxUTrakerNumber][MaxMessageSize]byte // cache data

	RegisterAddress [MaxUTrakerNumber]int // address of query
	RegisterNumber  [MaxUTrakerNumber]int //number of query
	RegisterIndex   [MaxUTrakerNumber]int // index of address

	ModbusData [MaxDataSize]byte // cache all data for query

	muid string
	suid string

	buttonPressLastStatus map[string]uint16

	colorTable         []color
	priorityColorTable []colorEx
	bSetColor          bool
	bSetColorResult    bool

	SampleCount [MaxUTrakerNumber]int
}

// NewSensorflowClient new modbus client
func NewSensorflowClient(port string, baudrate, timeout, keynumber int, muid, suid, wanift, wifiift string) (*SensorflowClient, error) {
	clients := []modbus.Client{}
	for i := 0; i < MaxUTrakerNumber; i++ {
		// new handler
		handler := modbus.NewRTUClientHandler(port)
		handler.BaudRate = baudrate
		handler.DataBits = 8
		handler.Parity = "N"
		handler.StopBits = 1
		handler.SlaveId = byte(i + 1)
		handler.Timeout = time.Millisecond * SampleTimeout
		// handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)

		if err := handler.Connect(); err != nil {
			return &SensorflowClient{}, fmt.Errorf("connect port[%s] failed, baudrate[%v], slaveid[%v], errmsg[%v]", port, baudrate, i, err)
		}

		// portlog.LOG.Infof("connect port[%s] success,  baudrate[%v], slaveid[%v]", port, baudrate, slaveid)
		clients = append(clients, modbus.NewClient(handler))
	}

	// check key number
	if keynumber <= 0 {
		keynumber = 6
	}

	// defer handler.Close()
	var client = SensorflowClient{
		Clients:    clients,
		ClientID:   SensorflowClientID,
		UKeyNumber: keynumber,
		muid:       muid,
		suid:       suid,
	}

	client.buttonPressLastStatus = make(map[string]uint16)
	client.bSetColor = false
	client.bSetColorResult = false

	// start poll
	client.Init()
	go client.Start()

	// post network param
	if err := client.PostNetworkInfo(wanift, wifiift); err != nil {
		fmt.Printf("post network info failed: %v\n", err)
		go func() {
			for {
				if err := client.PostNetworkInfo(wanift, wifiift); err != nil {
					time.Sleep(time.Second)
					continue
				}

				fmt.Printf("post network info success\n")
				break
			}
		}()
	}

	// read color data
	filename := "data" + muid + ".json"
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			var data []color

			var d color
			d.Red = 1
			d.Green = 0
			d.Blue = 0

			for i := 0; i < 48; i++ {
				data = append(data, d)
			}

			client.colorTable = data

			// write data
			bytedata, _ := json.MarshalIndent(data, "", "\t")
			ioutil.WriteFile(filename, bytedata, 0666)
		}
	} else {
		// file exist
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("read color table failed: %v\n", err)
			portlog.LOG.Warning("read color table failed: %v\n", err)
		} else {
			if err := json.Unmarshal(data, &client.colorTable); err != nil {
				fmt.Printf("unmarshal data failed: %v\n", err)
				portlog.LOG.Warning("unmarshal data failed: %v\n", err)
			}
		}
	}

	portlog.LOG.Info(client.colorTable)

	go client.checkColor()

	return &client, nil
}

func (sc *SensorflowClient) checkColor() {
	// sleep for a moment when at first
	time.Sleep(time.Second * CommandPrepareTimeConsume)

	for {
		l := len(sc.priorityColorTable)
		if l > 0 {
			i := sc.priorityColorTable[0].UWeiIndex
			j := sc.priorityColorTable[0].LEDIndex
			r := sc.priorityColorTable[0].Red
			g := sc.priorityColorTable[0].Green
			b := sc.priorityColorTable[0].Blue

			log.Println("set:", r, g, b)

			// resetRGBValue(&r, &g, &b)
			sc.SetUTrackerLED(i, j, r, g, b)

			for !sc.bSetColor {
				time.Sleep(time.Millisecond * 50)
			}

			log.Println("set color return")
			sc.bSetColor = false

			if !sc.bSetColorResult {
				log.Println("set color fail")
				continue
			}

			log.Println("set ok")

			if l > 1 {
				sc.priorityColorTable = sc.priorityColorTable[1:]
			} else {
				sc.priorityColorTable = []colorEx{}
			}

			continue
		}

		// query for one tracker
		for i := 0; i < sc.OnlineModuleNumber; i++ {
			// rgb start at 7th byte
			// query for every key, get color value
			for j := 0; j < sc.UKeyNumber; j++ {
				begin := (6 + j*3) * 2
				br := sc.Data[i][begin : begin+2]
				bg := sc.Data[i][begin+2 : begin+4]
				bb := sc.Data[i][begin+4 : begin+6]

				r := int(binary.BigEndian.Uint16(br))
				g := int(binary.BigEndian.Uint16(bg))
				b := int(binary.BigEndian.Uint16(bb))

				cindex := i*sc.UKeyNumber + j
				if len(sc.colorTable) < cindex {
					break
				}

				if (sc.colorTable[cindex].Red != r) || (sc.colorTable[cindex].Green != g) || (sc.colorTable[cindex].Blue != b) {
					log.Println("index:", i, j, "read:", r, g, b)

					r = sc.colorTable[cindex].Red
					g = sc.colorTable[cindex].Green
					b = sc.colorTable[cindex].Blue

					log.Println("set:", r, g, b)

					// resetRGBValue(&r, &g, &b)
					sc.SetUTrackerLED(i, j, r, g, b)

					for !sc.bSetColor {
						time.Sleep(time.Millisecond * 50)
					}

					log.Println("set color return")
					sc.bSetColor = false

					if !sc.bSetColorResult {
						log.Println("set fail")
						continue
					}

					log.Println("set ok")
				}

				// time.Sleep(time.Millisecond * 500)
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}

type networkPayload struct {
	WAN  string `json:"wan"`
	WIFI string `json:"wifi"`
	MAC  string `json:"mac"`
}

// PostNetworkInfo post network info
func (sc *SensorflowClient) PostNetworkInfo(wanname, wifiname string) error {
	wanip, wanmac, _ := public.QueryInterfaceInfoByName(wanname)
	wifiip, _, _ := public.QueryInterfaceInfoByName(wifiname)

	var np networkPayload
	np.WAN = wanip
	np.WIFI = wifiip
	np.MAC = wanmac

	topic := "sample-values/" + sc.muid + "/_/network"

	payload := public.MessagePayload{
		MonitoringUnitID: sc.muid,
		SampleUnitID:     "_",
		ChannelID:        "network",
		Name:             "network",
		Value:            np,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}

// PostButtonPressInfo post button press info
func (sc *SensorflowClient) PostButtonPressInfo(suid string) error {
	topic := "sample-values/" + sc.muid + "/" + suid + "/led"

	payload := public.MessagePayload{
		MonitoringUnitID: sc.muid,
		SampleUnitID:     suid,
		ChannelID:        "led",
		Name:             "led",
		Value:            "0-1000-0",
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}

// PostCurrentLedStatus post button press info
func (sc *SensorflowClient) PostCurrentLedStatus(suid string, uindex, index int) error {
	topic := "sample-values/" + sc.muid + "/" + suid + "/led"

	address := 6 + 3*index
	start := address * 2
	br := sc.Data[uindex][start : start+2]
	bg := sc.Data[uindex][start+2 : start+4]
	bb := sc.Data[uindex][start+4 : start+6]

	r := binary.BigEndian.Uint16(br)
	g := binary.BigEndian.Uint16(bg)
	b := binary.BigEndian.Uint16(bb)

	const remote = 0x8000
	if r >= remote {
		r = r - remote
	}

	if g >= remote {
		g = g - remote
	}

	if b >= remote {
		b = b - remote
	}

	value := strconv.Itoa(int(r)) + "-" + strconv.Itoa(int(g)) + "-" + strconv.Itoa(int(b))

	payload := public.MessagePayload{
		MonitoringUnitID: sc.muid,
		SampleUnitID:     suid,
		ChannelID:        "led",
		Name:             "led",
		Value:            value,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}

// PostCurrentModules post current modules
func (sc *SensorflowClient) PostCurrentModules() error {
	topic := "sample-values/" + sc.muid + "/_/modules"

	payload := public.MessagePayload{
		MonitoringUnitID: sc.muid,
		SampleUnitID:     "_",
		ChannelID:        "modules",
		Name:             "modules",
		Value:            sc.OnlineModuleNumber,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}

// PostCurrentUCount post current u count
func (sc *SensorflowClient) PostCurrentUCount() error {
	topic := "sample-values/" + sc.muid + "/_/ucount"

	payload := public.MessagePayload{
		MonitoringUnitID: sc.muid,
		SampleUnitID:     "_",
		ChannelID:        "ucount",
		Name:             "ucount",
		Value:            sc.OnlineModuleNumber * sc.UKeyNumber,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	p, _ := json.Marshal(payload)

	return portnet.DefaultNotify(topic, string(p))
}

// Init init
func (sc *SensorflowClient) Init() {
	for i := 0; i < MaxUTrakerNumber; i++ {
		sc.Overtime[i] = 2
	}

	// tag overtime counts
	for i := 0; i < MaxUTrakerNumber; i++ {
		sc.Data[i][49] = 255
		sc.Data[i][51] = 255
		sc.Data[i][53] = 255
		sc.Data[i][55] = 255
		sc.Data[i][57] = 255
		sc.Data[i][59] = 255
	}

	for i := 0; i < MaxUTrakerNumber; i++ {
		sc.RegisterNumber[i] = 32
	}

	sc.PollNumber = 2
}

// Poll poll
func (sc *SensorflowClient) Poll() {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()

	n := 0
	for i := 0; i < MaxUTrakerNumber; i++ {
		if sc.Overtime[i] < 2 {
			n++
		} else {
			break
		}
	}

	sc.OnlineModuleNumber = n

	if n >= 0 && n < MaxUTrakerNumber {
		sc.PollNumber = n + 1
	} else if n >= MaxUTrakerNumber {
		sc.PollNumber = MaxUTrakerNumber
	}

	// portlog.LOG.Info("online: %d, poll: %d, index: %d,", n, sc.PollNumber, sc.Index)

	var r []byte
	var err error

	if sc.LEDSet[sc.Index] && sc.SampleCount[sc.Index] >= MaxSampleCount {
		log.Println("set led, sample count:", sc.SampleCount[sc.Index])
		sc.SampleCount[sc.Index] = MaxSampleCount - 1

		// set led
		sc.LEDSet[sc.Index] = false

		addr := uint16(6 + sc.LEDIndex[sc.Index]*3)
		quan := uint16(3)

		val := []byte{}
		tmp := make([]byte, 2)

		binary.BigEndian.PutUint16(tmp, uint16(sc.LEDRed[sc.Index]))
		val = append(val, tmp...)
		binary.BigEndian.PutUint16(tmp, uint16(sc.LEDGreen[sc.Index]))
		val = append(val, tmp...)
		binary.BigEndian.PutUint16(tmp, uint16(sc.LEDBlue[sc.Index]))
		val = append(val, tmp...)

		r, err = sc.Clients[sc.Index].WriteMultipleRegisters(addr, quan, val)
		if err != nil {
			sc.SampleCount[sc.Index] = 0
			sc.bSetColorResult = false
		} else {
			sc.bSetColorResult = true
		}

		sc.bSetColor = true
	} else if sc.TransmissionData[sc.Index][0] > 0 {
		log.Println("set label")

		// set label
		sc.TransmissionData[sc.Index][0] = 0

		addr := uint16(220)
		quan := uint16(sc.TransmissionNumber[sc.Index])
		val := sc.TransmissionData[sc.Index][1 : quan*2+1]

		log.Println(val, addr, quan)
		r, err = sc.Clients[sc.Index].WriteMultipleRegisters(addr, quan, val)
	} else if sc.LEDPollTest[sc.Index] {
		log.Println("set all led")

		// set all led
		sc.LEDPollTest[sc.Index] = false

		addr := uint16(6)
		quan := uint16(3 * 6)
		val := []byte{}

		for i := 0; i < 6; i++ {
			val = append(val, 0x80, 0x00, 0x80, 0x00, 0x80, 0x64)
		}

		r, err = sc.Clients[sc.Index].WriteMultipleRegisters(addr, quan, val)
	} else {
		// get values
		sc.RegisterAddress[sc.Index] = RegisterAddressIndex[sc.RegisterIndex[sc.Index]]
		sc.RegisterNumber[sc.Index] = RegisterNumberIndex[sc.RegisterIndex[sc.Index]]

		addr := sc.RegisterAddress[sc.Index]
		quan := sc.RegisterNumber[sc.Index]

		log.Printf("get val, sc.Index: %d, register index: %d, addr: %d, quan: %d\n", sc.Index, sc.RegisterIndex[sc.Index], addr, quan)

		sc.RegisterIndex[sc.Index]++
		if sc.RegisterIndex[sc.Index] >= RegisterAddressListNumber {
			sc.RegisterIndex[sc.Index] = 0
		}

		r, err = sc.Clients[sc.Index].ReadHoldingRegisters(uint16(addr), uint16(quan))
		if err == nil {
			begin := addr * 2
			size := len(r)

			for i, j := begin, 0; j < size; i, j = i+1, j+1 {
				sc.Data[sc.Index][i] = r[j]
			}

			log.Println(begin, begin+size, sc.Index, sc.Data[sc.Index][:280])

			// count when get success
			if sc.SampleCount[sc.Index] < MaxSampleCount {
				sc.SampleCount[sc.Index]++
			}
		}
	}

	if err != nil {
		// do not increase continuously when reach limitation
		if sc.Overtime[sc.Index] < 100 {
			sc.Overtime[sc.Index]++
		}

		// portlog.LOG.Info("read error: %v, overtime: %v", err, sc.Overtime[sc.Index], sc.Index)
	} else {
		sc.Overtime[sc.Index] = 0
	}

	sc.Index++
	if sc.Index >= sc.PollNumber {
		sc.Index = 0
	}
}

// SetUTrackerLED set u tracker led, index: which u tracker, start at 0; id: which led in u tracker, start at 0
func (sc *SensorflowClient) SetUTrackerLED(index, id, r, g, b int) {
	sc.LEDSet[index] = true
	sc.LEDIndex[index] = id
	sc.LEDRed[index] = r
	sc.LEDGreen[index] = g
	sc.LEDBlue[index] = b
}

// SetTagLED set tag led, id start at 1, max 6
func (sc *SensorflowClient) SetTagLED(index, id, r, g, b int) {
	sc.TransmissionData[index][0] = 0xAA

	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(id))

	// id
	sc.TransmissionData[index][1] = tmp[0]
	sc.TransmissionData[index][2] = tmp[1]

	// address, r, g, b in 9, 10, 11
	sc.TransmissionData[index][3] = 0x00
	sc.TransmissionData[index][4] = 0x09

	// quantity, r, g, b has 6 bytes, 3 registers
	sc.TransmissionData[index][5] = 0x00
	sc.TransmissionData[index][6] = 0x03

	binary.BigEndian.PutUint16(tmp, uint16(r))
	sc.TransmissionData[index][7] = tmp[0]
	sc.TransmissionData[index][8] = tmp[1]

	binary.BigEndian.PutUint16(tmp, uint16(g))
	sc.TransmissionData[index][9] = tmp[0]
	sc.TransmissionData[index][10] = tmp[1]

	binary.BigEndian.PutUint16(tmp, uint16(b))
	sc.TransmissionData[index][11] = tmp[0]
	sc.TransmissionData[index][12] = tmp[1]

	// number: id(1) + address(1) + quantity(1) + rgb(3) = 6
	sc.TransmissionNumber[index] = 6
}

// SetTagID set tag id, param id represent which tag, value represent the tag id that want to be set
func (sc *SensorflowClient) SetTagID(index, id int, value []byte) {
	sc.TransmissionData[index][0] = 0xAA

	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(id))

	// id
	sc.TransmissionData[index][1] = tmp[0]
	sc.TransmissionData[index][2] = tmp[1]

	// address, start at 17
	sc.TransmissionData[index][3] = 0x00
	sc.TransmissionData[index][4] = 0x11

	// quantity, id has 13 bytes
	sc.TransmissionData[index][5] = 0x00
	sc.TransmissionData[index][6] = 0x07

	vlen := len(value)
	for i := 0; i < vlen; i++ {
		sc.TransmissionData[index][7+i] = value[i]
	}

	// number: id(1) + address(1) + quantity(1) + id(7) = 10
	sc.TransmissionNumber[index] = 10
}

// SetTagOutput set tag output, value use 0 or 100
func (sc *SensorflowClient) SetTagOutput(index, id, value int) {
	sc.TransmissionData[index][0] = 0xAA

	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(id))

	// id
	sc.TransmissionData[index][1] = tmp[0]
	sc.TransmissionData[index][2] = tmp[1]

	// address, start at 17
	sc.TransmissionData[index][3] = 0x00
	sc.TransmissionData[index][4] = 0x0F

	// quantity, id has 13 bytes
	sc.TransmissionData[index][5] = 0x00
	sc.TransmissionData[index][6] = 0x02

	// value
	sc.TransmissionData[index][7] = 0x00
	sc.TransmissionData[index][8] = 0x64
	sc.TransmissionData[index][9] = 0x00
	sc.TransmissionData[index][10] = byte(value)

	// number: id(1) + address(1) + quantity(1) + id(2) = 5
	sc.TransmissionNumber[index] = 5
}

// SetLampWith set lamp with id start at 1, max 6
func (sc *SensorflowClient) SetLampWith(index, id, m, c, t, r, g, b int) {
	sc.TransmissionData[index][0] = 0xAA

	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(id))

	// id
	sc.TransmissionData[index][1] = tmp[0]
	sc.TransmissionData[index][2] = tmp[1]

	// address, start at 9
	sc.TransmissionData[index][3] = 0x00
	sc.TransmissionData[index][4] = 0x09

	// quantity, 3 registers
	sc.TransmissionData[index][5] = 0x00
	sc.TransmissionData[index][6] = 0x03

	// high byte represent mode, low represent amount or position
	sc.TransmissionData[index][7] = byte(m)
	sc.TransmissionData[index][8] = byte(c)

	// high byte represent green, low represent red
	sc.TransmissionData[index][9] = byte(g)
	sc.TransmissionData[index][10] = byte(r)

	// high byte represent blue, low represent interval
	sc.TransmissionData[index][11] = byte(b)
	sc.TransmissionData[index][12] = byte(t)

	// number: id(1) + address(1) + quantity(1) + rgb(3) = 6
	sc.TransmissionNumber[index] = 6
}

// Start poll and receive command
func (sc *SensorflowClient) Start() {
	for {
		sc.Poll()
		sc.Parse()

		sc.PostCurrentModules()
		sc.PostCurrentUCount()

		time.Sleep(time.Millisecond * IntervalPoll)
	}
}

// Parse parse data
func (sc *SensorflowClient) Parse() {
	unum := sc.OnlineModuleNumber

	// set online u tracker, 1st register
	binary.BigEndian.PutUint16(sc.ModbusData[:2], uint16(unum))

	// set online tag (2nd register) and status (3rd to 6th register, total 8 bytes)
	s, ln := sc.TagFlag(unum)
	binary.BigEndian.PutUint16(sc.ModbusData[2:4], ln)
	binary.BigEndian.PutUint64(sc.ModbusData[4:12], s)

	// key status start at 9th, end at 12th
	ks := sc.UkeyStatus(unum)
	binary.BigEndian.PutUint64(sc.ModbusData[16:24], ks)

	// from 17th to 1040th, every u include 16 registers
	sc.SaveOtherData(unum)

	// log.Println(sc.ModbusData[0:64])
	// log.Println(sc.Data[0][:100])
}

// TagFlag get tag flag, tag status and number of lines
func (sc *SensorflowClient) TagFlag(unum int) (uint64, uint16) {
	s := uint64(0)
	ln := uint16(0)

	for i := 0; i < unum; i++ {
		// tag status at 24 ~ 29, mapping U0 ~ U5 status, total 6
		for j := 0; j < sc.UKeyNumber; j++ {
			start := (j + 24) * 2
			end := (j + 25) * 2

			d := sc.Data[i][start:end]
			t := binary.BigEndian.Uint16(d)

			// tag disconnect times, larger than 2 means do not exist
			if t > 2 {
				s = s & (^(1 << uint(i*6+j)))
			} else {
				s = s | (1 << uint(i*6+j))
				ln++
			}
		}
	}

	// transfer data
	tmp := s
	s = ((tmp & 0xFFFF) << 48) + (((tmp >> 16) & 0xFFFF) << 32) + (((tmp >> 32) & 0xFFFF) << 16) + ((tmp >> 48) & 0xFFFF)

	// log.Printf("%X, number: %d", s, ln)
	return s, ln
}

// UkeyStatus u key status
func (sc *SensorflowClient) UkeyStatus(unum int) uint64 {
	s := uint64(0)

	for i := 0; i < unum; i++ {
		// key status at 0 ~ 5, mapping U0 ~ U5 key status, total 6
		for j := 0; j < sc.UKeyNumber; j++ {
			start := j * 2
			end := (j + 1) * 2

			d := sc.Data[i][start:end]
			t := binary.BigEndian.Uint16(d)

			// key press times, larger than 0 means pressing
			if t == 0 {
				s = s & (^(1 << uint(i*6+j)))
			} else {
				s = s | (1 << uint(i*6+j))
			}
		}
	}

	// transfer data
	tmp := s
	s = ((tmp & 0xFFFF) << 48) + (((tmp >> 16) & 0xFFFF) << 32) + (((tmp >> 32) & 0xFFFF) << 16) + ((tmp >> 48) & 0xFFFF)

	// log.Printf("%X, %X", s, tmp)
	return s
}

// SaveOtherData including tag data and u tracker data
/**
in modbus data, arrange mode (total 16 registers, start at 16th register):
# tag id(7) + tag rgb(3) + rt(1) + reserve(1) + uled rgb(3) + ukey val(1) // original definition, reserve
tag type(1(th)): tag id(7) + tag rgb(3) + temperature(1) + humidity(1) + uled rgb(3) + ukey val(1)
tag type(*): tag id(7) + tag rgb(3) + rt(1) + data(1) + uled rgb(3) + ukey val(1)

in u tracker data, arrange mode (total 18 registers, start at 32nd register):
tag id(7) + tag key status(1) + tag vibration(1) + tag rgb(3) + rt(1) + h(1) + t(1) + extend mode(1) + extend data(1) + reserve(1)

u led rgb start at 6th register, end at 23rd register;
key status start at 0th register, end at 5th register;
**/
func (sc *SensorflowClient) SaveOtherData(unum int) {
	for i := 0; i < unum; i++ {
		for j := 0; j < sc.UKeyNumber; j++ {
			mstart := 32 + (i*sc.UKeyNumber+j)*32
			ustart := 64 + j*36

			// tag id, 14 bytes
			for k := 0; k < 14; k++ {
				mi := mstart + k
				ui := ustart + k
				sc.ModbusData[mi] = sc.Data[i][ui]
			}

			// tag rgb, 6 bytes
			for k := 0; k < 6; k++ {
				mi := mstart + 14 + k
				ui := ustart + 18 + k
				sc.ModbusData[mi] = sc.Data[i][ui]
			}

			/** original difinition, reserve
			// resistance temperature, 2 bytes
			// for k := 0; k < 2; k++ {
			// 	mi := mstart + 20 + k
			// 	ui := ustart + 24 + k
			// 	sc.ModbusData[mi] = sc.Data[i][ui]
			// }

			// reserve, 2 bytes, pass
			**/

			// custom difinition, save data by tag type, type at id's last byte
			ti := ustart + 13
			ltype := sc.Data[i][ti]

			switch ltype {
			case TagTypeRFID:
				// save resistance temperature (2 bytes), reserve (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 20 + k
					ui := ustart + 24 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}
			case TagTypeTH:
				// save humidity and temperature, 4 bytes
				for k := 0; k < 4; k++ {
					mi := mstart + 20 + k
					ui := ustart + 26 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}
			case TagTypeDoorMagnetic:
				// save resistance temperature (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 20 + k
					ui := ustart + 24 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}

				// save sensor data (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 22 + k
					ui := ustart + 32 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}
			case TagTypeAcoustoOptic:
				// TODO
			case TagTypeSmoke:
				// TODO
			case TagTypePT100:
				// save resistance temperature (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 20 + k
					ui := ustart + 24 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}

				// compute pt100
				di := ustart + 32
				d := uint(sc.Data[i][di])
				fA := float64(39)
				fI := 0.00125
				fR2 := 78.7
				fd := public.PT100(d, fA, fI, fR2)
				md := uint16(fd * 100)

				tmp := make([]byte, 2)
				binary.BigEndian.PutUint16(tmp, md)

				// save sensor data (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 22 + k
					sc.ModbusData[mi] = tmp[k]
				}
			case TagTypeElectricity:
				// save resistance temperature (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 20 + k
					ui := ustart + 24 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}

				// compute electricity
				di := ustart + 32
				d := uint(sc.Data[i][di])
				fd := public.Electricity(d)
				md := uint16(fd * 100)

				tmp := make([]byte, 2)
				binary.BigEndian.PutUint16(tmp, md)

				// save sensor data (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 22 + k
					sc.ModbusData[mi] = tmp[k]
				}
			case TagTypeVoltage:
				// save resistance temperature (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 20 + k
					ui := ustart + 24 + k
					sc.ModbusData[mi] = sc.Data[i][ui]
				}

				// compute voltage
				di := ustart + 32
				d := uint(sc.Data[i][di])
				fd := public.Voltage(d)
				md := uint16(fd * 100)

				tmp := make([]byte, 2)
				binary.BigEndian.PutUint16(tmp, md)

				// save sensor data (2 bytes)
				for k := 0; k < 2; k++ {
					mi := mstart + 22 + k
					sc.ModbusData[mi] = tmp[k]
				}
			}

			// u led rgb, 6 bytes
			ustart = 12 + j*3
			for k := 0; k < 6; k++ {
				mi := mstart + 24 + k
				ui := ustart + k
				sc.ModbusData[mi] = sc.Data[i][ui]
			}

			// u key value
			ustart = j
			for k := 0; k < 2; k++ {
				mi := mstart + 30 + k
				ui := ustart + k
				sc.ModbusData[mi] = sc.Data[i][ui]
			}
		}
	}
}

// ID client's id
func (sc *SensorflowClient) ID() string {
	return sc.ClientID
}

// Sample sensorflow sample
func (sc *SensorflowClient) Sample(payload string) (string, error) {
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	slaveid := req.Slaveid

	var result []byte
	if slaveid == 0 {
		// use modbus data
		address := req.Address
		quantity := req.Quantity
		result = sc.ModbusData[address*2 : (address+quantity)*2]
	} else {
		// check online or not
		unum := sc.OnlineModuleNumber
		if int(slaveid-1)/sc.UKeyNumber >= unum {
			return "", fmt.Errorf("u tracker offline")
		}

		uindex := int(slaveid-1) / sc.UKeyNumber

		if uindex >= MaxUTrakerNumber {
			return "", fmt.Errorf("slave id too large")
		}

		address := req.Address
		quantity := req.Quantity
		id := int32(int(slaveid-1) % sc.UKeyNumber)

		switch address {
		case 0:
			// u tracker button
			address = id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			count := binary.BigEndian.Uint16(data)

			r := uint16(0)
			if count > 0 {
				r = 1 // represent press
			} else {
				r = 0
			}

			// post button press payload
			last, ok := sc.buttonPressLastStatus[req.SUID]
			if !ok {
				// do not exist, set value
				sc.buttonPressLastStatus[req.SUID] = r
			} else {
				// exist, set value and send status
				sc.buttonPressLastStatus[req.SUID] = r

				if last != r {
					switch r {
					case 0:
						// post button press payload
						sc.PostCurrentLedStatus(req.SUID, uindex, int(id))
					case 1:
						sc.PostButtonPressInfo(req.SUID)
					}
				}
			}

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, r)
			result = tmp
		case 1:
			// u tracker led
			address = 6 + 3*id
			start := address * 2
			br := sc.Data[uindex][start : start+2]
			bg := sc.Data[uindex][start+2 : start+4]
			bb := sc.Data[uindex][start+4 : start+6]

			r := binary.BigEndian.Uint16(br)
			g := binary.BigEndian.Uint16(bg)
			b := binary.BigEndian.Uint16(bb)

			const remote = 0x8000
			if r >= remote {
				r = r - remote
			}

			if g >= remote {
				g = g - remote
			}

			if b >= remote {
				b = b - remote
			}

			tmp := make([]byte, 2)

			binary.BigEndian.PutUint16(tmp, r)
			result = append(result, tmp...)

			binary.BigEndian.PutUint16(tmp, g)
			result = append(result, tmp...)

			binary.BigEndian.PutUint16(tmp, b)
			result = append(result, tmp...)
		case 4:
			// tag status
			address = 24 + id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			count := binary.BigEndian.Uint16(data)

			r := uint16(0)
			if count < 2 {
				r = 1 // represent online
			} else {
				r = 0 // offline
			}

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, r)
			result = tmp
		case 10:
			// tag id
			address = 32 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 17:
			// tag button status
			address = 39 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 18:
			// tag vibration status
			address = 40 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 19:
			// tag led
			address = 41 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 22:
			// thermistor temperature
			address = 44 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 23:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypeTH {
				return "", fmt.Errorf("tag type not support")
			}

			// environment humidity
			address = 45 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 24:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypeTH {
				return "", fmt.Errorf("tag type not support")
			}

			// environment temperature
			address = 46 + 18*id
			result = sc.Data[uindex][address*2 : (address+quantity)*2]
		case 25:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypeDoorMagnetic {
				return "", fmt.Errorf("tag type not support")
			}

			// door
			address = 48 + 18*id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			count := binary.BigEndian.Uint16(data)

			r := uint16(0)
			if count > 0 {
				r = 1 // door close
			} else {
				r = 0 // open
			}

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, r)
			result = tmp
		case 26:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypeAcoustoOptic {
				return "", fmt.Errorf("tag type not support")
			}

			// acousto-optic
			address = 48 + 18*id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			count := binary.BigEndian.Uint16(data)

			r := uint16(0)
			if count > 0 {
				r = 0 // close
			} else {
				r = 1 // open
			}

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, r)
			result = tmp
		case 27:
			// smoke
		case 28:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypePT100 {
				return "", fmt.Errorf("tag type not support")
			}

			// PT100
			address = 45 + 18*id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			d := binary.BigEndian.Uint16(data)

			fA := float64(39)
			fI := 0.00125
			fR2 := 78.7
			r := public.PT100(uint(d), fA, fI, fR2)

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, uint16(r*100))
			result = tmp
		case 29:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypeElectricity {
				return "", fmt.Errorf("tag type not support")
			}

			// 4-20mA
			address = 45 + 18*id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			d := binary.BigEndian.Uint16(data)

			r := public.Electricity(uint(d))

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, uint16(r*100))
			result = tmp
		case 30:
			address = (38+18*id)*2 + 1
			t := sc.Data[uindex][address]
			if t != TagTypeVoltage {
				return "", fmt.Errorf("tag type not support")
			}

			// 0-10V
			address = 45 + 18*id
			data := sc.Data[uindex][address*2 : (address+quantity)*2]
			d := binary.BigEndian.Uint16(data)

			r := public.Voltage(uint(d))

			tmp := make([]byte, 2)
			binary.BigEndian.PutUint16(tmp, uint16(r*100))
			result = tmp
		default:
			return "", fmt.Errorf("address no support")
		}

		// log.Println(sc.Data[uindex][:120], address, quantity, uindex, id, result)
	}

	return fmt.Sprintf("%X", result), nil
}

// Command sensorflow command
func (sc *SensorflowClient) Command(payload string) (string, error) {
	// decode payload
	req, err := DecodeModbusPayload(payload)
	if err != nil {
		return "", fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	if req.Mode == public.SensorCommandModeSync {
		return sc.sync()
	}

	slaveid := uint16(req.Slaveid)
	address := uint16(req.Address)
	quantity := uint16(req.Quantity)
	// log.Printf("slaveid: %v, addr: %v, quan: %v", slaveid, address, quantity)

	bval := []byte{}

	switch v := req.Value.(type) {
	case int:
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(v))
		bval = append(bval, b...)
	case float64:
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(v))
		bval = append(bval, b...)
	case string:
		sep := ","
		vals := strings.Split(v, sep)

		l := len(vals)
		if l <= 0 {
			return "", fmt.Errorf("value null")
		}

		if vals[0] == "id" {
			// id
			if l >= 2 {
				bval = []byte(vals[1])
			} else {
				bval = []byte{}
			}
		} else {
			// led
			if len(vals) != int(quantity) {
				return "", fmt.Errorf("number of value not match, want[%v], input[%v], split by ','", quantity, len(vals))
			}

			// tranfer to bytes
			for _, v := range vals {
				i, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("invalid param [%v]", v)
				}

				// reset value
				if i < 0 {
					i = 0
				}

				b := make([]byte, 2)
				binary.BigEndian.PutUint16(b, uint16(i))
				bval = append(bval, b...)
			}
		}
	default:
		return "", fmt.Errorf("value type not support")
	}

	// compute which client to be invoked, which register to be set
	return sc.setCommandInfo(slaveid, address, quantity, bval, req.Mode, req.ColorTable)
}

// reset rgb value
func resetRGBValue(r, g, b *int) {
	const remote = 0x8000
	log.Printf("\n-----%d, %d, %d-----\n+++++++\n", *r, *g, *b)
	if ((*r == 0) && (*g == 0) && (*b == 0)) || ((*r == -1) && (*g == -1) && (*b == -1)) {
		*r = 1
		*g = 0
		*b = 0
	} else {
		*r = *r + remote
	}
}

// tag start at 24, u tracker start at 29, number of interval registers are 16
func (sc *SensorflowClient) setCommandInfo(slaveid, address, quantity uint16, value []byte, mode, colortable string) (string, error) {
	// slaveid equal to 0 use modbus data
	if slaveid == 0 {
		pos := (address - 23) % 16
		uindex := (address - 23) / 16

		index := int(uindex) / sc.UKeyNumber
		id := int(uindex) % sc.UKeyNumber

		// log.Printf("pos: %v, uindex: %v, index: %v, id: %v", pos, uindex, index, id)

		switch pos {
		case 0:
			// tag start at red
			switch quantity {
			case 1:
				r := int(binary.BigEndian.Uint16(value))
				sc.SetTagLED(index, id+1, r, 0, 0)
			case 2:
				r := int(binary.BigEndian.Uint16(value[:2]))
				g := int(binary.BigEndian.Uint16(value[2:4]))
				sc.SetTagLED(index, id+1, r, g, 0)
			case 3:
				r := int(binary.BigEndian.Uint16(value[:2]))
				g := int(binary.BigEndian.Uint16(value[2:4]))
				b := int(binary.BigEndian.Uint16(value[4:6]))
				sc.SetTagLED(index, id+1, r, g, b)
			default:
				return "", fmt.Errorf("invalid value")
			}
		case 1:
			// tag start at green
			switch quantity {
			case 1:
				g := int(binary.BigEndian.Uint16(value))
				sc.SetTagLED(index, id+1, 0, g, 0)
			case 2:
				g := int(binary.BigEndian.Uint16(value[:2]))
				b := int(binary.BigEndian.Uint16(value[2:4]))
				sc.SetTagLED(index, id+1, 0, g, b)
			default:
				return "", fmt.Errorf("invalid value")
			}
		case 2:
			// tag start at blue
			switch quantity {
			case 1:
				b := int(binary.BigEndian.Uint16(value))
				sc.SetTagLED(index, id+1, 0, 0, b)
			default:
				return "", fmt.Errorf("invalid value")
			}
		case 6:
			// u tracker start at red
			switch quantity {
			case 1:
				r := int(binary.BigEndian.Uint16(value))
				sc.SetUTrackerLED(index, id, r, 0, 0)
			case 2:
				r := int(binary.BigEndian.Uint16(value[:2]))
				g := int(binary.BigEndian.Uint16(value[2:4]))
				sc.SetUTrackerLED(index, id, r, g, 0)
			case 3:
				r := int(binary.BigEndian.Uint16(value[:2]))
				g := int(binary.BigEndian.Uint16(value[2:4]))
				b := int(binary.BigEndian.Uint16(value[4:6]))
				sc.SetUTrackerLED(index, id, r, g, b)
			default:
				return "", fmt.Errorf("invalid value")
			}
		case 7:
			// u tracker start at green
			switch quantity {
			case 1:
				g := int(binary.BigEndian.Uint16(value))
				sc.SetUTrackerLED(index, id, 0, g, 0)
			case 2:
				g := int(binary.BigEndian.Uint16(value[:2]))
				b := int(binary.BigEndian.Uint16(value[2:4]))
				sc.SetUTrackerLED(index, id, 0, g, b)
			default:
				return "", fmt.Errorf("invalid value")
			}
		case 8:
			// u tracker start at blue
			switch quantity {
			case 1:
				b := int(binary.BigEndian.Uint16(value))
				sc.SetUTrackerLED(index, id, 0, 0, b)
			default:
				return "", fmt.Errorf("invalid value")
			}
		}
	} else { // slaveid unequal to 0 use u tracker data
		index := int(slaveid-1) / sc.UKeyNumber
		id := int(slaveid-1) % sc.UKeyNumber

		switch address {
		case 1:
			// set u tacker led
			r := int(binary.BigEndian.Uint16(value[:2]))
			g := int(binary.BigEndian.Uint16(value[2:4]))
			b := int(binary.BigEndian.Uint16(value[4:6]))

			resetRGBValue(&r, &g, &b)
			// sc.SetUTrackerLED(index, id, r, g, b)

			var cex colorEx
			cex.UWeiIndex = index
			cex.LEDIndex = id
			cex.Red = r
			cex.Green = g
			cex.Blue = b

			sc.priorityColorTable = append(sc.priorityColorTable, cex)

			sc.colorTable[index*sc.UKeyNumber+id].Red = r
			sc.colorTable[index*sc.UKeyNumber+id].Green = g
			sc.colorTable[index*sc.UKeyNumber+id].Blue = b

			// set cache
			tmp := make([]byte, 6)
			binary.BigEndian.PutUint16(tmp[:2], uint16(r))
			binary.BigEndian.PutUint16(tmp[2:4], uint16(g))
			binary.BigEndian.PutUint16(tmp[4:6], uint16(b))

			start := (6 + id*3) * 2
			copy(sc.Data[index][start:start+6], tmp)

			// write data
			filename := "data" + sc.muid + ".json"
			bytedata, _ := json.MarshalIndent(sc.colorTable, "", "\t")
			ioutil.WriteFile(filename, bytedata, 0666)
		case 19:
			// set tag led
			r := int(binary.BigEndian.Uint16(value[:2]))
			g := int(binary.BigEndian.Uint16(value[2:4]))
			b := int(binary.BigEndian.Uint16(value[4:6]))

			sc.SetTagLED(index, id+1, r, g, b)
		case 10:
			// set tag id
			// log.Printf("set tag id: %v", value)
			sc.SetTagID(index, id+1, value)
		case 26:
			// log.Println("set output", index, id, value)
			v := int(binary.BigEndian.Uint16(value[:2]))
			sc.SetTagOutput(index, id+1, v)
		case 31:
			if mode != "" {
				switch mode {
				case DisplayMarquee:
					go func() {
						m := 6
						c := int(binary.BigEndian.Uint16(value[2:4]))
						t := int(binary.BigEndian.Uint16(value[4:6]))
						r := int(binary.BigEndian.Uint16(value[6:8]))
						g := int(binary.BigEndian.Uint16(value[8:10]))
						b := int(binary.BigEndian.Uint16(value[10:12]))

						for i := 1; i <= c; i++ {
							// lighten
							sc.SetLampWith(index, id+1, m, i, 0, r, g, b)
							time.Sleep(time.Millisecond * time.Duration(t))
						}

						for i := 1; i <= c; i++ {
							// extinguish
							sc.SetLampWith(index, id+1, m, i, 0, 0, 0, 0)
							time.Sleep(time.Millisecond * time.Duration(t))
						}
					}()
				}
			} else if colortable != "" {
				go func() {
					vals := strings.Split(colortable, ",")
					l := len(vals)
					r := l % 4
					switch r {
					case 1:
						vals = append(vals, "0", "0", "0")
					case 2:
						vals = append(vals, "0", "0")
					case 3:
						vals = append(vals, "0")
					}

					l = len(vals)
					m := 6
					t := 0
					for i := 0; i < l; i = i + 4 {
						c, _ := strconv.Atoi(strings.TrimSpace(vals[i]))
						r, _ := strconv.Atoi(strings.TrimSpace(vals[i+1]))
						g, _ := strconv.Atoi(strings.TrimSpace(vals[i+2]))
						b, _ := strconv.Atoi(strings.TrimSpace(vals[i+3]))

						sc.SetLampWith(index, id+1, m, c, t, r, g, b)
						time.Sleep(time.Millisecond * 500)
					}
				}()
			} else {
				m := int(binary.BigEndian.Uint16(value[:2]))
				c := int(binary.BigEndian.Uint16(value[2:4]))
				t := int(binary.BigEndian.Uint16(value[4:6]))
				r := int(binary.BigEndian.Uint16(value[6:8]))
				g := int(binary.BigEndian.Uint16(value[8:10]))
				b := int(binary.BigEndian.Uint16(value[10:12]))

				sc.SetLampWith(index, id+1, m, c, t, r, g, b)
			}
		}
	}

	return "", nil
}

// sync
func (sc *SensorflowClient) sync() (string, error) {
	var data public.SensorSync

	data.ID = sc.muid
	data.Modules = sc.OnlineModuleNumber
	data.UCount = sc.OnlineModuleNumber * sc.UKeyNumber

	for i := 0; i < data.Modules; i++ {
		for j := 0; j < sc.UKeyNumber; j++ {
			var d public.SensorData

			d.UIndex = i*sc.UKeyNumber + j + 1

			// led
			start := (6 + 3*j) * 2
			r := int(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
			g := int(binary.BigEndian.Uint16(sc.Data[i][start+2 : start+4]))
			b := int(binary.BigEndian.Uint16(sc.Data[i][start+4 : start+6]))
			d.LED = strconv.Itoa(r) + "-" + strconv.Itoa(g) + "-" + strconv.Itoa(b)

			// button
			d.Button = int(binary.BigEndian.Uint16(sc.Data[i][j*2 : (j+1)*2]))

			// tag status
			start = (24 + j) * 2
			ts := int(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
			if ts > 2 {
				d.Tag.State = 0
			} else {
				d.Tag.State = 1

				// type and id
				start := (32 + 18*j) * 2
				id := sc.Data[i][start : start+14]

				d.Tag.Type = int(id[13])
				d.Tag.Asset = string(id[:13])

				// button
				start = (32 + 18*j + 7) * 2
				btn := int(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
				if btn > 0 {
					d.Tag.Button = 1
				} else {
					d.Tag.Button = 0
				}

				// vibration
				start = (32 + 18*j + 8) * 2
				d.Tag.Vibration = int(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))

				// led
				start = (32 + 18*j + 9) * 2
				r := int(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
				g := int(binary.BigEndian.Uint16(sc.Data[i][start+2 : start+4]))
				b := int(binary.BigEndian.Uint16(sc.Data[i][start+4 : start+6]))
				d.Tag.LED = strconv.Itoa(r) + "-" + strconv.Itoa(g) + "-" + strconv.Itoa(b)

				// temperature
				start = (32 + 18*j + 12) * 2
				d.Tag.Temperature = float64(binary.BigEndian.Uint16(sc.Data[i][start:start+2])) / 10

				switch d.Tag.Type {
				case TagTypeRFID:
				case TagTypeTH:
					{
						start = (32 + 18*j + 13) * 2
						d.Tag.Humidity = float64(binary.BigEndian.Uint16(sc.Data[i][start:start+2])) / 10
						d.Tag.TemperatureP = float64(binary.BigEndian.Uint16(sc.Data[i][start+2:start+4])) / 10
					}
				case TagTypeDoorMagnetic, TagTypeAcoustoOptic, TagTypeSmoke:
					{
						start = (32 + 18*j + 16) * 2
						tmp := int(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
						if tmp > 0 {
							d.Tag.Door = 1
						} else {
							d.Tag.Door = 0
						}
					}
				case TagTypePT100:
					{
						start = (32 + 18*j + 13) * 2
						tmp := uint(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
						fA := float64(39)
						fI := 0.00125
						fR2 := 78.7
						d.Tag.PT100 = public.PT100(tmp, fA, fI, fR2)
					}
				case TagTypeElectricity:
					{
						start = (32 + 18*j + 13) * 2
						tmp := uint(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
						d.Tag.Electricity = public.Electricity(tmp)
					}
				case TagTypeVoltage:
					{
						start = (32 + 18*j + 13) * 2
						tmp := uint(binary.BigEndian.Uint16(sc.Data[i][start : start+2]))
						d.Tag.Voltage = public.Voltage(tmp)
					}
				case TagTypeLampWith:
				}
			}

			data.US = append(data.US, d)
		}
	}

	bytedata, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytedata), nil
}
