/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/08/28
 * Despcription: test file
 *
 */

package appmanager_test

import (
	"testing"

	"clc.hmu/app/public"
)

func TestDownloadUpgradePackage(t *testing.T) {
	testcase := []struct {
		url      string
		savepath string
		want     bool
	}{
		{"", "", false},
		{"http://127.0.0.1:6789/aggregation.tar.gz", "aggregation.tar.gz", true},
		{"http://127.0.0.1:6789/agg", "aggregation", false},
	}

	for i, v := range testcase {
		err := public.DownloadUpgradePackage(v.url, v.savepath)
		if (err == nil) != v.want {
			t.Errorf("[%d] test fail, url[%s], savepath[%s], want[%v], errmsg[%v]", i, v.url, v.savepath, v.want, err)
		}
	}
}

func TestDecompress(t *testing.T) {
	testcase := []struct {
		src  string
		dest string
		want bool
	}{
		{"./aggregation.tar.gz", "./", true},
	}

	for i, v := range testcase {
		err := public.Decompress(v.src, v.dest)
		if (err == nil) != v.want {
			t.Errorf("[%d] test fail, src[%s], dest[%s], want[%v], errmsg[%v]", i, v.src, v.dest, v.want, err)
		}
	}
}

func TestCheckFileSHA256(t *testing.T) {
	testcase := []struct {
		filepath string
		sum      string
		want     bool
	}{
		{"", "", false},
		{"aggregation", "9bc20bf43dc5e52bb68b330e46b7d3a6fc74badcd8f5213eb57e0b4aeec0d042", true},
		{"aggregation.tar.gz", "120a6c80ae34b649fba02a7951b8f0b779ecd5a58611b18d4c8c6d9e10dc7f68", true},
	}

	for i, v := range testcase {
		if public.CheckFileSHA256(v.filepath, v.sum) != v.want {
			t.Errorf("[%d] test fail, filepath[%s], sum[%s], want[%v]", i, v.filepath, v.sum, v.want)
		}
	}
}
