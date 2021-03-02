package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
)

func TestParseHTTPResponse(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Write([]byte("GET / HTTP/1.0\r\nHost: localhost\r\n\r\n")); err != nil {
		t.Fatal(err)
	}
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Header)
	defer res.Body.Close()
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		t.Fatal(err)
	}
}
