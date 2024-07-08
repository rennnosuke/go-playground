package main

import (
	"testing"
)

// Recursive 1.20ではcompile error
type Recursive interface {
	f() interface{ Recursive }
}

type RecursiveImpl struct{}

func (r *RecursiveImpl) f() interface{ Recursive } {
	return r
}

func TestIFCycleBuildErr(t *testing.T) {
	var r RecursiveImpl
	t.Logf("%+v\n", r.f())
}
