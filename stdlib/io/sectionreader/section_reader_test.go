package sectionreader

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestSectionReader(t *testing.T) {
	reader := strings.NewReader("hogehoge hogehoge fugafuga")
	sectionReader := io.NewSectionReader(reader, 10, 3)
	if _, err := io.Copy(os.Stdout, sectionReader); err != nil {
		t.Fatal(err)
	}
}
