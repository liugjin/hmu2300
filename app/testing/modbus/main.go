package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/goburrow/modbus"
)

type Cfg struct {
	DevPort  string `json:"dev_port"`
	BaudRate int    `json:"baud_rate"`
	Parity   string `json:"parity"`
	DataBit  int    `json:"data_bit"`
	StopBit  int    `json:"stop_bit"`
	SlaveId  byte   `json:"slave_id"`
	FnCode   byte   `json:"fn_code"`
	Timeout  int    `json:"timeout"` // milliSecond

}

func main() {
	data, err := ioutil.ReadFile("./cfg.json")
	if err != nil {
		panic(err)
	}
	cfg := &Cfg{}
	if err := json.Unmarshal(data, cfg); err != nil {
		panic(err)
	}
	handler := modbus.NewRTUClientHandler(cfg.DevPort)
	handler.BaudRate = cfg.BaudRate
	handler.DataBits = cfg.DataBit
	handler.Parity = cfg.Parity
	handler.StopBits = cfg.StopBit
	handler.SlaveId = cfg.SlaveId
	handler.Timeout = time.Duration(cfg.Timeout) * time.Millisecond

	if err := handler.Connect(); err != nil {
		panic(err)
	}
	defer handler.Close()
	fmt.Println("connected")

	client := modbus.NewClient(handler)

	switch cfg.FnCode {
	case 1:
		fmt.Println("ReadCoils")
		results, err := client.ReadCoils(0x0000, 6)
		if err != nil {
			panic(err)
		}
		fmt.Println(results)
	case 2:
		fmt.Println("ReadDiscreteInputs")
		results, err := client.ReadDiscreteInputs(0x0000, 6)
		if err != nil {
			panic(err)
		}
		fmt.Println(results)
	case 3:
		fmt.Println("ReadHoldingRegisters")
		results, err := client.ReadHoldingRegisters(0x0000, 6)
		if err != nil {
			panic(err)
		}
		if len(results) != 12 {
			panic(fmt.Sprintf("error protocal:%v", results))
		}

		for j := 0; j < 12; j += 2 {
			d := int32(0)
			d |= int32(results[j]) << 8
			d |= int32(results[j+1])

			fmt.Printf("s:%d: data:%d\n", j/2, d)
		}
	case 4:
		fmt.Println("ReadInputRegisters")
		results, err := client.ReadInputRegisters(0x0000, 6)
		if err != nil {
			panic(err)
		}
		if len(results) != 12 {
			panic(fmt.Sprintf("error protocal:%v", results))
		}

		for j := 0; j < 12; j += 2 {
			d := int32(0)
			d |= int32(results[j]) << 8
			d |= int32(results[j+1])

			fmt.Printf("s:%d: data:%d\n", j/2, d)
		}
	case 24:
		fmt.Println("ReadFIFOQueue")
		results, err := client.ReadFIFOQueue(0x0000)
		if err != nil {
			panic(err)
		}
		if len(results) != 12 {
			panic(fmt.Sprintf("error protocal:%v", results))
		}

		for j := 0; j < 12; j += 2 {
			d := int32(0)
			d |= int32(results[j]) << 8
			d |= int32(results[j+1])

			fmt.Printf("s:%d: data:%d\n", j/2, d)
		}

	}
}
