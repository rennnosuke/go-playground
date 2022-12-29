package main

import (
	"fmt"
	"testing"
)

func TestBytes_Clone(t *testing.T) {
	b := []byte("short")
	//b2 := bytes.Clone(b)
	fmt.Println(b)
	//fmt.Println(b2)
}
