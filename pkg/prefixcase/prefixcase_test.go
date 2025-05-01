package prefixcase

import (
	"encoding/json"
	"testing"
)

func TestPrefixCase(t *testing.T) {
	cache := NewPrefixCase[string]("/")
	cache.WithValue("a/b/c", "hello world")
	b, _ := json.MarshalIndent(cache, "", "  ")
	t.Log(string(b))
	c := cache.Value("a/b/c")
	if c[0] != "hello world" {
		t.Errorf("expect \"hello world\", got %s", c[0])
	}
}
