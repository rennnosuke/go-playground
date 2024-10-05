package conn

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_open(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf mySQLConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		{
			name: "connection to the database",
			args: args{
				ctx: context.Background(),
				conf: mySQLConfig{
					userName: "root",
					password: "root",
					dbName:   "test",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := open(tt.args.ctx, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := got.Ping(); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("open() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
