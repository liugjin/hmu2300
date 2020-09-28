/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: web handler implement
 *
 */

package web

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"

	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/sys"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gwaylib/errors"
	"fmt"
)

// message definition
var (
	MSGOK                  = "ok"
	MSGPONG                = "pong"
	MSGInvalidJSON         = "invalid json"
	MSGWriteFileFail       = "write file failed"
	MSGInvalidParam        = "invalid param"
	MSGReadPublicKeyFail   = "read public key failed"
	MSGUnknownErrorOccured = "unknown error occured"

	MSGMonitoringUnitIDRepeat = "monitoring unit id repeat"
	MSGMonitoringUnitNotFound = "monitoring unit not found"
	MSGSamplePortIDRepeat     = "sample port id repeat"
	MSGSamplePortRepeat       = "sample port repeat"
	MSGSamplePortNotFound     = "sample port not found"
	MSGSampleUnitIDRepeat     = "sample unit id repeat"
	MSGSampleUnitNotFound     = "sample unit not found"

	MSGQueryElementLibraryListFail = "query element library list failed"
	MSGReadElementLibraryFail      = "read element library failed"
	MSGElementLibraryIDRepeat      = "element library id repeat"
	MSGElementLibraryNotFound      = "element library not found"

	MSGConnectSystemServerFail = "connect system server fail"
	MSGQueryWifiConfigFail     = "query wifi config fail"
	MSGSetWifiModeFail         = "set wifi mode fail"
	MSGSetLTEModeFail          = "set lte mode fail"
	MSGSetETHStaticFail        = "set eth static mode fail"
	MSGSetETHDHCPFail          = "set eth dhcp mode fail"
	MSGSetAPModeFail           = "set ap mode fail"
	MSGSetLANFail              = "set lan fail"
	MSGQueryMUIDFail           = "query mu id fail"
	MSGRebootFail              = "reboot fail"
	MSGFactoryResetFail        = "factory reset fail"
	MSGQueryInternetConfigFail = "query internet config fail"
	MSGQueryLANConfigFail      = "query lan config fail"

	MSGLoginFail       = "login fail"
	MSGUserInvalid     = "invalid user"
	MSGPasswordInvalid = "invalid password"
	MSGHasNotLogin     = "has not login"
	MSGHasExpired      = "has expired"

	MSGCameraHasExist  = "camera has exist"
	MSGCameraNotFound  = "camera not found"
	MSGReachUpperLimit = "number of camera reach the upper limit"
	MSGRTSPHasBeenUsed = "rtsp url has been used"

	MSGRestartRtspClientFail = "restart rtsp client failed"

	MSGUnknownRecordType      = "unknown record operation type"
	MSGIllegalStartTime       = "illegal start time format"
	MSGIllegalEndTime         = "illegal end time"
	MSGIllegalStorageLimit    = "illegal storage limit"
	MSGReadRecordFileListFail = "read record file list failed"
	MSGIllegalTimeFormat      = "illegal time format"

	MSGSaveFrpConfigFail = "save frp config failed"
)

// code definition
var (
	ErrOK                  = "0"
	ErrInvalidJSON         = "101"
	ErrWriteFileFail       = "102"
	ErrInvalidParam        = "103"
	ErrReadPublicKeyFail   = "104"
	ErrUnknownErrorOccured = "105"

	ErrMonitoringUnitIDRepeat = "201"
	ErrMonitoringUnitNotFound = "202"
	ErrSamplePortIDRepeat     = "301"
	ErrSamplePortRepeat       = "301"
	ErrSamplePortNotFound     = "302"
	ErrSampleUnitIDRepeat     = "401"
	ErrSampleUnitNotFound     = "402"

	ErrQueryElementLibraryListFail = "501"
	ErrReadElementLibraryFail      = "502"
	ErrElementLibraryIDRepeat      = "503"
	ErrElementLibraryNotFound      = "504"

	ErrConnectSystemServerFail = "601"
	ErrQueryWifiConfigFail     = "602"
	ErrSetWifiModeFail         = "603"
	ErrSetLTEModeFail          = "604"
	ErrSetETHStaticFail        = "605"
	ErrSetETHDHCPFail          = "606"
	ErrSetAPModeFail           = "607"
	ErrSetLANFail              = "608"
	ErrQueryMUIDFail           = "609"
	ErrRebootFail              = "610"
	ErrFactoryResetFail        = "611"
	ErrQueryInternetConfigFail = "612"
	ErrQueryLANConfigFail      = "613"

	ErrLoginFail       = "701"
	ErrPasswordInvalid = "702"
	ErrUserInvalid     = "703"
	ErrHasNotLogin     = "704"
	ErrHasExpired      = "705"

	ErrCameraHasExist  = "801"
	ErrCameraNotFound  = "802"
	ErrReachUpperLimit = "803"
	ErrRTSPHasBeenUsed = "804"

	ErrRestartRtspClientFail = "901"

	ErrUnknownRecordType      = "1001"
	ErrIllegalStartTime       = "1002"
	ErrIllegalEndTime         = "1003"
	ErrIllegalStorageLimit    = "1004"
	ErrReadRecordFileListFail = "1005"
	ErrIllegalTimeFormat      = "1006"

	ErrSaveFrpConfigFail = "1101"
)

// file name
// var (
// 	FileNameApp = config.Configuration.Web.MonitoringUnit.Path
//
// // FileNameBus = "busmanager.json"
// )

// NormalGinH normal response
func NormalGinH(msg, status string) gin.H {
	return gin.H{"message": msg, "status": status}
}

// DataGinH response with data
func DataGinH(msg, status string, data interface{}) gin.H {
	return gin.H{"message": msg, "status": status, "data": data}
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandlePing ping
func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, NormalGinH(MSGPONG, ErrOK))
}

// HandleLogin login
func HandleLogin(c *gin.Context) {
	var req loginReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	if req.Username == cfg.Web.Username && req.Password == cfg.Web.Password {
		// add cookie
		session := sessions.Default(c)
		deadline := time.Now().Add(validityPeriod).Format(TimeFormatLayout)
		session.Set(CookieExpireTime, deadline)
		session.Save()

		c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
	} else {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGLoginFail, ErrLoginFail))
	}
}

