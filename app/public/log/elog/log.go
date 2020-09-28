package elog

import (
	"os"

	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/supd/bakfile"
)

// LOG log
var LOG = log.NewLog("elog")

func init() {
	if err := LOG.SetFile(
		os.ExpandEnv(etc.Etc.String("public/log", "elog")),
		bakfile.StrToSize(etc.Etc.String("public/log", "elog_max_size")),
		int(etc.Etc.Int64("public/log", "elog_bak_files")),
	); err != nil {
		// panic(err)
		log.Warning("elog set file failed: %s", err)
	}
}
