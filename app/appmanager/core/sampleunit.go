/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: sample unit
 *
 */

package core

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"clc.hmu/app/appmanager/appnet"
	"clc.hmu/app/appmanager/element"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/applog"
	"clc.hmu/app/public/sys"
	"github.com/Knetic/govaluate"
	"github.com/gwaylib/errors"
)

// SampleUnit sampleUnit
type SampleUnit struct {
	sys.SampleUnit
}

// ChannelValue channel
type ChannelValue struct {
	ID        string
	Name      string
	LastValue interface{}
}

// SampleUnitEx extend sample unit
type SampleUnitEx struct {
	SU                      *SampleUnit
	Elem                    element.Element
	Channels                []ChannelValue
	Status                  ChannelValue
	CommunicationErrorCount int32
}

// Start start
func (su *SampleUnitEx) Start() error {
	// find element mapping file
	libdir := sys.ElementLibDir
	suffix := ".json"
	filename := ""

	if strings.HasSuffix(su.SU.Element, suffix) {
		filename = libdir + su.SU.Element
	} else {
		filename = libdir + su.SU.Element + suffix
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read element library [%s] fail, errmsg [%v]", filename, err)
	}

	if err := json.Unmarshal(data, &su.Elem); err != nil {
		return fmt.Errorf("unmarshal data fail, errmsg [%v]", err)
	}

	for _, ch := range su.Elem.Channels {
		var channel ChannelValue
		channel.ID = ch.ID
		channel.Name = ch.Name
		// channel.LastValue = 0

		su.Channels = append(su.Channels, channel)
	}

	// set status channel
	su.Status.ID = "_state"
	su.Status.Name = "采集单元状态"
	su.Status.LastValue = -2

	su.CommunicationErrorCount = 0

	return nil
}

