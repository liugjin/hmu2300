package muid

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"sync"

	"clc.hmu/app/public/store"
)

var (
	id    = ""
	mutex = sync.Mutex{}
)

func GetMuID() string {
	mutex.Lock()
	defer mutex.Unlock()

	if len(id) > 0 {
		return id
	}
	busData, err := ioutil.ReadFile(store.GetRootDir() + "/app/aggregation/busmanager.json")
	if err != nil {
		panic(err)
	}
	busCfg := map[string]interface{}{}
	if err := json.Unmarshal(busData, &busCfg); err != nil {
		panic(err)
	}
	model := busCfg["model"].(string)
	switch model {
	case "hmu2000":
		data, err := ioutil.ReadFile("/tmp/dxs/snfile")
		if err != nil {
			ioutil.WriteFile("/tmp/dxs/snfile", []byte("0"), 0644)
			// panic("You need to create /tmp/dxs/snfile for developing or changing busmanager to 'pc' mode:" + err.Error())
		} else {
			id = strings.TrimSpace(string(data))
		}
	case "hmu2300":
		data, err := ioutil.ReadFile("/usrfs/app/uuid")
		if err != nil {
			id = "0"
		} else {
			id = strings.TrimSpace(string(data))
		}
	default:
		data, err := ioutil.ReadFile(store.GetRootDir() + "/app/aggregation/monitoring-units.json")
		if err != nil {
			panic(err)
		}
		mus := []map[string]interface{}{}
		if err := json.Unmarshal(data, &mus); err != nil {
			panic(err)
		}
		id = strings.TrimSpace(mus[0]["id"].(string))
	}
	return id
}

func SetMemID(ID string) {
	mutex.Lock()
	defer mutex.Unlock()
	id = ID
}
