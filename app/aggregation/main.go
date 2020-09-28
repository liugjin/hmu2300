/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/06/30
 * Despcription: aggregation
 *
 */

package main

import (
	"flag"
	"fmt"
	"os/signal"
	"time"

	"os"

	"clc.hmu/app/appmanager"
	"clc.hmu/app/busmanager"
	"clc.hmu/app/frp/frpc"
	"clc.hmu/app/portmanager"
	"clc.hmu/app/public"
	"clc.hmu/app/public/appver"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/applog"
	"clc.hmu/app/public/log/buslog"
	"clc.hmu/app/public/log/elog"
	"clc.hmu/app/public/log/portlog"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/supd/bakfile"
	"github.com/gwaylib/errors"
)

// build
var (
	BuildVersion = appver.BuildVersion
	BuildTime    = appver.BuildTime
	BuildName    = appver.BuildName
	CommitID     = appver.CommitID
)

func main() {

	showVerPtr := flag.Bool("v", false, "show build version")

	flag.Parse()

	if *showVerPtr {
		fmt.Println(BuildVersion)

		// exit
		os.Exit(0)
	}

	// print help
	fmt.Println("use '-h' or '--help' for more help")
	if err := log.SetFile(
		os.ExpandEnv(etc.Etc.String("public/log", "dlog")),
		bakfile.StrToSize(etc.Etc.String("public/log", "dlog_max_size")),
		int(etc.Etc.Int64("public/log", "dlog_bak_files")),
	); err != nil {
		// panic(err)
		log.Warning("log set file failed: %s", err)
	}

	log.Debug("App starting, [ctrl+c] to exit.")

	// start port
	go portmanager.Start()

	// start bus
	go busmanager.Start()

	// wait for a moment
	time.Sleep(time.Millisecond * 100)

	appmanager.Version = BuildVersion

	// start app
	go appmanager.Start()

	go func() {
		if err := frpc.Start(
			os.ExpandEnv(etc.Etc.String("public", "frpc_ini")),
		); err != nil {
			log.Warning(errors.As(err))
		}
	}()

	// record start log
	elog.LOG.Info(public.AppBoot)

	end := make(chan os.Signal, 2)
	signal.Notify(end, os.Interrupt, os.Kill)
	s := <-end

	elog.LOG.Infof("App down by signal:%s", s)
	elog.LOG.Close()
	portlog.LOG.Close()
	buslog.LOG.Close()
	applog.LOG.Close()
	log.Close()
}