func (su *SampleUnitEx) publishValue(muid, chanid, channame string, value interface{}) error {
	payload := public.MessagePayload{
		MonitoringUnitID: muid,
		SampleUnitID:     su.SU.ID,
		ChannelID:        chanid,
		Name:             channame,
		Value:            value,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	return appnet.PublishSampleValues(payload)
}

func (su *SampleUnitEx) modbusSample(mapp element.Mapping, muid, port string, baudrate int32) error {
	slaveid := su.SU.Setting.Address
	timeout := su.SU.Timeout

	// find channel definition
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		code := int32(0)
		switch cm.Code.(type) {
		case int:
			code = int32(cm.Code.(int))
		case float64:
			code = int32(cm.Code.(float64))
		}

		result, err := appnet.Sample(port, baudrate, code, slaveid, cm.Address, cm.Quantity, timeout, muid, su.SU.ID)
		if err != nil || result == "" {
			// log.Printf("sample failed, port {%s}, baudrate{%v}, code {%v}, slaveid {%v}, address {%v}, quantity {%v}, errmsg [%v], result{%v}", port, baudrate, cm.Code, slaveid, cm.Address, cm.Quantity, err, result)
			if err != nil && strings.Contains(err.Error(), "tag type not support") {
				continue
			}

			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)

					log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
				}

				// exit when error count reach threshold
				return errors.New("communication errors reach threashold").As(muid, su.SU.ID)
			}

			continue
		}

		// log.Printf("sample success, port {%s}, baudrate{%v}, code {%v}, slaveid {%v}, address {%v}, quantity {%v}", port, baudrate, cm.Code, slaveid, cm.Address, cm.Quantity)

		// sample success
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)

			log.Printf("sample unit [%v] connect\n", su.SU.ID)
		}

		// according to datatype, deal with result
		switch datatype {
		case public.DataTypeInt, public.DataTypeFloat:
			bval, err := hex.DecodeString(result)
			if err != nil {
				log.Printf("decode result failed, result [%s], errmsg [%v]", result, err)
				continue
			}

			// applog.LOG.Infof("id: %v, value: %v", cm.ChannelID, bval)

			var val int16

			vl := len(bval)
			switch vl {
			case 1:
				val = int16(bval[0])
			case 2:
				val = int16(binary.BigEndian.Uint16(bval))
			case 4:
				val = int16(binary.BigEndian.Uint32(bval))
			case 8:
				val = int16(binary.BigEndian.Uint64(bval))
			default:
				val = 0
			}

			ans := float64(0)

			if cm.Expression != "" {
				expression, err := govaluate.NewEvaluableExpression(cm.Expression)

				parameters := make(map[string]interface{}, 8)
				parameters["val"] = val

				result, err := expression.Evaluate(parameters)
				if err != nil {
					log.Printf("evaluate expression failed, errmsg {%v}", err)
					continue
				}

				ans = result.(float64)
			} else {
				ans = float64(val)
			}

			// special deal float value
			if datatype == public.DataTypeFloat && vl == 4 {
				ans = float64(math.Float32frombits(binary.BigEndian.Uint32(bval)))
			}

			lastvalue := float64(0)
			for i, ch := range su.Channels {
				if ch.ID == cm.ChannelID {
					// get last value
					if ch.LastValue != nil {
						lastvalue = ch.LastValue.(float64)
					} else {
						lastvalue = -1
					}

					// log.Printf("channel id: [%v], current value: [%v]; last value: [%v], cov: [%v]", cm.ChannelID, ans, lastvalue, mapp.Setting.COV)

					// check value, contrast by channel cov then by setting
					change := math.Abs(ans - lastvalue)
					if cm.COV != 0 && change < cm.COV {
						break
					} else if cm.COV == 0 && change < mapp.Setting.COV {
						break
					}

					// save value
					su.Channels[i].LastValue = ans

					// publish
					su.publishValue(muid, cm.ChannelID, ch.Name, ans)

					break
				}
			}
		case public.DataTypeString:
			bytes, err := hex.DecodeString(result)
			if err != nil {
				log.Printf("decode result failed, result [%s], errmsg [%v]", result, err)
				continue
			}

			value := ""
			if cm.Expression == "-" {
				l := len(bytes)
				num := l / 2
				s := []string{}

				for i := 0; i < num; i++ {
					v := binary.BigEndian.Uint16(bytes[i*2 : (i+1)*2])
					c := strconv.Itoa(int(v))
					s = append(s, c)
				}

				value = strings.Join(s, "-")
			} else {
				// cut bytes
				reg := regexp.MustCompile(`\d+`)
				s := reg.FindAllString(cm.Expression, -1)
				if len(s) > 0 {
					length, err := strconv.Atoi(s[0])
					if err == nil && length > 0 && length < len(bytes) {
						bytes = bytes[:length]
					}
				}

				value = string(bytes)
			}

			// save
			lastvalue := ""
			for i, ch := range su.Channels {
				if ch.ID == cm.ChannelID {
					// check last value
					if ch.LastValue != nil {
						lastvalue = ch.LastValue.(string)
					}

					// set value
					su.Channels[i].LastValue = value

					// publish only when it changes
					if lastvalue == value {
						break
					}

					// change, publish
					su.publishValue(muid, cm.ChannelID, ch.Name, value)

					break
				}
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) pmbusSample(mapp element.Mapping, protocol, muid, port string) error {
	type PMBusData struct {
		CID1   byte
		CID2   byte
		LENID  uint16
		Result string
	}

	var data []PMBusData

	// get for different cid
	for _, cm := range mapp.ChannelMappings {
		// check if is exist, pmbus use CID1 and CID2 for data
		isexist := false
		for _, d := range data {
			if d.CID1 == cm.CID1 && d.CID2 == cm.CID2 {
				isexist = true
				break
			}
		}

		// no exist, append
		if !isexist {
			newd := PMBusData{
				CID1:  cm.CID1,
				CID2:  cm.CID2,
				LENID: cm.COMMAND,
			}

			data = append(data, newd)
		}
	}

	// start sample
	for i, d := range data {
		result := ""
		var err error
		if protocol == public.ProtocolOilMachine {
			result, err = appnet.OilMachineSample(port, d.CID1, d.CID2, byte(su.SU.Setting.Address))
		} else {
			result, err = appnet.PMBusSample(port, d.CID1, d.CID2, byte(su.SU.Setting.Address), d.LENID)
		}

		if err != nil || result == "" {
			log.Printf("pmbus sample failed, errmsg {%v}", err)

			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)
				}

				// exit when error count reach threshold
				return fmt.Errorf("communication errors reach threashold")
			}

			continue
		}

		data[i].Result = result

		// sample success
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)
		}

		// sleep for a while
		time.Sleep(time.Millisecond * 100)
	}

	for _, cm := range mapp.ChannelMappings {
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		result := ""
		for _, d := range data {
			if d.CID1 == cm.CID1 && d.CID2 == cm.CID2 {
				result = d.Result
				break
			}
		}

		strlength := len(result)
		offset := cm.Offset * 2
		length := cm.Length * 2

		if offset+length > strlength {
			continue
		}

		result = result[offset : offset+length]

		switch datatype {
		case public.DataTypeInt, public.DataTypeFloat:
			bval, err := hex.DecodeString(result)
			if err != nil {
				log.Printf("decode result failed, result [%s], channel [%s], errmsg [%v]", result, cm.ChannelID, err)
				continue
			}

			log.Printf("channelid: %s; value: %X\n", cm.ChannelID, bval)

			var val int16

			switch len(bval) {
			case 1:
				val = int16(bval[0])
			case 2:
				val = int16(binary.BigEndian.Uint16(bval))
			case 4:
				val = int16(binary.BigEndian.Uint32(bval))
			case 5:
				bval = append([]byte{0x00, 0x00, 0x00}, bval...)
				val = int16(binary.BigEndian.Uint64(bval))
			case 8:
				val = int16(binary.BigEndian.Uint64(bval))
			default:
				val = 0
			}

			ans := float64(0)
			if cm.Expression != "" {
				expression, err := govaluate.NewEvaluableExpression(cm.Expression)

				parameters := make(map[string]interface{}, 8)
				parameters["val"] = val

				r, err := expression.Evaluate(parameters)
				if err != nil {
					log.Printf("evaluate expression failed, errmsg {%v}", err)
					continue
				}

				ans = r.(float64)
			} else {
				ans = float64(val)
			}

			// special deal float value, use little endian
			if datatype == public.DataTypeFloat && len(bval) == 4 {
				ans = float64(math.Float32frombits(binary.LittleEndian.Uint32(bval)))
			}

			log.Printf("channelid: %s; bvalue: %X, cvalue: %v\n", cm.ChannelID, bval, ans)

			lastvalue := float64(0)
			name := ""
			for i, ch := range su.Channels {
				if ch.ID == cm.ChannelID {
					// check last value
					if ch.LastValue != nil {
						lastvalue = ch.LastValue.(float64)
					} else {
						lastvalue = -1
					}

					name = ch.Name

					// set value
					su.Channels[i].LastValue = ans
				}
			}

			// log.Printf("sample unit name {%s}, value {%v}", name, ans)

			change := math.Abs(ans - lastvalue)
			if cm.COV != 0 && change < cm.COV {
				continue
			} else if cm.COV == 0 && change < mapp.Setting.COV {
				continue
			}

			su.publishValue(muid, cm.ChannelID, name, ans)
		case public.DataTypeString:
			r, err := hex.DecodeString(result)
			if err != nil {
				log.Printf("decode result failed, result [%s], errmsg [%v]", result, err)
				continue
			}

			// save
			name := ""
			lastvalue := []byte{}
			for i, ch := range su.Channels {
				if ch.ID == cm.ChannelID {
					// check last value
					if ch.LastValue != nil {
						lastvalue = ch.LastValue.([]byte)
					}

					name = ch.Name

					// set value
					su.Channels[i].LastValue = r
				}
			}

			if string(lastvalue) == string(r) {
				continue
			}

			su.publishValue(muid, cm.ChannelID, name, string(r))
		}
	}

	return nil
}

