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
			name:    "select#01",
			builder: sq.Select("*").From("users").Where(sq.Eq{"id": 1}),
			want:    "SELECT * FROM users WHERE id = ?",
		},
		{
			name:    "select#02",
			builder: sq.Select("id").From("users").Where(sq.Eq{"id": 1}),
			want:    "SELECT id FROM users WHERE id = ?",
		},
		{
			name:    "select#03 - `in` phrase",
			builder: sq.Select("*").From("users").Where(sq.Eq{"id": []int{1, 2, 3}}),
			want:    "SELECT * FROM users WHERE id IN (?,?,?)",
		},
		{
			// See: https://github.com/Masterminds/squirrel/pull/180#issuecomment-778403270
			name:    "select#04 - `between` phrase",
			builder: sq.Select("*").From("users").Where(sq.Expr("id BETWEEN ? AND ?", 1, 10)),
			want:    "SELECT * FROM users WHERE id BETWEEN ? AND ?",
		},
		{
			name:    "select#05 - count",
			builder: sq.Select("COUNT(id OR NULL)").From("users").Where(sq.Expr("name like 'hoge%'", 1, 10)),
			want:    "SELECT COUNT(id OR NULL) FROM users WHERE name like 'hoge%'",
		},
		{
			name:    "insert#1",
			builder: sq.Insert("users").Columns("id").Values("1"),
			want:    "INSERT INTO users (id) VALUES (?)",
		},
		{
			name:    "insert#1",
			builder: sq.Insert("users").Values("1"),
			want:    "INSERT INTO users VALUES (?)",
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
