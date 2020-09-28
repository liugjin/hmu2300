/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: util
 *
 */

package public

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"clc.hmu/app/public/log/bootflag"
	"clc.hmu/app/public/log/elog"
	"clc.hmu/app/public/store"
	"clc.hmu/app/public/sys"
	"github.com/gwaylib/errors"
)

// restart result
const (
	AppBoot                         = "App start"
	CommandStart                    = "command start"
	RestartBySetMuID                = "Restart by set muid"
	RestartByUpgrade                = "Restart by upgrade"
	RestartByWEB                    = "Restart by web"
	RestartByMQTT                   = "Restart by MQTT command"
	RestartByCommunicationInterrupt = "Restart by net timeout"
	RestartByAlarmLinkage           = "The alarm linkage"
	RestartBySelfCommand            = "Restart by SelfCommand"
	RebootByWEB                     = "Reboot by web"
	RebootByCommunicationInterrupt  = "Reboot by net timeout"
)

// UTCTimeStamp utc now string
func UTCTimeStamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// TimeToSTring specified time to string
func TimeToSTring(timestamp time.Time) string {
	return timestamp.Format("2006-01-02T15:04:05.000Z")
}

// ParseUTCTimeStamp parse string to time
func ParseUTCTimeStamp(timestamp string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", timestamp)
	return t
}

func UpgradeApp(model, callPath string) error {
	reader := strings.NewReader(fmt.Sprintf(`{"Args":["%s"]}`, callPath))
	req, err := http.NewRequest("PUT", "http://127.0.0.1:9001/exec", reader)
	if err != nil {
		return errors.As(err, model, callPath)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.As(err, model, callPath)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.As(err, model, callPath)
	}
	if resp.StatusCode != 200 {
		return errors.As(err, model, callPath, resp.StatusCode, string(respData))
	}
	log.Println(string(respData))
	return nil
}

// RestartApp restart
func RestartApp(model string, err error) error {
	if err := bootflag.WriteFlag(); err != nil {
		return errors.As(err)
	}
	elog.LOG.Info(err)

	callPath := "/usr/local/clc.hmu/app/supd/supd"
	switch model {
	case sys.MODEL_DEFAULT, sys.MODEL_HMU2300, sys.MODEL_HMU2400:
		callPath = store.GetRootDir() + "/app/supd/supd"
	case sys.MODEL_HMU2200:
		callPath = "/usr/data/clc.hmu/app/supd/supd"
	}

	reader := strings.NewReader(fmt.Sprintf(`{"Args":["%s","ctl","restart","clc.hmu.app.aggregation"]}`, callPath))
	req, err := http.NewRequest("PUT", "http://127.0.0.1:9001/exec", reader)
	if err != nil {
		return errors.As(err, model, callPath)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.As(err, model, callPath)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.As(err, model, callPath)
	}
	if resp.StatusCode != 200 {
		return errors.As(err, model, callPath, resp.StatusCode, string(respData))
	}
	log.Println(string(respData))
	return nil
}

// Reboot reboot
func Reboot(err error) error {
	if err := bootflag.WriteFlag(); err != nil {
		return errors.As(err)
	}
	elog.LOG.Info(err)

	// execute reboot
	cmd := exec.Command("/sbin/reboot")
	if err := cmd.Run(); err != nil {
		return errors.As(err)
	}

	// exit current process
	os.Exit(0)

	return nil
}

// LegalTimeFormat legal time format
func LegalTimeFormat(t string) bool {
	// exp := `^(20|21|22|23|[0-1]\d):[0-5]\d:[0-5]\d$`
	exp := `^((20|21|22|23|[0-1]\d):[0-5]\d)|24:00$`
	reg, err := regexp.Compile(exp)
	if err != nil {
		return false
	}

	return reg.MatchString(t)
}

// LegalQueryTimeFormat query format
func LegalQueryTimeFormat(t string) bool {
	exp := `^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\s((20|21|22|23|[0-1]\d):[0-5]\d)|24:00$`
	reg, err := regexp.Compile(exp)
	if err != nil {
		return false
	}

	return reg.MatchString(t)
}

// TransferQueryTimeFormat transfer format
func TransferQueryTimeFormat(t string) string {
	if !LegalQueryTimeFormat(t) {
		return ""
	}

	day := t[:10]
	hour := t[11:13]
	minute := t[14:]

	return day + "T" + hour + "-" + minute + "-00"
}

