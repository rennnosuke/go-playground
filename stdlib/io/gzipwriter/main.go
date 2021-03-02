package main

import (
	gzip2 "compress/gzip"
	"io"
	"os"
)

func main() {
	file, err := os.Create("file.txt.gz")
	if err != nil {
		panic(err)
	}
	gzip := gzip2.NewWriter(file)
	defer gzip.Close()
	if _, err := io.WriteString(gzip, "write string to gzip file"); err != nil {
		panic(err)
	}
}
