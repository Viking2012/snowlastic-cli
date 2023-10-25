package snowflake

import "testing"

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
			name: "uppercase",
			args: args{"TEST"},
			want: true,
		},
		{
			name: "lowercase",
			args: args{"test"},
			want: false,
		},
		{
			name: "Proper",
			args: args{"Test"},
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

func Test_needsQuoting(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "uppercase, no spaces",
			args:    args{"TESTFIELD"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "lowercase, no spaces",
			args:    args{"testfield"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "proper case, no spaces",
			args:    args{"TestField"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "uppercase, with spaces",
			args:    args{"TEST FIELD"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "lowercase, with spaces",
			args:    args{"test field"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "proper case, with spaces",
			args:    args{"Test Field"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := needsQuoting(tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("needsQuoting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("needsQuoting() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoteValue(t *testing.T) {
	type args struct {
		i any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int",
			args: args{1},
			want: "1",
		},
		{
			name: "int8",
			args: args{int8(1)},
			want: "1",
		},
		{
			name: "int16",
			args: args{int16(1)},
			want: "1",
		},
		{
			name: "int32",
			args: args{int32(1)},
			want: "1",
		},
		{
			name: "int64",
			args: args{int64(1)},
			want: "1",
		},
		{
			name: "float32, short decimal",
			args: args{float32(2.3)},
			want: "2.300000",
		},
		{
			name: "float64, short decimal",
			args: args{float64(2.3)},
			want: "2.300000",
		},
		{
			name: "float32, long decimal",
			args: args{float32(2.3000001)},
			want: "2.3000001",
		},
		{
			name: "float64, long decimal",
			args: args{float64(2.3000001)},
			want: "2.3000001",
		},
		{
			name: "string, no spaces",
			args: args{"test"},
			want: "'test'",
		},
		{
			name: "string, spaces",
			args: args{"test value"},
			want: "'test value'",
		},
		{
			name: "datetime as string",
			args: args{"2006-01-02"},
			want: "'2006-01-02'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuoteValue(tt.args.i); got != tt.want {
				t.Errorf("QuoteValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoteIdentifier(t *testing.T) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "uppercase, no spaces",
			args: args{"TESTFIELD"},
			want: "TESTFIELD",
		},
		{
			name: "lowercase, no spaces",
			args: args{"testfield"},
			want: `"testfield"`,
		},
		{
			name: "proper case, no spaces",
			args: args{"TestField"},
			want: `"TestField"`,
		},
		{
			name: "uppercase, with spaces",
			args: args{"TEST FIELD"},
			want: `"TEST FIELD"`,
		},
		{
			name: "lowercase, with spaces",
			args: args{"test field"},
			want: `"test field"`,
		},
		{
			name: "proper case, with spaces",
			args: args{"Test Field"},
			want: `"Test Field"`,
		},
		// later
		{
			name: "parenthetical, uppercase, no spaces, unquoted",
			args: args{"YEAR(TESTFIELD)"},
			want: "YEAR(TESTFIELD)",
		},
		{
			name: "parenthetical, lowercase, no spaces, unquoted",
			args: args{"YEAR(testfield)"},
			want: `YEAR("testfield")`,
		},
		{
			name: "parenthetical, uppercase, with spaces, unquoted",
			args: args{"YEAR(TEST FIELD)"},
			want: `YEAR("TEST FIELD")`,
		},
		{
			name: "parenthetical, lowercase, with spaces, unquoted",
			args: args{"YEAR(test field)"},
			want: `YEAR("test field")`,
		},
		{
			name: "parenthetical, uppercase, with spaces, quoted",
			args: args{`YEAR("TEST FIELD")`},
			want: `YEAR("TEST FIELD")`,
		},
		{
			name: "parenthetical, lowercase, with spaces, quoted",
			args: args{`YEAR("test field")`},
			want: `YEAR("test field")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuoteIdentifier(tt.args.identifier); got != tt.want {
				t.Errorf("QuoteIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
