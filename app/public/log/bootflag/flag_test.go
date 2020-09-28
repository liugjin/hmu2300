package bootflag

import (
	"fmt"
	"testing"
)

func TestGetFlag(t *testing.T) {
	flag, err := GetFlag()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(flag)
}
