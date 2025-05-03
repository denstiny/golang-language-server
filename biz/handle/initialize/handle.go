package initialize

import (
	"context"
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/denstiny/golang-language-server/biz/handle/progress"
	"github.com/denstiny/golang-language-server/pkg/engine"
	"log"
	"pkg.nimblebun.works/go-lsp"
)

func Handle(ctx context.Context, c *engine.LspService, params *lsp.InitializeParams) (lsp.InitializeResult, error) {
	for _, fold := range params.WorkspaceFolders {
		c.Config.WorkFolds = append(c.Config.WorkFolds, fold.Name)
	}
	c.Config.ClientInfo = params.ClientInfo

	// 创建进度条打印初始化进度
	progres := progress.NewProgress("initialize:golang-language-server", "golang-language-server")
	err := progres.Begin(ctx, "init workspace index", false)
	if err != nil {
		log.Fatalf("begin progres: init workspace index error: %s", err)
	}
	// TODO: 实现递归遍历project文件，通过 `file` 包解析索引存储到数据库中

	return lsp.InitializeResult{
		Capabilities: conts.ServerCapabilities,
		ServerInfo: lsp.ServerInfo{
			Version: conts.VERSION,
			Name:    conts.SERVICE_NAME,
		},
	}, nil
}