func (su *SampleUnitEx) snmpSample(mapp element.Mapping, muid, port string) error {
	all := []string{}
	for _, cm := range mapp.ChannelMappings {
		all = append(all, cm.OID)
	}

	m := make(map[string]string)

	// deal for error: oid count (x) is greater than MaxOids (60), every 60 channels for one query
	max := 60
	for {
		oids := []string{}
		lo := len(all)
		if lo > max {
			oids = all[:60]
			all = all[60:]
		} else {
			oids = all
		}

		result, err := appnet.SNMPSample(port, su.SU.Setting.Target, oids)
		if err != nil {
			fmt.Printf("snmp get failed: %v\n", err)

			// send off status
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)
				}

				// exit when error count reach threshold
				return fmt.Errorf("communication errors reach threashold")
			}

			return err
		}

		// send on status
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)
		}

		// parse data, result likes "oid1,val1,oid2,val2,oid3,val3"
		data := strings.Split(result, ",")

		l := len(data)
		for i := 0; i < l; i = i + 2 {
			m[data[i]] = data[i+1]
		}

		if lo <= max {
			break
		}
	}

	for _, cm := range mapp.ChannelMappings {
		if !strings.HasPrefix(cm.OID, ".") {
			cm.OID = "." + cm.OID
		}

		result := m[cm.OID]

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// check last value
				lastvalue := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
				}

				// applog.LOG.Infof("current channel {%s}, last value {%s}, current value {%s}", cm.ChannelID, lastvalue, result)

				// value same as last, do not dispose
				if result == lastvalue {
					continue
				}

				// change, check threshold when data type is int or float
				datatype := ""
				for _, ch := range su.Elem.Channels {
					if ch.ID == cm.ChannelID {
						datatype = ch.DataType
					}
				}

				// check data type
				var val interface{}
				if datatype == public.DataTypeInt || datatype == public.DataTypeFloat {
					cv, err := strconv.ParseFloat(result, 64)
					if err != nil {
						log.Printf("format current value [%v] to float failed: %v", result, err)
						continue
					}

					// set value
					val = cv

					if lastvalue != "" {
						lv, err := strconv.ParseFloat(lastvalue, 64)
						if err != nil {
							log.Printf("format last value [%v] to float failed: %v", lastvalue, err)
							continue
						}

						// do not reach threshold, continue
						change := math.Abs(lv - cv)
						if cm.COV != 0 && change < cm.COV {
							continue
						} else if cm.COV == 0 && change < mapp.Setting.COV {
							continue
						}
					}
				} else {
					val = result
				}

				// save current value
				su.Channels[i].LastValue = result

				// reach, compute expression
				if cm.Expression != "" {
					expression, err := govaluate.NewEvaluableExpression(cm.Expression)

					parameters := make(map[string]interface{}, 8)

					tmp, _ := strconv.Atoi(result)
					parameters["val"] = tmp

					r, err := expression.Evaluate(parameters)
					if err != nil {
						log.Printf("evaluate expression failed, errmsg {%v}", err)
						continue
					}

					val = r.(float64)
				}

				// publish
				su.publishValue(muid, cm.ChannelID, ch.Name, val)
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) selfSample(mapp element.Mapping, protocol, muid, port string) error {
	quantity := len(mapp.ChannelMappings)
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		var result string
		var err error
		if protocol == public.ProtocolHYIOTMU {
			result, err = appnet.SelfSample(port, su.SU.Setting.Model, cm.ChannelID, quantity)
			if err != nil || result == "" {
				continue
			}
		} else if protocol == public.ProtocolLuMiGateway {
			result, err = appnet.LuMiGatewaySample(port, su.SU.Setting.Model, su.SU.Setting.SID, cm.ChannelID)
			if err != nil || result == "" {
				continue
			}
		}

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// check last value
				lastvalue := ""
				name := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
					name = ch.Name
				}

				// set value
				su.Channels[i].LastValue = result

				// applog.LOG.Infof("current channel {%s}, last value {%s}, current value {%s}", cm.ChannelID, lastvalue, result)

				var val interface{}
				if cm.Expression != "" {
					expression, err := govaluate.NewEvaluableExpression(cm.Expression)

					parameters := make(map[string]interface{}, 8)

					tmp, _ := strconv.Atoi(result)
					parameters["val"] = tmp

					r, err := expression.Evaluate(parameters)
					if err != nil {
						log.Printf("evaluate expression failed, errmsg {%v}", err)
						continue
					}

					val = r.(float64)
				} else {
					val = result
				}

				// publish
				if result != lastvalue {
					// applog.LOG.Infof("current channel {%s}, value {%s}", cm.ChannelID, result)

					var suid string
					if protocol == public.ProtocolHYIOTMU {
						suid = "_"
					} else {
						suid = su.SU.ID
					}

					payload := public.MessagePayload{
						MonitoringUnitID: muid,
						SampleUnitID:     suid,
						ChannelID:        cm.ChannelID,
						Name:             name,
						Value:            val,
						Timestamp:        public.UTCTimeStamp(),
						Cov:              true,
						State:            0,
					}

					if err := appnet.PublishSampleValues(payload); err != nil {
						su.Channels[i].LastValue = ""
					}
				}
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) entrySample(mapp element.Mapping, muid, port string) error {
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		// start sample
		result, err := appnet.EntrySample(port, cm.Sequence, cm.Code.(string), cm.Group)
		if err != nil {
			fmt.Printf("entry sample failed: %v\n", err)

			// send off status
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)
				}

				// exit when error count reach threshold
				return fmt.Errorf("communication errors reach threashold")
			}

			return err
		}

		// send on status
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)
		}

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// check last value
				lastvalue := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
				}

				// applog.LOG.Infof("current channel {%s}, last value {%s}, current value {%s}", cm.ChannelID, lastvalue, result)

				// value same as last, do not dispose
				if result == lastvalue {
					continue
				}

				// save current value
				su.Channels[i].LastValue = result

				// publish
				su.publishValue(muid, cm.ChannelID, ch.Name, result)
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) didoSample(mapp element.Mapping, muid, port string) error {
	// find channel definition
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		result, err := appnet.DIDOSample(port)
		if err != nil || result == "" {
			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)

					log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
				}

				// exit when error count reach threshold
				return fmt.Errorf("[%s] communication errors reach threashold", su.SU.ID)
			}

			continue
		}

		// send on status
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)
		}

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// check last value
				lastvalue := ""
				name := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
					name = ch.Name
				}

				// set value
				su.Channels[i].LastValue = result

				// publish
				if result != lastvalue {
					payload := public.MessagePayload{
						MonitoringUnitID: muid,
						SampleUnitID:     su.SU.ID,
						ChannelID:        cm.ChannelID,
						Name:             name,
						Value:            result,
						Timestamp:        public.UTCTimeStamp(),
						Cov:              true,
						State:            0,
					}

					if err := appnet.PublishSampleValues(payload); err != nil {
						su.Channels[i].LastValue = ""
					}
				}
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) es5200Sample(mapp element.Mapping, muid, port string) error {
	// find channel definition
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		result, err := appnet.ES5200Sample(port, cm.CID1, cm.CID2, byte(su.SU.Setting.Address), cm.CommandGroup, cm.CommandType, cm.Length)
		if err != nil {
			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)

					log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
				}

				// exit when error count reach threshold
				return fmt.Errorf("communication errors reach threashold")
			}

			continue
		}

		if result == "" {
			continue
		}

		// send on status
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)
		}

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// check last value
				lastvalue := ""
				name := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
					name = ch.Name
				}

				// set value
				su.Channels[i].LastValue = result

				// publish
				if result != lastvalue {
					payload := public.MessagePayload{
						MonitoringUnitID: muid,
						SampleUnitID:     su.SU.ID,
						ChannelID:        cm.ChannelID,
						Name:             name,
						Value:            result,
						Timestamp:        public.UTCTimeStamp(),
						Cov:              true,
						State:            0,
					}

					if err := appnet.PublishSampleValues(payload); err != nil {
						su.Channels[i].LastValue = ""
					}
				}
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) cameraSample(mapp element.Mapping, muid, port string) error {
	// TODO sample state only, common should read element library
	_, err := appnet.CameraSample(port, su.SU.Setting.Host, "state")
	if err != nil {
		// sample failed
		su.CommunicationErrorCount++
		if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
			su.CommunicationErrorCount = 0

			off := -1
			if su.Status.LastValue != off {
				if err := su.publishValue(muid, su.Status.ID, su.Status.Name, off); err != nil {
					return err
				}

				su.Status.LastValue = off

				log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
			}

			// exit when error count reach threshold
			return fmt.Errorf("communication errors reach threashold")
		}

		return err
	}

	su.CommunicationErrorCount = 0

	on := 0
	if su.Status.LastValue != on {
		if err := su.publishValue(muid, su.Status.ID, su.Status.Name, on); err != nil {
			return err
		}

		su.Status.LastValue = on
	}

	return nil
}

