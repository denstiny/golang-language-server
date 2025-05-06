package didOpen

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/dal/cache"
	"github.com/denstiny/golang-language-server/pkg/file"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"pkg.nimblebun.works/go-lsp"
	"strings"
)

func Handle(ctx context.Context, param *lsp.DidOpenTextDocumentParams) error {
	filePath := strings.Trim(string(param.TextDocument.URI), "file://")
	gofile, err := file.ParseGoBuffer(filePath, []byte(param.TextDocument.Text))
	if err == nil {
		cache.OpenedFile.Store(filePath, gofile)
	}
	return nil
}

func ParseMod(ctx context.Context, filePath string) (*modfile.File, error) {
	if file.Exists(filePath) {
		return nil, fmt.Errorf("file %s exists", filePath)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fest, err := modfile.Parse(filePath, b, nil)
	if err != nil {
		return nil, err
	}
	return fest, nil
}
