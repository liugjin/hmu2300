package bootflag

import (
	"io/ioutil"
	"os"
	"strings"

	"clc.hmu/app/public/store/etc"
	"github.com/gwaylib/errors"
)

var flagPath = os.ExpandEnv(etc.Etc.String("public/log", "boot_flag"))

func WriteFlag() error {
	return errors.As(ioutil.WriteFile(flagPath, []byte("1"), 0666))
}

func GetFlag() (string, error) {
	data, err := ioutil.ReadFile(flagPath)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			return "0", nil
		}
		return "", errors.As(err)
	}
	return string(data), nil
}

func CleanFlag() error {
	return errors.As(ioutil.WriteFile(flagPath, []byte("0"), 0666))
}
