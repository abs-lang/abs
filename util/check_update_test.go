package util

import (
	"testing"
)

func TestUpdateAvailable(t *testing.T) {
	_, outdated := UpdateAvailable("1.0")
	if !outdated {
		t.Fatalf("expected 1.0 to be outdated")
	}
}
