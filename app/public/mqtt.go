/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: mqtt
 *
 */

package public

// phase
const (
	PhaseExcuting = "executing"
	PhaseError    = "error"
	PhaseTimeout  = "timeout"
	PhaseComplete = "complete"
)

// PublishPayload pubilsh payload
type PublishPayload struct {
	MonitoringUnitID string  `json:"monitoringUnitId"`
	SampleUnitID     string  `json:"sampleUnitId"`
	ChannelID        string  `json:"channelId"`
	Value            float64 `json:"value"`
	Timestamp        string  `json:"timestamp"`
	Cov              bool    `json:"cov"`
	State            int     `json:"state"`
}

// CommandPayload command payload
/**
{
	"monitoringUnit": "u-wei", 			//监控单元ID
	"sampleUnit": "rack-2", 				//采集单元ID
	"channel": "led-3", 				//控制通道ID
	"parameters": {					//控制参数
		"value": 0					//led灯状态
	},
	"phase": "executing",				//控制阶段
	"timeout": 1000,					//控制超时时间，单位：毫秒
	"operator": "admin",				//控制执行者
	"startTime": "2014-10-01T12:00:00Z",	//控制命令执行时间（UTC）
	"endTime": "2014-10-01T12:00:01Z", 	//控制命令完成时间（UTC），由控制器返回
	"retryTimes": 0,
	"result": "ok",						//执行结果，由控制器返回
	"_phase": "executing"
}
**/
type CommandPayload struct {
	MonitoringUnit string      `json:"monitoringUnit"`
	SampleUnit     string      `json:"sampleUnit"`
	Channel        string      `json:"channel"`
	Parameters     interface{} `json:"parameters"`
	Phase          string      `json:"phase"`
	Timeout        int         `json:"timeout"`
	Operator       string      `json:"operator"`
	StartTime      string      `json:"startTime"`
	EndTime        string      `json:"endTime"`
	Result         interface{} `json:"result"`
	UnderlinePhase string      `json:"_phase"`
}

// CommandParameter command parameter
type CommandParameter struct {
	// common
	Value interface{} `json:"value"`

	// sensorflow
	Red   int `json:"r"`
	Green int `json:"g"`
	Blue  int `json:"b"`

	// lamp with
	Mode       string `json:"mode"`
	ColorTable string `json:"colorTable"`

	// upgrade
	Path     string `json:"path"`
	CheckSum string `json:"checksum"`

	Channel string `json:"channel"`
}