func (su *SampleUnitEx) elecfireSample(mapp element.Mapping, muid, port string) error {
	serialnum := su.SU.Setting.SerialNumber

	// find channel definition
	for _, cm := range mapp.ChannelMappings {
		result, err := appnet.ElecFireSample(port, serialnum, int(cm.Address), int(cm.Quantity))
		if err != nil || result == "" {
			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)

					log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
				}

				// exit when error count reach threshold
				return fmt.Errorf("communication errors reach threashold")
			}

			continue
		}

		// sample success
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)

			log.Printf("sample unit [%v] connect\n", su.SU.ID)
		}

		bval, err := hex.DecodeString(result)
		if err != nil {
			log.Printf("decode result failed, result [%s], errmsg [%v]", result, err)
			continue
		}

		var val int16

		vl := len(bval)
		switch vl {
		case 1:
			val = int16(bval[0])
		case 2:
			val = int16(binary.BigEndian.Uint16(bval))
		case 4:
			val = int16(binary.BigEndian.Uint32(bval))
		case 8:
			val = int16(binary.BigEndian.Uint64(bval))
		default:
			val = 0
		}

		ans := float64(0)

		if cm.Expression != "" {
			expression, err := govaluate.NewEvaluableExpression(cm.Expression)

			parameters := make(map[string]interface{}, 8)
			parameters["val"] = val

			result, err := expression.Evaluate(parameters)
			if err != nil {
				log.Printf("evaluate expression failed, errmsg {%v}", err)
				continue
			}

			ans = result.(float64)
		} else {
			ans = float64(val)
		}

		lastvalue := float64(0)
		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// get last value
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(float64)
				} else {
					lastvalue = -1
				}

				// check value, contrast by channel cov then by setting
				change := math.Abs(ans - lastvalue)
				if cm.COV != 0 && change < cm.COV {
					break
				} else if cm.COV == 0 && change < mapp.Setting.COV {
					break
				}

				// save value
				su.Channels[i].LastValue = ans

				// publish
				su.publishValue(muid, cm.ChannelID, ch.Name, ans)

				break
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) virtualAntennaSample(mapp element.Mapping, muid, port string) error {
	// find channel definition
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		result, err := appnet.VirtualAntennaSample(port, cm.ChannelID)
		if err != nil || result == "" {
			fmt.Println("err:", err, "result：", result, "ChannelId:", cm.ChannelID)
			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					su.publishValue(muid, su.Status.ID, su.Status.Name, off)

					log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
				}

				// exit when error count reach threshold
				return fmt.Errorf("communication errors reach threashold")
			}

			continue
		}

		// send on status
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			su.publishValue(muid, su.Status.ID, su.Status.Name, on)
		}

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				// check last value
				lastvalue := ""
				name := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
					name = ch.Name
				}

				// set value
				su.Channels[i].LastValue = result

				// publish
				if result != lastvalue {
					payload := public.MessagePayload{
						MonitoringUnitID: muid,
						SampleUnitID:     su.SU.ID,
						ChannelID:        cm.ChannelID,
						Name:             name,
						Value:            result,
						Timestamp:        public.UTCTimeStamp(),
						Cov:              true,
						State:            0,
					}

					if err := appnet.PublishSampleValues(payload); err != nil {
						su.Channels[i].LastValue = ""
					}
				}
			}
		}
	}

	return nil
}

