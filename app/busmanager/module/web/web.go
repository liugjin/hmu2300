/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: web router
 *
 */

package web

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"os/exec"
	"time"

	"clc.hmu/app/public"
	"clc.hmu/app/public/appver"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/store/memcfg"
	"clc.hmu/app/public/sys"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gwaylib/errors"
	"sync"
)

var videoconfig public.CloudVideo

var muid string
var mac sys.HMUSystemResp

// WSHub hub
var WSHub = newHub()

// AgencyConfig agency config
var AgencyConfig IniParser

// Status status
type Status struct {
	SUID  string `json:"suid"`
	State int    `json:"state"`
}

// DeviceStatus device's status, key: suid, value: state
var DeviceStatus map[string]int

//实时数据
var PayloadChan = make(chan []byte)
//实时数据
var PayloadMap sync.Map

// difinition
const (
	validityPeriod   = time.Second * 3600
	CookieExpireTime = "expiretime"
	TimeFormatLayout = "2006-01-02 15:04:05"
)

// StartRouter start router
func StartRouter() {
	// go SyncRecord()

	cfg := sys.GetBusManagerCfg()

	r := gin.Default()

	// make log to file.
	logFile, err := buslog.LOG.GetFile()
	if err != nil {
		log.Panic(errors.As(err))
		return
	}
	r.Use(gin.LoggerWithWriter(logFile))

	// cors config
	c := cors.DefaultConfig()
	c.AllowAllOrigins = true
	c.AllowCredentials = true
	c.AddAllowMethods("OPTIONS")
	r.Use(cors.New(c))

	// session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("token", store))

	agencyPath := os.ExpandEnv(etc.Etc.String("public", "frpc_ini"))
	// read ini file
	if err := AgencyConfig.Load(agencyPath); err != nil {
		buslog.LOG.Errorf("load agency config failed: %s", errors.As(err, agencyPath))
	}

	// write config to frpc config
	// muid := getUUID()
	// if len(muid) == 13 {
	// 	remoteport := muid[6:8] + muid[10:13]
	// 	AgencyConfig.AddFrpcAgencySection(sectionSSH, "127.0.0.1", "22", remoteport, agencyPath)
	// 	RestartFrpc()
	// }

	// ping
	r.GET("/ping", HandlePing)
	r.POST("/reset", HandleFactoryReset)

	// Register hmu-upgrade interface, and use port of webapi
	r.Any("/hacheck", func(c *gin.Context) {
		c.String(200, "1")
		return
	})
	r.Any("/version", func(c *gin.Context) {
		c.String(200, appver.BuildVersion)
		return
	})

	// static file
	path := cfg.Web.Pages.Path.ExpandEnv()
	r.StaticFS("/p/", http.Dir(path))

	user := r.Group("/user")
	{
		user.POST("/login", HandleLogin)
		user.PUT("/passwd", HandleModifyPassword)
	}

	// router
	mu := r.Group("/mu", AuthorizationMiddle)
	{
		mu.GET("/", HandleGetConfig)
		mu.GET("/id", HandleGetMUID)
		mu.POST("/id", HandleSetMUID)
		// mu.POST("/", HandleAddMonitoringUnit)
		mu.PUT("/", HandleModifyMonitoringUnit)
		// mu.DELETE("/", HandleDeleteMonitoringUnit)
		mu.PUT("/reboot", HandleReboot)
		mu.PUT("/restart", HandleRestart)
		mu.GET("/pubkey", HandleGetPublicKey)
		mu.GET("/upgrade", HandleGetUpgrade)
		mu.POST("/upgrade", HandlePostUpgrade)
		mu.POST("/chooseProjectForm",HandlePostChooseProjectForm)
	}

	sp := r.Group("/sp", AuthorizationMiddle)
	{
		sp.POST("/", HandleAddSamplePort)
		sp.PUT("/", HandleModifySamplePort)
		sp.DELETE("/", HandleDeleteSamplePort)
	}

	su := r.Group("/su", AuthorizationMiddle)
	{
		su.POST("/", HandleAddSampleUnit)
		su.PUT("/", HandleModifySampleUnit)
		su.DELETE("/", HandleDeleteSampleUnit)
	}

	pl := r.Group("/pl", AuthorizationMiddle)
	{
		pl.GET("/", HandleGetProtocolLibrary)
	}

	el := r.Group("/el", AuthorizationMiddle)
	{
		el.GET("/", HandleGetElementLibraryList)
		el.GET("/:elname", HandleGetElementLibrary)
		el.POST("/", HandleAddElementLibrary)
		el.PUT("/:elname", HandleModifyElementLibrary)
		el.DELETE("/:elname", HandleDeleteElementLibrary)
	}

	mqtt := r.Group("/mqtt", AuthorizationMiddle)
	{
		mqtt.GET("/", HandleGetMQTTConfig)
		mqtt.PUT("/", HandleModifyMQTTConfig)
	}

	ntp := r.Group("/ntp", AuthorizationMiddle)
	{
		ntp.GET("/", HandleGetNTP)
		ntp.PUT("/", HandleModifyNTP)
	}

	internet := r.Group("/internet", AuthorizationMiddle)
	{
		internet.GET("/", HandleGetInternetConfig)
	}

	wifi := r.Group("/wifi", AuthorizationMiddle)
	{
		wifi.GET("/", HandleGetWifiConfig)
		wifi.PUT("/", HandleSetWifiConfig)
	}

	ap := r.Group("/ap", AuthorizationMiddle)
	{
		ap.GET("/", HandleGetAPConfig)
		ap.PUT("/", HandleSetAPMode)
	}

	lan := r.Group("lan", AuthorizationMiddle)
	{
		lan.GET("/", HandleGetLANConfig)
		lan.PUT("/", HandleSetLANConfig)
	}

	lte := r.Group("lte", AuthorizationMiddle)
	{
		lte.PUT("/", HandleSetLTE)
	}

	eth := r.Group("eth", AuthorizationMiddle)
	{
		eth.GET("/", HandleGetWANConfig)
		eth.PUT("/static", HandleSetWANStatic)
		eth.PUT("/dhcp", HandleSetWANDHCP)
	}

	video := r.Group("video", AuthorizationMiddle)
	{
		video.GET("/", HandleGetVideoInfo)
		video.GET("/:id", HandleGetVideoByCameraID)
		video.POST("/", HandleAddCamera)
		video.PUT("/", HandleModifyCamera)
		video.DELETE("/:id", HandleDeleteCamera)
		video.PUT("/restart", HandleRestartRtspClient)
	}

	record := r.Group("record", AuthorizationMiddle)
	{
		record.GET("/:id", HandleQueryRecordFileByID)
		record.PUT("/", HandleModifyRecordInfoByID)

	}

	storage := r.Group("storage", AuthorizationMiddle)
	{
		storage.GET("/", HandleGetStorageInfo)
		storage.PUT("/", HandleSetStorageLimit)
	}

	mapcfg := r.Group("map", AuthorizationMiddle)
	{
		mapcfg.GET("/", HandleGetRemoteMapConfig)
		mapcfg.POST("/", HandleSetRemoteMapConfig)
		mapcfg.PUT("/restart", HandleRestartFrpc)
	}
	// web socket
	go WSHub.run()

	r.GET("/ws", func(c *gin.Context) {
		serveWs(WSHub, c.Writer, c.Request)
	})


	//实时数据
	r.GET("/real", func(c *gin.Context) {
		realTimeDataFunc(WSHub, c.Writer, c.Request)
	})

	r.POST("/commandForm",CommmandForm)
	//realTimeData := r.Group("real",AuthorizationMiddle)
	//{
	//	realTimeData.GET("/",func(c *gin.Context) {
	//		realTimeDataFunc(WSHub, c.Writer, c.Request)
	//	})
	//}

	// element lib files mapping
	ellib := r.Group("ellib", AuthorizationMiddle)
	{
		ellib.Static("/", os.ExpandEnv(etc.Etc.String("public", "element-dir")))
		ellib.POST("/", HandleUploadElementLibrary)
	}

	port := ":" + cfg.Web.Port
	r.Run(port)
}

