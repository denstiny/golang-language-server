package prefixcase

import (
	"strings"
	"testing"
)

func TestStringSplit(t *testing.T) {
	a := strings.Split("a/b/c", "/")
	if len(a) != 3 {
		t.Errorf("a = %v; want 3", len(a))
	}

	b := strings.Split("a", "/")
	if len(b) != 1 {
		t.Errorf("a = %v; want 1", len(a))
	}

	c := []int{1}
	if len(c[1:]) != 0 {
		t.Errorf("a = %v; want 0", len(a))
	}
}
