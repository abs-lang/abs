package util

import (
	"os"
	"testing"
)

func TestUnaliasPath(t *testing.T) {
	tests := []struct {
		path     string
		aliases  map[string]string
		expected string
	}{
		{"test", map[string]string{}, "test"},
		{"test" + string(os.PathSeparator) + "sample.abs", map[string]string{}, "test" + string(os.PathSeparator) + "sample.abs"},
		{"test" + string(os.PathSeparator) + "sample.abs", map[string]string{"test": "path"}, "path" + string(os.PathSeparator) + "sample.abs"},
		{"test", map[string]string{"test": "path"}, "path" + string(os.PathSeparator) + "index.abs"},
		{"." + string(os.PathSeparator) + "test", map[string]string{"test": "path"}, "." + string(os.PathSeparator) + "test"},
	}

	for _, tt := range tests {
		res := UnaliasPath(tt.path, tt.aliases)

		if res != tt.expected {
			t.Fatalf("error unaliasing path, expected %s, got %s", tt.expected, res)
		}
	}
}
