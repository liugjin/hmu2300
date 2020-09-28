/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: video
 *
 */

package public

// record mode
const (
	RecordModeDisable  = "disable"
	RecordModeEnable   = "enable"
	RecordModeWholeDay = "wholeday"
)

// RecordInfo record info
type RecordInfo struct {
	Mode      string `json:"mode"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// SyncInfo cloud sync info
type SyncInfo struct {
	Mode      string `json:"mode"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// Camera camera info
type Camera struct {
	CameraID   string     `json:"cameraId"`
	CameraName string     `json:"cameraName"`
	RtspURL    string     `json:"rtspUrl"`
	StreamID   string     `json:"streamId"`
	ServerURL  string     `json:"serverUrl"`
	StreamName string     `json:"streamName"`
	RTMPURL    string     `json:"rtmpUrl"`
	HLSURL     string     `json:"hlsUrl"`
	UserName   string     `json:"userName"`
	Password   string     `json:"password"`
	Record     RecordInfo `json:"record"`
	Sync       SyncInfo   `json:"sync"`
}

// CloudVideo video info
type CloudVideo struct {
	MQTT struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"mqtt"`
	MUID         string   `json:"muid"`
	StorageLimit string   `json:"storageLimit"`
	Cameras      []Camera `json:"cameras"`
}
