package loggers

import "testing"

func TestNewReturnsLogger(t *testing.T) {
	if logger := New(); logger == nil {
		t.Fatal("expected logger")
	}
}
