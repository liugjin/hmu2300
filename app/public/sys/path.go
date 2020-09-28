package sys

import "os"

type Path string

func (p Path) ExpandEnv() string {
	return os.ExpandEnv(string(p))
}
