/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: bus server implement
 *
 */

package module

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "clc.hmu/app/busmanager/buspb"
	"clc.hmu/app/busmanager/module/web"
	"clc.hmu/app/extend"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/bootflag"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// BusServer is used to implement busmanager.BusServer.
type BusServer struct {
	MqttClient MQTTClient

	enableCache bool
	cache       []string
}

// LEDSetInterval set led interval
var LEDSetInterval = 50

// status
const (
	Offline       = -1 // mqtt确认了下线
	Online        = 0  // mqtt确认了上线
	WaitingOnline = 1  // 检测到网络正常，等待mqtt上线
)

// NetworkStatus network status, 0 for connect, -1 for disconnect
var (
	networkStatus     = -1
	networkStatusSync = sync.Mutex{}
)

func GetNetworkStatus() int {
	networkStatusSync.Lock()
	defer networkStatusSync.Unlock()
	return networkStatus
}

func SetNetworkStatus(status int, mqtt bool) {
	networkStatusSync.Lock()
	defer networkStatusSync.Unlock()

	if networkStatus == Online && !mqtt {
		// 当mqtt在线时，只能由mqtt处理
		return
	}
	networkStatus = status
}

// Init do some init operation
func (s *BusServer) Init() {
	cfg := sys.GetBusManagerCfg()

	// download element library
	downloadDependentDeviceLibrary()

	if err := web.ReadVideoConfig(); err != nil {
		log.Printf("open video config file failed, errmsg {%v}", err)
	}

	s.enableCache = true

	// check directory exist or not
	_, err := os.Stat(sys.GetBusManagerCfg().Cache.Directory)
	if err != nil {
		if os.IsNotExist(err) {
			// do not exist, create
			if err := os.Mkdir(cfg.Cache.Directory, os.ModeDir); err != nil {
				log.Printf("create directory failed: %s", err)
				s.enableCache = false
			}
		} else {
			s.enableCache = false
		}
	}

	log.Printf("enable cache: %v", s.enableCache)

	muid := getUUID()
	willtopic := "sample-values/" + muid + "/_/_state"
	conntopic := willtopic

	payload := public.MessagePayload{
		MonitoringUnitID: muid,
		SampleUnitID:     "_",
		ChannelID:        "_state",
		Name:             "采集器连接状态",
		Value:            -1,
		// Timestamp:        public.UTCTimeStamp(),// 此值因与服务器发生了冲突，估不上报
		Cov:   true,
		State: 0,
	}

	willpayload, _ := json.Marshal(payload)

	payload.Value = 0
	connpayload, _ := json.Marshal(payload)

	s.MqttClient = NewMQTTClient(SubMessageHandler, willtopic, string(willpayload), conntopic, string(connpayload))

	if err := s.MqttClient.ConnectServer(); err != nil {
		log.Printf("connect mqtt server failed, errmsg {%v}, start reconnect...", err)

		// start to reconnect
		go s.MqttClient.ReconnectServer()
	}

	// check network status
	go checkNetworkStatus()

	s.MqttClient.Subscribe("sample-values/+/_/upgrade")
	s.MqttClient.Subscribe("command/" + muid + "/#")

	// init status
	web.DeviceStatus = make(map[string]int)

	mu := sys.GetMonitoringUnitCfg()
	for _, sp := range mu.SamplePorts {
		for _, su := range sp.SampleUnits {
			web.DeviceStatus[su.ID] = 0
		}
	}

	// set led status
	go controlAppLEDStatus()

	// check start log
	go func() {
		topic := "sample-values/" + getUUID() + "/_/restart"
		payload := public.MessagePayload{
			MonitoringUnitID: cfg.MQTT.ClientID,
			SampleUnitID:     "_",
			ChannelID:        "restart",
			Name:             "",
			Value:            0,
			Timestamp:        public.UTCTimeStamp(),
			Cov:              true,
			State:            0,
		}

		flag, err := bootflag.GetFlag()
		if err != nil {
			log.Warning(errors.As(err))
			flag = "-1"
		}
		switch flag {
		case "0":
			payload.Value = 1
		case "1":
			payload.Value = 2
		default:
			// using 0
		}
		bp, _ := json.Marshal(payload)
		for {
			if GetNetworkStatus() == Online {
				s.MqttClient.PublishSampleValues(topic, string(bp))
				break
			}
			time.Sleep(time.Second)
		}
		if err := bootflag.CleanFlag(); err != nil {
			log.Warning(errors.As(err))
		}
	}()
}

