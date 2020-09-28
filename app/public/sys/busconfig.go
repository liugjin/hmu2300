package sys

import (
	"log"
	"os"

	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/store/memcfg"
	"github.com/gwaylib/errors"
)

var (
	BusmanagerCfgPath = os.ExpandEnv(etc.Etc.String("public", "buscfg"))

	busCfg *BusConfig
)

// 加载并获取monitoring-unit.json的数据, 该数据由memcfg维护。
func GetBusManagerCfg() *BusConfig {
	// reading from cache
	if busCfg != nil {
		return busCfg
	}
	cfg := &BusConfig{}
	if err := memcfg.GetJsonCfg(BusmanagerCfgPath, cfg); err != nil {
		log.Panic(errors.As(err))
	}

	// default one hour
	if cfg.Web.Restart.Duration == 0 {
		cfg.Web.Restart.Duration = 3600
	}

	// default for 6 times
	if cfg.Web.Restart.Max == 0 {
		cfg.Web.Restart.Max = 6
	}

	if cfg.Cache.Directory == "" {
		cfg.Cache.Directory = "/mnt/sda1/cache/"
	}

	busCfg = cfg

	return busCfg
}

// 更新并存储monitoring-unit.json的数据, 该存储由memcfg维护。
func SaveBusManagerCfg(cfg *BusConfig) error {
	busCfg = cfg
	if err := memcfg.WriteJsonCfg(BusmanagerCfgPath, cfg); err != nil {
		return errors.As(err)
	}
	return nil
}

// 以下是详细的配置文件结构

// MQTTOption mqtt option
type MQTTOption struct {
	// Broker       string `json:"broker"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	ClientID     string `json:"clientid"`
	CleanSession bool   `json:"cleansession"`
	Store        string `json:"store"`
	Qos          byte   `json:"qos"`
}

// WebOption web option
type WebOption struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Port           string `json:"port"`
	MonitoringUnit struct {
		Path Path `json:"path"`
	} `json:"monitoringunit"`
	ElementLib struct {
		Server string `json:"server"`
		// Path   string `json:"path"` // move to etc.ini
	} `json:"elementlib"`
	Pages struct {
		Path Path `json:"path"` // Need to os.ExpandEnv
	} `json:"pages"`
	Video struct {
		Path Path `json:"path"` // Need to os.ExpandEnv
		Max  int  `json:"max"`  // limitation of cameras
		Sync struct {
			User string `json:"user"`
			Host string `json:"host"`
			Path Path   `json:"path"`
		}
	} `json:"video"`
	NetChecking struct {
		Timeout int `json:"timeout"`  // 单位秒
		DoTimes int `json:"do_times"` // 启动执行的次数，若需要长时间执行，请填写int32最大值
		// 如果配置数据为零，不启用网络自动检查功能, 若配置，自动增加mqtt路径进行并检测
		Hosts []string `json:"hosts"`
	} `json:"net_checking"`
	Restart struct {
		Duration int `json:"duration"` // means how long for software restart when disconnect (uint: s(second))
		Times    int `json:"times"`    // record software restart times
		Max      int `json:"max"`      // harware restart when software restart reaches max times
	} `json:"restart"`
}

// SystemServerOption system server option
type SystemServerOption struct {
	Uri  string `json:"uri"`
	Vals string `json:"vals"`
}

// CacheOption cache
type CacheOption struct {
	Directory  string `json:"dir"`
	MaxFile    int    `json:"maxFile"`
	MaxMessage int    `json:"maxMsg"`
}

// Signal signal
type Signal struct {
	Topic string `json:"topic"` // format 'suid/chid', like 'di1/val'
	Value string `json:"value"` // value
}

// CaptureSubOption capture sub option
type CaptureSubOption struct {
	SUID    string   `json:"su"`      // video's sample unit id
	Signals []Signal `json:"signals"` // signals want to set
}

// AutoLockSubOption auto lock sub option
type AutoLockSubOption struct {
	Topic    string `json:"topic"`    // format 'suid/chid', like 'di1/val'
	SetValue int    `json:"setvalue"` // value want to set
}

// BusConfig bus config
type BusConfig struct {
	MQTT          MQTTOption          `json:"mqtt"`
	Web           WebOption           `json:"web"`
	SystemServer  SystemServerOption  `json:"systemserver"`
	Model         string              `json:"model"`
	Cache         CacheOption         `json:"cache"`
	StartLogPath  Path                `json:"startlogpath"`
	CaptureOption []CaptureSubOption  `json:"capture"`
	AutoLock      []AutoLockSubOption `json:"autolock"`
}
