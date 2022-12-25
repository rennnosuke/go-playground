package main

import (
	"io"
	"testing"
)

func F[T comparable]() {}

type I[T comparable] interface {
	G(T)
}

func TestComparable(t *testing.T) {
	F[any]()
	F[io.Reader]()
	F[struct{ any }]()
	F[[1]any]()
	F[I[any]]()
}
