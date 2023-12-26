package crud

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFindAll(t *testing.T) {
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
				db: getDB(t, func(m sqlmock.Sqlmock) {
					m.ExpectQuery("SELECT (.+) FROM products").
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at"}).
								AddRow(1, "test1", 1000, time.Now(), time.Now()).
								AddRow(2, "test2", 2000, time.Now(), time.Now()),
						)
				}),
			},
			want: []Product{
				{ID: 1, Name: "test1", Price: 1000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{ID: 2, Name: "test2", Price: 2000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
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
			got, err := FindAll(context.Background(), db)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateApproxTime(time.Second)); diff != "" {
				t.Errorf("FindAll() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *sql.DB
		p   *Product
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				db: getDB(t, func(m sqlmock.Sqlmock) {
					m.ExpectExec("INSERT INTO products (.+) VALUES (.+)").
						WithArgs("test", 1000, ApproxTime{}, ApproxTime{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}),
				p: &Product{Name: "test", Price: 1000},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "db is nil",
			args: args{
				ctx: context.Background(),
				db:  nil,
				p:   nil,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.ctx, tt.args.db, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateApproxTime(time.Second)); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *sql.DB
		p   *Product
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				db: getDB(t, func(m sqlmock.Sqlmock) {
					m.ExpectExec("UPDATE products SET name = (.+), price = (.+), updated_at = (.+) WHERE id = (.+)").
						WithArgs("updated", 10000, ApproxTime{}, 1).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}),
				p: &Product{
					ID:    1,
					Name:  "updated",
					Price: 10000,
				},
			},
			wantErr: false,
		},
		{
			name: "db is nil",
			args: args{
				ctx: context.Background(),
				db:  nil,
				p:   nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Update(tt.args.ctx, tt.args.db, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		db  *sql.DB
		id  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				db: getDB(t, func(m sqlmock.Sqlmock) {
					m.ExpectExec("DELETE FROM products WHERE id = (.+)").
						WithArgs(1).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}),
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "db is nil",
			args: args{
				ctx: context.Background(),
				db:  nil,
				id:  0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.ctx, tt.args.db, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getDB(t *testing.T, fn func(m sqlmock.Sqlmock)) *sql.DB {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	fn(mock)
	return db
}

type ApproxTime time.Time

// Match satisfies sqlmock.Argument interface
func (a ApproxTime) Match(v driver.Value) bool {
	t, ok := v.(time.Time)
	if !ok {
		return false
	}
	return cmp.Diff(a, t, cmpopts.EquateApproxTime(time.Second)) != ""
}
