/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/03
 * Despcription: common definition
 *
 */

package public

// status
const (
	StatusOK  = 0
	StatusErr = 1
)

// message
const (
	MessageOK         = "ok"
	MessageErrUnknown = "unknown err occured"
)

// operation
const (
	OperateSample  = "sample"
	OperateCommand = "command"
)

// datatype
const (
	DataTypeInt     = "int"
	DataTypeFloat   = "float"
	DataTypeString  = "string"
	DataTypeCommand = "command"
)

// element type
const (
	ElementTypeModbus      = "ModbusElement"
	ElementTypePMBus       = "PMBusElement"
	ElementTypeOilMachine  = "OilMachineElement"
	ElementTypeHMU         = "HMUElement"
	ElementTypeCamera      = "CameraElement"
	ElementTypeLuMiGateway = "LuMiGatewayElement"
	ElementTypeSNMP        = "SnmpManagerElement"
)

// protocol type
const (
	ProtocolTypeModbusSerial    = "ModbusSerialProtocol"
	ProtocolTypeModbusTCPClient = "ModbusTcpClientProtocol"
	ProtocolTypeYDN23SerialPort = "YDN23SerialPortProtocol"
	ProtocolTypePMBUS           = "PMBusProtocol"
	ProtocolTypeDeltaUPS        = "DeltaUPSProtocol"
	ProtocolTypeOilMachine      = "OilMachineProtocol"
	ProtocolTypeHMU             = "HMUProtocol"
	ProtocolTypeCamera          = "CameraProtocol"
	ProtocolTypeLampWith        = "LampWithProtocol"
	ProtocolTypeCreditCard      = "CreditCardProtocol"
	ProtocolTypeLuMiGateway     = "LuMiGatewayProtocol"
	ProtocolTypeFaceIPC         = "FaceIPCProtocol"
	ProtocolTypeSNMP            = "SnmpManagerProtocol"
	ProtocolTypeEntry           = "EntryProtocol"
)

// protocol
const (
	ProtocolModbusSerial    = "protocol-modbus-serial"
	ProtocolModbusTCP       = "protocol-modbus-tcp"
	ProtocolHYIOTMU         = "protocol-hyiot-mu"
	ProtocolPMBUS           = "protocol-pmbus"
	ProtocolYDN23           = "protocol-ydn23"
	ProtocolDeltaUPS        = "protocol-delta-ups"
	ProtocolOilMachine      = "protocol-oil-machine"
	ProtocolSensorflow      = "protocol-sensorflow"
	ProtocolCamera          = "protocol-camera"
	ProtocolCameraRstp      = "protocol-camera-rstp"
	ProtocolLampWith        = "protocol-lamp-with"
	ProtocolCreditCard      = "protocol-credit-card"
	ProtocolLuMiGateway     = "protocol-lumi-gateway"
	ProtocolFaceIPC         = "protocol-face-ipc"
	ProtocolSNMP            = "protocol-snmp-manager"
	ProtocolWeiGengEntry    = "protocol-weigeng-entry"
	ProtocolDIDO            = "protocol-dido"
	ProtocolHSJRFID         = "protocol-hsj-rfid"
	ProtocolES5200          = "protocol-es5200"
	ProtocolElecFire        = "protocol-elec-fire"
	ProtocolLTE             = "protocol-lte"
	ProtocolVirtualAntenna  = "protocol-virtual-antenna"
	ProtocolAutoLockDIDO    = "protocol-autolock-dido"
	ProtocolAutoLockIOBlock = "protocol-autolock-ioblock"
)

// 用于后台页面的端口配置的协议列表
var ProtocolList = []string{
	ProtocolModbusSerial,
	ProtocolModbusTCP,
	ProtocolSensorflow,
	ProtocolPMBUS,
	ProtocolYDN23,
	ProtocolOilMachine,
	ProtocolCamera,
	ProtocolCameraRstp,
	ProtocolLampWith,
	ProtocolLuMiGateway,
	ProtocolFaceIPC,
	ProtocolSNMP,
	ProtocolWeiGengEntry,
	ProtocolDIDO,
	ProtocolAutoLockDIDO,
	ProtocolAutoLockIOBlock,
	ProtocolHSJRFID,
	ProtocolES5200,
	ProtocolElecFire,
	ProtocolLTE,
	ProtocolVirtualAntenna,
}

// caller
const (
	CallerBusServer  = "busserver"
	CallerPortServer = "portserver"
)
