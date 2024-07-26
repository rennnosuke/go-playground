# Shutting down server gracefully

goのHTTPサーバーをGraceful shutdownします。

## Graceful Shutdownとは

Graceful shutdownは、サーバーが終了する際に、リクエストを完了させるための方法です。

サーバープロセスに対してプロセスをKILLするシグナルが発行されても、リクエストを安全に完了させてからシステムを終了できます。
これにより、サーバー終了のたびに接続中のクライアントがエラーレスポンスを受け取ってしまうことを防ぎます。

近年使用されることの多いオーケストレーションシステム下のコンテナフリート環境では、リソース使用状況の変化やデプロイによって頻繁にサーバープロセスが終了します。
このような環境下では プロセスの終了をGracefulにすることで、クライアントへのエラーを減らしより信頼性の高いシステムを構築できます。

（アプリケーション前段のリバースプロキシやサービスメッシュなどがそのあたりを代理してくれる場合もありますが、アプリケーションでもやっておくと安全です）

## サンプルコード

```go
	// catch SIGTERM, SIGINT, SIGKILL
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	// sync until signal received
	<-ctx.Done()

	// shutting down server within 1min - otherwise force shutdown
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	fmt.Printf("Shutdown server...\n")

	if err := svr.Shutdown(ctx); err != nil {
		fmt.Printf("failed to shutdown server: %v\n", err)
	}
```

### 検知したいシグナルを設定する

`signal.NotifyContext` は指定したシグナルを検知するためのコンテキストを返します。

`stop` を呼び出すと、シグナルの検知を終了することができます。

下記コードでは以下の3つのシグナルを検知します。

- SIGTERM
- SIGINT
- SIGKILL

```go
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
    defer stop()
```

### シグナルを検知する

`ctx.Done()` から返礼されるチャネルを待ち受けることで、シグナルを検知します。

```go
    <-ctx.Done()
```

### シャットダウン処理

シグナルを検知した後は、シグナルがプロセスを終了させるプロセスであっても実行プログラムを終了せず、シャットダウン処理に移行します。

`svr.Shutdown(ctx)` は、接続を全てcloseした上でサーバーをシャットダウンします。
現在接続している全てのリクエストが処理を終了するまで待機します。

ただし、コンテキストに設定したタイムアウトを超過した場合、接続がアイドル状態でなくても切断します。

> シャットダウンは、アクティブな接続を中断することなく、サーバーを正常にシャットダウンします。シャットダウンは、まず開いているすべてのリスナーを閉じ、次にアイドル状態のすべての接続を閉じ、接続がアイドル状態に戻るまで無期限に待機してからシャットダウンすることで機能します。シャットダウンが完了する前に提供されたコンテキストの有効期限が切れた場合、シャットダウンはコンテキストのエラーを返します。それ以外の場合は、サーバーの基になるリスナーを閉じることによって返されたエラーを返します。

```go
    ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()

    fmt.Printf("Shutdown server...\n")

    if err := svr.Shutdown(ctx); err != nil {
        fmt.Printf("failed to shutdown server: %v\n", err)
    }
```









