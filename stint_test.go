package stint

import (
	"strings"
	"testing"
)

func TestTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		env      map[string]string
		expected string
		wantErr  bool
	}{
		{
			name:     "int and len",
			input:    "val: %bar:int:len(2)%",
			env:      map[string]string{"bar": "52"},
			expected: "val: 52",
			wantErr:  false,
		},
		{
			name:     "len range",
			input:    "%code:len(3, 5)%",
			env:      map[string]string{"code": "ABCD"},
			expected: "ABCD",
			wantErr:  false,
		},
		{
			name:     "not an integer",
			input:    "%bar:int%",
			env:      map[string]string{"bar": "abc"},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "too short",
			input:    "%bar:len(5)%",
			env:      map[string]string{"bar": "123"},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "variable missing",
			input:    "hello %name%",
			env:      map[string]string{"foo": "bar"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := NewTemplate(tt.input)
			if err != nil {
				t.Fatalf("NewTemplate() error = %v", err)
			}

			var out strings.Builder
			err = tmpl.Execute(&out, tt.env)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("Execute() got = %q, want %q", out.String(), tt.expected)
			}
		})
	}
}
