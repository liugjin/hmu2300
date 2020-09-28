package main

import (
	"clc.hmu/app/appmanager/applog"
	"clc.hmu/app/appmanager/applog/log"
	"github.com/op/go-logging"
)

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func main() {
	lFile, err := applog.OpenLogFile("./testing.log")
	if err != nil {
		panic(err)
	}
	defer lFile.Close()

	log.Debugf("debug %s", Password("secret"))
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
}