func (su *SampleUnitEx) commonSample(timeout time.Duration, mapp element.Mapping, muid, port, suid string) error {
	for _, cm := range mapp.ChannelMappings {
		// filter command channel
		datatype := ""

		for _, ch := range su.Elem.Channels {
			if ch.ID == cm.ChannelID {
				datatype = ch.DataType
			}
		}

		if datatype == public.DataTypeCommand {
			continue
		}

		// find channel definitio
		result, err := appnet.CommonSample(time.Duration(su.SU.Timeout)*time.Millisecond, port, suid)
		if err != nil || result == "" {
			// sample failed
			su.CommunicationErrorCount++
			if su.CommunicationErrorCount > su.SU.MaxCommunicationErrors {
				su.CommunicationErrorCount = 0

				off := -1
				if su.Status.LastValue != off {
					su.Status.LastValue = off
					if err := su.publishValue(muid, su.Status.ID, su.Status.Name, off); err != nil {
						applog.LOG.Warning(errors.As(err))
					}

					log.Printf("sample unit [%v] disconnect\n", su.SU.ID)
				}

				// exit when error count reach threshold
				return errors.As(err, "communication errors reach threashold")
			}

			continue
		}

		// send on status
		su.CommunicationErrorCount = 0

		on := 0
		if su.Status.LastValue != on {
			su.Status.LastValue = on
			if err := su.publishValue(muid, su.Status.ID, su.Status.Name, on); err != nil {
				applog.LOG.Warning(errors.As(err, suid))
			}
		}

		// 解析出需要上报的数据
		sampleResult, err := public.ParseSamplePayload(result)
		if err != nil {
			return errors.As(err, suid)
		}
		if !sampleResult.Send {
			return nil
		}

		for i, ch := range su.Channels {
			if ch.ID == cm.ChannelID {
				val, ok := sampleResult.Data[ch.ID]
				if !ok {
					continue
				}
				// check last value
				lastvalue := ""
				name := ""
				if ch.LastValue != nil {
					lastvalue = ch.LastValue.(string)
					name = ch.Name
				}

				// set value
				su.Channels[i].LastValue = val

				// publish
				if val != lastvalue {
					payload := public.MessagePayload{
						MonitoringUnitID: muid,
						SampleUnitID:     su.SU.ID,
						ChannelID:        cm.ChannelID,
						Name:             name,
						Value:            val,
						Timestamp:        public.UTCTimeStamp(),
						Cov:              true,
						State:            0,
					}

					if err := appnet.PublishSampleValues(payload); err != nil {
						su.Channels[i].LastValue = ""
					}
				}

				// 已找到
				break
			}
		}
	}

	return nil
}