type modifyPasswdReq struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

// HandleModifyPassword change password
func HandleModifyPassword(c *gin.Context) {
	var req modifyPasswdReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}
	cfg := sys.GetBusManagerCfg()

	if req.Username != cfg.Web.Username {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGUserInvalid, ErrUserInvalid))
		return
	}

	if req.OldPassword != cfg.Web.Password {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGPasswordInvalid, ErrPasswordInvalid))
		return
	}

	cfg.Web.Password = req.NewPassword

	// write to file
	if err := sys.SaveBusManagerCfg(cfg); err != nil {
		buslog.LOG.Warningf("write file failed, filename {%s}, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

type opSamplePortReq struct {
	MUID string `json:"muid"`
	sys.SamplePort
}

// HandleAddSamplePort add sp
func HandleAddSamplePort(c *gin.Context) {
	var req opSamplePortReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// sample port id not found
	if req.SamplePort.ID == "" {
		log.Debug(errors.New("SamplePort.ID Invalid"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	if req.SamplePort.Setting.Port == "" {
		log.Debug(errors.New("SamplePort.Setting.Port Invalid"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// find in config
	mu := sys.GetMonitoringUnitCfg()
	// check sample port exist or not
	if mu.HasSamplePort(req.SamplePort.Setting.Port) {
		log.Debug(errors.New("sp.Port == req.SamplePort.Setting.Port").As(req.SamplePort.Setting.Port))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortIDRepeat, ErrSamplePortIDRepeat))
		return
	}

	for _, sp := range mu.SamplePorts {
		if sp.ID == req.SamplePort.ID {
			log.Debug(errors.New("sp.ID == req.SamplePort.ID").As(sp.ID, req.SamplePort.ID))
			c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortIDRepeat, ErrSamplePortIDRepeat))
			return
		}
	}

	// add to mu port
	mu.SamplePorts = append(mu.SamplePorts, req.SamplePort)

	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		log.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*mu}))
}

// HandleModifySamplePort modify sp
func HandleModifySamplePort(c *gin.Context) {
	var req opSamplePortReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// sample port id not found
	if req.SamplePort.ID == "" {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}
	if req.SamplePort.Setting.Port == "" {
		log.Debug(errors.New("SamplePort.Setting.Port Invalid"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}
	// find in config
	mu := sys.GetMonitoringUnitCfg()

	// check sample port exist or not
	findsp := false
	spIdx := -1
	for isp, sp := range mu.SamplePorts {
		if sp.ID == req.SamplePort.ID {
			findsp = true
			req.SamplePort.SampleUnits = sp.SampleUnits
			spIdx = isp
		} else {
			// 遍历其他端口是否被占用, 已占用的会冲突
			if sp.Enable && sp.Setting.Port == req.SamplePort.Setting.Port {
				c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortRepeat, ErrSamplePortRepeat))
				return
			}
		}
	}
	if !findsp {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// modify mu port
	mu.SamplePorts[spIdx] = req.SamplePort
	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*mu}))
}

