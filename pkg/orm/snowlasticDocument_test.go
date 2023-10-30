package types

import (
	"database/sql"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func Test_keysToLower(t *testing.T) {
	type args struct {
		m map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "string",
			args: args{map[string]any{
				"TEST": "TEST",
			}},
			want: map[string]any{
				"test": "TEST",
			},
		},
		{
			name: "int",
			args: args{map[string]any{
				"TEST": 1,
			}},
			want: map[string]any{
				"test": 1,
			},
		},
		{
			name: "bool",
			args: args{map[string]any{
				"TEST": false,
			}},
			want: map[string]any{
				"test": false,
			},
		},
		{
			name: "nested",
			args: args{map[string]any{
				"TEST": map[string]any{
					"NESTED": "nested",
				},
			}},
			want: map[string]any{
				"test": map[string]any{
					"nested": "nested",
				},
			},
		},
		{
			name: "deeply nested",
			args: args{map[string]any{
				"MAP_OF_MAPS": map[string]any{
					"LEVEL1": map[string]any{
						"LEVEL-2": map[string]any{
							"FLOOR": nil,
						},
					},
				},
			}},
			want: map[string]any{
				"map_of_maps": map[string]any{
					"level1": map[string]any{
						"level-2": map[string]any{
							"floor": nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keysToLower(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keysToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_Get(t *testing.T) {
	var baseMap = map[string]any{
		"testString": "test",
		"testInt":    1,
		"testBool":   true,
	}
	type fields struct {
		m m
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
	}{
		{
			name:   "string",
			fields: fields{baseMap},
			args:   args{"testString"},
			want:   "test",
		},
		{
			name:   "int",
			fields: fields{baseMap},
			args:   args{"testInt"},
			want:   1,
		},
		{
			name:   "bool",
			fields: fields{baseMap},
			args:   args{"testBool"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				m: tt.fields.m,
			}
			if got := d.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_GetID(t *testing.T) {
	viper.Set("identifier", "id")
	var stringID = map[string]any{
		"id":         "1",
		"testString": "test",
		"testInt":    1,
		"testBool":   true,
	}
	var intID = map[string]any{
		"id":         1,
		"testString": "test",
		"testInt":    1,
		"testBool":   true,
	}
	type fields struct {
		m m
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "string id",
			fields: fields{stringID},
			want:   "1",
		},
		{
			name:   "int id",
			fields: fields{intID},
			want:   "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				m: tt.fields.m,
			}
			if got := d.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_MarshalJSON(t *testing.T) {
	type fields struct {
		m m
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "simple",
			fields: fields{map[string]any{
				"id":   "1",
				"test": "TEST",
			}},
			want:    []byte(`{"id":"1","test":"TEST"}`),
			wantErr: false,
		},
		{
			name: "exclude nil",
			fields: fields{map[string]any{
				"id":   "1",
				"test": nil,
			}},
			want:    []byte(`{"id":"1"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				m: tt.fields.m,
			}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestDocument_New(t *testing.T) {
	type fields struct {
		m m
	}
	tests := []struct {
		name   string
		fields fields
		want   SnowlasticDocument
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				m: tt.fields.m,
			}
			if got := d.New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocument_ScanFrom(t *testing.T) {
	type fields struct {
		m m
	}
	type args struct {
		rows *sql.Rows
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				m: tt.fields.m,
			}
			if err := d.ScanFrom(tt.args.rows); (err != nil) != tt.wantErr {
				t.Errorf("ScanFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDocument(t *testing.T) {
	tests := []struct {
		name string
		want SnowlasticDocument
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDocument(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDocumentFromMap(t *testing.T) {
	type args struct {
		m map[string]any
	}
	tests := []struct {
		name string
		args args
		want SnowlasticDocument
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDocumentFromMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocumentFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
