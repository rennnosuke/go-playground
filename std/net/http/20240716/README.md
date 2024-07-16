# Inspect HTTP Client fields

```go
type Client struct {
    Transport     RoundTripper
    CheckRedirect func(req *Request, via []*Request) error
    Jar           CookieJar
    Timeout       time.Duration
}
```

## Transport

単一のHTTPリクエストを実行するときの処理を定義する、 `http.RoundTripper` インターフェース型フィールドです。

（ほとんどの場合） `Transport` はTCPコネクションをキャッシュしているため、再利用される必要があります。
`Transport` はインターフェースなので、具象型のフィールドでこのキャッシュが管理される必要があります。

`Transport` が `nil` の場合、`http.DefaultTransport`が使用されます。

### DefaultTransport

`http.DefaultTransport` は、デフォルトのHTTPクライアントのTransportです。

`http.Client` で `Transport` フィールドが指定されていない場合、このTransportが使用されます。

また `http.DefaultClient` でも使用されます。

内部ではTCPコネクションをキャッシュします。キャッシュ可能な最大数は `MaxIdleConns` で設定される100ですが、
ホストあたりの接続数は `MaxIdleConnsPerHost` で設定される2であるため注意が必要です。

```go
var DefaultTransport RoundTripper = &Transport{
	Proxy: ProxyFromEnvironment,
	DialContext: defaultTransportDialContext(&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}),
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
```


## CheckRedirect

HTTPクライアントのリダイレクトポリシーを定義します。

`CheckRedirect` がnullでない場合、HTTPリダイレクトの前に指定した関数が実行されます。

`req` は最新のリクエストを、 `via` はリダイレクト途中に形成されたリクエストを表します。

`CheckRedirect` がエラーを返す時、以下のパターンがあります。

- `GET` メソッドの場合、リダイレクトは中止され、エラー一つ前のclose済み `Response` と `url.Error` でwrapされたエラーを返します。
- `ErrUseLastResponse` を返す場合、最後のリクエストの `Response` が返されます。このレスポンスはcloseされていません。

`CheckRedirect` が `nil` の場合 `http.defaultCheckRedirect` が使用され、10回までリダイレクト可能になります。

## Jar

HTTPクッキーを管理するための `http.CookieJar` インターフェース型フィールドです。

外部リクエストで使用するクッキーを管理するのに使用します。

クッキーはHTTPクライアントのインバウンドレスポンスごとに更新されます。

またリダイレクトのたびに `Jar` は更新されます。

`Jar` がnilの場合、 `Request` に明示的にクッキーが設定された場合にのみクッキーが送信されるようになります。

## Timeout

クライアントリクエストのタイムリミットを設定します。

タイムアウトの範囲は `Client.Do` 実行後から `Response.Body` の読み取りまでです。 

`Timeout` がゼロの場合、タイムアウトは無効になります。

デフォルトでは、 `Client.Timeout` はゼロです。すなわち、タイムアウトは無効になります。

タイムアウトした場合、クライアントはコンテキストキャンセルと同じようにTransport接続を終了します。


以前は `CancelRequest` メソッドを使用してキャンセル処理を実装していましたが、すでに非推奨となっています。
代わりにリクエストに設定した `context.Context` を経由したキャンセルが推奨されています。