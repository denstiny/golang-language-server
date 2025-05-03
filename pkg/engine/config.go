package engine

import "github.com/sourcegraph/go-lsp"

type Config struct {
	ServerPort      int
	ProjectRoot     string
	ServerConfigDir string
	Trace           bool
	ClientInfo      lsp.ClientInfo
}
