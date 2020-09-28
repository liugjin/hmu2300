package log

// 构建默认的日志器
var defaultLog = NewLog("default")

// 关闭日志
func Close() error {
	return defaultLog.Close()
}

// 默认日志不输出到文件，需要自行启用设置
func SetFile(fileName string, maxSize int64, bakNum int) error {
	return defaultLog.SetFile(fileName, maxSize, bakNum)
}

// implements go log
func Fatal(v ...interface{}) {
	defaultLog.Fatal(v...)
}

// implements go log
func Fatalf(format string, v ...interface{}) {
	defaultLog.Fatalf(format, v...)
}

// implements go log
func Fatalln(v ...interface{}) {
	defaultLog.Fatal(v...)
}

// implements go log
func Panic(v ...interface{}) {
	defaultLog.Panic(v...)
}

// implements go log
func Panicf(format string, v ...interface{}) {
	defaultLog.Panicf(format, v...)
}

// implements go log
func Panicln(v ...interface{}) {
	defaultLog.Panic(v...)
}

// implements go log
func Print(v ...interface{}) {
	defaultLog.Debug(v...)
}

// implements go log
func Printf(format string, v ...interface{}) {
	defaultLog.Debugf(format, v...)
}

// implements go log
func Println(v ...interface{}) {
	defaultLog.Debug(v...)
}

// implements go-log
func Critical(args ...interface{}) {
	defaultLog.Critical(args...)
}

// implements go-log
func Criticalf(format string, args ...interface{}) {
	defaultLog.Criticalf(format, args...)
}

// implements go-log
func Error(args ...interface{}) {
	defaultLog.Error(args...)
}

// implements go-log
func Errorf(format string, args ...interface{}) {
	defaultLog.Errorf(format, args...)
}

// implements go-log
func Warning(args ...interface{}) {
	defaultLog.Warning(args...)
}

// implements go-log
func Warningf(format string, args ...interface{}) {
	defaultLog.Warningf(format, args...)
}

// implements go-log
func Notice(args ...interface{}) {
	defaultLog.Notice(args...)
}

// implements go-log
func Noticef(format string, args ...interface{}) {
	defaultLog.Noticef(format, args...)
}

// implements go-log
func Info(args ...interface{}) {
	defaultLog.Info(args...)
}

// implements go-log
func Infof(format string, args ...interface{}) {
	defaultLog.Infof(format, args...)
}

// implements go-log
func Debug(args ...interface{}) {
	defaultLog.Debug(args...)
}

// implements go-log
func Debugf(format string, args ...interface{}) {
	defaultLog.Debugf(format, args...)
}
