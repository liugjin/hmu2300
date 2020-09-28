/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/09/19
 * Despcription: video implement
 *
 */

package protocol

import (
	"encoding/json"
	"fmt"
	"os"

	"clc.hmu/app/public"
	"clc.hmu/app/public/ffmpeg"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/sys"
	"github.com/google/uuid"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolCameraRstp, generalCameraRstpDriverProtocol)
}

// Implement DriverProtocol
type cameraRstpDriverProtocol struct {
	uri     string
	suid    string // 采集单元的配置文件ID
	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit
}

func (dp *cameraRstpDriverProtocol) Payload() interface{} {
	return dp.suid
}

func (dp *cameraRstpDriverProtocol) ClientID() string {
	return dp.suid
}

func (dp *cameraRstpDriverProtocol) NewInstance() (PortClient, error) {
	return NewCameraRstpClient(dp.uri, dp.muCfg, dp.portCfg, dp.unitCfg)
}

func generalCameraRstpDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	cfg := sys.GetMonitoringUnitCfg()
	portCfg := cfg.GetSamplePort(uri)
	unitCfg := portCfg.GetSampleUnit(suid)
	return &cameraRstpDriverProtocol{
		uri:     uri,
		suid:    suid,
		muCfg:   cfg,
		portCfg: portCfg,
		unitCfg: unitCfg,
	}, nil
}

// ============= register driver end ==========================

// CameraRstpClient video client
type CameraRstpClient struct {
	clientID string
	uri      string

	muCfg   *sys.MonitoringUnit
	portCfg *sys.SamplePort
	unitCfg *sys.SampleUnit

	cmd *ffmpeg.FFMPEG
}

// NewCameraRstpClient new video client
func NewCameraRstpClient(uri string, muCfg *sys.MonitoringUnit, portCfg *sys.SamplePort, unitCfg *sys.SampleUnit) (PortClient, error) {
	cmd := portCfg.Setting.ExecCmd
	if len(cmd) == 0 {
		cmd = "ffmpeg"
	}
	return &CameraRstpClient{
		clientID: unitCfg.ID,
		uri:      uri,
		muCfg:    muCfg,
		portCfg:  portCfg,
		unitCfg:  unitCfg,
		cmd:      ffmpeg.NewCmd(cmd),
	}, nil
}

// Start start
func (vc *CameraRstpClient) Start() {
	log.Debug("CameraRstp Ignore CameraRstpClient Start:" + vc.clientID)
}

// ID client id
func (vc *CameraRstpClient) ID() string {
	return vc.clientID
}

// Sample get values
func (vc *CameraRstpClient) Sample(payload string) (string, error) {
	log.Debug("Ignore CameraRstpClient Sample:" + vc.clientID)
	return public.NewSamplePayload(false).Serial(), nil
}

func (vc CameraRstpClient) getCacheFile(fileName string) string {
	return os.TempDir() + "/" + fileName
}

func (vc *CameraRstpClient) removeCacheFile(fileName string) error {
	return errors.As(os.Remove(vc.getCacheFile(fileName)))
}

func (vc *CameraRstpClient) upload(fileName string) (string, error) {
	uploadFileName := vc.muCfg.ID + "_" + fileName
	host := vc.unitCfg.Setting.UploadUrl
	if len(host) == 0 {
		host = "lab.huayuan-iot.com"
	}
	author := "admin"
	project := "vidoe"
	token := vc.unitCfg.Setting.UploadToken
	if len(token) == 0 {
		token = "b4da0ed0-d1b7-11e8-b75e-435a751a1801"
	}
	user := "admin"
	file := vc.getCacheFile(fileName)
	dlUrl, err := public.UploadFile(
		file,
		uploadFileName, host, author, project, token, user,
	)
	if err != nil {
		return "", errors.As(err, fileName)
	}

	return dlUrl, nil
}

// CameraRstpOperationPayload camera operation payload
type CameraRstpOperationParam struct {
	Type string `json:"type"` // 支持两种:image(截图),video(截取)

	// 参数从配置文件, 详见public/monitoring-unit.go
}

// Command set values
func (vc *CameraRstpClient) Command(payload string) (string, error) {
	op := &CameraRstpOperationParam{}
	if err := json.Unmarshal([]byte(payload), op); err != nil {
		return "", errors.As(err, payload)
	}

	id := uuid.New().String()

	switch op.Type {
	case "image":
		fileName := id + ".jpg"
		filePath := vc.getCacheFile(fileName)

		if err := vc.cmd.FetchImage(
			vc.uri,
			filePath,
			vc.unitCfg.Setting.Size,
			"mjpeg",
			vc.unitCfg.Setting.Time,
			1,
		); err != nil {
			return "", errors.As(err)
		}
		defer vc.removeCacheFile(fileName)

		// upload
		dlUrl, err := vc.upload(fileName)
		if err != nil {
			return "", errors.As(err)
		}
		return dlUrl, nil

	case "video":
		fileName := id + ".mp4"
		if err := vc.cmd.FetchVideo(
			vc.uri,
			os.TempDir()+fmt.Sprintf("/%s", fileName),
			vc.unitCfg.Setting.Size,
			"mp4",
			vc.unitCfg.Setting.Time,
			vc.unitCfg.Setting.Duration,
			20,
		); err != nil {
			return "", errors.As(err)
		}
		defer vc.removeCacheFile(fileName)

		// upload
		dlUrl, err := vc.upload(fileName)
		if err != nil {
			return "", errors.As(err)
		}
		return dlUrl, nil

	}
	return "", errors.New("Unknow type").As(op.Type)
}
