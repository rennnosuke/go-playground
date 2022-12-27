# Go 1.20 New Feaature

## Skip `go test`

[change](https://go-review.googlesource.com/c/go/+/421439/4/src/cmd/go/testdata/script/test_skip.txt#21)

```shell
go test -v -run Test -skip skip_test.go
```

## Skip `go generate`

[change](https://go-review.googlesource.com/c/go/+/421440)

```shell
go generate -skip th..sand './generate/flag.go'
```

## `vet` detect loop var in multiple nests

```shell
go vet ./vetloopvar
```

### ./vetloopvar/vetloopvar.go

```go
package main

import "fmt"

func main() {
	seq := []int{1, 2, 3}
	for k, v := range seq {
		func() {
			fmt.Println(k, v)
		}()
	}
}
```