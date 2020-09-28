package memcfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestJsonCfg(t *testing.T) {
	// Create tmp data
	tmpfn := filepath.Join(os.TempDir(), uuid.New().String())
	if err := ioutil.WriteFile(tmpfn, []byte(`{"testing":"1"}`), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfn)

	dst := map[string]string{}

	// for no cache
	if err := GetJsonCfg(tmpfn, &dst); err != nil {
		t.Fatal(err)
	}
	val, _ := dst["testing"]
	if val != "1" {
		t.Fatalf("expect:1, but:%s", val)
	}

	// for cache
	dst = map[string]string{}
	// 移除测试文件，若未命中缓存，则读取文件并出错
	if err := os.Remove(tmpfn); err != nil {
		t.Fatal(err)
	}
	if err := GetJsonCfg(tmpfn, &dst); err != nil {
		t.Fatal(err)
	}
	if val != "1" {
		t.Fatalf("expect:1, but:%s", val)
	}

	// for write
	if err := WriteJsonCfg(tmpfn, dst); err != nil {
		t.Fatal(err)
	}
	// clean cached, and reread file data.
	dst = map[string]string{}
	CleanJsonCache(tmpfn)
	if err := GetJsonCfg(tmpfn, &dst); err != nil {
		t.Fatal(err)
	}
	if val != "1" {
		t.Fatalf("expect:1, but:%s", val)
	}

	// for update cache
	dst["testing"] = "2"
	if err := WriteJsonCfg(tmpfn, dst); err != nil {
		t.Fatal(err)
	}
	if err := GetJsonCfg(tmpfn, &dst); err != nil {
		t.Fatal(err)
	}
	val = dst["testing"]
	if val != "2" {
		t.Fatalf("expect:2, but:%s", val)
	}
}
