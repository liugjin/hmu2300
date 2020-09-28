package muid

import (
	"fmt"
	"testing"
)

func TestGetMuID(t *testing.T) {
	id := GetMuID()
	fmt.Println(id)
}
