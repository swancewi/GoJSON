package jsonParser

import (
	"reflect"
	"testing"
)

// TestParseJSON tests valid JSON inputs and expected Go values.
func TestParseJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "simple string object",
			input:    `{"name":"Alice","age":"30"}`,
			expected: map[string]any{"name": "Alice", "age": "30"},
		},
		{
			name:     "simple string no quotes",
			input:    `name`,
			expected: "name",
		},
		{
			name:     "simple string with quotes",
			input:    `"name"`,
			expected: "name",
		},
		{
			name:     "simple object",
			input:    `{"name":"Alice","age":30}`,
			expected: map[string]any{"name": "Alice", "age": float64(30)},
		},
		{
			name:     "array of numbers",
			input:    `[1, 2, 3]`,
			expected: []any{float64(1), float64(2), float64(3)},
		},
		{
			name:  "nested structure",
			input: `{"user":{"name":"Bob","tags":["dev","go"]}}`,
			expected: map[string]any{
				"user": map[string]any{
					"name": "Bob",
					"tags": []any{"dev", "go"},
				},
			},
		},
		{
			name:     "booleans and null",
			input:    `{"ok":true,"fail":false,"none":null}`,
			expected: map[string]any{"ok": true, "fail": false, "none": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJSON(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseJSON(%q) = %#v, want %#v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestParseJSONErrors ensures invalid JSON fails with an error.
func TestParseJSONErrors(t *testing.T) {
	tests := []string{
		`{`,               // unterminated object
		`[1,2,`,           // unterminated array
		`{"key": "value"`, // missing closing brace
		`{"key": tru`,     // incomplete literal
		`{"x": 09}`,       // invalid number
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseJSON(input)
			if err == nil {
				t.Errorf("expected error for input %q, got nil", input)
			}
		})
	}
}
