/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/09/11
 * Despcription: generate clients
 *
 */

package portmanager

import (
	"fmt"
	"strconv"

	"clc.hmu/app/portmanager/protocol"
	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// HasClientExist check client exist or not
func (s *PortServer) HasClientExist(port, id string) (protocol.PortClient, bool) {
	clients := s.Clients[port]
	for _, client := range clients {
		if client.ID() == id {
			return client, true
		}
	}

	return nil, false
}

// get client, pro: protocol
func (s *PortServer) getClient(pro, port, suid, payload string) (protocol.PortClient, error) {
	id := ""
	switch pro {
	case public.ProtocolModbusSerial, public.ProtocolModbusTCP:
		// decode payload
		req, err := protocol.DecodeModbusPayload(payload)
		if err != nil {
			return nil, err
		}

		id = strconv.Itoa(int(req.Slaveid)) + strconv.Itoa(int(req.BaudRate))
	case public.ProtocolHYIOTMU:
		// decode payload
		req, err := protocol.DecodeSystemOperationRequest(payload)
		if err != nil {
			return nil, err
		}

		id = req.Model
	case public.ProtocolPMBUS, public.ProtocolYDN23:
		// decode payload
		req, err := protocol.DecodePMBUSBindingPayload(payload)
		if err != nil {
			return nil, err
		}

		id = protocol.PMBUSClientID + strconv.Itoa(int(req.ADR))
	case public.ProtocolOilMachine:
		id = protocol.OilMachineClientID
	case public.ProtocolDeltaUPS:
		id = protocol.DeltaUPSClientID
	case public.ProtocolSensorflow:
		id = protocol.SensorflowClientID
	case public.ProtocolLampWith:
		id = protocol.LampWithClientID
	case public.ProtocolCreditCard:
		id = protocol.CreditCardClientID
	case public.ProtocolLuMiGateway:
		id = protocol.LuMiGatewayClientID
	case public.ProtocolFaceIPC:
		id = protocol.FaceIPCClientID
	case public.ProtocolSNMP:
		// decode payload
		req, err := protocol.DecodeSNMPOperatePayload(payload)
		if err != nil {
			return nil, err
		}

		id = protocol.SNMPClientID + req.Target
	case public.ProtocolWeiGengEntry:
		id = protocol.WeiGengEntryClientID
	case public.ProtocolDIDO:
		id = protocol.DIDOClientID + port
	case public.ProtocolHSJRFID:
		id = protocol.HSJRFIDClientID
	case public.ProtocolES5200:
		id = protocol.ES5200ClientID
	case public.ProtocolCamera:
		// decode payload
		req, err := protocol.DecodeCameraBindingPayload(payload)
		if err != nil {
			return nil, err
		}

		id = protocol.VideoClientID + req.Host
	case public.ProtocolElecFire:
		id = protocol.ElecFireClientID
	default:
		// TODO:后续统一使用此接口获取ID
		drvPro, err := protocol.GetDriverProtocol(pro, port, suid, payload)
		if err != nil {
			return nil, errors.As(err)
		}
		id = drvPro.ClientID()
	}

	// find client
	client, exist := s.HasClientExist(port, id)
	if !exist {
		return nil, fmt.Errorf("client not exist")
	}

	return client, nil
}

// generate modbus rtu client
func (s *PortServer) generateModbusRTUClient(payload, port string) (protocol.PortClient, error) {
	// decode payload
	req, err := protocol.DecodeModbusPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := strconv.Itoa(int(req.Slaveid)) + strconv.Itoa(int(req.BaudRate))
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewModbusSerialClient(port, req.BaudRate, req.Timeout, byte(req.Slaveid))
}

// generate modbus tcp client
func (s *PortServer) generateModbusTCPClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeModbusPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := strconv.Itoa(int(req.Slaveid)) + strconv.Itoa(int(req.BaudRate))
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewModbusTCPClient(port, req.BaudRate, req.Timeout, byte(req.Slaveid))
}

