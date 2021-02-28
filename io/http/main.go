package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "json")

	gz := gzip.NewWriter(w)
	defer func() {
		if err := gz.Flush(); err != nil {
			panic(err)
		}
	}()

	mw := io.MultiWriter(gz, os.Stdout)
	je := json.NewEncoder(mw)

	source := map[string]string{
		"hello": "'world",
	}

	je.SetIndent("", "  ")
	if err := je.Encode(source); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
