/*
 *
 * Copyright 2019 huayuan-iot
 *
 * Author: shu
 * Date: 2019/06/19
 * Despcription: log
 *
 */
package log

import (
	"os"

	"github.com/op/go-logging"
)

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:.4s} %{id:03x} %{module} %{color:reset} %{message}`,
)

type Log struct {
	*logging.Logger

	module string

	consoleLog logging.LeveledBackend

	file    *LogFile
	fileLog logging.LeveledBackend
}

func (f *Log) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func NewLog(module string) *Log {
	l := logging.MustGetLogger(module)

	consoleLog := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
	consoleLog.SetLevel(logging.DEBUG, module) // Echo to console
	l.SetBackend(consoleLog)

	return &Log{
		Logger:     l,
		module:     module,
		consoleLog: consoleLog,
	}
}
