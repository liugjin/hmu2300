/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/10/30
 * Despcription: credit card machine implement
 *
 */

package protocol

import (
	"encoding/json"
	"sync"

	"clc.hmu/app/public"
	"github.com/gwaylib/errors"
)

// ============= register driver start ==========================

// 驱动注册
func init() {
	RegDriverProtocol(public.ProtocolCreditCard, generalCreditCardDriverProtocol)
}

// Implement DriverProtocol
type creditCardDriverProtocol struct {
	req *public.CreditCardBindingPayload
	uri string
}

func (dp *creditCardDriverProtocol) Payload() interface{} {
	return dp.req
}

func (dp *creditCardDriverProtocol) ClientID() string {
	return CreditCardClientID
}

func (dp *creditCardDriverProtocol) NewInstance() (PortClient, error) {
	return NewCreditCardClient(
		dp.req.Username, dp.req.Password,
		dp.req.SerialNumber,
	)
}

func generalCreditCardDriverProtocol(uri, suid, payload string) (DriverProtocol, error) {
	// decode payload
	req, err := DecodeCreditCardBindingPayload(payload)
	if err != nil {
		return nil, errors.As(err, payload)
	}
	return &creditCardDriverProtocol{
		req: &req,
		uri: uri,
	}, nil
}

// ============= register driver end ==========================

// CreditCardClientID id
var CreditCardClientID = "credit-card-machine-client-id"

const ccSuccess = 0

// UserLoginInfo login info
type UserLoginInfo struct {
	dwSize         int32
	dwTimeout      int32
	dwVID          int32
	dwPID          int32
	szUserName     [32]byte
	szPassword     [16]byte
	szSerialNumber [48]byte
	byRes          [80]byte
}

// DeviceRegisterResult device register result
type DeviceRegisterResult struct {
	dwSize            int32
	szDeviceName      [32]byte
	szSerialNumber    [48]byte
	dwSoftwareVersion int32
	byRes             [40]byte
}

// WaitSecond wait second
type WaitSecond struct {
	dwSize int32
	byWait byte // 0: always run until response; else: seconds
	byRes  [27]byte
}

// ConfigInputInfo config input info
type ConfigInputInfo struct {
	lpCondBuffer     *WaitSecond
	dwCondBufferSize int32
	lpInBuffer       *WaitSecond
	dwInBufferSize   int32
	byRes            [48]byte
}

// BeepFlicker beep and flicker
type BeepFlicker struct {
	dwSize         int32
	byBeepType     byte // 0 disable; 1 continue; 2 slow; 3 fast; 4 stop
	byBeepCount    byte // only enable for slow and fast, can not be zero
	byFlickerType  byte // 0 disable; 1 continue; 2 error; 3 correct; 4 stop
	byFlickerCount byte // only enable for error and correct, can not be zero
	byRes          [24]byte
}

// BeepFlickerInputInfo input info
type BeepFlickerInputInfo struct {
	lpCondBuffer     *BeepFlicker
	dwCondBufferSize int32
	lpInBuffer       *BeepFlicker
	dwInBufferSize   int32
	byRes            [48]byte
}

// ActivateCardResult activate card result
type ActivateCardResult struct {
	dwSize            int32
	byCardType        byte     // card type （0-TypeA m1 card，1-TypeA cpu card, 2-TypeB card, 3-125kHz Id card）
	bySerialLen       byte     // length of card's serial number
	bySerial          [10]byte // card's serial number
	bySelectVerifyLen byte     // length of select verify
	bySelectVerify    [3]byte  // select verify
	byRes             [12]byte
}

// ConfigOutputInfo config oputput info
type ConfigOutputInfo struct {
	lpOutBuffer     *ActivateCardResult // output cache
	dwOutBufferSize int32               // cache size
	byRes           [56]byte
}

// DeviceInfo device info
type DeviceInfo struct {
	dwSize         int32
	dwVID          int32
	dwPID          int32
	szManufacturer [32]byte
	szDeviceName   [32]byte
	szSerialNumber [48]byte
	byRes          [68]byte
}

// CreditCardClient credit card client
type CreditCardClient struct {
	ClientID string

	// login info
	username     string
	password     string
	serialNumber string

	userID uintptr     // distribute by login
	cardID chan []byte // cache card id

	mtx sync.Mutex

	// initFunc       *syscall.Proc
	// getLastErrFunc *syscall.Proc
	// enumDeviceFunc *syscall.Proc
	// loginFunc      *syscall.Proc
	// getConfigFunc  *syscall.Proc
	// setConfigFunc  *syscall.Proc
	// cleanupFunc    *syscall.Proc
	// logoutFunc     *syscall.Proc
}

