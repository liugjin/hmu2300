package httptry

import (
	"fmt"
	"net/http/httputil"
	"testing"
)

func TestHttpsClient(t *testing.T) {
	resp, err := HttpsClient.Get("https://baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}
