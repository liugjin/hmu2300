/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: server implement
 *
 */

package appmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"clc.hmu/app/appmanager/appnet"
	"clc.hmu/app/appmanager/core"
	"clc.hmu/app/public/log/applog"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"

	pb "clc.hmu/app/appmanager/apppb"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
)

// Version version
var Version = ""

type initiativePushPayload struct {
	hasPushed bool
	data      interface{}
}

// appserver is used to implement busmanager.BusServer.
type appserver struct {
	InitiativePushValue map[string]initiativePushPayload
}

// Notify notify implement
func (s *appserver) Notify(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyReply, error) {
	switch in.Caller {
	case public.CallerBusServer:
		if err := s.DealUpperMessage(in.Topic, in.Payload); err != nil {
			log.Printf("deal upper message failed, errmsg: %v", errors.As(err))
			return &pb.NotifyReply{Status: public.StatusErr, Message: err.Error()}, nil
		}
	case public.CallerPortServer:
		if err := s.HandlePushMessage(in.Topic, in.Payload); err != nil {
			log.Printf("handle active push message failed, errmsg: %v", err)
			return &pb.NotifyReply{Status: public.StatusErr, Message: err.Error()}, nil
		}
	}

	return &pb.NotifyReply{Status: public.StatusOK, Message: public.MessageOK}, nil
}

// DealUpperMessage deal upper message
func (s *appserver) DealUpperMessage(topic, payload string) error {
	// applog.LOG.Infof("receive topic: %s, payload: %s", topic, payload)

	// topic like 'command/mu/su/channel', split by '/' for mu, su, ch
	ts := strings.Split(topic, "/")
	if len(ts) < 4 {
		return errors.New("unknown command topic").As(topic)
	}

	muid := ts[1]
	suid := ts[2]
	chid := ts[3]

	// applog.LOG.Info(muid, suid, chid)

	// parse
	var p public.CommandPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return errors.As(err, topic, payload)
	}

	switch p.Phase {
	case public.PhaseExcuting:
		// pass
	case public.PhaseComplete:
		// 由终端发出的complete协议, 直接忽略
		return nil
	default:
		return errors.New("Unknow Phase").As(p.Phase)
	}

	start := public.ParseUTCTimeStamp(p.StartTime)
	truestart := time.Now()

	switch ts[0] {
	case "sample-values":
		// update muid
		p.MonitoringUnit = sys.GetMonitoringUnitCfg().ID

		interval := time.Since(truestart)
		p.EndTime = public.TimeToSTring(start.Add(interval))
		p.Phase = public.PhaseComplete
		// p.UnderlinePhase = public.PhaseComplete

		appnet.ReplyCommand(topic, p)

		// restart app
		time.Sleep(time.Second)
		return nil
		// return public.RestartApp()

	case "command":
		// deal update mu config
		if suid == "_" && chid == "mucfg" {
			if err := s.HandleUpdateMUConfig(p.Parameters); err != nil {
				log.Debug(errors.As(err))
				interval := time.Since(truestart)
				p.EndTime = public.TimeToSTring(start.Add(interval))
				p.Result = err.Error()
				p.Phase = public.PhaseError

				return errors.As(appnet.ReplyCommand(topic, p))
			}

			interval := time.Since(truestart)
			p.EndTime = public.TimeToSTring(start.Add(interval))
			p.Result = "ok"
			p.Phase = public.PhaseComplete

			return errors.As(appnet.ReplyCommand(topic, p))
		}

		// deal remote command
		if suid == "_" && chid == "remote" {
			if err := s.HandleRemoteCommand(payload); err != nil {
				log.Debug(errors.As(err))
				interval := time.Since(truestart)
				p.EndTime = public.TimeToSTring(start.Add(interval))
				p.Result = err.Error()
				p.Phase = public.PhaseError

				return errors.As(appnet.ReplyCommand(topic, p))
			}

			interval := time.Since(truestart)
			p.EndTime = public.TimeToSTring(start.Add(interval))
			p.Result = "ok"
			p.Phase = public.PhaseComplete

			return errors.As(appnet.ReplyCommand(topic, p))
		}

		// special for sensorflow
		if suid == "_" && chid == "leds" {
			if _, err := s.HandleSensorflowSetLEDS(p); err != nil {
				interval := time.Since(truestart)
				p.EndTime = public.TimeToSTring(start.Add(interval))
				p.Result = err.Error()
				p.Phase = public.PhaseError

				return errors.As(appnet.ReplyCommand(topic, p))
			}

			interval := time.Since(truestart)
			p.EndTime = public.TimeToSTring(start.Add(interval))
			p.Result = "ok"
			p.Phase = public.PhaseComplete

			return errors.As(appnet.ReplyCommand(topic, p))
		}

		if suid == "_" && chid == "sync" {
			result, err := s.HandleSensorflowSync()
			if err != nil {
				log.Debug(errors.As(err))
				interval := time.Since(truestart)
				p.EndTime = public.TimeToSTring(start.Add(interval))
				p.Result = err.Error()
				p.Phase = public.PhaseError

				return errors.As(appnet.ReplyCommand(topic, p))
			}

			interval := time.Since(truestart)
			p.EndTime = public.TimeToSTring(start.Add(interval))
			p.Result = result
			p.Phase = public.PhaseComplete

			return errors.As(appnet.ReplyCommand(topic, p))
		}

		// deal self command
		if suid == "_" {
			return errors.As(s.HandleSelfCommand(topic, payload, chid))
		}

		// deal common command
		result, err := s.DealCommonCommand(muid, suid, chid, payload, p)
		if err != nil {
			interval := time.Since(truestart)
			p.EndTime = public.TimeToSTring(start.Add(interval))
			p.Result = err.Error()
			p.Phase = public.PhaseError
			// p.UnderlinePhase = public.PhaseError

			// log.Printf("command fail, err: %v", err)
			log.Warning(errors.As(err))
			return errors.As(appnet.ReplyCommand(topic, p))
		}

		interval := time.Since(truestart)
		p.EndTime = public.TimeToSTring(start.Add(interval))
		p.Result = result
		p.Phase = public.PhaseComplete
		// p.UnderlinePhase = "complete"

		return errors.As(appnet.ReplyCommand(topic, p))
	}

	return nil
}

