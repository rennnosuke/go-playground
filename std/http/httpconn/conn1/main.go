package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "localhost:8888")
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
				conn.Close()
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			resp := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Body:       io.NopCloser(strings.NewReader("{\"msg\": \"message\"}")),
			}
			resp.Write(conn)
		}()
	}
}