// NewCreditCardClient new client
func NewCreditCardClient(username, password, serailnumber string) (PortClient, error) {
	var client = CreditCardClient{
		ClientID:     CreditCardClientID,
		username:     username,
		password:     password,
		serialNumber: serailnumber,
	}

	if err := client.Init(); err != nil {
		return nil, err
	}

	// if err := client.Start(); err != nil {
	// 	return nil, err
	// }

	return &client, nil
}

// DecodeCreditCardBindingPayload decode binding payload
func DecodeCreditCardBindingPayload(payload string) (public.CreditCardBindingPayload, error) {
	var p public.CreditCardBindingPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return p, err
	}

	return p, nil
}

// DeviceCallBack enum device call back
func DeviceCallBack(pDevceInfo *DeviceInfo, pUser *CreditCardClient) *interface{} {
	// // compare serial number, login if exist
	// serialnum := string(bytes.TrimRight(pDevceInfo.szSerialNumber[:], "\x00"))
	// if serialnum == pUser.serialNumber {
	// 	var info UserLoginInfo
	// 	var res DeviceRegisterResult

	// 	user := []byte(pUser.username)
	// 	pass := []byte(pUser.password)

	// 	l := len(user)
	// 	for i := 0; i < l; i++ {
	// 		info.szUserName[i] = user[i]
	// 	}

	// 	l = len(pass)
	// 	for i := 0; i < l; i++ {
	// 		info.szPassword[i] = pass[i]
	// 	}

	// 	info.dwSize = int32(unsafe.Sizeof(info))
	// 	info.dwTimeout = 5000
	// 	info.dwVID = pDevceInfo.dwVID
	// 	info.dwPID = pDevceInfo.dwPID
	// 	info.szSerialNumber = pDevceInfo.szSerialNumber

	// 	res.dwSize = int32(unsafe.Sizeof(res))

	// 	pUser.userID, _, _ = pUser.loginFunc.Call(uintptr(unsafe.Pointer(&info)), uintptr(unsafe.Pointer(&res)))
	// 	r1, _, _ := pUser.getLastErrFunc.Call()
	// 	if r1 != ccSuccess {
	// 		log.Printf("login fail, errcode: %v", r1)
	// 		return nil
	// 	}

	// 	log.Printf("login success, userid: %v", pUser.userID)
	// }

	return nil
}

// Init init
func (cc *CreditCardClient) Init() error {
	// dllname := "HCUsbSDK.dll"
	// sdk, err := syscall.LoadDLL(dllname)
	// if err != nil {
	// 	return fmt.Errorf("load dll fail: %v", err)
	// }

	// cc.initFunc, err = sdk.FindProc("USB_SDK_Init")
	// if err != nil {
	// 	return fmt.Errorf("find init failed: %v", err)
	// }

	// cc.getLastErrFunc, err = sdk.FindProc("USB_SDK_GetLastError")
	// if err != nil {
	// 	return fmt.Errorf("find get last error failed: %v", err)
	// }

	// cc.enumDeviceFunc, err = sdk.FindProc("USB_SDK_EnumDevice")
	// if err != nil {
	// 	return fmt.Errorf("find enum device failed: %v", err)
	// }

	// cc.loginFunc, err = sdk.FindProc("USB_SDK_Login")
	// if err != nil {
	// 	return fmt.Errorf("find login failed: %v", err)
	// }

	// cc.getConfigFunc, err = sdk.FindProc("USB_SDK_GetDeviceConfig")
	// if err != nil {
	// 	return fmt.Errorf("find get device config failed: %v", err)
	// }

	// cc.setConfigFunc, err = sdk.FindProc("USB_SDK_SetDeviceConfig")
	// if err != nil {
	// 	return fmt.Errorf("find set device config failed: %v", err)
	// }

	// cc.cleanupFunc, err = sdk.FindProc("USB_SDK_Cleanup")
	// if err != nil {
	// 	return fmt.Errorf("find cleanup failed: %v", err)
	// }

	// cc.logoutFunc, err = sdk.FindProc("USB_SDK_Logout")
	// if err != nil {
	// 	return fmt.Errorf("find logout failed: %v", err)
	// }

	return nil
}

// Start start
func (cc *CreditCardClient) Start() error {
	// // init
	// cc.initFunc.Call()

	// r1, _, _ := cc.getLastErrFunc.Call()
	// if r1 != ccSuccess {
	// 	return fmt.Errorf("init failed, errcode: %v", r1)
	// }

	// enum device, login device in callback
	return cc.EnumDevice()
}