// Cleanup cleanup
func (s *BusServer) Cleanup() {
	s.MqttClient.DisconnectServer()
}

func (s *BusServer) publishCacheFile() error {
	files, err := filepath.Glob(filepath.Join(sys.GetBusManagerCfg().Cache.Directory, "*"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	topic := "sample-block/" + getUUID()
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			continue
		}

		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			continue
		}

		sd := string(data)
		ms := strings.Split(sd, "\n")

		sp := []string{}
		for _, m := range ms {
			ps := strings.Split(m, "&")
			if len(ps) < 2 {
				continue
			}

			p := ps[1]
			sp = append(sp, p)
		}

		d := strings.Join(sp, ",")
		d = "[" + d + "]"

		log.Printf("publish file: %s", filename)

		// publish data
		if err := s.MqttClient.PublishSampleValues(topic, d); err != nil {
			return err
		}

		// remove cache file
		os.Remove(filename)
	}

	return nil
}

func saveCacheToFile(cache []string) error {
	cfg := sys.GetBusManagerCfg()
	files, err := filepath.Glob(filepath.Join(cfg.Cache.Directory, "*"))
	if err != nil {
		return err
	}

	var ifl []int
	for _, f := range files {
		fn := filepath.Base(f)
		i, _ := strconv.Atoi(fn)
		ifl = append(ifl, i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(ifl)))
	log.Println(ifl)

	l := len(ifl)
	if l > cfg.Cache.MaxFile {
		rfs := ifl[cfg.Cache.MaxFile:]

		// remove files
		for _, f := range rfs {
			os.Remove(filepath.Join(cfg.Cache.Directory, strconv.Itoa(f)))
		}
	}

	var nf int
	if l == 0 {
		nf = 0
	} else {
		nf = ifl[0] + 1
	}

	filepath := filepath.Join(cfg.Cache.Directory, strconv.Itoa(nf))
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer f.Close()

	data := strings.Join(cache, "\n")

	if _, err = f.Write([]byte(data)); err != nil {
		return fmt.Errorf("write file [%s] failed: %s", filepath, err)
	}

	return nil
}

// Publish publish implement
func (s *BusServer) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishReply, error) {
	// check message and boradcast to clients when necessary
	checkMessage(in.Topic, in.Payload)

	go func(a string, b []byte) {
		web.PayloadMap.Store(a, b)
		web.PayloadChan <- b
	}(in.Topic, []byte(in.Payload))

	// check whether should capture or not
	buscfg := sys.GetBusManagerCfg()
	for _, cap := range buscfg.CaptureOption {
		match := false
		for _, signal := range cap.Signals {
			if strings.Contains(in.Topic, signal.Topic) {
				// topic coincident, check value
				var p public.MessagePayload
				if err := json.Unmarshal([]byte(in.Payload), &p); err != nil {
					continue
				}

				val := ""
				switch p.Value.(type) {
				case int:
					val = strconv.Itoa(p.Value.(int))
				case float64:
					val = strconv.Itoa(int(p.Value.(float64)))
				case string:
					val = p.Value.(string)
				}

				if signal.Value == val {
					match = true
					break
				}
			}
		}

		// match, capture
		if match {
			var p public.CommandPayload
			var para public.CommandParameter

			muid := getUUID()
			chid := "capture"

			para.Channel = chid

			p.MonitoringUnit = muid
			p.SampleUnit = cap.SUID
			p.Channel = chid
			p.StartTime = public.UTCTimeStamp()
			p.Phase = public.PhaseExcuting
			p.Parameters = para

			topic := "command/" + muid + "/" + cap.SUID + "/" + chid
			msg, err := json.Marshal(p)
			if err != nil {
				continue
			}

			// publish
			s.MqttClient.PublishSampleValues(topic, string(msg))
		}
	}

	// enable cache
	if s.enableCache {
		// online or not
		if GetNetworkStatus() == Online {
			// check cache files exist or not, send cache files first
			if err := s.publishCacheFile(); err != nil {
				log.Printf("publish failed: %s", err)
				return &pb.PublishReply{Status: public.StatusOK, Message: public.MessageOK}, nil
			}

			// check cache exist, publish
			if len(s.cache) > 0 {
				for _, m := range s.cache {
					ms := strings.Split(m, "&")
					if len(ms) == 2 {
						topic := ms[0]
						payload := ms[1]

						log.Printf("publish cache: %s", m)
						if err := s.MqttClient.PublishSampleValues(topic, payload); err != nil {
							return &pb.PublishReply{Status: public.StatusOK, Message: public.MessageOK}, nil
						}
					}
				}

				s.cache = []string{}
			}

			// then publish current message
			s.MqttClient.PublishSampleValues(in.Topic, in.Payload)

		} else {
			// offline, save data to cache, check cache quantity
			if len(s.cache) < sys.GetBusManagerCfg().Cache.MaxMessage {
				log.Printf("save to cache, current number: %d", len(s.cache))
				s.cache = append(s.cache, in.Topic+"&"+in.Payload)
			} else {
				log.Printf("save to file")
				// save to file
				if err := saveCacheToFile(s.cache); err != nil {
					log.Printf("save cache faield: %s", err)
				}

				s.cache = []string{}
			}
		}
	} else {
		if err := s.MqttClient.PublishSampleValues(in.Topic, in.Payload); err != nil {
			return &pb.PublishReply{Status: public.StatusErr, Message: public.MessageErrUnknown}, nil
		}
	}

	return &pb.PublishReply{Status: public.StatusOK, Message: public.MessageOK}, nil
}