// Sample sample
func (su *SampleUnitEx) Sample(muid, port, protocol string, baudrate int32) error {
	findprotocol := false
	var mapp element.Mapping

	// sensorflow protocol same as modbus rtu protocol
	if protocol == public.ProtocolSensorflow {
		protocol = public.ProtocolModbusSerial
	}

	for _, mapping := range su.Elem.Mappings {
		if mapping.Protocol == protocol {
			findprotocol = true
			mapp = mapping
			break
		}
	}

	if !findprotocol {
		log.Printf("protocol [%s] not found in mapping file", protocol)
		return fmt.Errorf("protocol [%s] not found in mapping file", protocol)
	}

	switch protocol {
	case public.ProtocolModbusSerial, public.ProtocolModbusTCP, public.ProtocolHSJRFID:
		return su.modbusSample(mapp, muid, port, baudrate)
	case public.ProtocolPMBUS, public.ProtocolYDN23, public.ProtocolOilMachine:
		return su.pmbusSample(mapp, protocol, muid, port)
	case public.ProtocolSNMP:
		return su.snmpSample(mapp, muid, port)
	case public.ProtocolHYIOTMU, public.ProtocolLuMiGateway:
		return su.selfSample(mapp, protocol, muid, port)
	case public.ProtocolFaceIPC:
	case public.ProtocolWeiGengEntry:
		return su.entrySample(mapp, muid, port)
	case public.ProtocolDIDO:
		return su.didoSample(mapp, muid, port)
	case public.ProtocolES5200:
		return su.es5200Sample(mapp, muid, port)
	case public.ProtocolCamera:
		return su.cameraSample(mapp, muid, port)
	case public.ProtocolElecFire:
		return su.elecfireSample(mapp, muid, port)
	case public.ProtocolVirtualAntenna:
		return su.virtualAntennaSample(mapp, muid, port)
	}

	// 默认使用透传指令，由采集单元自行处理逻辑。
	// TODO: 统一使用此透传，以便用逻辑不超出单元的范围。
	// 只需传输端口号与采集单元ID过去，协议自行再取出配置文件进行处理
	// 配置文件读取:public.GetMonitoringUnitCfg()
	return errors.As(su.commonSample(time.Duration(su.SU.Timeout)*time.Millisecond, mapp, muid, port, su.SU.ID))
}

