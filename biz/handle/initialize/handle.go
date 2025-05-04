package initialize

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/denstiny/golang-language-server/biz/dal/cache"
	"github.com/denstiny/golang-language-server/biz/handle/progress"
	"github.com/denstiny/golang-language-server/pkg/engine"
	"github.com/denstiny/golang-language-server/pkg/file"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"pkg.nimblebun.works/go-lsp"
	"time"
)

func Handle(ctx context.Context, c *engine.LspService, params *lsp.InitializeParams) (lsp.InitializeResult, error) {
	InitializeService(c, params)
	// 创建进度条打印初始化进度
	progres := progress.NewProgress("initialize:golang-language-server", "golang-language-server")
	err := progres.Begin(ctx, "init workspace index", false)
	if err != nil {
		log.Error().Msg("begin progres: init workspace index error: %s" + err.Error())
	}

	// TODO: 实现递归遍历project文件，通过 `file` 包解析索引存储到数据库中
	for _, folder := range params.WorkspaceFolders {
		err = LoadGoMod(ctx, folder.Name)
		if err != nil {
			log.Error().Msg("load go modules failed: %s" + err.Error())
			return lsp.InitializeResult{}, err
		}

		err = LoadGoCodeFile(ctx, folder.Name)
		if err != nil {
			log.Error().Msg("load go modules failed: %s" + err.Error())
			return lsp.InitializeResult{}, err
		}
	}
	time.Sleep(10000 * time.Millisecond)

	progres.End(ctx, "golsp", "golang-language-server init workspace index")

	return lsp.InitializeResult{
		Capabilities: conts.ServerCapabilities,
		ServerInfo: lsp.ServerInfo{
			Version: conts.VERSION,
			Name:    conts.SERVICE_NAME,
		},
	}, nil
}

func InitializeService(c *engine.LspService, param *lsp.InitializeParams) {
	// 将初始化信息暂存到service中
	for _, fold := range param.WorkspaceFolders {
		c.Config.WorkFolds = append(c.Config.WorkFolds, fold.Name)
	}
	c.Config.ClientInfo = param.ClientInfo
}

func LoadGoMod(ctx context.Context, p string) error {
	filePath := filepath.Join(p, "go.mod")
	if !file.Exists(filePath) {
		return fmt.Errorf("load go mod err: %v go.mod not found", filePath)
	}

	f, err := file.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	fest, err := modfile.Parse(f.Name(), data, nil)
	if err != nil {
		return err
	}

	pacakges, err := cache.GetPackage(fest.Module.Mod.Path, fest.Module.Mod.Version)
	if err != nil {
		return err
	}
	if pacakges != nil {
		log.Info().Str("module", pacakges.PackageName).Msg("found pacakges")
	}

	return nil
}

func LoadGoCodeFile(ctx context.Context, path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return LoadGoCodeFile(ctx, path)
		}
		if filepath.Ext(info.Name()) != ".go" {
			f, err := file.Open(info.Name())
			if err != nil {
				return err
			}
			defer f.Close()

			gf, err := file.ParseGoFile(f)
			if err != nil {
				log.Error().Msg(err.Error())
				return err
			}
			log.Info().Interface("gf", gf).Msg("parse ok")
		}
		return nil
	})
	return err
}
