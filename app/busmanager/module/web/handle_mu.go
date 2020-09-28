package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"io"
	"strings"

	"clc.hmu/app/public"
	"clc.hmu/app/public/appver"
	"clc.hmu/app/public/httptry"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/bootflag"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/log/elog"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/sys"
	"github.com/gin-gonic/gin"
	"github.com/gwaylib/errors"
)

// HandleGetConfig get mu config
func HandleGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*(sys.GetMonitoringUnitCfg())}))
}

// HandleGetMUID get muid
func HandleGetMUID(c *gin.Context) {
	if muid != "" {
		c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.HMUSystemResp{UUID: muid}))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.UUID()
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryMUIDFail, ErrQueryMUIDFail))
		return
	}

	muid = resp.UUID

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

type setUUIDReq struct {
	UUID string `json:"uuid"`
}

// HandleSetMUID set muid
func HandleSetMUID(c *gin.Context) {
	var req setUUIDReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetUUID(req.UUID)
	if err != nil {
		buslog.LOG.Warningf("system set uuid failed, errmsg {%v}", errors.As(err))
		// c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryMUIDFail, ErrQueryMUIDFail))
		// return
	}

	muid = req.UUID

	mu := sys.GetMonitoringUnitCfg()

	// modify mu-units json file's id, by default modify the first mu
	mu.ID = req.UUID

	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	// write to bus manager config
	cfg.MQTT.ClientID = req.UUID
	if err := sys.SaveBusManagerCfg(cfg); err != nil {
		buslog.LOG.Warningf("set mqtt clientid failed, errmsg {%v}", errors.As(err))
	}

	// write to video config
	videoconfig.MUID = req.UUID
	if err := WriteConfigToFile(cfg.Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
		buslog.LOG.Warningf("set video muid failed, errmsg {%v}", errors.As(err))
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))

	// restart app
	go func() {
		// sleep for one second, giving time for response
		time.Sleep(time.Second * 1)

		// then restart
		if err := public.RestartApp(cfg.Model, errors.New(public.RestartBySetMuID)); err != nil {
			buslog.LOG.Errorf("restart app failed, errmsg {%v}", errors.As(err))
			return
		}

		buslog.LOG.Infof("set uuid restart app success")
	}()
}

// HandleModifyMonitoringUnit modify mu
func HandleModifyMonitoringUnit(c *gin.Context) {
	var req sys.MonitoringUnit
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// write to file
	if err := sys.SaveMonitoringUnitCfg(&req); err != nil {
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{req}))
}

// HandleReboot reboot
func HandleReboot(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// log restart
	if err := bootflag.WriteFlag(); err != nil {
		log.Warning(errors.As(err))
		return
	}
	elog.LOG.Info(errors.New(public.RebootByWEB))

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.Reboot()
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGRebootFail, ErrRebootFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleRestart restart, actually exit the program, which watchdog wake up it
func HandleRestart(c *gin.Context) {
	// response ok
	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))

	go func() {
		time.Sleep(1e9)
		cfg := sys.GetBusManagerCfg()
		public.RestartApp(cfg.Model, errors.New(public.RestartByWEB))
	}()
}