// ReadVideoConfig read video config
func ReadVideoConfig() error {
	filename := sys.GetBusManagerCfg().Web.Video.Path.ExpandEnv()
	if err := memcfg.GetJsonCfg(filename, &videoconfig); err != nil {
		return errors.As(err)
	}
	return nil
}

// OpenElementFile read mapping file
func OpenElementFile(filename string) (public.Element, error) {
	var el public.Element

	if err := memcfg.GetJsonCfg(filename, &el); err != nil {
		return el, errors.As(err)
	}
	return el, nil
}

// WriteConfigToFile write config
func WriteConfigToFile(filename string, data interface{}) error {
	log.Debug("Write ConfigFile:", filename, data)

	// 更新内存中的配置文件，并写入配置文件
	return errors.As(memcfg.WriteJsonCfg(filename, data))
}

// NewFile new file
func NewFile(filename string, data interface{}) error {
	log.Debug("New ConfigFile:", filename, data)
	// 更新内存中的配置文件，并写入配置文件
	return errors.As(memcfg.WriteJsonCfg(filename, data))
}

// DeviceStatusToBytes device;'s status to bytes
func DeviceStatusToBytes() ([]byte, error) {
	var allstatus []Status
	for k, v := range DeviceStatus {
		status := Status{
			SUID:  k,
			State: v,
		}

		allstatus = append(allstatus, status)
	}

	bs, err := json.Marshal(allstatus)
	if err != nil {
		log.Printf("marshal data failed, err: %v", err)
		return nil, err
	}

	return bs, nil
}

