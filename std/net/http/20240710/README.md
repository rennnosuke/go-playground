# ListenAndServe without helper

## Goal
- `http.ListenAndServe` を使わずに、`net.Listener` と `http.Server` を使ってサーバを起動する。

## 手順

### 1. `net.Listener` を使用して、リクエストを受け付ける


