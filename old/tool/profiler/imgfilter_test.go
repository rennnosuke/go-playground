package profiler

import (
	"reflect"
	"testing"
)

func TestFlipFilter(t *testing.T) {
	type args struct {
		img [][][3]uint8
	}
	tests := []struct {
		name string
		args args
		want [][][3]uint8
	}{
		{
			name: "flip",
			args: args{
				img: [][][3]uint8{
					{
						{255, 0, 0},
						{0, 255, 0},
						{0, 0, 255},
					},
				},
			},
			want: [][][3]uint8{
				{
					{0, 255, 255},
					{255, 0, 255},
					{255, 255, 0},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FlipFilter(tt.args.img); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlipFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFlipFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FlipFilter([][][3]uint8{
			{
				{255, 0, 0},
				{0, 255, 0},
				{0, 0, 255},
			},
		})
	}
}