// DealCommonCommand parse command
func (s *appserver) DealCommonCommand(muid, suid, chid, payload string, p public.CommandPayload) (string, error) {
	// search mu
	mu := sys.GetMonitoringUnitCfg()

	// find channel
	for _, sus := range portmap {
		for _, suex := range sus {
			if suex.SU.ID == suid {
				for _, ch := range suex.Channels {
					if ch.ID == chid {
						port := ""
						protocol := ""
						var baudRate int32

						for _, sp := range mu.SamplePorts {
							for _, su := range sp.SampleUnits {
								if su.ID == suid {
									port = sp.Setting.Port
									baudRate = sp.Setting.BaudRate
									protocol = sp.Protocol
								}
							}
						}

						// send command
						result, err := suex.CommonCommand(port, protocol, suid, chid, baudRate, p.Parameters)
						if err != nil {
							return "", errors.As(err, port, protocol, suid, chid, baudRate, p.Parameters)
						}

						// check for lock singal
						buscfg := sys.GetBusManagerCfg()
						for _, signal := range buscfg.AutoLock {
							ids := strings.Split(signal.Topic, "/")
							if len(ids) != 2 {
								continue
							}

							// check id
							if suid == ids[0] && chid == ids[1] {
								// coincident
								go func() {
									time.Sleep(time.Second * 2)

									// set value
									var para public.CommandParameter
									para.Value = signal.SetValue

									r, err := suex.CommonCommand(port, protocol, suid, chid, baudRate, para)
									log.Debugf("auto lock set value result: %s, %s", r, err)
								}()

								break
							}
						}

						log.Debugf("port {%s}, protocol {%s}, channel {%s}, value {%v} command result: {%s}", port, protocol, chid, p.Parameters, result)
						return result, nil
					}
				}
			}
		}
	}

	return "", errors.New("channel not found").As(muid, suid, chid)
}

