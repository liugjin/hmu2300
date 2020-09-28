package log

import "testing"

//
func TestFileLog(t *testing.T) {
	l := NewLog("testing")
	if err := l.SetFile("./test.log", 10*1024, 20); err != nil {
		t.Fatal(err)
	}
	for i := 10000; i > 0; i-- {
		l.Warning("Hello")
	}
	l.Close()
}
