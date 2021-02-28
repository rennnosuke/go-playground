package main

import (
	"io"
	"os"
)

func main() {
	file, err := os.Create("file.txt")
	if err != nil {
		panic(err)
	}
	writer := io.MultiWriter(file, os.Stdout)
	if _, err := io.WriteString(writer, "multi writing"); err != nil {
		panic(err)
	}
}
