/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/04
 * Despcription: entry point
 *
 */

package busmanager

import (
	"net"
	"os"

	pb "clc.hmu/app/busmanager/buspb"
	"clc.hmu/app/busmanager/module"
	"clc.hmu/app/busmanager/module/web"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Start ture main
func Start() {
	port := etc.Etc.String("busmanager", "rpc_server")
	logpath := os.ExpandEnv(etc.Etc.String("busmanager", "blog"))
	log.Debug("BusManger start by:", port, logpath)

	// connect app server
	module.ConnectAppServer()
	defer module.DisconnectAppServer()

	//connect web server
	// connect app server
	web.ConnectAppServer()
	defer web.DisconnectAppServer()

	// start web server
	go web.StartRouter()

	// go module.MappingStaticFiles()

	// listen
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// new server
	var server module.BusServer
	server.Init()
	defer server.Cleanup()

	s := grpc.NewServer()
	pb.RegisterBusServer(s, &server)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
