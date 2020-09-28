package at

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gwaylib/errors"
)

type FileATCmd struct {
	dev     string
	timeout time.Duration

	mutex      sync.Mutex
	resultChan chan string
	errChan    chan error
}

func NewFileATCmd(dev string, timeout time.Duration) AT {
	return &FileATCmd{
		dev:     dev,
		timeout: timeout,

		resultChan: make(chan string, 1),
		errChan:    make(chan error, 1),
	}
}

func (at *FileATCmd) do(cmd string) (string, error) {
	file, err := os.OpenFile(at.dev, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return "", errors.As(err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(cmd + "\r\n")); err != nil {
		return "", errors.As(err)
	}

	// TODO:make timeout
	result := []byte{}
	for {
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				return strings.TrimSpace(string(result)), nil
			}
			return "", errors.As(err)
		}
		// echo ok\n
		if n == 3 && buf[0] == 79 && buf[1] == 75 && buf[2] == 10 {
			return strings.TrimSpace(string(result)), nil
		}
		result = append(result, buf[:n]...)
	}
	panic("Not reach here")
	return "", nil
}

func (at *FileATCmd) Do(cmd string) (string, error) {
	at.mutex.Lock()
	defer at.mutex.Unlock()

	go func() {
		result, err := at.do(cmd)
		if err != nil {
			at.errChan <- err
		} else {
			at.resultChan <- result
		}
	}()
	select {
	case <-time.After(at.timeout):
		return "", errors.New("timeout").As(at.dev)
	case result := <-at.resultChan:
		return result, nil
	case err := <-at.errChan:
		return "", errors.As(err)
	}
	panic("Not reach here")
}

// 检查网络模块是否正常
func (at *FileATCmd) AT() error {
	_, err := at.Do("AT")
	return errors.As(err)
}

// 检查sim卡是否存在
func (at *FileATCmd) CPIN() (string, error) {
	// TODO: 解析数据
	return at.Do("AT+CPIN?")
}

// 取信号强度
// 33最高
func (at *FileATCmd) CSQ() (string, error) {
	// TODO: 解析数据
	return at.Do("AT+CSQ")
}
