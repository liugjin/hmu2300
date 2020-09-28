/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: buslog
 *
 */

package buslog

import (
	"os"

	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/supd/bakfile"
)

// LOG log
var LOG = log.NewLog("buslog")

func init() {
	if err := LOG.SetFile(
		os.ExpandEnv(etc.Etc.String("public/log", "blog")),
		bakfile.StrToSize(etc.Etc.String("public/log", "blog_max_size")),
		int(etc.Etc.Int64("public/log", "blog_bak_files")),
	); err != nil {
		// panic(err)
		log.Warning("bus log set file failed: %s", err)
	}
}
