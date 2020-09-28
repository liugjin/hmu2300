/*
 *
 * Copyright 2018 huayuan-iot
 *
 * Author: lynn
 * Date: 2018/07/27
 * Despcription: test file
 *
 */

package public_test

import (
	"testing"

	"clc.hmu/app/public"
)

func TestVerifyTimeFormat(t *testing.T) {
	var testcase = []struct {
		value string
		want  bool
	}{
		{"", false},
		{"::", false},
		{"0:0:0", false},
		{"00:00:00", true},
		{"23:59:59", true},
		{"24:00:00", false},
		{"23:60:00", false},
		{"23:00:60", false},
		{"12:12:12", true},
		{"123:12:00", false},
	}

	for i, v := range testcase {
		if public.LegalTimeFormat(v.value) != v.want {
			t.Errorf("[%v] test fail, want [%v]", i+1, v.want)
		}
	}
}

func TestVerifyQueryTimeFormat(t *testing.T) {
	var testcase = []struct {
		value string
		want  bool
	}{
		{"", false},
		{"2018-10-1012:12", false},
		{"2018-10-10 12:12", true},
		{"2018-10-10  12:12", false},
	}

	for i, v := range testcase {
		if public.LegalQueryTimeFormat(v.value) != v.want {
			t.Errorf("[%v] test fail, want [%v]", i+1, v.want)
		}
	}
}

func TestTransferQueryTimeFormat(t *testing.T) {
	var testcase = []struct {
		value string
		want  string
	}{
		{"", ""},
		{"2018-10-10 12:12", "2018-10-10T12-12-00"},
	}

	for i, v := range testcase {
		if public.TransferQueryTimeFormat(v.value) != v.want {
			t.Errorf("[%v] test fail, want [%v]", i+1, v.want)
		}
	}
}

func TestAppVersion(t *testing.T) {
	testcase := []struct {
		app string
		ver string
	}{
		{"/tmp/aggregation", "v1.1.1"},
	}

	for i, v := range testcase {
		output, err := public.AppVersion(v.app)
		if err != nil || output != v.ver {
			t.Errorf("[%d] test failed, app[%s], output[%s], want version[%s]", i, v.app, output, v.ver)
		}
	}
}

func TestUploadFile(t *testing.T) {
	filename := "cindy.jpg"
	host := "lab.huayuan-iot.com"
	author := "hyiot"
	project := "video"
	token := "b2b8ec80-8a3a-11e8-9083-afae74b81b2b"
	user := "hyiot"

	url, err := public.UploadFile(filename, filename, host, author, project, token, user)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(url)
}

func TestQueryInterfaceInfoByName(t *testing.T) {
	name := "lo"
	ip, mac, err := public.QueryInterfaceInfoByName(name)
	if err != nil {
		t.Error(err)
	}

	t.Log(ip, mac)
}
