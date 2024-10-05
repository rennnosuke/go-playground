package conn

import (
	"context"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

func Test_open(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf mysql.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success to connect to the database",
			args: args{
				ctx: context.Background(),
				conf: mysql.Config{
					User:      "testuser",
					Passwd:    "testpass",
					Addr:      "localhost:3306",
					DBName:    "test",
					Collation: "utf8mb4_general_ci",
					Loc:       time.UTC,
					Timeout:   time.Second * 30,
				},
			},
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
		})
	}
}