// HandleDeleteSamplePort delete sp
func HandleDeleteSamplePort(c *gin.Context) {
	var req opSamplePortReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// sample port id not found
	if req.SamplePort.ID == "" {
		log.Debug(errors.New("Missing ID"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// find in config
	mu := sys.GetMonitoringUnitCfg()
	// check sample port exist or not
	findsp := false
	length := len(mu.SamplePorts)
	for isp, sp := range mu.SamplePorts {
		if sp.ID == req.SamplePort.ID {
			findsp = true

			// delete mu port
			if (isp + 1) == length {
				mu.SamplePorts = append(mu.SamplePorts[:isp])
			} else {
				mu.SamplePorts = append(mu.SamplePorts[:isp], mu.SamplePorts[isp+1:]...)
			}
			break
		}
	}

	if !findsp {
		log.Debug(errors.New("Not Found Sample Port").As(req.SamplePort.ID, mu.SamplePorts))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		log.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*mu}))
}

type opSampleUnitReq struct {
	MUID string `json:"muid"`
	SPID string `json:"spid"`
	sys.SampleUnit
}

// HandleAddSampleUnit add su
func HandleAddSampleUnit(c *gin.Context) {
	var req opSampleUnitReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// sample port id not found
	if req.SPID == "" {
		log.Debug(errors.New("spid not found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// sample unit id not found
	if req.SampleUnit.ID == "" {
		log.Debug(errors.New("SampleUnit.ID not found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSampleUnitNotFound, ErrSampleUnitNotFound))
		return
	}

	mu := sys.GetMonitoringUnitCfg()
	findsp := false
	for isp, sp := range mu.SamplePorts {
		if sp.ID == req.SPID {
			findsp = true

			// check su exit or not
			for _, su := range sp.SampleUnits {
				if su.ID == req.SampleUnit.ID {
					log.Debug(errors.New("Same SampleUnit ID"))
					c.JSON(http.StatusBadRequest, NormalGinH(MSGSampleUnitIDRepeat, ErrSampleUnitIDRepeat))
					return
				}
			}

			mu.SamplePorts[isp].SampleUnits = append(mu.SamplePorts[isp].SampleUnits, req.SampleUnit)
		}
	}

	if !findsp {
		log.Debug(errors.New("SamplePort Not Found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		log.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*mu}))
}

// HandleModifySampleUnit modify su
func HandleModifySampleUnit(c *gin.Context) {
	var req opSampleUnitReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// sample port id not found
	if req.SPID == "" {
		log.Debug(errors.New("spid not found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// sample unit id not found
	if req.SampleUnit.ID == "" {
		log.Debug(errors.New("SampleUnit.ID not found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSampleUnitNotFound, ErrSampleUnitNotFound))
		return
	}

	mu := sys.GetMonitoringUnitCfg()
	findsp := false
	for isp, sp := range mu.SamplePorts {
		if sp.ID == req.SPID {
			findsp = true

			// check su exit or not
			findsu := false
			for isu, su := range sp.SampleUnits {
				if su.ID == req.SampleUnit.ID {
					findsu = true
					mu.SamplePorts[isp].SampleUnits[isu] = req.SampleUnit
					break
				}
			}

			if !findsu {
				log.Debug(errors.New("SamplePort Not Found"))
				c.JSON(http.StatusBadRequest, NormalGinH(MSGSampleUnitNotFound, ErrSampleUnitNotFound))
				return
			}

			// has found
			break
		}
	}

	if !findsp {
		log.Debug(errors.New("SamplePort Not Found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		log.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*mu}))
}

// HandleDeleteSampleUnit delete su
func HandleDeleteSampleUnit(c *gin.Context) {
	var req opSampleUnitReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// sample port id not found
	if req.SPID == "" {
		log.Debug(errors.New("spid not found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// sample unit id not found
	if req.SampleUnit.ID == "" {
		log.Debug(errors.New("SampleUnit.ID not found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSampleUnitNotFound, ErrSampleUnitNotFound))
		return
	}

	mu := sys.GetMonitoringUnitCfg()
	findsp := false
	for isp, sp := range mu.SamplePorts {
		if sp.ID == req.SPID {
			findsp = true
			findsu := false
			length := len(sp.SampleUnits)
			for isu, su := range sp.SampleUnits {
				if su.ID == req.SampleUnit.ID {
					findsu = true

					if isu == length {
						mu.SamplePorts[isp].SampleUnits = mu.SamplePorts[isp].SampleUnits[:isu]
					} else {
						mu.SamplePorts[isp].SampleUnits = append(mu.SamplePorts[isp].SampleUnits[:isu], mu.SamplePorts[isp].SampleUnits[isu+1:]...)
					}
					break
				}
			}

			if !findsu {
				log.Debug(errors.New("SampleUnit Not Found"))
				c.JSON(http.StatusBadRequest, NormalGinH(MSGSampleUnitNotFound, ErrSampleUnitNotFound))
				return
			}

			// has found
			break
		}
	}

	if !findsp {
		log.Debug(errors.New("SamplePort Not Found"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGSamplePortNotFound, ErrSamplePortNotFound))
		return
	}

	// write to file
	if err := sys.SaveMonitoringUnitCfg(mu); err != nil {
		log.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, sys.MUS{*mu}))
}

// HandleGetProtocolLibrary query protocol library
func HandleGetProtocolLibrary(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol != "" {
		c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, []string{}))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, public.ProtocolList))
}

// ELResp element library response
type ELResp []string

// HandleGetElementLibraryList get element library list
func HandleGetElementLibraryList(c *gin.Context) {
	dir := sys.ElementLibDir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Warning(errors.As(err, dir))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryElementLibraryListFail, ErrQueryElementLibraryListFail))
		return
	}

	var resp ELResp
	for _, v := range files {
		resp = append(resp, v.Name())
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleGetElementLibrary get specified element library
func HandleGetElementLibrary(c *gin.Context) {
	elname := c.Param("elname")

	suffix := ".json"
	elpath := ""
	dir := sys.ElementLibDir
	if strings.HasSuffix(elname, suffix) {
		elpath = dir + elname
	} else {
		elpath = dir + elname + suffix
	}

	el, err := OpenElementFile(elpath)
	if err != nil {
		log.Debug(errors.As(err, elpath))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGReadElementLibraryFail, ErrReadElementLibraryFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, el))
}

type opElementLibraryReq struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Version     string  `json:"version"`
	Description string  `json:"description"`
	Cov         float64 `json:"cov"`
	Channels    []struct {
		ChannelID string      `json:"chid"`
		Name      string      `json:"name"`
		DataType  string      `json:"type"`
		Value     interface{} `json:"value"`

		// modbus
		Code     int32       `json:"code,omitempty"`
		Address  int32       `json:"address,omitempty"`
		Quantity int32       `json:"quantity,omitempty"`
		Format   interface{} `json:"format,omitempty"`

		// pmbus
		CID1    int32 `json:"cid1,omitempty"`
		CID2    int32 `json:"cid2,omitempty"`
		Command int32 `json:"command,omitempty"`
		Offset  int   `json:"offset,omitempty"`
		Length  int   `json:"length,omitempty"`

		// snmp
		OID string `json:"oid,omitempty"`

		Expression string  `json:"expression"`
		Cov        float64 `json:"cov"`
	} `json:"channels"`
}

func reqToEL(elPath string, req opElementLibraryReq) public.Element {
	var el public.Element
	el.ID = req.ID
	el.Name = req.Name
	el.Path = elPath + el.ID + ".json"
	el.Type = req.Type
	el.Version = req.Version
	el.Description = req.Description

	switch el.Type {
	case public.ElementTypeModbus:
		// modbus serial
		var elmap public.Mapping
		elmap.Protocol = public.ProtocolModbusSerial
		elmap.Setting.COV = req.Cov
		elmap.Type = public.ProtocolTypeModbusSerial

		for _, ch := range req.Channels {
			// channel array
			var elch public.Channel
			elch.ID = ch.ChannelID
			elch.Name = ch.Name
			elch.DataType = ch.DataType
			elch.Value = ch.Value

			el.Channels = append(el.Channels, elch)

			// mapping array
			var chmap public.ChannelMapping
			chmap.ChannelID = ch.ChannelID
			chmap.Code = ch.Code
			chmap.Address = ch.Address
			chmap.Quantity = ch.Quantity
			chmap.Format = ch.Format
			chmap.Expression = ch.Expression
			chmap.COV = ch.Cov
			elmap.ChannnelMappings = append(elmap.ChannnelMappings, chmap)
		}

		el.Mappings = append(el.Mappings, elmap)

		// modbus tcp
		elmap.Protocol = public.ProtocolModbusTCP
		elmap.Setting.COV = req.Cov
		elmap.Type = public.ProtocolTypeModbusTCPClient
		elmap.ChannnelMappings = []public.ChannelMapping{}

		for _, ch := range req.Channels {
			var chmap public.ChannelMapping
			chmap.ChannelID = ch.ChannelID
			chmap.Code = ch.Code
			chmap.Address = ch.Address
			chmap.Quantity = ch.Quantity
			chmap.Format = ch.Format
			chmap.Expression = ch.Expression
			chmap.COV = ch.Cov
			elmap.ChannnelMappings = append(elmap.ChannnelMappings, chmap)
		}

		el.Mappings = append(el.Mappings, elmap)

	case public.ElementTypePMBus:
		// pmbus serial
		var elmap public.Mapping
		elmap.Protocol = public.ProtocolPMBUS
		elmap.Setting.COV = req.Cov
		elmap.Type = public.ProtocolTypePMBUS

		for _, ch := range req.Channels {
			// channel array
			var elch public.Channel
			elch.ID = ch.ChannelID
			elch.Name = ch.Name
			elch.DataType = ch.DataType
			elch.Value = ch.Value

			el.Channels = append(el.Channels, elch)

			// mapping array
			var chmap public.ChannelMapping
			chmap.ChannelID = ch.ChannelID
			chmap.CID1 = byte(ch.CID1)
			chmap.CID2 = byte(ch.CID2)
			chmap.COMMAND = uint16(ch.Command)
			chmap.Offset = ch.Offset
			chmap.Length = ch.Length
			chmap.Format = ch.Format
			chmap.Expression = ch.Expression
			chmap.COV = ch.Cov
			elmap.ChannnelMappings = append(elmap.ChannnelMappings, chmap)
		}

		el.Mappings = append(el.Mappings, elmap)

	case public.ElementTypeOilMachine:
		var elmap public.Mapping
		elmap.Protocol = public.ProtocolOilMachine
		elmap.Setting.COV = req.Cov
		elmap.Type = public.ProtocolTypeOilMachine

		for _, ch := range req.Channels {
			// channel array
			var elch public.Channel
			elch.ID = ch.ChannelID
			elch.Name = ch.Name
			elch.DataType = ch.DataType
			elch.Value = ch.Value

			el.Channels = append(el.Channels, elch)

			// mapping array
			var chmap public.ChannelMapping
			chmap.ChannelID = ch.ChannelID
			chmap.CID1 = byte(ch.CID1)
			chmap.CID2 = byte(ch.CID2)
			chmap.COMMAND = uint16(ch.Command)
			chmap.Offset = ch.Offset
			chmap.Length = ch.Length
			chmap.Format = ch.Format
			chmap.Expression = ch.Expression
			chmap.COV = ch.Cov
			elmap.ChannnelMappings = append(elmap.ChannnelMappings, chmap)
		}

		el.Mappings = append(el.Mappings, elmap)

	case public.ElementTypeHMU:
		// hmu
		var elmap public.Mapping
		elmap.Protocol = public.ProtocolHYIOTMU
		elmap.Setting.COV = req.Cov
		elmap.Type = public.ProtocolTypeHMU

		for _, ch := range req.Channels {
			// channel array
			var elch public.Channel
			elch.ID = ch.ChannelID
			elch.Name = ch.Name
			elch.DataType = ch.DataType
			elch.Value = ch.Value

			el.Channels = append(el.Channels, elch)

			// mapping array
			var chmap public.ChannelMapping
			chmap.ChannelID = ch.ChannelID
			elmap.ChannnelMappings = append(elmap.ChannnelMappings, chmap)
		}

		el.Mappings = append(el.Mappings, elmap)
	case public.ElementTypeSNMP:
		// snmp
		var elmap public.Mapping
		elmap.Protocol = public.ProtocolSNMP
		elmap.Setting.COV = req.Cov
		elmap.Type = public.ProtocolTypeSNMP

		for _, ch := range req.Channels {
			// channel array
			var elch public.Channel
			elch.ID = ch.ChannelID
			elch.Name = ch.Name
			elch.DataType = ch.DataType
			elch.Value = ch.Value

			el.Channels = append(el.Channels, elch)

			// mapping array
			var chmap public.ChannelMapping
			chmap.ChannelID = ch.ChannelID
			chmap.OID = ch.OID
			chmap.Format = ch.Format
			chmap.Expression = ch.Expression
			chmap.COV = ch.Cov
			elmap.ChannnelMappings = append(elmap.ChannnelMappings, chmap)
		}

		el.Mappings = append(el.Mappings, elmap)
	}

	return el
}

// HandleAddElementLibrary add element library
func HandleAddElementLibrary(c *gin.Context) {
	var req opElementLibraryReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// check id
	if req.ID == "" {
		log.Debug(errors.New("req.ID NOT FOUND"))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGElementLibraryNotFound, ErrElementLibraryNotFound))
		return
	}

	// transform and add
	el := reqToEL(sys.ElementLibDir, req)
	if err := NewFile(el.Path, el); err != nil {
		log.Warning(errors.As(err, el.Path))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, el))
}

// HandleModifyElementLibrary modify element library
func HandleModifyElementLibrary(c *gin.Context) {
}

// HandleDeleteElementLibrary delete element library
func HandleDeleteElementLibrary(c *gin.Context) {
	elname := c.Param("elname")
	filepath := sys.ElementLibDir + elname + ".json"
	if err := os.Remove(filepath); err != nil {
		log.Warning(errors.As(err, filepath))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

type opMqttConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// HandleGetMQTTConfig get mqtt config
func HandleGetMQTTConfig(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, opMqttConfig{Host: cfg.MQTT.Host, Port: cfg.MQTT.Port}))
}

// HandleModifyMQTTConfig modify mqtt config
func HandleModifyMQTTConfig(c *gin.Context) {
	var req opMqttConfig
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()
	cfg.MQTT.Host = req.Host
	cfg.MQTT.Port = req.Port

	// write to file
	if err := sys.SaveBusManagerCfg(cfg); err != nil {
		log.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	// write to video file
	videoconfig.MQTT.Host = req.Host
	videoconfig.MQTT.Port = req.Port
	if err := WriteConfigToFile(cfg.Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
		buslog.LOG.Warningf("modify video's mqtt config failed, errmsg {%v}", errors.As(err))
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, req))
}

type ntpReq struct {
	NTP1 string `json:"ntp1"`
	NTP2 string `json:"ntp2"`
	NTP3 string `json:"ntp3"`
}

// HandleGetNTP get ntp
func HandleGetNTP(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.Time()
	if err != nil {
		buslog.LOG.Warningf("get ntp config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGRebootFail, ErrRebootFail))
		return
	}

	var ntp ntpReq

	ns := strings.Split(resp.TimeServer, " ")
	l := len(ns)
	if l > 3 {
		l = 3
	}

	switch l {
	case 1:
		ntp.NTP1 = ns[0]
	case 2:
		ntp.NTP1 = ns[0]
		ntp.NTP2 = ns[1]
	case 3:
		ntp.NTP1 = ns[0]
		ntp.NTP2 = ns[1]
		ntp.NTP3 = ns[2]
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, ntp))
}

// HandleModifyNTP modify ntp
func HandleModifyNTP(c *gin.Context) {
	var req ntpReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	ts := req.NTP1 + " " + req.NTP2 + " " + req.NTP3

	resp, err := client.SetTimeServer(ts)
	if err != nil {
		buslog.LOG.Warningf("set ntp config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGRebootFail, ErrRebootFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

type setAPReq struct {
	SSID       string `json:"ssid"`
	Encryption string `json:"encryption"`
	Key        string `json:"key"`
	Channel    string `json:"channel"`
	Hide       string `json:"hide"`
}

// HandleGetAPConfig get ap config
func HandleGetAPConfig(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.Wireless()
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetAPModeFail, ErrSetAPModeFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleSetAPMode set ap
func HandleSetAPMode(c *gin.Context) {
	var req setAPReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetAP(req.SSID, req.Encryption, req.Key, req.Channel, req.Hide)
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetAPModeFail, ErrSetAPModeFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleGetInternetConfig get internet config
func HandleGetInternetConfig(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.Internet()
	if err != nil {
		buslog.LOG.Warningf("get internet config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryInternetConfigFail, ErrQueryInternetConfigFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleGetWifiConfig get wifi config
func HandleGetWifiConfig(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.Wireless()
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryWifiConfigFail, ErrQueryWifiConfigFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

type setWifiReq struct {
	SSID string `json:"ssid"`
	Key  string `json:"key"`
}

// HandleSetWifiConfig set wifi
func HandleSetWifiConfig(c *gin.Context) {
	var req setWifiReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetWifi(req.SSID, req.Key)
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetWifiModeFail, ErrSetWifiModeFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

type opLANReq struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

// HandleGetLANConfig get LAN
func HandleGetLANConfig(c *gin.Context) {
	if mac.LANIP != "" {
		c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, mac))
		return
	}
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.LAN()
	if err != nil {
		buslog.LOG.Warningf("get lan config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryLANConfigFail, ErrQueryLANConfigFail))
		return
	}

	mac.LANIP = resp.LANIP
	mac.LANMask = resp.LANMask

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleSetLANConfig set LAN
func HandleSetLANConfig(c *gin.Context) {
	var req opLANReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetLAN(req.IP, req.Mask)
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetLANFail, ErrSetLANFail))
		return
	}

	// update
	mac.LANIP = req.IP
	mac.LANMask = req.Mask

	// 已在不需要修改js的ip地址
	// if err := modifyJSHost(req.IP); err != nil {
	// 	buslog.LOG.Warningf("modify js host failed, errmsg {%v}", errors.As(err))
	// }

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleSetLTE set lte
func HandleSetLTE(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetLTE()
	if err != nil {
		buslog.LOG.Warningf("set lte mode failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetLTEModeFail, ErrSetLTEModeFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleGetWANConfig get WAN
func HandleGetWANConfig(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.WAN()
	if err != nil {
		buslog.LOG.Warningf("get lan config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryLANConfigFail, ErrQueryLANConfigFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

type setWANStaticReq struct {
	IP      string `json:"ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
	PDNS    string `json:"pdns"`
	SDNS    string `json:"sdns"`
}

// HandleSetWANStatic set WAN static mode
func HandleSetWANStatic(c *gin.Context) {
	var req setWANStaticReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetEthStatic(req.IP, req.Mask, req.Gateway, req.PDNS, req.SDNS)
	if err != nil {
		buslog.LOG.Warningf("set eth static mode failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetETHStaticFail, ErrSetETHStaticFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleSetWANDHCP set WAN dhcp mode
func HandleSetWANDHCP(c *gin.Context) {

	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SetEthDHCP()
	if err != nil {
		buslog.LOG.Warningf("get wifi config failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSetETHDHCPFail, ErrSetETHDHCPFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleFactoryReset factory reset
func HandleFactoryReset(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.FactoryReset("all")
	if err != nil {
		buslog.LOG.Warningf("factory reset failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGFactoryResetFail, ErrFactoryResetFail))
		return
	}

	// parse data
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, resp))
}

// HandleGetVideoInfo get video config
func HandleGetVideoInfo(c *gin.Context) {
	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, videoconfig.Cameras))
}

// HandleGetVideoByCameraID get by id
func HandleGetVideoByCameraID(c *gin.Context) {
	id := c.Param("id")

	for _, camera := range videoconfig.Cameras {
		// find
		if id == camera.StreamName {
			c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, camera))
			return
		}
	}

	log.Debug(errors.New("NOT FOUND").As(id))
	// not found
	c.JSON(http.StatusBadRequest, NormalGinH(MSGCameraNotFound, ErrCameraNotFound))
}

// HandleAddCamera add cemera
func HandleAddCamera(c *gin.Context) {
	// check number of cameras
	cfg := sys.GetBusManagerCfg()
	if len(videoconfig.Cameras) >= cfg.Web.Video.Max {
		log.Debug(errors.New("NOT MATCH").As(len(videoconfig.Cameras), cfg.Web.Video.Max))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGReachUpperLimit, ErrReachUpperLimit))
		return
	}

	var req public.Camera
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// id not found
	if req.StreamName == "" {
		log.Debug(errors.ErrNoData.As(req.StreamName))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGCameraNotFound, ErrCameraNotFound))
		return
	}

	// if stream name exist
	for _, camera := range videoconfig.Cameras {
		if camera.StreamName == req.StreamName {
			log.Debug(errors.New("Same stream:").As(camera.StreamName, req.StreamName))
			c.JSON(http.StatusBadRequest, NormalGinH(MSGCameraHasExist, ErrCameraHasExist))
			return
		}
	}

	// if rtsp url has been used
	for _, camera := range videoconfig.Cameras {
		if camera.RtspURL == req.RtspURL {
			log.Debug(errors.New("Same stream:").As(camera.RtspURL, req.RtspURL))
			c.JSON(http.StatusBadRequest, NormalGinH(MSGRTSPHasBeenUsed, ErrRTSPHasBeenUsed))
			return
		}
	}

	req.CameraID = req.StreamName
	videoconfig.Cameras = append(videoconfig.Cameras, req)

	// write to file
	if err := WriteConfigToFile(cfg.Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
		buslog.LOG.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	// restart rtsp client
	_, err := RestartRtspClient()
	if err != nil {
		buslog.LOG.Warningf("restart rtsp client failed, errmsg {%v}", errors.As(err))
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, videoconfig.Cameras))
}

// HandleModifyCamera modify camera
func HandleModifyCamera(c *gin.Context) {
	var req public.Camera
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// if rtsp url has been used
	for _, camera := range videoconfig.Cameras {
		if camera.RtspURL == req.RtspURL && camera.CameraID != req.CameraID {
			log.Debug(errors.New("NOT MATCH").As(camera, req))
			c.JSON(http.StatusBadRequest, NormalGinH(MSGRTSPHasBeenUsed, ErrRTSPHasBeenUsed))
			return
		}
	}

	cfg := sys.GetBusManagerCfg()
	// if exist
	for i, camera := range videoconfig.Cameras {
		if camera.CameraID == req.CameraID {
			// change request cameraid
			req.CameraID = req.StreamName

			// modify
			videoconfig.Cameras[i] = req

			// write to file
			if err := WriteConfigToFile(cfg.Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
				buslog.LOG.Warning(errors.As(err))
				c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
				return
			}

			// restart rtsp client
			_, err := RestartRtspClient()
			if err != nil {
				buslog.LOG.Warningf("restart rtsp client failed, errmsg {%v}", errors.As(err))
			}

			c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, videoconfig.Cameras))
			return
		}
	}

	log.Debug(errors.New("MsgCameraNotFound"))
	// not found
	c.JSON(http.StatusBadRequest, NormalGinH(MSGCameraNotFound, ErrCameraNotFound))
}

// HandleDeleteCamera delete camera
func HandleDeleteCamera(c *gin.Context) {
	id := c.Param("id")

	// if exist
	length := len(videoconfig.Cameras)
	cfg := sys.GetBusManagerCfg()
	for i, camera := range videoconfig.Cameras {
		if camera.CameraID == id {
			if (i + 1) == length {
				videoconfig.Cameras = videoconfig.Cameras[:i]
			} else {
				videoconfig.Cameras = append(videoconfig.Cameras[:i], videoconfig.Cameras[i+1:]...)
			}

			// write to file
			if err := WriteConfigToFile(cfg.Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
				buslog.LOG.Warning(errors.As(err))
				c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
				return
			}

			// restart rtsp client
			_, err := RestartRtspClient()
			if err != nil {
				buslog.LOG.Warningf("restart rtsp client failed, errmsg {%v}", errors.As(err))
			}

			c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, videoconfig.Cameras))
			return
		}
	}

	log.Debug(errors.New("NOT FOUND"))
	// not found
	c.JSON(http.StatusBadRequest, NormalGinH(MSGCameraNotFound, ErrCameraNotFound))
}

// get record true file time
func getRecordFileTime(name string) string {
	// file name like '9H200A1710001_camera1_2018-08-02T05-20-24.mp4'
	length := len(name)
	if length < 23 {
		return ""
	}

	filetime := name[length-23 : length-4]
	return filetime
}

// record filename satisfy condition
func verifyRecordFileName(name, starttime, endtime string) bool {
	if !strings.HasSuffix(name, "mp4") {
		return false
	}

	// ignore time
	if starttime == "" || endtime == "" {
		return true
	}

	// with time
	if public.LegalQueryTimeFormat(starttime) && public.LegalQueryTimeFormat(endtime) {
		// transfer starttime and endtime to specify format like '2018-08-02T05-20-00'
		start := public.TransferQueryTimeFormat(starttime)
		end := public.TransferQueryTimeFormat(endtime)

		// get record file time
		filetime := getRecordFileTime(name)

		if start <= filetime && end >= filetime {
			return true
		}
	}

	return false
}

// HandleQueryRecordFileByID query record file
func HandleQueryRecordFileByID(c *gin.Context) {
	cameraid := c.Param("id")
	rootdir := "/mnt/sda1/record/"
	dir := rootdir + cameraid + "/"

	// format like "2018-08-01 00:01"
	starttime := c.Query("startTime")
	endtime := c.Query("endTime")

	// if !public.LegalQueryTimeFormat(starttime) || !public.LegalQueryTimeFormat(endtime) {
	// 	c.JSON(http.StatusBadRequest, NormalGinH(MSGIllegalTimeFormat, ErrIllegalTimeFormat))
	// 	return
	// }

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		buslog.LOG.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGReadRecordFileListFail, ErrReadRecordFileListFail))
		return
	}

	// filter
	filelist := []string{}
	for _, f := range files {
		name := f.Name()
		if verifyRecordFileName(name, starttime, endtime) {
			filelist = append(filelist, name)
		}
	}

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, filelist))
}

type modifyRecordInfoReq struct {
	CameraID string `json:"cameraId"`
	Record   struct {
		Enable    bool   `json:"enable"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"record"`
	Sync struct {
		Enable    bool   `json:"enable"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"sync"`
}

func verifyModifyRecordInfoRequest(req modifyRecordInfoReq) (string, string) {
	if !public.LegalTimeFormat(req.Record.StartTime) {
		return MSGIllegalStartTime, ErrIllegalStartTime
	}

	if !public.LegalTimeFormat(req.Record.EndTime) {
		return MSGIllegalEndTime, ErrIllegalEndTime
	}

	if !public.LegalTimeFormat(req.Sync.StartTime) {
		return MSGIllegalStartTime, ErrIllegalStartTime
	}

	if !public.LegalTimeFormat(req.Sync.EndTime) {
		return MSGIllegalEndTime, ErrIllegalEndTime
	}

	return MSGOK, ErrOK
}

// HandleModifyRecordInfoByID modify record info
func HandleModifyRecordInfoByID(c *gin.Context) {
	var req modifyRecordInfoReq
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// verify request
	// errmsg, errcode := verifyModifyRecordInfoRequest(req)
	// if errcode != ErrOK {
	// 	c.JSON(http.StatusBadRequest, NormalGinH(errmsg, errcode))
	// 	return
	// }

	// find camera
	for i, camera := range videoconfig.Cameras {
		if camera.CameraID == req.CameraID {
			// record
			if !req.Record.Enable {
				videoconfig.Cameras[i].Record.Mode = public.RecordModeDisable
			} else if req.Record.StartTime == "00:00" && req.Record.EndTime == "24:00" {
				videoconfig.Cameras[i].Record.Mode = public.RecordModeWholeDay
			} else {
				videoconfig.Cameras[i].Record.Mode = public.RecordModeEnable
				videoconfig.Cameras[i].Record.StartTime = req.Record.StartTime
				videoconfig.Cameras[i].Record.EndTime = req.Record.EndTime
			}

			// sync
			if !req.Sync.Enable {
				videoconfig.Cameras[i].Sync.Mode = public.RecordModeDisable
			} else if req.Sync.StartTime == "00:00" && req.Sync.EndTime == "24:00" {
				videoconfig.Cameras[i].Sync.Mode = public.RecordModeWholeDay
			} else {
				videoconfig.Cameras[i].Sync.Mode = public.RecordModeEnable
				videoconfig.Cameras[i].Sync.StartTime = req.Sync.StartTime
				videoconfig.Cameras[i].Sync.EndTime = req.Sync.EndTime
			}

			// write to file
			cfg := sys.GetBusManagerCfg()
			if err := WriteConfigToFile(cfg.Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
				buslog.LOG.Warning(errors.As(err))
				c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
				return
			}

			// restart video
			RestartRtspClient()

			c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, videoconfig))
			return
		}
	}

	log.Debug(errors.New("NOT FOUND"))
	c.JSON(http.StatusBadRequest, NormalGinH(MSGCameraNotFound, ErrCameraNotFound))
}

type setStorageLimit struct {
	Limit string `json:"limit"`
}

type storageInfo struct {
	Total      string `json:"total"`
	RecordUsed string `json:"recordUsed"`
	OtherUsed  string `json:"otherUsed"`
	Limit      string `json:"limit"`
}

func sumFilesSize(dir string) int64 {
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		buslog.LOG.Warningf("get file list failed, dir {%v}, errmsg {%v}", dir, errors.As(err))
		return 0
	}

	sum := int64(0)
	for _, fi := range filelist {
		if fi.IsDir() {
			sum += sumFilesSize(dir + "/" + fi.Name())
		} else {
			sum += fi.Size()
		}
	}

	return sum
}

// HandleGetStorageInfo get storage info
func HandleGetStorageInfo(c *gin.Context) {
	cfg := sys.GetBusManagerCfg()

	// get info from hmu
	client := sys.ConnectSystemDaemon(cfg.Model, &cfg.SystemServer)
	defer client.Disconnect()

	resp, err := client.SystemInfo()
	if err != nil {
		buslog.LOG.Warningf("get system info failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGQueryWifiConfigFail, ErrQueryWifiConfigFail))
		return
	}

	total, _ := strconv.Atoi(resp.SDTotal)
	free, _ := strconv.Atoi(resp.SDFree)

	recordused := int(sumFilesSize("/mnt/sda1/record") / 1024 / 1024)
	otherused := total - free - recordused

	var r storageInfo
	r.Total = resp.SDTotal
	r.RecordUsed = strconv.Itoa(recordused)
	r.OtherUsed = strconv.Itoa(otherused)
	r.Limit = videoconfig.StorageLimit

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, r))
}

// HandleSetStorageLimit set storge limit
func HandleSetStorageLimit(c *gin.Context) {
	var req setStorageLimit
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	// if req.Limit < 0 {
	// 	c.JSON(http.StatusBadRequest, NormalGinH(MSGIllegalStorageLimit, ErrIllegalStorageLimit))
	// 	return
	// }

	videoconfig.StorageLimit = req.Limit

	// write to file
	if err := WriteConfigToFile(sys.GetBusManagerCfg().Web.Video.Path.ExpandEnv(), videoconfig); err != nil {
		buslog.LOG.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGWriteFileFail, ErrWriteFileFail))
		return
	}

	// restart video
	RestartRtspClient()

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

// HandleGetRemoteMapConfig get remote map config
func HandleGetRemoteMapConfig(c *gin.Context) {
	cfg := AgencyConfig.Read()

	c.JSON(http.StatusOK, DataGinH(MSGOK, ErrOK, cfg))
}

// HandleSetRemoteMapConfig set remote map config
func HandleSetRemoteMapConfig(c *gin.Context) {
	var req FrpcMapConfig
	if err := c.BindJSON(&req); err != nil {
		log.Debug(errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGInvalidJSON, ErrInvalidJSON))
		return
	}

	agencyPath := os.ExpandEnv(etc.Etc.String("public", "frpc_ini"))
	if err := AgencyConfig.SaveFrpConfig(req, agencyPath); err != nil {
		buslog.LOG.Warning(errors.As(err))
		c.JSON(http.StatusInternalServerError, NormalGinH(MSGSaveFrpConfigFail, ErrSaveFrpConfigFail))
		return
	}

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

// HandleRestartFrpc restart frpc
func HandleRestartFrpc(c *gin.Context) {
	RestartFrpc()

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

// RestartFrpc restart frpc
func RestartFrpc() error {
	cfg := sys.GetBusManagerCfg()

	if err := public.RestartApp(cfg.Model, errors.New("Restart Frpc")); err != nil {
		return errors.As(err)
	}

	return nil
}

// RestartRtspClient restart rtspclient
func RestartRtspClient() (string, error) {
	pro := "/etc/init.d/rtspclient"
	param := "restart"

	cmd := exec.Command(pro, param)
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// HandleRestartRtspClient restart rtspclient
func HandleRestartRtspClient(c *gin.Context) {
	result, err := RestartRtspClient()
	if err != nil {
		buslog.LOG.Errorf("restart rtsp client failed, errmsg {%v}", errors.As(err))
		c.JSON(http.StatusBadRequest, NormalGinH(MSGRestartRtspClientFail, ErrRestartRtspClientFail))
		return
	}

	buslog.LOG.Infof("restart rtsp client success, result {%s}", result)

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

// generateKey genarate key
func generateKey() (string, error) {
	// /usr/bin/dropbearkey -t rsa -s 1024 -f /root/rsa_key

	// execute generate private key
	cmd := exec.Command("/usr/bin/dropbearkey", "-t", "rsa", "-s", "1024", "-f", "/root/rsa_key")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// get public key, the second line
	rd := bufio.NewReader(bytes.NewBuffer(output))
	if _, err = rd.ReadString('\n'); err != nil || io.EOF == err {
		return "", err
	}

	line, err := rd.ReadString('\n')
	if err != nil || io.EOF == err {
		return "", err
	}

	key := line[:len(line)-1]
	return key, nil
}

// readKeyFromFile read key from file
func readKeyFromFile(filepath string) (string, error) {
	// /usr/bin/dropbearkey -y -f /root/rsa_key

	// execute read public key
	cmd := exec.Command("/usr/bin/dropbearkey", "-y", "-f", "/root/rsa_key")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// get public key, the second line
	rd := bufio.NewReader(bytes.NewBuffer(output))
	if _, err = rd.ReadString('\n'); err != nil || io.EOF == err {
		return "", err
	}

	line, err := rd.ReadString('\n')
	if err != nil || io.EOF == err {
		return "", err
	}

	key := line[:len(line)-1]
	return key, nil
}

// HandleUploadElementLibrary upload element library
func HandleUploadElementLibrary(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	for _, file := range files {
		err := c.SaveUploadedFile(file, sys.ElementLibDir+file.Filename)
		if err != nil {
			log.Printf("upload file failed: %v", file.Filename)
			continue
		}
	}

	c.JSON(http.StatusOK, NormalGinH(MSGOK, ErrOK))
}

var appclient AppWebClient

//实时数据中的控制
func CommmandForm(c *gin.Context){
	commandMuid := c.PostForm("commandMuid")
	commandUnitId := c.PostForm("commandUnitId")
	commandChannelId := c.PostForm("commandChannelId")
	commandKey := c.PostForm("commandKey")
	commandValue := c.PostForm("commandValue")
	commandType := c.PostForm("commandType")

	topic := "command/" + commandMuid + "/" + commandUnitId + "/" + commandChannelId

	var payload string

	if commandType == "string" {
		payload = "{\"phase\":\"executing\",\"monitoringUnit\":\"" + commandMuid + "\",\"sampleUnit\":\"" + commandUnitId + "\",\"channel\":\"" + commandChannelId + "\",\"parameters\":{\"" + commandKey + "\":\"" +
			commandValue + "\"},\"timeout\":5000,\"operator\":\"admin\",\"startTime\":\"2018-06-13T07:14:23.135Z\",\"retryTimes\":0,\"endTime\":null,\"result\":" +
			"null,\"_phase\":\"executing\"}"
	} else {
		payload = "{\"phase\":\"executing\",\"monitoringUnit\":\"" + commandMuid + "\",\"sampleUnit\":\"" + commandUnitId + "\",\"channel\":\"" + commandChannelId + "\",\"parameters\":{\"" + commandKey + "\":" +
			commandValue + "},\"timeout\":5000,\"operator\":\"admin\",\"startTime\":\"2018-06-13T07:14:23.135Z\",\"retryTimes\":0,\"endTime\":null,\"result\":" +
			"null,\"_phase\":\"executing\"}"
	}

	//fmt.Println("控制topic:",topic)
	fmt.Println("控制payload:",payload)


	// check appclient available
	if appclient.Client == nil {
		// buslog.LOG.Warningf("app client unavaliable, new client")

		appclient = NewAppWebClient()
		if appclient.Client == nil {
			// buslog.LOG.Warningf("new app client failed")
			return
		}

		buslog.LOG.Infof("new app client success")
	}

	//var testTopic = "command/9H200A2400093/lan2/ch183"
	//var testPayload = "{\"phase\":\"executing\",\"monitoringUnit\":\"9H200A2400093\",\"sampleUnit\":\"lan2\",\"channel\":\"ch183\",\"parameters\":{\"value\":\"OFF\"},\"timeout\":5000,\"operator\":\"admin\",\"startTime\":\"2018-06-13T07:14:23.135Z\",\"retryTimes\":0,\"endTime\":null,\"result\":null,\"_phase\":\"executing\"}"

	// send message to app
	if err := appclient.Notify(topic, payload); err != nil {
		buslog.LOG.Warningf("notify app server failed, errmsg {%v}\n", errors.As(err))
	}

	c.JSON(http.StatusOK,gin.H{
		"topic" : topic,
		"payload" : payload,
	})
}