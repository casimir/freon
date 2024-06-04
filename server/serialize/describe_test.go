package serialize

import (
	"reflect"
	"testing"
)

func TestDescribe(t *testing.T) {
	type args struct {
		obj any
	}
	type S struct {
		ID string `desc:"readonly"`
	}
	tests := []struct {
		name    string
		args    args
		want    []Field
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				obj: struct {
					Name string `desc:""`
					Age  int    `desc:"readonly"`
					Pet  string
				}{
					Name: "Alice",
					Age:  30,
					Pet:  "Bibao",
				},
			},
			want: []Field{
				{Name: "Name", Value: "Alice", Readonly: false, Obfuscate: false},
				{Name: "Age", Value: 30, Readonly: true, Obfuscate: false},
				{Name: "Pet", Value: "Bibao"},
			},
			wantErr: false,
		},
		{
			name: "obfuscate",
			args: args{
				obj: struct {
					Password string `desc:"obfuscate"`
				}{
					Password: "secret",
				},
			},
			want: []Field{
				{Name: "Password", Value: "secret", Readonly: false, Obfuscate: true},
			},
			wantErr: false,
		},
		{
			name: "hidden",
			args: args{
				obj: struct {
					Password string `desc:"hidden"`
				}{
					Password: "secret",
				},
			},
			want:    []Field{},
			wantErr: false,
		},
		{
			name: "unknown tag",
			args: args{
				obj: struct {
					Name string `desc:"unknown"`
				}{
					Name: "Alice",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "embedded",
			args: args{
				obj: struct {
					S
					Name string
				}{
					S:    S{ID: "123"},
					Name: "Alice",
				},
			},
			want: []Field{
				{Name: "ID", Value: "123", Readonly: true, Obfuscate: false},
				{Name: "Name", Value: "Alice", Readonly: false, Obfuscate: false},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Describe(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Describe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Describe() = %v, want %v", got, tt.want)
			}
		})
	}
}
