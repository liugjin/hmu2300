/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: applog
 *
 */

package applog

import (
	"os"

	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/supd/bakfile"
)

// LOG log
var LOG = log.NewLog("applog")

func init() {
	if err := LOG.SetFile(
		os.ExpandEnv(etc.Etc.String("public/log", "alog")),
		bakfile.StrToSize(etc.Etc.String("public/log", "alog_max_size")),
		int(etc.Etc.Int64("public/log", "alog_bak_files")),
	); err != nil {
		// panic(err)
		log.Warning("app log set file failed: %s", err)
	}
}
