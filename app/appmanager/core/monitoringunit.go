/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: monitoring unit
 *
 */

package core

import (
	"encoding/json"
	"time"

	"clc.hmu/app/public"

	"clc.hmu/app/appmanager/appnet"
	"clc.hmu/app/public/log/applog"
	"clc.hmu/app/public/sys"
)

// MonitoringUnit monitoring unit
type MonitoringUnit struct {
	sys.MonitoringUnit
}

// PortMap port map
type PortMap map[string][]SampleUnitEx

// Start start
func (mu *MonitoringUnit) Start() PortMap {
	portmap := make(PortMap)

	for _, sp := range mu.SamplePorts {
		if !sp.Enable {
			applog.LOG.Infof("port [%s] has been disable", sp.ID)
			continue
		}

		sPort := &SamplePort{sp}

		sus, err := sPort.Start(mu.ID)
		if err != nil {
			applog.LOG.Warningf("start port failed, errmsg[%v]", err)
			continue
		}

		portmap[sp.ID] = sus

		applog.LOG.Infof("sample port [%s] start success", sp.ID)
	}

	if err := mu.SendDiscovery(); err != nil {
		applog.LOG.Warningf("send discovery failed, errmsg {%v}", err)

		go func() {
			for {
				if err := mu.SendDiscovery(); err != nil {
					time.Sleep(time.Second * 1)
					continue
				}

				applog.LOG.Info("send discovery success")
				break
			}
		}()
	}

	return portmap
}

// GivingStatus giving status
func (mu *MonitoringUnit) GivingStatus(value float64) error {
	payload := public.MessagePayload{
		MonitoringUnitID: mu.ID,
		SampleUnitID:     "_",
		ChannelID:        "_state",
		Name:             "采集器连接状态",
		Value:            value,
		Timestamp:        public.UTCTimeStamp(),
		Cov:              true,
		State:            0,
	}

	return appnet.PublishSampleValues(payload)
}

// SendDiscovery send discovery
func (mu *MonitoringUnit) SendDiscovery() error {
	prefix := "{"
	suffix := "}"
	sep := ","
	data := prefix

	for _, sp := range mu.SamplePorts {
		for _, su := range sp.SampleUnits {
			var d public.Discovery
			d.Name = su.Name
			d.Model = "hmu2000"
			d.Project = "hmu-manager"
			d.Station = "shenzhen"
			d.StationName = "shenzhen"
			d.Type = "hmu"
			d.User = "admin"
			d.Vendor = "clc"

			bd, err := json.Marshal(d)
			if err != nil {
				// applog.LOG.Warningf("marshal sample unit failed, errmsg {%v}", err)
				continue
			}

			data += "\"" + su.ID + "\":" + string(bd) + sep

		}
	}

	tmp := data[:len(data)-1]
	data = tmp + suffix

	return appnet.PublishDiscovery(mu.ID, data)
}
