package initialize

import (
	"context"
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/denstiny/golang-language-server/pkg/engine"
	"github.com/sourcegraph/go-lsp"
)

func Handle(ctx context.Context, c *engine.LspService, params *lsp.InitializeParams) (lsp.InitializeResult, error) {
	c.Config.ProjectRoot = params.RootPath
	c.Config.ClientInfo = params.ClientInfo

	return lsp.InitializeResult{
		Capabilities: conts.ServerCapabilities,
	}, nil
}
