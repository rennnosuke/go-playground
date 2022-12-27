package main

import (
	"fmt"
	"testing"
)

func TestVetDetectLoopVarInFunc(t *testing.T) {
	seq := []int{1, 2, 3}
	for k, v := range seq {
		func() {
			fmt.Println(k, v)
		}()
	}
}
