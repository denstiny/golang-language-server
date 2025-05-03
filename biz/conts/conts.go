package conts

import (
	"pkg.nimblebun.works/go-lsp"
)

const (
	VERSION      = "0.0.1"
	SERVICE_NAME = "golang-language-server"
)

// 设置lsp.Server默认功能全部关闭
var ServerCapabilities = lsp.ServerCapabilities{
	HoverProvider: &lsp.HoverOptions{
		WorkDoneProgressOptions: lsp.WorkDoneProgressOptions{
			WorkDoneProgress: true,
		},
	},
	CompletionProvider: &lsp.CompletionOptions{
		ResolveProvider: true,
		TriggerCharacters: []string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", ".",
		},
	},
}

const CacheFileName = "go_lsp_cahce"