// Subscribe subscribe implement
func (s *BusServer) Subscribe(ctx context.Context, in *pb.SubscribeRequest) (*pb.SubscribeReply, error) {
	if err := s.MqttClient.Subscribe(in.Topic); err != nil {
		return &pb.SubscribeReply{Status: public.StatusErr, Message: err.Error()}, nil
	}

	return &pb.SubscribeReply{Status: public.StatusOK, Message: public.MessageOK}, nil
}

// get uuid
func getUUID() string {
	// address := config.Configuration.SystemServer.Host + ":" + config.Configuration.SystemServer.Port

	// // get info from hmu
	// var client public.SystemClient
	// if err := client.ConnectSystemDaemon(address); err != nil {
	// 	log.Fatalf("connect system server failed, errmsg {%v}", err)
	// }

	// resp, err := client.UUID()
	// if err != nil {
	// 	log.Fatalf("get uuid failed, errmsg {%v}", err)
	// }
	// defer client.Disconnect()

	// return resp.UUID

	// read id from config file
	return sys.GetMonitoringUnitCfg().ID
}

// contorl app led status
func controlAppLEDStatus() {
	var appled extend.AppLED
	if err := appled.Prepare(sys.GetBusManagerCfg().Model); err != nil {
		buslog.LOG.Warningf("prepare app led failed, errmsg: %v", err)
		return
	}
	defer appled.CleanUp()

	status := 0

	// loop
	for {
		// toggle status
		status = status ^ 1

		// sleep for a moment, interval set by mqtt connect/disconnect handler
		time.Sleep(time.Millisecond * time.Duration(LEDSetInterval))

		if err := appled.SetLEDStatus(status); err != nil {
			// log.Printf("set appled %v", err)
		}
	}
}

func checkMessage(topic, payload string) {
	s := strings.Split(topic, "/")
	if len(s) != 4 {
		return
	}

	suid := s[2]
	channelid := s[3]

	if channelid != "_state" {
		return
	}

	// parse payload, get value
	var p public.MessagePayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		log.Printf("parse payload fail, payload: %s, errmsg: %v", payload, err)
		return
	}

	v := int(p.Value.(float64))
	if v == -1 {
		v = 0
	} else {
		v = 1
	}

	// set status
	lastvalue, ok := web.DeviceStatus[suid]
	if !ok {
		log.Printf("channel id `%s` do not exist", suid)
		return
	}

	if v != lastvalue {
		// update status, broadcast
		web.DeviceStatus[suid] = v

		bs, _ := web.DeviceStatusToBytes()
		web.WSHub.BroadcastMessage(bs)
	}
}

