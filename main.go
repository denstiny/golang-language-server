package main

import (
	"github.com/denstiny/golang-language-server/biz/flags"
	"github.com/denstiny/golang-language-server/pkg/engine"
)

func main() {
	client := engine.NewClient(Handles())
	client.SetConfig(engine.Config{
		ServerPort:      flags.SERVICE_PROT,
		ServerConfigDir: flags.SERVICE_CONFIG_DIR,
	})
	client.Start()
}
