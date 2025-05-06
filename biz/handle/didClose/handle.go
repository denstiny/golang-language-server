package didClose

import (
	"context"
	"pkg.nimblebun.works/go-lsp"
)

func Handle(ctx context.Context, param *lsp.DidCloseTextDocumentParams) error {
	return nil
}
