package entities

import (
	"reflect"
	"testing"
)

func TestNewVS_Snippet(t *testing.T) {
	type args struct {
		name        string
		body        string
		prefix      string
		description string
		scope       string
	}
	tests := []struct {
		name string
		args args
		want *VsSnippet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVS_Snippet(tt.args.name, tt.args.body, tt.args.prefix, tt.args.description, tt.args.scope); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVS_Snippet() = %v, want %v", got, tt.want)
			}
			got := NewVS_Snippet(tt.args.name, tt.args.body, tt.args.prefix, tt.args.description, tt.args.scope)
			if got.Output() != tt.want.Output() {
				t.Errorf("NewVS_Snippet() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSnippet_Output(t *testing.T) {
	type fields struct {
		Name        string
		Body        string
		Prefix      string
		Description string
		Scope       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snippet{
				Name:        tt.fields.Name,
				Body:        tt.fields.Body,
				Prefix:      tt.fields.Prefix,
				Description: tt.fields.Description,
				Scope:       tt.fields.Scope,
			}
			if got := s.Output(); got != tt.want {
				t.Errorf("Snippet.Output() = %v, want %v", got, tt.want)
			}
		})
	}
}