// StartSync start
func StartSync(currentdir, remotedir string) error {
	if err := os.Setenv("RSYNC_RSH", "/usr/bin/ssh -i /root/rsa_key"); err != nil {
		log.Println("set env, err: ", err)
		return err
	}

	cmd := exec.Command("/usr/bin/rsync", "-r", currentdir, remotedir)
	log.Println(cmd.Args)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	return nil
}

// StopSync stop
func StopSync() error {
	cmd := exec.Command("/usr/bin/killall", "rsync")
	if _, err := cmd.Output(); err != nil {
		log.Println("killall error: ", err)
		return err
	}

	return nil
}

// SyncRunning is runnig or not
func SyncRunning() bool {
	cmd := exec.Command("/usr/bin/pgrep", "rsync")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return false
	}

	log.Println(string(output))

	if string(output) != "" {
		return true
	}

	return false
}

// SyncRecord sync record files
func SyncRecord() {
	retry := time.NewTicker(3 * time.Second)
	cfg := sys.GetBusManagerCfg()

	for {
		select {
		case <-retry.C:
			// check config
			for _, camera := range videoconfig.Cameras {
				if camera.Sync.Mode == public.RecordModeEnable {
					// enable, check current time
					currenttime := time.Now().Format("2006-01-02 15:04:05")
					starttime := currenttime[:11] + camera.Sync.StartTime + ":00"
					endtime := currenttime[:11] + camera.Sync.EndTime + ":00"

					if currenttime < starttime || currenttime > endtime {
						// stop sync
						log.Println("stop sync")
						StopSync()
					} else {
						// check
						if !SyncRunning() {
							log.Println("start sync")

							localdir := "/mnt/sda1/record/" + camera.CameraID
							remoteuser := cfg.Web.Video.Sync.User
							remotehost := cfg.Web.Video.Sync.Host
							remotepath := cfg.Web.Video.Sync.Path
							remotedir := remoteuser + "@" + remotehost + ":" + remotepath.ExpandEnv()

							StartSync(localdir, remotedir)
						}
					}
				} else if camera.Sync.Mode == public.RecordModeWholeDay {
					if !SyncRunning() {
						log.Println("start sync wholeday")

						localdir := "/mnt/sda1/record/" + camera.CameraID
						remoteuser := cfg.Web.Video.Sync.User
						remotehost := cfg.Web.Video.Sync.Host
						remotepath := cfg.Web.Video.Sync.Path
						remotedir := remoteuser + "@" + remotehost + ":" + remotepath.ExpandEnv()

						StartSync(localdir, remotedir)
					}
				}
			}
		}
	}
}

// AuthorizationMiddle auth middle
func AuthorizationMiddle(c *gin.Context) {
	// session := sessions.Default(c)
	// et := session.Get(CookieExpireTime)
	// if et == nil {
	// 	c.JSON(http.StatusBadRequest, NormalGinH(MSGHasNotLogin, ErrHasNotLogin))
	// 	c.Abort()
	// 	return
	// }

	// last := et.(string)
	// now := time.Now().Format(TimeFormatLayout)

	// if now > last {
	// 	c.JSON(http.StatusBadRequest, NormalGinH(MSGHasExpired, ErrHasExpired))
	// 	c.Abort()
	// 	return
	// }

	// // update deadline
	// deadline := time.Now().Add(validityPeriod).Format(TimeFormatLayout)
	// session.Set(CookieExpireTime, deadline)
	// session.Save()
}