// generate system client
func (s *PortServer) generateSystemClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeSystemBindingRequest(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := req.Model
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewSystemClient(req)
}

// generate pm bus client
func (s *PortServer) generatePMBusClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodePMBUSBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// find client
	id := protocol.PMBUSClientID + strconv.Itoa(int(req.ADR))
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewPMBusClient(port, req.BaudRate, req.Timeout, req.SOI, req.VER, req.ADR, req.CID1, req.EOI)
}

// generate oil machine client
func (s *PortServer) generateOilMachineClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeOilMachineBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// find client
	id := protocol.OilMachineClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewOilMachineClient(port, req.BaudRate, req.Timeout, req.SOI, req.EOI)
}

// generate delta ups client
func (s *PortServer) generateDeltaUPSClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeDeltaUPSBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// find client
	id := protocol.DeltaUPSClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewDeltaUPSClient(port, req.BaudRate, req.Timeout, req.Header, uint16(req.ID))
}

// generate sensorflow client
func (s *PortServer) generateSensorflowClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeModbusPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.SensorflowClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewSensorflowClient(port, int(req.BaudRate), int(req.Timeout), int(req.KeyNumber), req.MUID, req.SUID, req.WANInterface, req.WifiInterface)
}

// generate video client
func (s *PortServer) generateVideoClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeCameraBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.VideoClientID + req.Host
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewVideoClient(req.Host, req.User, req.Password)
}

// generate lamp with client
func (s *PortServer) generateLampWithClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeLampWithBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.LampWithClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewLampWithClient(port, int(req.BaudRate), req.Timeout)
}

// generate credit card client
func (s *PortServer) generateCreditCardClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeCreditCardBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.CreditCardClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewCreditCardClient(req.Username, req.Password, req.SerialNumber)
}

// generate lumi gateway client
func (s *PortServer) generateLuMigatewayClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeLuMiGatewayBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.LuMiGatewayClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewLuMiGatewayClient(req.SID, req.Password, req.NetInterface)
}

// generate face ipc client
func (s *PortServer) generateFaceIPCClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeFaceIPCBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.FaceIPCClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewFaceIPCClient(req)
}

// generate snmp client
func (s *PortServer) generateSNMPClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeSNMPBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.SNMPClientID + req.Target
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewSNMPClient(req.Version, req.Target, req.Port, req.ReadCommunity, req.WriteCommunity)
}

// generate wei geng entry client
func (s *PortServer) generateWeiGengEntryClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeWeiGengEntryBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.WeiGengEntryClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewWeiGengEntryClient(req.LocalAddress, req.LocalPort, req.DoorAddress, req.DoorPort, req.SerialNumber)
}

// generate dido client
func (s *PortServer) generateDIDOClient(payload, port string) (protocol.PortClient, error) {
	// check client exist or not
	id := protocol.DIDOClientID + port
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewDIDOClient(port)
}

// generate hsj rfid client
func (s *PortServer) generateHSJRFIDClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeModbusPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.HSJRFIDClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}
	return protocol.NewHSJRFIDClient(port, int(req.BaudRate), int(req.Timeout))
}

// generate ES5200 client
func (s *PortServer) generateES5200Client(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeES5200BindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.ES5200ClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}
	return protocol.NewES5200Client(port, int(req.BaudRate), int(req.Timeout))
}

// generate elec fire client
func (s *PortServer) generateElecFireClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeElecFireBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.ElecFireClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewElecFireClient(req)
}

func (s *PortServer) generateVirtualAntennaClient(payload, port string) (protocol.PortClient, error) {
	req, err := protocol.DecodeVirtualAntennaBindingPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload failed, errmsg [%v]", err)
	}

	// check client exist or not
	id := protocol.VirtualAntennaClientID
	if _, exist := s.HasClientExist(port, id); exist {
		return nil, fmt.Errorf("client has exist, id: %v", id)
	}

	return protocol.NewVirtualAntennaClient(req)
}
