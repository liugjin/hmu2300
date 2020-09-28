package at

type AT interface {
	Do(cmd string) (string, error)

	// 检查网络模块是否正常
	AT() error

	// 检查sim卡是否存在
	CPIN() (string, error)

	// 取信号强度
	CSQ() (string, error)
}
