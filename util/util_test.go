package util

import (
	"os"
	"testing"

	"github.com/abs-lang/abs/object"
)

func TestUnaliasPath(t *testing.T) {
	tests := []struct {
		path     string
		aliases  map[string]string
		expected string
	}{
		{"test", map[string]string{}, "test" + string(os.PathSeparator) + "index.abs"},
		{"test" + string(os.PathSeparator) + "sample.abs", map[string]string{}, "test" + string(os.PathSeparator) + "sample.abs"},
		{"test" + string(os.PathSeparator) + "sample.abs", map[string]string{"test": "path"}, "path" + string(os.PathSeparator) + "sample.abs"},
		{"test", map[string]string{"test": "path"}, "path" + string(os.PathSeparator) + "index.abs"},
		{"." + string(os.PathSeparator) + "test", map[string]string{"test": "path"}, "test" + string(os.PathSeparator) + "index.abs"},
	}

	for _, tt := range tests {
		res := UnaliasPath(tt.path, tt.aliases)

		if res != tt.expected {
			t.Fatalf("error unaliasing path, expected %s, got %s", tt.expected, res)
		}
	}
}

func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		strings []string
		len     int
	}{
		{[]string{"a", "b", "c"}, 3},
		{[]string{"a", "a", "a"}, 1},
	}

	for _, tt := range tests {
		if len(UniqueStrings(tt.strings)) != tt.len {
			t.Fatalf("expected %d, got %d", tt.len, len(UniqueStrings(tt.strings)))
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		strings  []string
		match    string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "a", "a"}, "d", false},
	}

	for _, tt := range tests {
		if tt.expected != Contains(tt.strings, tt.match) {
			t.Fatalf("expected %v", tt.expected)
		}
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		number   string
		expected bool
	}{
		{"12", true},
		{"12a", false},
		{"12.2", true},
	}

	for _, tt := range tests {
		if tt.expected != IsNumber(tt.number) {
			t.Fatalf("expected %v (%s)", tt.expected, tt.number)
		}
	}
}

func TestInterpolateStringVars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string", "string"},
		{"string $string string", "string test string"},
		{"string $string", "string test"},
		{"$string", "test"},
		{"${string}", "test"},
		{"\\$string", "$string"},
		{"\\${string}", "${string}"},
		{"_$string", "_test"},
		{"string$string\\string", "stringtest\\string"},
		{"$string_", ""},
		{"xy\\z", "xy\\z"},
		{"${string}_", "test_"},
		{"${string x", "${string x"},
	}

	env := object.NewEnvironment(os.Stdout, os.Stderr, "", "dev", false)
	env.Set("string", &object.String{Value: "test"})

	for _, tt := range tests {
		output := InterpolateStringVars(tt.input, env)
		if tt.expected != output {
			t.Fatalf("expected '%v', got '%v' (original: %s)", tt.expected, output, tt.input)
		}
	}
}

func TestMapify(t *testing.T) {
	elements := []object.Object{}
	first := &object.String{Value: "x"}
	second := &object.Number{Value: 10}
	elements = append(elements, first, second)

	m := Mapify(elements)

	if len(m) != 2 {
		t.Fatalf("expected len '%d', got '%d'", 2, len(m))
	}

	if m["STRING:x"] != first {
		t.Fatalf("string element not found")
	}

	if m["NUMBER:10"] != second {
		t.Fatalf("number element not found")
	}
}
