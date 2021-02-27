package main

import (
	"bufio"
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
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if _, err := io.WriteString(conn, scanner.Text()); err != nil {
				panic(err)
			}
		}
	}
}
