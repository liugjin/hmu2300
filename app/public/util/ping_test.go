package util

import "testing"

func TestTCPPing(t *testing.T) {
	var tests = []struct {
		Host   string
		Result bool
	}{
		{"baidu.com:80", true},
		{"127.0.0.1:888", false},
	}
	for index, val := range tests {
		result := TCPPing(val.Host)
		if result != val.Result {
			t.Fatalf("Expect:%t,but:%t, index:%d", val.Result, result, index)
		}
	}
}