// HandleRemoteCommand remote
func (s *appserver) HandleRemoteCommand(payload string) error {
	var p public.RemotePayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return err
	}

	if p.Phase != "executing" {
		return nil
	}

	// get param
	var ri public.RemoteInfo
	if err := json.Unmarshal([]byte(p.Parameters.Value), &ri); err != nil {
		return err
	}

	// execute, `ssh -i rsa_key -N -f -R port:localhost:22 root@remotehost`
	if ri.Operation == "start" {
		// start remote
		local := ri.Port + ":localhost:22"
		remote := "root@" + ri.Host

		cmd := exec.Command("ssh", "-y", "-i", "/root/rsa_key", "-N", "-f", "-R", local, remote)
		applog.LOG.Infof("execute command: %v", cmd.Args)

		if err := cmd.Run(); err != nil {
			applog.LOG.Warningf("command failed, errmsg {%v}", errors.As(err))
			return err
		}
	} else if ri.Operation == "stop" {
		// stop remote
		cmd := exec.Command("killall", "ssh")
		applog.LOG.Infof("execute command: %v", cmd.Args)

		if err := cmd.Run(); err != nil {
			applog.LOG.Warningf("command failed, errmsg {%v}", errors.As(err))
			return err
		}
	}

	applog.LOG.Info("remote command success")

	return nil
}

// HandleSelfCommand self
func (s *appserver) HandleSelfCommand(topic, payload, channel string) error {
	var p public.RemotePayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return err
	}

	if p.Phase != "executing" {
		return nil
	}

	applog.LOG.Infof("handle self command, topic {%v}, payload {%v}", topic, payload)

	start := public.ParseUTCTimeStamp(p.StartTime)
	truestart := time.Now()
	interval := time.Since(truestart)
	p.EndTime = public.TimeToSTring(start.Add(interval))
	p.Phase = "complete"

	appnet.ReplyCommand(topic, p)

	// send command
	port := "/dev/self"
	model := "hmu2000"
	result, err := appnet.SelfCommand(port, model, p.Parameters.Value, channel)
	if err != nil {
		// applog.LOG.Warning("self command failed, channel {%s}, value {%s}, errmsg {%v}", channel, p.Parameters.Value, err)
		applog.LOG.Warning(errors.As(err, channel, p.Parameters.Value))

		// interval := time.Since(truestart)
		// p.EndTime = public.TimeToSTring(start.Add(interval))
		// p.Phase = "error"

		// appnet.ReplyCommand(topic, p)

		// return public.RestartApp()

		return public.RestartApp(model, errors.New(public.RestartBySelfCommand))
	}

	applog.LOG.Infof("port {%s}, channel {%s}, value {%v} command result: {%s}", port, channel, p.Parameters.Value, result)

	// interval := time.Since(truestart)
	// p.EndTime = public.TimeToSTring(start.Add(interval))
	// p.Phase = "complete"

	// appnet.ReplyCommand(topic, p)

	return public.RestartApp(model, errors.New(public.RestartBySelfCommand))
}

// HandlePushMessage active push
func (s *appserver) HandlePushMessage(topic, payload string) error {
	var p public.MessagePayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return fmt.Errorf("parse payload failed: %v", err)
	}

	// sample unit is "_"
	if "_" == p.SampleUnitID {
		return s.HandleInitiativePushMessage(topic, p)
	}

	mu := sys.GetMonitoringUnitCfg()
	// find su
	for _, sp := range mu.SamplePorts {
		for _, su := range sp.SampleUnits {
			if p.SampleUnitID == su.ID {
				// get protocol
				if sp.Protocol == public.ProtocolSensorflow {
					return s.HandleInitiativePushMessage(topic, p)
				}
			}
		}
	}

	return appnet.PublishSampleValues(p)
}

