package initialize

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/conts"
	"github.com/denstiny/golang-language-server/biz/dal/sqlite"
	"github.com/denstiny/golang-language-server/biz/dal/sqlite/model"
	"github.com/denstiny/golang-language-server/biz/handle/progress"
	"github.com/denstiny/golang-language-server/pkg/engine"
	"github.com/denstiny/golang-language-server/pkg/file"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"pkg.nimblebun.works/go-lsp"
	"strings"
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
		ctx, err = LoadGoMod(ctx, folder.Name)
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

func LoadGoMod(ctx context.Context, p string) (context.Context, error) {
	filePath := filepath.Join(p, "go.mod")
	if !file.Exists(filePath) {
		return ctx, fmt.Errorf("load go mod err: %v go.mod not found", filePath)
	}

	f, err := file.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return ctx, err
	}

	fest, err := modfile.Parse(f.Name(), data, nil)
	if err != nil {
		return ctx, err
	}

	packInfo, err := sqlite.GetPackage(fest.Module.Mod.Path, fest.Module.Mod.Version)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			pageNames := strings.Split(strings.ReplaceAll(fest.Module.Mod.Path, "\"", ""), "/")
			name := pageNames[len(pageNames)-1]
			packInfo = &model.Package{
				Name:        name,
				PackagePath: fest.Module.Mod.Path,
				Version:     fest.Module.Mod.Version,
			}
			err = sqlite.CreatePackage(packInfo)
			if err != nil {
				return ctx, err
			}
		} else {
			return ctx, err
		}
	}

	log.Info().Str("module", packInfo.PackagePath).Msg("found pacakges")
	ctx = context.WithValue(ctx, "packageName", packInfo.Name)
	ctx = context.WithValue(ctx, "packagePath", packInfo.PackagePath)

	// 存储依赖的包
	for _, require := range fest.Require {
		packInfo, err := sqlite.GetPackage(require.Mod.Path, require.Mod.Version)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				pageNames := strings.Split(strings.ReplaceAll(require.Mod.Path, "\"", ""), "/")
				name := pageNames[len(pageNames)-1]
				packInfo = &model.Package{
					Name:        name,
					PackagePath: require.Mod.Path,
					Version:     require.Mod.Version,
				}
				err = sqlite.CreatePackage(packInfo)
				if err != nil {
					return ctx, err
				}
			} else {
				return ctx, err
			}
		}
	}

	// 存储当前项目的所有子包
	return ctx, nil
}

func LoadGoCodeFile(ctx context.Context, root string) error {
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if filePath == root {
			return nil
		}
		if info.IsDir() {
			return LoadGoCodeFile(ctx, filePath)
		}
		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		f, err := file.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		gf, err := file.ParseGoFile(f)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}

		packageName := ctx.Value("packageName").(string)
		packagePath := ctx.Value("packagePath").(string)
		if packageName != gf.PackageName() {
			ctx = context.WithValue(ctx, "packageName", packageName)
			ctx = context.WithValue(ctx, "packagePath", strings.Join([]string{packagePath, gf.PackageName()}, "/"))
		}
		sqlite.FindPackage(sqlite.PackageFindParams{})

		return nil
	})
	return err
}
