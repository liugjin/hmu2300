package memcfg

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"sync"

	"github.com/gwaylib/errors"
)

var (
	jsonCfgMap  = map[string]interface{}{}
	jsonCfgSync = sync.Mutex{}
)

// 像json.Umarshal一样获取，若需要实时的配置数据，应每次使用重新读取缓存的数据。
// 若缓存在，将直接返回缓存;
// 若没有，加载配置文件的数据到内存，并填充请求的值。
// 如果是多进程访问，请构建微服务通过网络(例如rpc)来取值。
func GetJsonCfg(fileName string, dst interface{}) error {
	jsonCfgSync.Lock()
	defer jsonCfgSync.Unlock()

	val, ok := jsonCfgMap[fileName]
	if ok {
		// Match cache.
		rv := reflect.ValueOf(dst)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return errors.New("non-pointer").As(fileName, rv.Type().String())
		}
		// 将原指针的值改到缓存的值
		dstVal := reflect.Indirect(rv)
		dstVal.Set(reflect.Indirect(reflect.ValueOf(val)))
		return nil
	}

	// Loading file to memory.
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.As(err, fileName)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return errors.As(err, fileName)
	}
	jsonCfgMap[fileName] = dst
	return nil
}

// 更改缓存中的数据，并将缓存写入持久性配置文件。
func WriteJsonCfg(fileName string, src interface{}) error {
	jsonCfgSync.Lock()
	defer jsonCfgSync.Unlock()

	jsonCfgMap[fileName] = src
	data, err := json.MarshalIndent(src, "", "	")
	if err != nil {
		return errors.As(err, fileName, src)
	}
	if err := ioutil.WriteFile(fileName, data, 0666); err != nil {
		return errors.As(err, fileName, src)
	}
	return nil
}

// 清空内存缓存，但不会清除文件上的数据，以便临时的配置文件可以释放缓存。
func CleanJsonCache(fileName string) {
	jsonCfgSync.Lock()
	defer jsonCfgSync.Unlock()
	delete(jsonCfgMap, fileName)
}
