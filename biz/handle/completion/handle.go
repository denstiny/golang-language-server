package completion

import (
	"context"
	"fmt"
	"github.com/denstiny/golang-language-server/biz/dal/cache"
	"github.com/denstiny/golang-language-server/biz/dal/sqlite"
	"github.com/denstiny/golang-language-server/pkg/file"
	"github.com/rs/zerolog/log"
	"go/ast"
	"pkg.nimblebun.works/go-lsp"
	"regexp"
	"strings"
)

func Handle(ctx context.Context, params *lsp.CompletionParams) (interface{}, error) {
	fileUrl := string(params.TextDocumentPositionParams.TextDocument.URI)

	filepath := strings.Trim(fileUrl, "file://")
	var completionItem []lsp.CompletionItem
	v, ok := cache.OpenedFile.Load(filepath)
	if ok {
		gf := v.(*file.GoFile)
		word := gf.GetCursorWord(file.Position{Line: params.Position.Line, Column: params.Position.Character - 1})
		keys := strings.Split(word, ".")
		log.Info().Str("word", word).Msg("Found word")
		if len(keys) > 1 {
			// 查找正在使用的包
			rootPacakgeSpec := keys[0]
			pkg, err := gf.FindPackage(rootPacakgeSpec)
			if err != nil {
				return nil, err
			}
			sqlite.GetPackage(sqlite.PackageFindParams{PackagePath: &pkg.Path, Name: &pkg.Name})

		} else {
			log.Info().Msg(fmt.Sprintf("word(%v:%v): %v", params.Position.Line, params.Position.Character-1, word))
			completionItem = append(completionItem, func(gf *file.GoFile) []lsp.CompletionItem {
				var resp []lsp.CompletionItem
				for _, importItem := range gf.Imports[file.Global] {
					resp = append(resp, buildCompletionItem(importItem.Name, lsp.CIKModule, bigScore))
				}

				for _, values := range gf.Variables[file.Global] {
					resp = append(resp, buildCompletionItem(*values.Name, lsp.CIKVariable, bigScore))
				}

				for _, fun := range gf.Functions[file.Global] {
					resp = append(resp, buildCompletionItem(fun.Name, lsp.CIKFunction, bigScore))
				}

				return resp
			}(gf)...)
		}
	}

	return lsp.CompletionList{
		IsIncomplete: true,
		Items: append(completionItem,
			buildCompletionItem(IMPORT, lsp.CIKKeyword, bigScore),
			buildCompletionItem(IF, lsp.CIKKeyword, bigScore),
			buildCompletionItem(CASE, lsp.CIKKeyword, bigScore),
			buildCompletionItem(DEFAULT, lsp.CIKKeyword, bigScore),
			buildCompletionItem(FUNC, lsp.CIKKeyword, bigScore),
			buildCompletionItem(SWITCH, lsp.CIKKeyword, bigScore),
		),
	}, nil
}

func buildCompletionItem(label string, kind lsp.CompletionItemKind, sort float64) lsp.CompletionItem {
	return lsp.CompletionItem{
		Label:      label,
		Kind:       kind,
		Detail:     label,
		InsertText: label,
		SortText:   fmt.Sprintf("%.2f", sort),
		Tags:       []lsp.CompletionItemTag{},
	}
}

const (
	lowScore  float64 = 0.01
	stdScore  float64 = 1.0
	bigScore  float64 = 1.5
	highScore float64 = 100.0
)

func fuzzyMatchRegex(text, pattern string) bool {
	// 构建正则表达式模式，允许匹配部分单词
	regexPattern := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(pattern))
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Printf("编译正则表达式出错: %v\n", err)
		return false
	}
	return re.MatchString(text)
}

func getKeywordsNode(gf *file.GoFile, workds []string) (ast.Node, bool) {
	var stat bool
	for _, word := range workds {
		for _, v := range gf.Imports {
			for _, importItem := range v {
				if importItem.Name == word {
					stat = true
					// 找到之后需要获取这个包的索引
					//v, ok := cache.OpenedFile.Load(importItem.Path)
				}
			}
		}
	}
	return gf, stat
}
