package format

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func BenchmarkFormatPrintInteger(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(ioutil.Discard, "format print %d", i)
	}
}

func BenchmarkPrintIntegerConcatWithItoa(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprint(ioutil.Discard, "format print "+strconv.Itoa(i))
	}
}

func BenchmarkPrintIntegerConcatWithFormatInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprint(ioutil.Discard, "format print "+strconv.FormatInt(int64(i), 10))
	}
}
