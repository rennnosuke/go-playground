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

デフォルトでは、`http.DefaultTransport`が使用されます。

## CheckRedirect

HTTPクライアントのリダイレクトポリシーを定義します。

`CheckRedirect` がnullでない場合、HTTPリダイレクトの前に指定した関数が実行されます。

`req` は最新のリクエストを、 `via` はリダイレクト途中に形成されたリクエストを表します。

`CheckRedirect` がエラーを返す時、以下のパターンがあります。

- `GET` メソッドの場合、リダイレクトは中止され、エラー一つ前のclose済み `Response` と `url.Error` でwrapされたエラーを返します。
- `ErrUseLastResponse` を返す場合、最後のリクエストの `Response` が返されます。このレスポンスはcloseされていません。

デフォルトでは `http.defaultCheckRedirect` が使用され、10回までリダイレクト可能になります。
