/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/05/17
 * Despcription: port server implement
 *
 */

package portmanager

import (
	"fmt"

	pb "clc.hmu/app/portmanager/portpb"
	"clc.hmu/app/portmanager/protocol"
	"clc.hmu/app/public"
	"clc.hmu/app/public/log"
	"clc.hmu/app/public/log/portlog"
	"github.com/gwaylib/errors"
	"golang.org/x/net/context"
)

// PortServer portserver is used to implement portmanager.PortServer.
type PortServer struct {
	Clients     map[string][]protocol.PortClient // one port binding multiple clients, for there may be multiple sample units in one port
	ProtocolMap map[string]string                // record port's mapping protocol, ex: "/dev/ttyS0" -> "protocol-modbus-serial"
}

// Binding implement binding interface
func (s *PortServer) Binding(ctx context.Context, in *pb.BindingRequest) (*pb.BindingReply, error) {
	// portlog.LOG.Infof("Binding request: %+v\n", in)

	drvName := in.Protocol
	uri := in.Port

	// Get driver template
	drv, err := protocol.GetDriverProtocol(drvName, in.Port, in.Suid, in.Payload)
	if err != nil {
		log.Warning(errors.As(err))
		errmsg := fmt.Sprintf("unimplement protocol [%v]", drvName)
		return &pb.BindingReply{Status: public.StatusErr, Message: errmsg}, nil
	}

	// Checking cache
	pid := drv.ClientID()
	if _, exist := s.HasClientExist(in.Port, pid); exist {
		errmsg := fmt.Sprintf("client has exist, id: %v", pid)
		portlog.LOG.Warning(errors.New(errmsg))
		return &pb.BindingReply{Status: public.StatusErr, Message: errmsg}, nil
	}

	// Make driver connect
	client, err := drv.NewInstance()
	if err != nil {
		errmsg := fmt.Sprintf("new %s client failed, errmsg [%v]", drvName, err)
		portlog.LOG.Warning(err)
		return &pb.BindingReply{Status: public.StatusErr, Message: errmsg}, nil
	}
	portlog.LOG.Infof("new client binding success, port [%s], protocol [%s], id [%v]", in.Port, in.Protocol, client.ID())

	// Make cache
	s.Clients[uri] = append(s.Clients[uri], client)
	s.ProtocolMap[uri] = drvName

	return &pb.BindingReply{Status: public.StatusOK, Message: public.MessageOK}, nil
}

// Release implement release interface
func (s *PortServer) Release(ctx context.Context, in *pb.ReleaseRequest) (*pb.ReleaseReply, error) {
	// LOG.Infof("Release request: %+v\n", in)

	// switch in.Protocol {
	// case ModbusSerial, ModbusTcp:
	// 	if err := ReleaseModbusClient(in.Port, in.Payload); err != nil {
	// 		errmsg := fmt.Sprintf("release modbus client failed, errmsg [%v]", err)
	// 		portlog.LOG.Warning(errmsg)

	// 		return &pb.ReleaseReply{Status: public.StatusErr, Message: errmsg}, nil
	// 	}

	// 	portlog.LOG.Infof("release modbus client failed, port [%s], protocol [%s]\n", in.Port, in.Protocol)
	// }

	return &pb.ReleaseReply{Status: public.StatusOK, Message: public.MessageOK}, nil
}

// Operate implements operate iterface
func (s *PortServer) Operate(ctx context.Context, in *pb.OperateRequest) (*pb.OperateReply, error) {
	// log.Printf("Operate request: %+v", in)

	drvName, ok := s.ProtocolMap[in.Port]
	if !ok {
		return &pb.OperateReply{Status: public.StatusErr, Message: errors.ErrNoData.As(in.Port).Error()}, nil
	}

	// execute operation
	switch in.Type {
	case public.OperateSample:
		client, err := s.getClient(drvName, in.Port, in.Suid, in.Payload)
		if err != nil {
			return &pb.OperateReply{Status: public.StatusErr, Message: err.Error()}, nil
		}

		result, err := client.Sample(in.Payload)
		if err != nil {
			return &pb.OperateReply{Status: public.StatusErr, Message: err.Error()}, errors.As(err)
		}

		return &pb.OperateReply{Status: public.StatusOK, Message: public.MessageOK, Data: result}, nil
	case public.OperateCommand:
		client, err := s.getClient(drvName, in.Port, in.Suid, in.Payload)
		if err != nil {
			return &pb.OperateReply{Status: public.StatusErr, Message: err.Error()}, nil
		}

		result, err := client.Command(in.Payload)
		if err != nil {
			return &pb.OperateReply{Status: public.StatusErr, Message: err.Error()}, err
		}

		return &pb.OperateReply{Status: public.StatusOK, Message: public.MessageOK, Data: result}, nil

	default:
		err := errors.New(public.MessageErrUnknown).As(in.Type)
		log.Warning(err)
		return &pb.OperateReply{Status: public.StatusErr, Message: err.Error()}, nil
	}
}
