package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var addr = "localhost:8090"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("start server %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if _, err := fmt.Fprintf(os.Stderr, "error occured: %s", err); err != nil {
				panic(err)
			}
			continue
		}

		go func(conn net.Conn) {
			fmt.Println("server is connected to client:")

			for {
				if _, err := io.Copy(os.Stdout, conn); err != nil {
					_ = fmt.Errorf("error occured: %s", err)
					break
				}
			}
		}(conn)
	}
}

func handleConnection(conn net.Conn) error {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buf[:n]))
	}
	return nil
}
