/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: entry point
 *
 */

package portmanager

import (
	"net"
	"os"

	"clc.hmu/app/portmanager/portnet"
	pb "clc.hmu/app/portmanager/portpb"
	"clc.hmu/app/portmanager/protocol"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/store/etc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Start entry point
func Start() {
	port := etc.Etc.String("portmanager", "rpc_server")
	logpath := os.ExpandEnv(etc.Etc.String("portmanager", "plog"))
	log.Debug("PortManager start by:", port, logpath)

	// start listen
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// init port server
	var ps PortServer
	ps.Clients = make(map[string][]protocol.PortClient)
	ps.ProtocolMap = make(map[string]string)

	portnet.ConnectAppServer()
	defer portnet.DisconnectAppServer()

	// register port server
	s := grpc.NewServer()
	pb.RegisterPortServer(s, &ps)

	// register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
