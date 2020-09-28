package frpc

import (
	"fmt"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	go func() {
		Start()
	}()
	time.Sleep(1 * 1e9) // 等待程序起来
	tpl, err := FrpcTemplateFromEtc()
	if err != nil {
		t.Fatal(err)
	}
	tpl.ServerAddr = "127.0.0.1"
	fmt.Println(tpl.Format())
	Reload(tpl)
	time.Sleep(10 * 1e9)
}
