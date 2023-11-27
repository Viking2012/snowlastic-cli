package _import

import "testing"

func Test_SegmentIsGiven(t *testing.T) {
	type args struct {
		segment       any
		givenSegments []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bool/true/lowercase",
			args: args{
				segment:       true,
				givenSegments: []string{"true"},
			},
			want: true,
		},
		{
			name: "bool/true/uppercase",
			args: args{
				segment:       true,
				givenSegments: []string{"TRUE"},
			},
			want: true,
		},
		{
			name: "bool/true/numeric",
			args: args{
				segment:       true,
				givenSegments: []string{"1"},
			},
			want: true,
		},
		{
			name: "bool/false",
			args: args{
				segment:       true,
				givenSegments: []string{"0", "false", "FALSE"},
			},
			want: false,
		},
		{
			name: "int/true",
			args: args{
				segment:       1,
				givenSegments: []string{"1", "2"},
			},
			want: true,
		},
		{
			name: "int/false",
			args: args{
				segment:       1,
				givenSegments: []string{"2", "3", "4"},
			},
			want: false,
		},
		{
			name: "uint/true",
			args: args{
				segment:       uint(1),
				givenSegments: []string{"1", "2"},
			},
			want: true,
		},
		{
			name: "uint/false",
			args: args{
				segment:       uint(1),
				givenSegments: []string{"2", "3", "4"},
			},
			want: false,
		},
		{
			name: "float/true",
			args: args{
				segment:       1.001,
				givenSegments: []string{"1.001"},
			},
			want: true,
		},
		{
			name: "float/false",
			args: args{
				segment:       1.001,
				givenSegments: []string{"1", "2", "1.002"},
			},
			want: false,
		},
		{
			name: "complex/true",
			args: args{
				segment:       complex(1.001, 2.002),
				givenSegments: []string{"1.001+2.002i"},
			},
			want: true,
		},
		{
			name: "complex/false",
			args: args{
				segment:       complex(1.001, 2.002),
				givenSegments: []string{"1.001+2.003i"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SegmentIsGiven(tt.args.segment, tt.args.givenSegments); got != tt.want {
				t.Errorf("segmentIsGiven() = %v, want %v", got, tt.want)
			}
		})
	}
}
