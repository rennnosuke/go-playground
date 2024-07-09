package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/net/websocket"
)

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Fatalln(err)
	}
}

func echoHandler(c *websocket.Conn) {
	if _, err := io.Copy(c, c); err != nil {
		log.Fatalln(err)
	}
}
