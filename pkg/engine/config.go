package engine

import "pkg.nimblebun.works/go-lsp"

type Config struct {
	ServerPort      int
	WorkFolds       []string
	ServerConfigDir string
	Trace           bool
	ClientInfo      lsp.ClientInfo
}
