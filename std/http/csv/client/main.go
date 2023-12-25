package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	tr := &http.Transport{
		DisableCompression:  true, // disable compression as gzip when no Accept-Encoding header
		MaxIdleConnsPerHost: 100,  // max idle keep-alive connections per host
		IdleConnTimeout:     30,   // keep-alive timeout
	}
	cli := http.Client{
		Transport: tr,
		Timeout:   time.Second,
	}
	csv := "1,2,3,4,5"
	body := io.NopCloser(strings.NewReader(csv))
	header := http.Header(map[string][]string{
		"Content-Type": {"text/csv"},
	})
	path, err := url.Parse("http://localhost:8080/")
	if err != nil {
		panic(err)
	}
	req := http.Request{
		Method:        "POST",
		URL:           path,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        header,
		Body:          body,
		ContentLength: int64(len(csv)),
	}
	resp, err := cli.Do(&req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b := make([][]string, 0)
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		panic(err)
	}

	// output response
	fmt.Printf("output:\n%s", b)
}