// Stop stop
func (cc *CreditCardClient) Stop() error {
	// // cc.logoutFunc.Call(cc.userID)
	// // r1, _, _ := cc.getLastErrFunc.Call()
	// // if r1 != ccSuccess {
	// // 	return fmt.Errorf("logout failed, errcode: %v", r1)
	// // }

	// cc.cleanupFunc.Call()
	// r1, _, _ := cc.getLastErrFunc.Call()
	// if r1 != ccSuccess {
	// 	return fmt.Errorf("cleanup failed, errcode: %v", r1)
	// }

	return nil
}

// EnumDevice enum device
func (cc *CreditCardClient) EnumDevice() error {
	// // enum device, login device in callback
	// fn := syscall.NewCallback(DeviceCallBack)
	// cc.enumDeviceFunc.Call(fn, uintptr(unsafe.Pointer(cc)))

	// r1, _, _ := cc.getLastErrFunc.Call()
	// if r1 != ccSuccess {
	// 	return fmt.Errorf("enum device failed, errcode: %v", r1)
	// }

	return nil
}

// BeepAndFlicker beep and flicker
func (cc *CreditCardClient) BeepAndFlicker() {
	// // slow beep one time, correct flicker one time
	// var bf BeepFlicker
	// bf.dwSize = int32(unsafe.Sizeof(bf))
	// bf.byBeepType = 2
	// bf.byBeepCount = 1
	// bf.byFlickerType = 3
	// bf.byFlickerCount = 1

	// var struBeep BeepFlickerInputInfo
	// struBeep.dwInBufferSize = int32(unsafe.Sizeof(bf))
	// struBeep.lpInBuffer = &bf

	// var struOutput ConfigOutputInfo
	// cc.setConfigFunc.Call(cc.userID, uintptr(0x0100), uintptr(unsafe.Pointer(&struBeep)), uintptr(unsafe.Pointer(&struOutput)))
}

// ID id
func (cc *CreditCardClient) ID() string {
	return cc.ClientID
}

// Sample sample, get values
func (cc *CreditCardClient) Sample(payload string) (string, error) {
	return "", nil
}

// Command command, set values
func (cc *CreditCardClient) Command(payload string) (string, error) {
	// cc.mtx.Lock()
	// defer cc.mtx.Unlock()

	// // enum device and login
	// if err := cc.Start(); err != nil {
	// 	fmt.Println("err")

	// 	// clean up
	// 	fmt.Println(cc.Stop())

	// 	return "", err
	// }

	// cc.BeepAndFlicker()

	// var struWaitSecond WaitSecond
	// struWaitSecond.dwSize = int32(unsafe.Sizeof(struWaitSecond))
	// struWaitSecond.byWait = 3

	// var struActivateRes ActivateCardResult
	// struActivateRes.dwSize = int32(unsafe.Sizeof(struActivateRes))

	// var struInput ConfigInputInfo
	// struInput.dwInBufferSize = int32(unsafe.Sizeof(struWaitSecond))
	// struInput.lpInBuffer = &struWaitSecond

	// var struOutput ConfigOutputInfo
	// struOutput.dwOutBufferSize = int32(unsafe.Sizeof(struActivateRes))
	// struOutput.lpOutBuffer = &struActivateRes

	// // get config
	// cc.getConfigFunc.Call(cc.userID, uintptr(0x0104), uintptr(unsafe.Pointer(&struInput)), uintptr(unsafe.Pointer(&struOutput)))
	// r1, _, _ := cc.getLastErrFunc.Call()
	// if r1 != ccSuccess {
	// 	return "", fmt.Errorf("get device config failed, errcode: %v", r1)
	// }

	// cardid := struActivateRes.bySerial[:struActivateRes.bySerialLen]

	// // reserve
	// l := len(cardid)
	// for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
	// 	cardid[i], cardid[j] = cardid[j], cardid[i]
	// }

	// var val int
	// switch l {
	// case 1:
	// 	val = int(cardid[0])
	// case 2:
	// 	val = int(binary.BigEndian.Uint16(cardid))
	// case 4:
	// 	val = int(binary.BigEndian.Uint32(cardid))
	// case 8:
	// 	val = int(binary.BigEndian.Uint64(cardid))
	// default:
	// 	val = 0
	// }

	// log.Println(val)

	// cc.BeepAndFlicker()

	// // clean up
	// cc.Stop()

	// return strconv.Itoa(val), nil

	return "", nil
}
