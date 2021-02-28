package main

import (
	"io"
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "http.ResponseWriter sample")
}

func main() {
	http.HandleFunc("/sample", handleFunc)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
