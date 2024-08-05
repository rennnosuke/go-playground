# `http.Client` リクエストの前処理・後処理を挿入できるRoundTripper

## Blog
本ドキュメントの内容を以下の記事にまとめました。

[Ren's Blog](https://rennnosukesann.hatenablog.com/entry/2024/08/03/162216?_gl=1*r1ln9n*_gcl_au*MTc1NTM1MTc1Mi4xNzIyNjA3MTMz)

## What
`http.Client` のリクエストの前処理・後処理を挿入できるようにするための `http.RoundTripper` を実装します。

## Why
HTTPクライアントリクエストに付随して処理を実行したいケースに便利です。

e.g.
- リクエスト時のログ
- 特定条件下で特定のヘッダー値を書き換え
- endpointごとのrate limiting


## How
`http.RoundTripper`を実装した構造体を作成し、 `RoundTrip` メソッド実装で前処理・後処理を挿入できるようにします。

## 使い方

```go
var cli = http.Client{
	Transport: &DecorateRoundTripper{
		// wrapするRoundTripper
		base: http.DefaultTransport,
		// 前処理
		before: func(r *http.Request) error {
			slog.InfoContext(r.Context(), "before request")
			return nil
		},
		// 後処理
		after: func(r *http.Response) error {
			slog.InfoContext(r.Request.Context(), "after response")
			return nil
		},
	},
	Timeout: 30,
}
```

## 解説

- ① `DecorateRoundTripper` 構造体は `RoundTripper` インターフェースを実装しています。
フィールドには前処理を行う `before` 関数と後処理を行う `after` 関数を持ちます。
- ② `RoundTrip` メソッドは、リクエストを送信する前に `before` 関数を実行し、レスポンスを受信した後に `after` 関数を実行します。


```go
// ①
type DecorateRoundTripper struct {
	base   http.RoundTripper
	before func(*http.Request) error
	after  func(*http.Response) error
}

// ②
func (d *DecorateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if d.before != nil {
		if err := d.before(req); err != nil {
			return nil, err
		}
	}
	resp, err := d.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if d.after != nil {
		if err := d.after(resp); err != nil {
			return nil, err
		}
	}
	return resp, nil
```

## リクエスト・レスポンスボディを読み取る場合の注意点
- 前処理・後処理でリクエストやレスポンスを消費してしまった場合、その後の処理でボディを再利用することができません（再読み取りに対して `io.EOF` エラーが返ります）。
- 前処理でリクエストボディを読み取った上でwrapしたRoundTripperに渡す場合は、読み取った内容を一時保存し、 wrapしたRoundTripper.RoundTripの前にリクエストボディを再度設定する必要があります。
- 後処理でレスポンスボディを読み取った上でレスポンスをreturnする場合は、読み取った内容を一時保存し、 after実行後にレスポンスボディを再度設定する必要があります。

## 解決案：前処理・後処理でリクエスト・レスポンスボディをコピーし、読み取り後コピーを再設定する

### ① `c.before` 呼び出し前に、ボディの内容を読み取り一次保存

`req.GetBody()` で `req.Body` のコピーを取得し、一時保存します。<br>
`req.Body` 読み取り後に `req.GetBody()` を呼び出すとコピーを取得することができないので、必ず `req.Body` 読み取り処理前に実行してください。

`c.before` でリクエストボディを読み取った後、一次保存したボディのコピーを `req.Body` に再度設定します。

### ② `c.after` 呼び出し前に、ボディの内容を読み取り一次保存
`io.ReadAll` で `resp.Body` のコピーを取得し、一時保存します。

`io.ReadAll` の読み取りで `resp.Body` が閉じられるため、 一度保存したレスポンスボディバイナリを使用して `bytes.NewReader` を初期化し、 `io.NopCloser` でラップして再度設定します。

```go
func (c *ConsumableDecorateRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.before != nil {
		// ①
		rb, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		if err := c.before(req); err != nil {
			return nil, err
		}
		req.Body = rb
	}
	resp, err := c.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if c.after != nil {
		// ②
		rb, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body = io.NopCloser(bytes.NewReader(rb))
		if err := c.after(resp); err != nil {
			return nil, err
		}
	}
	return resp, nil
}
```

## Tips
### 後処理のあとは、`resp.Body`をすべて読み取り、接続を閉じることを忘れないようにする
通常のリクエスト同様、 `resp.Body` を全て読み取り、closeすることを忘れないようにしましょう。keep-alive TCP接続が再利用されない可能性があります（ `http.DefaultTransport` の場合）。

> Body を閉じるのは呼び出し側の責任です
。デフォルトの HTTP クライアントのトランスポートは、Body が
最後まで読み取られて閉じられていない場合、HTTP/1.x の "keep-alive" TCP 接続を再利用しない場合があります。

https://pkg.go.dev/net/http#Response.Body

```go
_, _ = io.ReadAll(io.Discard, resp.Body)
_ = resp.Body.Close()
```