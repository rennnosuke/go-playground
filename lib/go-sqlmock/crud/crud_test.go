package crud

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func TestGetProducts(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    []Product
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				db: getDB(t, "SELECT (.+) FROM products",
					sqlmock.NewRows([]string{"id", "name", "price"}).
						AddRow(1, "test1", 1000).
						AddRow(2, "test2", 2000),
				),
			},
			want: []Product{
				{ID: 1, Name: "test1", Price: 1000},
				{ID: 2, Name: "test2", Price: 2000},
			},
			wantErr: false,
		},
		{
			name: "db is nil",
			args: args{
				ctx: context.Background(),
				db:  nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.args.db
			defer func() {
				if db != nil {
					db.Close()
				}
			}()
			got, err := GetProducts(context.Background(), db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetProducts() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func getDB(t *testing.T, expectQuery string, returnRows *sqlmock.Rows) *sql.DB {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(expectQuery).WillReturnRows(returnRows)
	return db
}
