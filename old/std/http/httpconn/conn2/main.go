package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer func() {
				_ = conn.Close()
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			_, err := io.WriteString(conn, strings.Join([]string{
				"HTTP/1.1 200 OK",
				"Content-Type: text/plain",
				"",
				"{\"message\":\"hogehoge\"}",
			}, "\r\n"))
			if err != nil {
				panic(err)
			}
		}()
	}
}
