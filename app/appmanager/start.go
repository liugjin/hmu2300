/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: entry point
 *
 */

package appmanager

import (
	"io"
	"net"
	"os"
	"time"

	"clc.hmu/app/appmanager/appnet"
	"clc.hmu/app/appmanager/appnet/pmnet"
	pb "clc.hmu/app/appmanager/apppb"
	"clc.hmu/app/appmanager/core"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"clc.hmu/app/public/sys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var portmap core.PortMap

// Start true main
func Start() {
	eldir := sys.ElementLibDir
	port := etc.Etc.String("appmanager", "rpc_server")
	logpath := os.ExpandEnv(etc.Etc.String("appmanager", "alog"))

	log.Debug("AppManger start by:", port, eldir, logpath)

	// set heartbeat
	// if len(config.GlobalConfiguration) > 0 {
	// 	appnet.SetHeartbeat(config.GlobalConfiguration[0].ID, "")
	// }

	// connect to port manager
	pmnet.ConnectPortManager(etc.Etc.String("portmanager", "rpc_client")) // config.GlobalSetting.PortManager.Address)
	defer pmnet.DisconnectPortManager()

	// connect to bus manager
	appnet.ConnectBusManager(etc.Etc.String("busmanager", "rpc_client")) // config.GlobalSetting.BusManager.Address)
	defer appnet.DisconnectBusManager()

	// start monitor

	mu := &core.MonitoringUnit{
		*sys.GetMonitoringUnitCfg(),
	}
	pm := mu.Start()
	portmap = pm

	// update config
	go updateConfigFromSDCard()

	// Start server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var appser appserver
	appser.InitiativePushValue = make(map[string]initiativePushPayload)

	s := grpc.NewServer()
	pb.RegisterAppServer(s, &appser)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func updateConfigFromSDCard() {
	for {
		// check duration
		time.Sleep(time.Second * 3)

		// check file exist or not
		newfile := "/mnt/sda1/.config/monitoring-units.json"
		fi, err := os.Stat(newfile)
		if err != nil {
			// if os.IsNotExist(err) {
			// 	log.Printf("file %s not exist", newfile)
			// } else {
			// 	log.Printf("other error occured")
			// }

			continue
		}

		oldfile := "/usr/bin/aggregation/monitoring-units.json"
		ofi, err := os.Stat(oldfile)
		if err != nil {
			log.Printf("stat config file failed, errmsg {%v}", err)
			continue
		}

		newfilemodtime := fi.ModTime()
		oldfilemodtime := ofi.ModTime()
		if oldfilemodtime.String() < newfilemodtime.String() {
			log.Println("copy file")

			// copy file
			copyFile(newfile, oldfile)

			// then restart app
			// public.RestartApp()

			continue
		}

		// log.Println("file latest")
	}
}

func copyFile(source, dest string) bool {
	if source == "" || dest == "" {
		log.Println("source or dest is null")
		return false
	}

	fs, err := os.Open(source)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer fs.Close()

	fd, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer fd.Close()

	_, copyerr := io.Copy(fd, fs)
	if copyerr != nil {
		log.Println(copyerr.Error())
		return false
	}

	return true
}
