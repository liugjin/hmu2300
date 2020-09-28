package log

import (
	"clc.hmu/app/supd/bakfile"
	"github.com/gwaylib/errors"
	"github.com/op/go-logging"
)

type LogFile struct {
	*bakfile.File
}

// 同一时间仅能输出到一个文件
// maxFiles -- 记录的最多文件数
// maxSize -- 每个文件最大字节数, 102
func (f *Log) SetFile(fileName string, maxSize int64, bakNum int) error {
	// Auto close the last file
	if f.file != nil {
		f.file.Close() // ignore error
	}
	file, err := bakfile.OpenFile(fileName, maxSize, bakNum)
	if err != nil {
		return errors.As(err)
	}
	f.file = &LogFile{file}
	f.fileLog = logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(f.file, "", 0), format))
	f.fileLog.SetLevel(logging.INFO, f.module) // Echo to console
	f.SetBackend(logging.MultiLogger(f.consoleLog, f.fileLog))
	return nil
}

func (f *Log) GetFile() (*LogFile, error) {
	if f.file != nil {
		return f.file, nil
	}
	return nil, errors.New("Not Init")
}
