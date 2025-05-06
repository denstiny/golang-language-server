package file

import "testing"

func TestFidnPackage(t *testing.T) {
	s, err := FindPackageName("/Users/bytedance/denstiny/golang-language-server/pkg/engine", "engine")
	if err != nil {
		t.Error(err)
	}
	if s != "github.com/denstiny/golang-language-server/pkg/engine" {
		t.Errorf("FindPackageName() returned %s", s)
	}
}

func TestFidnPackage2(t *testing.T) {
	s, err := FindPackageName("/Users/bytedance/denstiny/golang-language-server/biz/handle/initialize", "initialize")
	if err != nil {
		t.Error(err)
	}
	if s != "github.com/denstiny/golang-language-server/biz/handle/initialize" {
		t.Errorf("FindPackageName() returned %s", s)
	}
}
