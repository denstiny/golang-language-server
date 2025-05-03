package completion

import (
	"context"
	"pkg.nimblebun.works/go-lsp"
)

func Handle(ctx context.Context, params *lsp.CompletionParams) (interface{}, error) {
	return lsp.CompletionList{
		IsIncomplete: true,
		Items: []lsp.CompletionItem{
			buildCompletionItem(IMPORT, lsp.CIKKeyword),
			buildCompletionItem(IF, lsp.CIKKeyword),
			buildCompletionItem(CASE, lsp.CIKKeyword),
			buildCompletionItem(DEFAULT, lsp.CIKKeyword),
			buildCompletionItem(FUNC, lsp.CIKKeyword),
			buildCompletionItem(SWITCH, lsp.CIKKeyword),
		},
	}, nil
}

func buildCompletionItem(label string, kind lsp.CompletionItemKind) lsp.CompletionItem {
	return lsp.CompletionItem{
		Label:      label,
		Kind:       kind,
		Detail:     label,
		InsertText: label,
		Tags:       []lsp.CompletionItemTag{},
	}
}

const (
	lowScore  float64 = 0.01
	stdScore  float64 = 1.0
	highScore float64 = 100.0
)
