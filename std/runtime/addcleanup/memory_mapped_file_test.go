package addcleanup

import (
	"reflect"
	"testing"
)

func TestNewMemoryMappedFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *MemoryMappedFile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMemoryMappedFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemoryMappedFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryMappedFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
