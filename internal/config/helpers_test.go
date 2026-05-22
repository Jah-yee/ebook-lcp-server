package config

import "testing"

func TestHelpers(t *testing.T) {
	t.Setenv("SOME_INT", "12")
	if got := envInt("SOME_INT", 1); got != 12 {
		t.Fatalf("unexpected envInt value %d", got)
	}
	t.Setenv("SOME_INT", "0")
	if got := envInt("SOME_INT", 7); got != 7 {
		t.Fatalf("expected fallback envInt value, got %d", got)
	}
	items := splitCSV(" one, two ,, three ")
	if len(items) != 3 || items[0] != "one" || items[2] != "three" {
		t.Fatalf("unexpected splitCSV result %+v", items)
	}
	if got := defaultString("  value ", "fallback"); got != "value" {
		t.Fatalf("unexpected defaultString result %q", got)
	}
	if got := defaultString(" ", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback defaultString result, got %q", got)
	}
}
