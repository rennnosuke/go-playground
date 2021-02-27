# Benchmark for fmt.Fprintf/fmt.Fprint

`fmt.Fprintf()` を使用して
formatへ値を埋め込み出力する場合と、 `fmt.Fprint()` を使用して、 `strconv.Itoa` で整数値の変換+連結行った文字列を出力する場合のベンチマークをとりました。

# Result
```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/rennnosuke/go-playground/measurement/format
cpu: Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz
BenchmarkFormatPrintString-12                  	22885459	        53.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkPrintConcatString-12                  	26645217	        41.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkPrintSixConcatString-12               	28030809	        42.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkFormatPrintInteger-12                 	14126896	        89.88 ns/op	       8 B/op	       0 allocs/op
BenchmarkPrintIntegerConcatWithItoa-12         	 8373324	       147.7 ns/op	      47 B/op	       3 allocs/op
BenchmarkPrintIntegerConcatWithFormatInt-12    	 7169640	       146.7 ns/op	      47 B/op	       3 allocs/op
PASS
ok  	github.com/rennnosuke/go-playground/measurement/format	8.833s
```

先3つは、文字列リテラルを `fmt.Fprintf` でformatで埋め込んで出力するか、 `fmt.Fprint` で文字列連結して出力するかの比較です。若干連結の方が高速ですが、そこまで大きな差はありませんでした。結合数を1つ/6つの場合で試してみましたが、この場合も実行速度はそこまで変わらないようです。

後3つでは、整数値を文字列に含めた場合を試してみました。整数値を文字列に連結するパターンでは整数→文字列変換を実施する関数を使う必要がありますが、これらを使うと低速なだけでなくメモリアロケーションも発生してしまうことがわかります。このような場合、素直にformatを使用したほうが良いことがわかります。

