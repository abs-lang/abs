package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestGenerateEqualityString(t *testing.T) {
	tests := []struct {
		input    Object
		expected string
	}{
		{&String{Value: "a"}, "STRING:a"},
		{&Array{Elements: []Object{FALSE}}, "ARRAY:[false]"},
		{&Array{Elements: []Object{&Number{Value: 1}}}, "ARRAY:[1]"},
		{&Array{Elements: []Object{&String{Value: "1"}}}, "ARRAY:[\"1\"]"},
		{&Hash{}, "HASH:{}"},
		{&Hash{Pairs: map[HashKey]HashPair{}}, "HASH:{}"},
		{&Hash{Pairs: map[HashKey]HashPair{HashKey{Value: "x", Type: STRING_OBJ}: HashPair{&String{Value: "x"}, &String{Value: "x"}}}}, "HASH:{\"x\": \"x\"}"},
	}

	for _, tt := range tests {
		output := GenerateEqualityString(tt.input)
		if tt.expected != output {
			t.Fatalf("expected '%v', got '%v' (original: %s)", tt.expected, output, tt.input.Inspect())
		}
	}
}
