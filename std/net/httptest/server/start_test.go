package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHTTPGetWithUnstartedServer tests the httptest.NewUnstartedServer function.
// httptest.NewUnstartedServer は、開始されていないサーバーを返す.
// そのためテストサーバーをproxyしたテストを実施する場合 ts.Start() を呼び出す必要がある.
func TestHTTPGetWithUnstartedServer(t *testing.T) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, World")
	}))
	ts.Start() // serverは起動していないため、明示的に開始する
	defer ts.Close()

	testHTTPGet(t, ts)
}

// TestHTTPGetWithStartedServer tests the httptest.NewServer function.
// httptest.NewServer は、開始されたサーバーを返す.
// そのためテストサーバーを明示的にStartする必要はない.
func TestHTTPGetWithStartedServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, World")
	}))
	defer ts.Close()

	testHTTPGet(t, ts)
}

func testHTTPGet(t *testing.T, ts *httptest.Server) {
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status code error: %v", res.StatusCode)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	if got, want := buf.String(), "Hello, World\n"; got != want {
		t.Fatalf("body = %q; want %q", got, want)
	}
}
