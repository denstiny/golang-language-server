package file

import "testing"

func TestParse(t *testing.T) {
	f, err := Open("/Users/bytedance/denstiny/golang-language-server/pkg/file/testgofile.gox")
	if err != nil {
		t.Error(err)
		return
	}
	gf, err := ParseGoFile(f)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	t.Log(gf)
}