// HandleGetPublicKey public key
func HandleGetPublicKey(c *gin.Context) {
	keypath := "/root/rsa_key"

	if _, err := os.Stat(keypath); err != nil {
		if os.IsNotExist(err) {
			// key file not exist
			key, err := generateKey()
			if err != nil {
				buslog.LOG.Errorf("generate private key failed, errmsg {%v}", errors.As(err))
				c.JSON(http.StatusInternalServerError, NormalGinH(MSGReadPublicKeyFail, ErrReadPublicKeyFail))
				return
			}

			c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, key))
			return
		}

		buslog.LOG.Error(errors.As(err))
		// other error occured
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGUnknownErrorOccured, ErrUnknownErrorOccured))
		return
	}

	// key file exist
	key, err := readKeyFromFile(keypath)
	if err != nil {
		buslog.LOG.Errorf("read public key failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGReadPublicKeyFail, ErrReadPublicKeyFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, key))
}

type UpgradeInfo struct {
	Ver         string `json:"ver"`
	DownloadUrl string `json:"dl_url"`
	CheckSum    string `json:"checksum"`
}

func HandleGetUpgrade(c *gin.Context) {
	upgradeUrl := etc.Etc.String("busmanager", "upgrade_url")
	busCfg := sys.GetBusManagerCfg()
	muCfg := sys.GetMonitoringUnitCfg()

	vals := url.Values{}
	vals.Add("uuid", muCfg.ID)

	resp, err := httptry.HttpsClient.PostForm(upgradeUrl+"/"+busCfg.Model, vals)
	if err != nil {
		log.Warning(errors.As(err, upgradeUrl, busCfg.Model))
		c.JSON(200, NormalGinH("升级服务器不可用", "500"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		log.Warning(errors.New("error upgrade protocol").As(resp.Status))
		c.JSON(200, NormalGinH("升级服务器不可用", "500"))
		return
	}
	if resp.StatusCode != 200 {
		log.Warning(errors.New("error upgrade protocol").As(resp.Status))
		c.JSON(200, NormalGinH("已是最新版本", "404"))
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warning(errors.As(err))
		c.JSON(200, NormalGinH("升级服务器不可用", "500"))
		return
	}
	info := &UpgradeInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		log.Warning(errors.As(err))
		c.JSON(200, NormalGinH("升级服务器不可用", "500"))
		return
	}
	if info.Ver == appver.BuildVersion {
		c.JSON(200, NormalGinH("已是最新版本", "404"))
		return
	}

	c.JSON(200, NormalGinH(fmt.Sprintf("发现新的版本(%s), 是否升级？", info.Ver), "200"))
	return
}

func HandlePostUpgrade(c *gin.Context) {
	// log restart to upgrade
	if err := bootflag.WriteFlag(); err != nil {
		log.Warning(errors.As(err))
		return
	}
	elog.LOG.Info(errors.New(public.RestartByUpgrade))

	callPath := os.ExpandEnv(etc.Etc.String("public", "upgrade_sh"))

	if err := public.UpgradeApp(sys.GetBusManagerCfg().Model, callPath); err != nil {
		log.Warning(errors.As(err))
		c.JSON(200, NormalGinH("升级服务不可用", "500"))
		return
	}

	// Do upgrade
	c.JSON(200, NormalGinH("正在升级，请稍候", "200"))
	return
}

func HandlePostChooseProjectForm(c *gin.Context) {
	scheme := c.PostForm("scheme1")
	//E:/git_local_dev/src/clc.hmu/app/aggregation/monitoring-units.json
	path := os.ExpandEnv(etc.Etc.String("public", "mucfg"))
	index := strings.LastIndex(path, "/")
	n1 := path[0:index+1] + "monitoring-units@v1.json" //电池柜方案
	n2 := path[0:index+1] + "monitoring-units@v2.json" //设备柜方案
	n3 := path[0:index+1] + "monitoring-units@v3.json" //双机柜方案
	n4 := path[0:index+1] + "temp.json"                //原始配置
	fileList := []string{n1, n2, n3, n4}

	var err error

	//如果temp.json不存在就不copy一份
	if boo, _ := pathExists(fileList[3]); !boo {
		_, err = copyFile(path, fileList[3])
	}

	switch scheme {
	case "1":
		os.Remove(path)
		_, err = copyFile(fileList[0], path)
	case "2":
		os.Remove(path)
		_, err = copyFile(fileList[1], path)
	case "3":
		os.Remove(path)
		_, err = copyFile(fileList[2], path)
	case "4":
		os.Remove(path)
		_, err = copyFile(fileList[3], path)
	}

	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
	}
	c.JSON(200, gin.H{
		"code": "0",
		"err":  err,
	})
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
