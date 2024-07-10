package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		panic(err)
	}
}
