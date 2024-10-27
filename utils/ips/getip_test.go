package ips

import (
	"fmt"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	addr, err := GetLocalIP()
	if err != nil {
		return
	}
	fmt.Println(addr)
}
