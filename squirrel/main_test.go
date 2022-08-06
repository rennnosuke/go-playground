package squirrel

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
)

func TestSquirrel(t *testing.T) {
	t.Log("start squirrel examples")

	// select
	q, _, err := sq.Select(
		"*",
	).From(
		"users",
	).Where(sq.Eq{"id": 1}).ToSql()
	if err != nil {
		t.Fatal(err)
	}
	want := "SELECT * FROM users WHERE id = ?"
	t.Log(want)
	if diff := cmp.Diff(want, q); diff != "" {
		t.Errorf("MyFunc() mismatch (-want +got):\n%s", diff)
	}

}
