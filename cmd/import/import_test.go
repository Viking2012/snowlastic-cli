package _import

import (
	"testing"
)

func Test_isUpper(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "All upper",
			args: args{s: "TEST"},
			want: true,
		},
		{
			name: "All lower",
			args: args{s: "test"},
			want: false,
		},
		{
			name: "Mixed",
			args: args{s: "Test"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUpper(tt.args.s); got != tt.want {
				t.Errorf("isUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quoteField(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Uppercase field",
			args: args{i: string("DATABASE")},
			want: "DATABASE",
		},
		{
			name: "Mixedcase field",
			args: args{i: string("Database")},
			want: `"Database"`,
		},
		{
			name: "Function with quoted argument",
			args: args{`YEAR("Closure Date")`},
			want: `YEAR("Closure Date")`,
		},
		{
			name: "Function with unquoted argument",
			args: args{`YEAR(Closure Date)`},
			want: `YEAR("Closure Date")`,
		},
		{
			name: "Function with proper argument",
			args: args{`YEAR(CLOSURE_DATE)`},
			want: `YEAR(CLOSURE_DATE)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteField(tt.args.i); got != tt.want {
				t.Errorf("quoteField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quoteParam(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{string("test")},
			want: "'test'",
		},
		{
			name: "int",
			args: args{int(1)},
			want: "1",
		},
		{
			name: "float",
			args: args{float32(1.1)},
			want: "1.100000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteParam(tt.args.i); got != tt.want {
				t.Errorf("quoteParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