// Command command
func (su *SampleUnitEx) Command(port, protocol, chid string, baudRate int32, value int) string {
	findprotocol := false
	var mapp element.Mapping

	for _, mapping := range su.Elem.Mappings {
		if mapping.Protocol == protocol {
			findprotocol = true
			mapp = mapping
			break
		}
	}

	if !findprotocol {
		applog.LOG.Warningf("protocol [%s] not found in mapping file", protocol)
		return ""
	}

	slaveid := su.SU.Setting.Address
	timeout := su.SU.Timeout

	// applog.LOG.Infof("%+v", mapp)
	for _, cm := range mapp.ChannelMappings {
		code := int32(0)
		switch cm.Code.(type) {
		case int:
			code = int32(cm.Code.(int))
		case float64:
			code = int32(cm.Code.(float64))
		}

		if cm.ChannelID == chid {
			result, err := appnet.Command(port, baudRate, code, slaveid, cm.Address, int32(value), timeout)
			if err != nil {
				applog.LOG.Warningf("sample failed, errmsg [%s]", err)
				continue
			}

			return result
		}
	}

	return "channel not found"
}

// CommonCommand value: string format of interface request
func (su *SampleUnitEx) CommonCommand(port, protocol, suid, chid string, baudRate int32, value interface{}) (string, error) {
	findprotocol := false
	var mapp element.Mapping

	// sensorflow same as modbus rtu
	pro := protocol
	if protocol == public.ProtocolSensorflow {
		pro = public.ProtocolModbusSerial
	}

	for _, mapping := range su.Elem.Mappings {
		if mapping.Protocol == pro {
			findprotocol = true
			mapp = mapping
			break
		}
	}

	if !findprotocol {
		applog.LOG.Warningf("protocol [%s] not found in mapping file", protocol)
		return "", fmt.Errorf("protocol not found")
	}

	// pasre parameter
	var para public.CommandParameter

	bytepara, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("invalid parameters")
	}

	if err := json.Unmarshal(bytepara, &para); err != nil {
		return "", fmt.Errorf("invalid parameters")
	}

	netTimeout := time.Duration(su.SU.Timeout) * time.Millisecond

	// log.Printf("parameter: %+v", para)
	// applog.LOG.Infof("%+v", mapp)
	for _, cm := range mapp.ChannelMappings {
		if cm.ChannelID == chid {
			result := ""
			switch protocol {
			case public.ProtocolModbusSerial, public.ProtocolModbusTCP, public.ProtocolSensorflow:
				{
					slaveid := su.SU.Setting.Address
					timeout := su.SU.Timeout

					if protocol == public.ProtocolSensorflow && slaveid > 0 {
						switch cm.Address {
						case 1, 19:
							para.Value = strconv.Itoa(para.Red) + "," + strconv.Itoa(para.Green) + "," + strconv.Itoa(para.Blue)
						case 10:
							para.Value = "id," + para.Value.(string)
						}
					}

					code := int32(0)
					switch cm.Code.(type) {
					case int:
						code = int32(cm.Code.(int))
					case float64:
						code = int32(cm.Code.(float64))
					}

					var p = public.ModbusPayload{
						Port:       port,
						BaudRate:   baudRate,
						Code:       code,
						Slaveid:    slaveid,
						Address:    cm.Address,
						Quantity:   cm.Quantity,
						Value:      para.Value,
						Timeout:    timeout,
						ColorTable: para.ColorTable,
						Mode:       para.Mode,
					}

					payload, err := json.Marshal(p)
					if err != nil {
						return "", fmt.Errorf("encode payload failed")
					}

					result, err = appnet.CommonCommand(netTimeout, port, suid, string(payload))
					if err != nil {
						return "", err
					}
				}
			case public.ProtocolLampWith:
				var p = public.LampWithOperationPayload{
					Value: para.Value,
				}

				payload, err := json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(payload))
				if err != nil {
					return "", err
				}
			case public.ProtocolCreditCard:
				result, err = appnet.CommonCommand(netTimeout, port, suid, "")
				if err != nil {
					return "", err
				}
			case public.ProtocolFaceIPC:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				var p public.FaceIPCOperationPayload
				if err := json.Unmarshal(b, &p); err != nil {
					return "", fmt.Errorf("payload illegal: %v", err)
				}

				p.Perpose = chid
				b, err = json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			case public.ProtocolSNMP:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				var p public.SNMPOperationPayload
				if err := json.Unmarshal(b, &p); err != nil {
					return "", fmt.Errorf("payload illegal: %v", err)
				}

				// set oid
				p.Target = su.SU.Setting.Target
				p.OID = cm.OID

				b, err = json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			case public.ProtocolWeiGengEntry:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				var p public.WeiGengEntryOperationPayload
				if err := json.Unmarshal(b, &p); err != nil {
					return "", fmt.Errorf("payload illegal: %v", err)
				}

				// set param
				p.FunctionID = cm.Code.(string)

				b, err = json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			case public.ProtocolES5200:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				var p public.ES5200OperationPayload
				if err := json.Unmarshal(b, &p); err != nil {
					return "", fmt.Errorf("payload illegal: %v", err)
				}

				// set oid
				p.ADR = byte(su.SU.Setting.Address)
				p.CID1 = cm.CID1
				p.CID2 = cm.CID2
				p.LENID = cm.Length
				p.COMMANDGROUP = cm.CommandGroup
				p.COMMANDTYPE = cm.CommandType

				b, err = json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			case public.ProtocolLuMiGateway, public.ProtocolDIDO:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			case public.ProtocolVirtualAntenna:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				var p public.VirtualAntennaOperationPayload
				if err := json.Unmarshal(b, &p); err != nil {
					return "", fmt.Errorf("payload illegal: %v", err)
				}

				p.Perpose = chid
				b, err = json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			case public.ProtocolCamera:
				b, err := json.Marshal(value)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				var p public.CameraOperationPayload
				if err := json.Unmarshal(b, &p); err != nil {
					return "", fmt.Errorf("payload illegal: %v", err)
				}

				p.Host = su.SU.Setting.Host
				b, err = json.Marshal(p)
				if err != nil {
					return "", fmt.Errorf("encode payload failed")
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", err
				}
			default:
				// 默认使用透传指令，由采集单元自行处理逻辑。
				// TODO: 统一使用此透传，以便用逻辑不超出单元的范围。
				b, err := json.Marshal(value)
				if err != nil {
					return "", errors.New("encode payload failed").As(value)
				}

				result, err = appnet.CommonCommand(netTimeout, port, suid, string(b))
				if err != nil {
					return "", errors.As(err, value)
				}
			}

			return result, nil
		}
	}

	return "", fmt.Errorf("channel [%v] not found", chid)
}
