package completion

import (
	"context"
	"encoding/json"
	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func Completion(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	var params lsp.CompletionParams
	if req.Params != nil {
		err := json.Unmarshal(*req.Params, &params)
		if err != nil {
			return nil, err
		}
	}

	return lsp.CompletionList{}, nil
}