func checkNetworkStatus() {
	timer := time.NewTicker(5 * time.Second)
	lastRestartTime := time.Now()
	cfg := sys.GetBusManagerCfg()
	netCheckList := []string{cfg.MQTT.Host + ":" + cfg.MQTT.Port}
	netCheckList = append(netCheckList, cfg.Web.NetChecking.Hosts...)
	doTimeout := cfg.Web.NetChecking.Timeout
	if doTimeout == 0 {
		doTimeout = 5
	}
	doTimes := cfg.Web.NetChecking.DoTimes

	for {
		select {
		case <-timer.C:
			status := GetNetworkStatus()
			// 大部分网络是正常的，优化走这个
			if status == Online {
				lastRestartTime = time.Now()
				continue
			}

			if status == WaitingOnline && len(cfg.Web.NetChecking.Hosts) > 0 {
				// 处理检测到网络正常时的ticker事件
				sysd := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
				if _, err := sysd.AutoCheckNetworking(netCheckList, time.Duration(doTimeout)*1e9); err == nil {
					SetNetworkStatus(WaitingOnline, false)
				} else {
					// 在等待mqtt上线的过程中发如果检查到网络又下线了，恢复到网络不可用的状态。
					SetNetworkStatus(Offline, false)
				}
				sysd.Disconnect()
				continue
			}

			if status == Offline {
				if doTimes > 0 && len(cfg.Web.NetChecking.Hosts) > 0 {
					doTimes--
					// 先尝试网络
					sysd := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
					if _, err := sysd.AutoCheckNetworking(netCheckList, time.Duration(doTimeout)*1e9); err == nil {
						// 检测到网络正常了, 执行等待mqtt上线的逻辑。
						sysd.Disconnect()
						SetNetworkStatus(WaitingOnline, false)
						continue
					}

					sysd.Disconnect()
					// 网络失败，走失败的逻辑
				}

				now := time.Now()
				d := now.Sub(lastRestartTime)
				rd := time.Duration(cfg.Web.Restart.Duration) * time.Second
				if d >= rd {
					lastRestartTime = now
					rt := cfg.Web.Restart.Times
					if rt < cfg.Web.Restart.Max {
						// add retart times
						cfg.Web.Restart.Times++
						if err := sys.SaveBusManagerCfg(cfg); err != nil {
							buslog.LOG.Warningf("save bus config failed, errmsg {%v}", err)
						}

						buslog.LOG.Infof("software restart: %d times", cfg.Web.Restart.Times)

						// software rstart
						if err := public.RestartApp(cfg.Model, errors.New(public.RestartByCommunicationInterrupt)); err != nil {
							buslog.LOG.Warning(errors.As(err))
						}
					} else {
						// clear times
						cfg.Web.Restart.Times = 0
						if err := sys.SaveBusManagerCfg(cfg); err != nil {
							buslog.LOG.Warningf("save bus config failed, errmsg {%v}", err)
						}

						buslog.LOG.Info("hardware restart")

						// hardware restart
						if err := public.Reboot(errors.New(public.RebootByCommunicationInterrupt)); err != nil {
							buslog.LOG.Warning(errors.As(err))
						}
					}
				}
			}
		}
	}
}

func downloadDependentDeviceLibrary() error {
	cfg := sys.GetBusManagerCfg()
	mu := sys.GetMonitoringUnitCfg()
	elementPath := os.ExpandEnv(etc.Etc.String("public", "element-dir"))
	for _, sp := range mu.SamplePorts {
		for _, su := range sp.SampleUnits {
			// check device library exist or not
			ep := filepath.Join(elementPath, su.Element)
			_, err := os.Stat(ep)
			if err == nil {
				// do not update when exist
				log.Debugf("element library [%s] exist", ep)
				continue
			}

			// other error occured
			if !os.IsNotExist(err) {
				log.Debugf("check element library [%s] existence fail", errors.As(err, ep))
				continue
			}

			// do not exist, download and save
			np := cfg.Web.ElementLib.Server + su.Element
			if err := public.HTTPDownloadFile(np, ep); err != nil {
				fmt.Printf("download or save element library [%s] failed: %s\n", su.Element, errors.As(err))
				continue
			}

			fmt.Printf("download or save element library [%s] success\n", su.Element)
		}
	}

	return nil
}
