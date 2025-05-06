package didChange

import (
	"context"
	"github.com/denstiny/golang-language-server/biz/dal/cache"
	"github.com/denstiny/golang-language-server/pkg/file"
	"github.com/rs/zerolog/log"
	"pkg.nimblebun.works/go-lsp"
	"strings"
)

func Handle(ctx context.Context, param *lsp.DidChangeTextDocumentParams) error {
	log.Info().Str("text: ", param.ContentChanges[0].Text).Msg("handling didChange textDocument")
	filePath := strings.Trim(string(param.TextDocument.URI), "file://")
	v, ok := cache.OpenedFile.Load(filePath)
	gofile, err := file.ParseGoBuffer(filePath, []byte(param.ContentChanges[0].Text))
	if err != nil {
		if ok {
			gofile = v.(*file.GoFile)
			gofile.UpdateBuffer([]byte(param.ContentChanges[0].Text))
		} else {
			return nil
		}
	}
	cache.OpenedFile.Store(filePath, gofile)
	return nil
}
