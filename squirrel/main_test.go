package squirrel

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
)

func hoge() {

}

func TestSquirrel(t *testing.T) {
	t.Log("start squirrel examples")

	tests := []struct {
		name    string
		builder interface {
			ToSql() (string, []interface{}, error)
		}
		want string
	}{
		{
			name:    "select#1",
			builder: sq.Select("*").From("users").Where(sq.Eq{"id": 1}),
			want:    "SELECT * FROM users WHERE id = ?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, args, err := tt.builder.ToSql()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("query: %s", got)
			t.Logf("args: %v", args)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("squirrel: built sql mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
