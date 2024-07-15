package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func main() {
	fmt.Println("start server...")

	// 1. net.Listenerを作成
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	for {
		// 2. リクエストを受け付ける
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {

			// 3. レスポンスを返す
			body := `{"status":"ok"}`
			resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: " + strconv.Itoa(len(body)) + "\r\n\r\n" + body)
			if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			// 4. 一度リクエストを処理したら、1秒待つ（並列リクエスト検証用）
			time.Sleep(1 * time.Second)
		}(conn)
	}
}
