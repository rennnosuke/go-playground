package main

import (
	"github.com/rennnosuke/skeleton_example"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(skeleton_example.Analyzer)
}
