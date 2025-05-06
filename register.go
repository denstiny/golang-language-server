package main

import (
	"context"
	"encoding/json"
	"github.com/denstiny/golang-language-server/biz/handle/completion"
	"github.com/denstiny/golang-language-server/biz/handle/didChange"
	"github.com/denstiny/golang-language-server/biz/handle/didOpen"
	"github.com/denstiny/golang-language-server/biz/handle/initialize"
	"github.com/denstiny/golang-language-server/biz/handle/initialized"
	"github.com/denstiny/golang-language-server/pkg/engine"
	"github.com/sourcegraph/jsonrpc2"
	"pkg.nimblebun.works/go-lsp"
)

func RpcHandles() map[string]engine.RouteFunc {
	return map[string]engine.RouteFunc{
		"initialized": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, initialized.Handle(ctx)
		},
		"initialize": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			var param lsp.InitializeParams
			err := json.Unmarshal(*req.Params, &param)
			if err != nil {
				return nil, err
			}
			return initialize.Handle(ctx, c, &param)
		},
		"shutdown": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		},
		"textDocument/didSave": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		},
		"textDocument/completion": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			var param lsp.CompletionParams
			err := json.Unmarshal(*req.Params, &param)
			if err != nil {
				return nil, err
			}
			return completion.Handle(ctx, &param)
		},
		"completionItem/resolve": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		},
		"workspace/executeCommand": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		},
		"$/cancelRequest": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		},
		"textDocument/didChange": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			var param lsp.DidChangeTextDocumentParams
			err := json.Unmarshal(*req.Params, &param)
			if err != nil {
				return nil, err
			}
			return nil, didChange.Handle(ctx, &param)
		},
		"textDocument/didOpen": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			var param lsp.DidOpenTextDocumentParams
			err := json.Unmarshal(*req.Params, &param)
			if err != nil {
				return nil, err
			}
			return nil, didOpen.Handle(ctx, &param)
		},
		"textDocument/didClose": func(ctx context.Context, c *engine.LspService, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
			return nil, nil
		},
	}
}
