package util

import (
	"fmt"
	"math/rand"
	"time"
)

// Return
// return six bit num code
func RandCodeInt(debug bool) int64 {
	if debug {
		return 1234567890
	}

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	code := r.Int63n(9999999999)
	if code < 10000000000 {
		code += 10000000000
	}
	return code
}

func RandCodeStr(debug bool) string {
	return fmt.Sprintf("%09d", RandCodeInt(debug))
}
