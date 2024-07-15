# Listen multi http requests concurrently

## Goal
`net` パッケージを利用して複数のHTTPリクエストを並列に受け付ける簡易なHTTPサーバーを実装していきます。

ただし、接続を受け付けた後のハンドリングは実装しません。

## 1. HTTPサーバーのクライアント接続の実装

クライアントの接続を受け付けるHTTPサーバーを実装していきます。

### 1.1. 単一のリクエストを一度だけ受け付ける

最初は `net` パッケージを使用して、単一のHTTPリクエストを受け付けるサーバを作成してみます。

サーバーはクライアントの接続を受け付けると、HTTP status code 200 OK と空のJSONを返し、そのまま処理を終了します。

また複数のクライアントの接続を並列に受け付けることもできません。

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
	resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: " + strconv.Itoa(len(body)) + "\r\n\r\n" + body)
	if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
```

#### リクエストの実行
```bash
$ go run server.go
start server...
---
$ curl -XGET "http://localhost:8080"
{"status":"ok"}
```

### 1.2. 単一のリクエストを複数回受け付ける

1. で作成したHTTPサーバーが、一度リクエストを受け付けた後も再度リクエストを受け取れるようにします。

#### server.go

`lis.Accept()` でリクエストを受け付けレスポンスを返す処理を for ループで囲みます。

並列リクエストができないことを検証するため、リクエストを受け付けた後1秒待機するようにしています。

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
		resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: " + strconv.Itoa(len(body)) + "\r\n\r\n" + body)
		if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		// 4. 一度リクエストを処理したら、1秒待つ（並列リクエスト検証用）
		time.Sleep(1 * time.Second)
	}
}
```

#### リクエストの実行

`xargs` で `curl` を並列実行してみます。 

サーバー側は一度に1リクエストしか受け付けていないため、並列実行しても直列に結果が出力されます。

```bash
$ go run server.go
start server...
---
$ time seq 1 10 | xargs -I XX -P 10 curl "http://localhost:8080" 
{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}seq 1 10  0.00s user 0.00s system 66% cpu 0.006 total
xargs -I ZZ -P 10 curl "http://localhost:8080"  0.18s user 0.16s system 3% cpu 9.059 total
```

### 1.3. 複数リクエストを並列に受け付ける

1. で作成したHTTPサーバーが、複数のクライアントの接続を並列に受け付けるようにします。

#### server.go

ポイントは2.でリクエストを受け付けた後の部分になります。
リクエストを受け付けたのち、goroutineを使用して非同期にレスポンスを処理するようにします。
非同期処理の間に次のリクエストの受付を開始でき、新たな接続を開始できるため、複数のリクエストを並列に処理できます。

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
		go func(conn net.Conn) {

			// 3. レスポンスを返す
			body := `{"status":"ok"}`
			resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: " + strconv.Itoa(len(body)) + "\r\n\r\n" + body)
			if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			// 4. 一度リクエストを処理したら、1秒待つ（並列リクエスト検証用）
			time.Sleep(1 * time.Second)
		}(conn)
	}
}

```

#### リクエストの実行

サーバー側は一度に複数リクエストを受け付けられるようになったため、並列リクエストに対して即座に結果が返ります。

```bash
$ go run server.go
start server...
---
$ time seq 1 10 | xargs -I XX -P 10 curl "http://localhost:8080" 
{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}{"status":"ok"}seq 1 10  0.00s user 0.00s system 66% cpu 0.006 total
xargs -I ZZ -P 10 curl "http://localhost:8080"  0.18s user 0.16s system 3% cpu 9.059 total
```

---

なお以下のように2.3.を `go func` 内で実行してしまうと、一つのファイルディスクリプタに対する同時実行数をオーバーしてしまい、以下のpanicが発生します。

```go
panic: too many concurrent operations on a single file or socket (max 1048575)
```

#### server.go

```go
func main() {
	fmt.Println("start server...")

	// 1. net.Listenerを作成
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	for {
		go func() {
			// 2. リクエストを受け付ける
			conn, err := lis.Accept()
			if err != nil {
				panic(err)
			}

			// 3. レスポンスを返す
			body := `{"status":"ok"}`
			resp := []byte("HTTP/1.1 200 OK\r\nContent-Length: " + strconv.Itoa(len(body)) + "\r\n\r\n" + body)
			if _, err := io.Copy(conn, bytes.NewBuffer(resp)); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}()
	}
}
```
