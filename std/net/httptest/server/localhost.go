package server

import (
	"fmt"
	"net"
	"testing"
)

// serveTestHTTP
// httptest.Server 初期化時に net.Listener を生成する関数 httptest.newLocalListener とほぼ同じルールで net.Listener を生成する
// 内部的にはループバックアドレスに対するlistenerを生成しているだけ
// See: https://github.com/golang/go/blob/608acff8479640b00c85371d91280b64f5ec9594/src/net/http/httptest/server.go#L60
func serveTestHTTP(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1")
	if err != nil {
		l, err = net.Listen("tcp6", "[::1]:0")
		if err != nil {
			panic(fmt.Sprintf("failed to listen on a port: %v", err))
		}
	}
	return l
}
