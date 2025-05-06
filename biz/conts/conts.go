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
	Workspace: &struct {
		WorkspaceFolders lsp.WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
		FileOperations   *struct {
			DidCreate  *lsp.FileOperationRegistrationOptions `json:"didCreate,omitempty"`
			WillCreate *lsp.FileOperationRegistrationOptions `json:"willCreate,omitempty"`
			DidRename  *lsp.FileOperationRegistrationOptions `json:"didRename,omitempty"`
			WillRename *lsp.FileOperationRegistrationOptions `json:"willRename,omitempty"`
			DidDelete  *lsp.FileOperationRegistrationOptions `json:"didDelete,omitempty"`
			WillDelete *lsp.FileOperationRegistrationOptions `json:"willDelete,omitempty"`
		} `json:"fileOperations,omitempty"`
	}{
		WorkspaceFolders: lsp.WorkspaceFoldersServerCapabilities{
			Supported: true,
		},
	},
	TextDocumentSync: &lsp.TextDocumentSyncOptions{
		OpenClose: true,
		Change:    lsp.TDSyncKindFull,
	},
	// 文件跟踪
	DefinitionProvider: &lsp.DefinitionRegistrationOptions{
		TextDocumentRegistrationOptions: lsp.TextDocumentRegistrationOptions{
			DocumentSelector: []lsp.DocumentFilter{
				{
					Language: "go",
					Scheme:   "file",
					Pattern:  "*.{go,mod}",
				},
			},
		},
	},
	// 光标跟踪
	DocumentSymbolProvider: &lsp.DocumentSymbolRegistrationOptions{
		DocumentHighlightOptions: lsp.DocumentHighlightOptions{
			WorkDoneProgressOptions: lsp.WorkDoneProgressOptions{
				WorkDoneProgress: true,
			},
		},
	},
}

const CacheFileName = "go_lsp_cahce.db"
