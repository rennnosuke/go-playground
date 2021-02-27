package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var addr = "localhost:8090"

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("client connected to %s:\n", addr)

	for {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			panic(err)
		}
	}
}
