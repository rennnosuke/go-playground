package crud

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func TestGetProducts(t *testing.T) {
	tests := []struct {
		name    string
		want    []Product
		wantErr bool
	}{
		{
			name: "success",
			want: []Product{
				{ID: 1, Name: "test1", Price: 1000},
				{ID: 2, Name: "test2", Price: 2000},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			resultRows := sqlmock.NewRows([]string{"id", "name", "price"})
			for _, w := range tt.want {
				resultRows.AddRow(w.ID, w.Name, w.Price)
			}
			mock.ExpectQuery("SELECT (.+) FROM products").
				WillReturnRows(resultRows)

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
