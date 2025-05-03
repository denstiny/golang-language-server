package completion

import (
	"context"
	"pkg.nimblebun.works/go-lsp"
)

func Handle(ctx context.Context, params *lsp.CompletionParams) (interface{}, error) {
	return lsp.CompletionList{}, nil
}

const (
	lowScore  float64 = 0.01
	stdScore  float64 = 1.0
	highScore float64 = 100.0
)
