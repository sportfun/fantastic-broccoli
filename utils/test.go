package utils

import "testing"

func AssertEquals(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}

	t.Fatalf("Expected '%v', but get '%v'", a, b)
}

func AssertNotEquals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		return
	}

	t.Fatalf("Expected something different than '%v'", a)
}
