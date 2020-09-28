package etc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gwaylib/conf/ini"
)

var Etc *ini.File

func init() {
	etcPath, err := findEtc()
	if err != nil {
		panic(err)
	}
	// 写死位置，以便用调用
	file, err := ini.GetFile(etcPath)
	if err != nil {
		panic(err)
	}
	Etc = file
}

func findEtc() (string, error) {
	possibleConf := []string{
		os.ExpandEnv("$PRJ_ROOT/app/aggregation/etc/etc.ini"),
		"./etc.ini",
		"./etc/etc.ini",
		"/etc/aggregation/etc.ini",
	}

	for _, file := range possibleConf {
		if _, err := os.Stat(file); err == nil {
			fmt.Println("etc file:", file)
			abs_file, err := filepath.Abs(file)
			if err == nil {
				return abs_file, nil
			} else {
				return file, nil
			}
		}
	}

	return "", fmt.Errorf("fail to find etc.ini")
}
