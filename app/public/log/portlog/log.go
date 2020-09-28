/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: portlog
 *
 */

package portlog

import (
	"os"

	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/supd/bakfile"
)

// LOG log
var LOG = log.NewLog("portlog")

func init() {
	if err := LOG.SetFile(
		os.ExpandEnv(etc.Etc.String("public/log", "plog")),
		bakfile.StrToSize(etc.Etc.String("public/log", "plog_max_size")),
		int(etc.Etc.Int64("public/log", "plog_bak_files")),
	); err != nil {
		// panic(err)
		log.Warning("port log set file failed: %s", err)
	}
}
