# ListenAndServe without helper

## Goal
`http.ListenAndServe` を使わずに、`net.Listener` と `http.Server` を使ってサーバを起動します。

## Implementation

### 1. 単一のリクエストを一度だけ受け付ける

最初は `net` パッケージを使用して、単一のHTTPリクエストを受け付けるサーバを作成してみます。

サーバーはクライアントの接続を受け付けると、HTTP status code 200 OK と空のJSONを返し、そのまま処理を終了します。

#### server.go

1. `net.Listen` で `net.Listener` を作成します。
<br> `net.Listener` は、クライアントからの接続を受け付けるためのインターフェースです。
<br> `net.Listen` は、指定されたネットワークとアドレスにリスナーを作成します。
<br> 今回は第一引数に `tcp` を指定しているため、TCPプロトコルを使用する `net.TCPListener` が返されます。


2. `net.Listener.Accept` でリクエストを受け付けます。
<br>リスエストを受け付けると、クライアント接続への接続 `net.Conn` が返されます。
<br>`net.Conn` はインターフェースであり、今回はtcp接続で `net.TCPListener` が使用されているので `net.TCPConn` が返されます。


3. クライアントにレスポンスを返します。 `io.Copy` で `net.Conn` に対してレスポンスを書き込みます。

```go
func main() {
	fmt.Println("start server...")

	// 1. net.Listenerを作成
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	// 2. リクエストを受け付ける
	conn, err := lis.Accept()
	if err != nil {
		panic(err)
	}

	// 3. レスポンスを返す
	body := `{"status":"ok"}`
	resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\n" + body)
	if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
```

#### client.go
立ち上げたサーバーに対するクライアントです。標準出力にレスポンスボディを書き込みます。
```go
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
```

#### 実行結果
```bash
$ go run server.go
start server...
---
$ go run client.go
{"status":"ok"}
```

### 2. 単一のリクエストを複数回受け付ける

1. で作成したHTTPサーバーが、一度リクエストを受け付けた後も再度リクエストを受け取れるようにします。

#### server.go

`lis.Accept()` でリクエストを受け付けレスポンスを返す処理を for ループで囲みます。

```go
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

		// 3. レスポンスを返す
		body := `{"status":"ok"}`
		resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\n" + body)
		if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

```
