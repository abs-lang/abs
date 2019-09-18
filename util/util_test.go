package util

import (
	"testing"
)

func TestUnaliasPath(t *testing.T) {
	tests := []struct {
		path     string
		aliases  map[string]string
		expected string
	}{
		{"test", map[string]string{}, "test"},
		{"test/sample.abs", map[string]string{}, "test/sample.abs"},
		{"test/sample.abs", map[string]string{"test": "path"}, "path/sample.abs"},
		{"test", map[string]string{"test": "path"}, "path/index.abs"},
		{"./test", map[string]string{"test": "path"}, "./test"},
	}

	for _, tt := range tests {
		res := UnaliasPath(tt.path, tt.aliases)

		if res != tt.expected {
			t.Fatalf("error unaliasing path, expected %s, got %s", tt.expected, res)
		}
	}
}