// FileSHA256 sum file sha256
func FileSHA256(filePath string) (string, error) {
	var hashValue string
	file, err := os.Open(filePath)
	if err != nil {
		return hashValue, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}

	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)

	return hashValue, nil
}

// AppVersion get app version
func AppVersion(app string) (string, error) {
	cmd := exec.Command(app, "-v")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	l := len(output)
	if l == 0 {
		return "", fmt.Errorf("exec fail")
	}

	// trim line feed
	return string(output[:l-1]), nil
}

// uploadBody upload response
type uploadBody struct {
	Error      string `json:"err"`
	Enable     bool   `json:"enable"`
	Visible    bool   `json:"visible"`
	ID         string `json:"_id"`
	User       string `json:"user"`
	Project    string `json:"project"`
	Resource   string `json:"resource"`
	Name       string `json:"name"`
	Extension  string `json:"extension"`
	Type       string `json:"type"`
	Author     string `json:"author"`
	Path       string `json:"path"`
	Size       int    `json:"size"`
	UpdateTime string `json:"updatetime"`
	CreateTime string `json:"createtime"`
	Index      int    `json:"_index"`
	V          int    `json:"__v"`
}

// UploadFile upload file, request url such as:  http://lab.huayuan-iot.com/resource/upload/img/public/9H200A1710004_camera1_2018-09-06T15-16-20.mp4?author=admin&project=video&token=9de47be0-a28f-11e8-b7c6-3575c6f21b3d&user=admin
func UploadFile(filepath, filename, host, author, project, token, user string) (string, error) {
	requrlprefix := "http://" + host + "/resource/upload/"
	requrl := requrlprefix + "img/public/" + filename + "?author=" + author
	requrl += "&project=" + project
	requrl += "&token=" + token
	requrl += "&user=" + user

	buf := bytes.NewBufferString("")
	writer := multipart.NewWriter(buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		return "", errors.As(err, filepath)
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filepath)
	if err != nil {
		return "", errors.As(err, filepath)
	}
	defer fh.Close()

	// need to know the boundary to properly close the part myself.
	boundary := writer.Boundary()
	closebuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	reqreader := io.MultiReader(buf, fh, closebuf)
	fi, err := fh.Stat()
	if err != nil {
		return "", errors.As(err, filepath)
	}

	req, err := http.NewRequest("POST", requrl, reqreader)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(buf.Len()) + int64(closebuf.Len())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.As(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code").As(resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.As(err)
	}

	var b uploadBody
	if err := json.Unmarshal(body, &b); err != nil {
		return "", errors.As(err, string(body))
	}
	if len(b.Error) > 0 {
		return "", errors.New(b.Error).As(requrl, string(body))
	}

	return requrlprefix + b.Path, nil
}

// QueryInterfaceInfoByName query interface ip and mac
func QueryInterfaceInfoByName(name string) (string, string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", "", fmt.Errorf("query interface info failed: %v", err)
	}

	for _, i := range interfaces {
		if i.Name == name {
			addrs, err := i.Addrs()
			if err != nil {
				return "", "", fmt.Errorf("query address failed: %v", err)
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), i.HardwareAddr.String(), nil
					}
				}
			}

			return "", "", fmt.Errorf("query address failed, ipv4 not found")
		}
	}

	return "", "", fmt.Errorf("interface [%v] not found", name)
}

// HTTPDownloadFile download file
func HTTPDownloadFile(netpath, localpath string) error {
	// get data first
	resp, err := http.Get(netpath)
	if err != nil {
		return err
	}

	// check code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code abnormal: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// open or create local file
	f, err := os.OpenFile(localpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	// get data from body
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// save
	if _, err := f.Write(d); err != nil {
		return err
	}

	return nil
}

// OnvifHTTPCapture onvif http capture, url like: http://onvif:hyiot123@192.168.1.64/onvif-http/snapshot?
func OnvifHTTPCapture(host, user, password string) ([]byte, error) {
	url := "http://" + user + ":" + password + "@" + host + "/onvif-http/snapshot?"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get resource failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("properties: %s error reading response. %s", url, err)
	}

	return body, nil
}
