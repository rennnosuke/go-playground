package format

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func BenchmarkFormatPrintString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(ioutil.Discard, "format print %s", "string")
	}
}

func BenchmarkPrintConcatString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprint(ioutil.Discard, "format print "+"string")
	}
}

func BenchmarkPrintSixConcatString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprint(ioutil.Discard, "format print "+"s"+"t"+"r"+"i"+"n"+"g")
	}
}
