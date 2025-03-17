package dir1

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestOSOpenTraversal(t *testing.T) {
	// 相対ディレクトリコンポーネントを利用して、目的のパスから逸脱できる
	trustedLocation := "."
	f, err := os.Open(filepath.Join(trustedLocation, "../dir2/dummy.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := io.ReadAll(f); string(b) != "Traversal!" {
		t.Errorf("want: Traversal!, got: %s", string(b))
	}
}

func TestOSOpenRootTraversal(t *testing.T) {
	root, err := os.OpenRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	// os.RootによるOpenは、ルートディレクトリからの相対パスを受け入れない
	_, err = root.Open("../dir2/dummy.txt")
	if wantErr := (*os.PathError)(nil); !errors.As(err, &wantErr) {
		t.Errorf("want: os.PathError, got: %T", err)
	}
}
