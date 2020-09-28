package sys

import (
	"log"
	"os"

	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/store/memcfg"
	"clc.hmu/app/public/store/muid"
	"github.com/gwaylib/errors"
)

var (
	// path for monitoring-units.json
	MonitoringUnitCfgPath = os.ExpandEnv(etc.Etc.String("public", "mucfg"))

	// ElementLibraryDir element library directory
	ElementLibDir = os.ExpandEnv(etc.Etc.String("public", "element-dir"))

	muCfg *MonitoringUnit
)

func SetMonitoringUnitCfgPath(filePath string) {
	MonitoringUnitCfgPath = filePath
}

// 加载并获取monitoring-unit.json的数据, 该数据由memcfg维护。
func GetMonitoringUnitCfg() *MonitoringUnit {
	// get from cache
	if muCfg != nil {
		return muCfg
	}

	mus := MUS{}
	if err := memcfg.GetJsonCfg(MonitoringUnitCfgPath, &mus); err != nil {
		log.Panic(errors.As(err))
		return nil
	}
	if len(mus) != 1 {
		panic(errors.New("Error configuration").As(len(mus)))
	}
	muCfg = &mus[0]

	// fix uuid
	sn := muid.GetMuID()
	if len(sn) > 0 {
		muCfg.ID = sn
	}
	return muCfg
}

// 更新并存储monitoring-unit.json的数据, 该存储由memcfg维护。
func SaveMonitoringUnitCfg(cfg *MonitoringUnit) error {
	muCfg = cfg
	if err := memcfg.WriteJsonCfg(MonitoringUnitCfgPath, &MUS{*cfg}); err != nil {
		return errors.As(err)
	}
	return nil
}

// =================================================================================
// 以下是监控单元配置的数据结构。
// =================================================================================

// MUS mu config
type MUS []MonitoringUnit

// MonitoringUnit monitoring unit
type MonitoringUnit struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Setting struct {
		Port     string `json:"port,omitempty"`
		BaudRate int32  `json:"baudRate,omitempty"`
	} `json:"setting"`
	SamplePorts []SamplePort `json:"ports"`
}

func (mu *MonitoringUnit) HasSamplePort(port string) bool {
	for _, sp := range mu.SamplePorts {
		if !sp.Enable {
			continue
		}
		if sp.Setting.Port == port {
			return true
		}
	}
	return false
}

func (mu *MonitoringUnit) GetSamplePort(port string) *SamplePort {
	for _, sp := range mu.SamplePorts {
		if !sp.Enable {
			continue
		}
		if sp.Setting.Port == port {
			return &sp
		}
	}
	panic("SamplePort Not Found:" + port)
}

type SamplePortSetting struct {
	// 注意，此字段每一个集单元必须写，用于端口寻址，且配置时应该是全局唯一的。
	// 格式如：rtsp://admin:hyiot123@192.168.1.64:554/Streaming/Channels/101; /dev/com1等
	Port string `json:"port"`

	QrcodePort     string `json:"qrcodePort,omitempty"`
	BaudRate       int32  `json:"baudRate,omitempty"`
	QrcodeBaudRate int32  `json:"qrcodeBaudRate,omitempty"`
	User           string `json:"user,omitempty"`
	Password       string `json:"password,omitempty"`
	ConnectStriong string `json:"connectString,omitempty"`

	// sensorflow
	KeyNumber     int32  `json:"keyNumber,omitempty"`
	WANInterface  string `json:"WANInterface,omitempty"`
	WifiInterface string `json:"WifiInterface,omitempty"`

	// lumi gateway, including password
	SID          string `json:"sid,omitempty"`
	NetInterface string `json:"netinterface,omitempty"`

	// face ipc, including port, user
	Host         string `json:"host,omitempty"`
	UploadServer string `json:"uploadServer,omitempty"` // upload param: server address
	Author       string `json:"author,omitempty"`       // upload param: author
	Project      string `json:"project,omitempty"`      // upload param: project
	Token        string `json:"token,omitempty"`        // upload param: token

	// weigeng entry
	LocalhostAddress string `json:"localhostAddress,omitempty"`
	LocalhostPort    string `json:"localhostPort,omitempty"`

	// for video
	ExecCmd string `json:"exec_cmd,omitempty"` // 执行的指令，空值默认为:ffmpeg, 填写时应该: ./ff
}

// SamplePort sample port
type SamplePort struct {
	ID          string             `json:"id"`
	Symbol      string             `json:"symbol"`
	Protocol    string             `json:"protocol"`
	Name        string             `json:"name"`
	Enable      bool               `json:"enable"`
	Setting     *SamplePortSetting `json:"setting"`
	SampleUnits []SampleUnit       `json:"sampleUnits"`
}

func (mu *SamplePort) GetSampleUnit(suid string) *SampleUnit {
	for _, su := range mu.SampleUnits {
		if !su.Enable {
			continue
		}

		if su.ID == suid {
			return &su
		}
	}
	panic("SampleUnit Not Found:" + suid)
}

type SampleUnitSetting struct {
	Address int32 `json:"address,omitempty"` // modbus同一端口的子地址

	// hmu
	Host  string `json:"host,omitempty"`
	Port  string `json:"port,omitempty"`
	Model string `json:"model,omitempty"`

	// camera
	IP       string `json:"ip,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	// credit card, including username and password
	SerialNumber string `json:"serialnum,omitempty"`

	// lumi gateway, including model
	SID string `json:"sid,omitempty"`

	// snmp, including port
	Version        string `json:"version,omitempty"`
	Target         string `json:"target,omitempty"`
	ReadCommunity  string `json:"readCommunity,omitempty"`
	WriteCommunity string `json:"writeCommunity,omitempty"`

	// weigeng entry
	DoorAddress string `json:"doorAddress,omitempty"` // door address
	DoorPort    string `json:"doorPort,omitempty"`    // door port
	SerialNo    string `json:"serialNo,omitempty"`    // device serial number

	// rstp视频
	Size     string `json:"size,omitempty"`            // 图片或视频的尺寸, 例如：640x480
	Time     string `json:"time,omitempty"`            // 截取的起始时间，例如：00:00:00
	Duration int    `json:"duration,string,omitempty"` // 截取的持续时间, 单位为秒(仅视频有效)，例如：5

	UploadUrl   string `json:"upload_url,omitempty"`   // 上传的地址，填写正式或测试服务器地址，例如:lab.huayuan-iot.com, 未填时默认为测试服地址
	UploadToken string `json:"upload_token,omitempty"` // 上传的地址，须为http上传接口, 格式如：上传服务器的token值, 不填写时使用系统内置默认值。

	// DidoLoc
	AutoLockDiValue string `json:"di_value,omitempty"`
	AutoLockDoValue string `json:"do_value,omitempty"`
	AutoLockNum     int32  `json:"num,omitempty"` // 设置modbus的状态读取长度
}

// SampleUnit sample unit
type SampleUnit struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Period                 int32  `json:"period"` // 需要注意，此值为0时，只会调用一遍，不会启动循环采集功能。
	Timeout                int32  `json:"timeout"`
	MaxCommunicationErrors int32  `json:"maxCommunicationErrors"`
	Element                string `json:"element"`
	Enable                 bool   `json:"enable"`

	Setting SampleUnitSetting `json:"setting"`
}