// HandleInitiativePushMessage handle initiative message
func (s *appserver) HandleInitiativePushMessage(topic string, p public.MessagePayload) error {
	// find in map
	val, ok := s.InitiativePushValue[topic]
	if !ok {
		// not found, append
		var value initiativePushPayload
		value.data = p
		value.hasPushed = false

		s.InitiativePushValue[topic] = value
		if err := appnet.PublishSampleValues(p); err == nil {
			v := s.InitiativePushValue[topic]
			v.hasPushed = true
			s.InitiativePushValue[topic] = v
		}

		return nil
	}

	// found, get data and check whether has been change
	v := val.data.(public.MessagePayload)

	haschange := false
	switch v.Value.(type) {
	case int:
		last := v.Value.(int)
		now := p.Value.(int)
		if last != now {
			haschange = true
		}
	case float64:
		last := v.Value.(float64)
		now := p.Value.(float64)
		if last != now {
			haschange = true
		}
	case string:
		last := v.Value.(string)
		now := p.Value.(string)
		if last != now {
			haschange = true
		}
	}

	// change, set value
	if haschange {
		var value initiativePushPayload
		value.data = p
		value.hasPushed = false
		s.InitiativePushValue[topic] = value
	}

	for _, payload := range s.InitiativePushValue {
		if !payload.hasPushed {
			if err := appnet.PublishSampleValues(payload.data.(public.MessagePayload)); err == nil {
				payload.hasPushed = true
				s.InitiativePushValue[topic] = payload
			}
		}
	}

	return nil
}

// HandleSensorflowSetLEDS handle sensorflow set leds
func (s *appserver) HandleSensorflowSetLEDS(p public.CommandPayload) (string, error) {
	var para public.SensorLEDParameter

	bytepara, err := json.Marshal(p.Parameters)
	if err != nil {
		return "", fmt.Errorf("invalid parameters")
	}

	if err := json.Unmarshal(bytepara, &para); err != nil {
		return "", fmt.Errorf("invalid parameters")
	}

	// find port
	mu := sys.GetMonitoringUnitCfg()
	for _, sp := range mu.SamplePorts {
		if sp.Protocol == public.ProtocolSensorflow {
			sus := portmap[sp.ID]
			for _, su := range sus {

				t := reflect.TypeOf(para)
				v := reflect.ValueOf(para)

				num := t.NumField()
				for k := 0; k < num; k++ {
					if su.SU.Setting.Address != int32(k+1) {
						continue
					}

					val := v.Field(k).Interface().(public.SensorRGB)
					if (val.Red == 0) && (val.Green == 0) && (val.Blue == 0) {
						continue
					}

					su.CommonCommand(sp.Setting.Port, sp.Protocol, su.SU.ID, "led", 0, val)
				}
			}
		}
	}

	return "ok", nil
	// return "", fmt.Errorf("port not found")
}

// HandleSensorflowSync handle sensorflow sync
func (s *appserver) HandleSensorflowSync() (interface{}, error) {
	// find port
	mu := sys.GetMonitoringUnitCfg()
	for _, sp := range mu.SamplePorts {
		if sp.Protocol == public.ProtocolSensorflow {
			// find su
			sus := portmap[sp.ID]
			if len(sus) > 0 {
				su := sus[0]

				var para public.CommandParameter
				para.Mode = public.SensorCommandModeSync
				result, err := su.CommonCommand(sp.Setting.Port, sp.Protocol, su.SU.ID, "led", 0, para)
				if err != nil {
					return "", errors.As(err, sp.Setting.Port, sp.Protocol)
				}

				var resp public.SensorSync
				if err := json.Unmarshal([]byte(result), &resp); err != nil {
					return "", errors.As(err, result)
				}

				return resp, nil
			}
		}
	}

	return "", fmt.Errorf("port not found")
}

// HandleUpdateMUConfig update mu config
func (s *appserver) HandleUpdateMUConfig(p interface{}) error {
	d, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal mu config failed: %s", err)
	}

	var cfg []core.MonitoringUnit
	if err := json.Unmarshal(d, &cfg); err != nil {
		return fmt.Errorf("unmarshal mu config failed: %s", err)
	}

	if len(cfg) == 0 {
		return fmt.Errorf("do not have sample units")
	}

	oldMu := sys.GetMonitoringUnitCfg()

	// check version
	nv := cfg[0].Version
	cv := oldMu.Version

	// new version do not larger than current version
	if nv <= cv {
		return fmt.Errorf("new version [%s] do not larger than current version [%s]", nv, cv)
	}

	// save new config
	if err := sys.SaveMonitoringUnitCfg(&cfg[0].MonitoringUnit); err != nil {
		return fmt.Errorf("save new config failed: %s", err)
	}

	// exit for restart
	go func() {
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	return nil
}
