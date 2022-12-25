package main

import "fmt"

// Recursive 1.20ではcompile error
type Recursive interface {
	f() interface{ Recursive }
}

type RecursiveImpl struct{}

func (r *RecursiveImpl) f() interface{ Recursive } {
	return r
}

func main() {
	var r RecursiveImpl
	fmt.Printf("%+v\n", r.f())
}
