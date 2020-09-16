package ch8

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

func TestGoRoutine_Clock(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		t.Fatalf("failed to listen tcp: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		now := time.Now().Format("15:04:05\n")
		_, err := io.WriteString(c, now)
		if err != nil {
			return
		}
		fmt.Printf("Send: %s", now)
		time.Sleep(1 * time.Second)
	}
}
