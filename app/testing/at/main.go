// testing AT for hmu2000
package main

import (
	"flag"
	"fmt"
	"time"

	"clc.hmu/app/public/at"
)

func main() {
	dev := flag.String("dev", "/dev/ttyUSB5", "device path of AT")
	cmd := flag.String("cmd", "AT", "cmd of AT")
	flag.Parse()

	c := at.NewHMU2000ATCmd(*dev, 60*1e9)
	result, err := c.Do(*cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println("cmd:", result)

	result, err = c.CPIN()
	if err != nil {
		panic(err)
	}
	fmt.Println("cpin:", result)

	for {
		result, err = c.CSQ()
		if err != nil {
			panic(err)
		}
		fmt.Println("AT+CSQ:", result)
		time.Sleep(1e9)
	}

}
