package file

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := Open("/Users/bytedance/denstiny/golang-language-server/pkg/file/testgofile.txt")
	if err != nil {
		t.Error(err)
		return
	}
	gf, err := ParseGoFile(f)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	t.Log(gf)
}

// findVariableDefinitions 函数用于在 AST 中查找变量定义
func findVariableDefinitions(file *ast.File) []string {
	var variables []string
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok == token.VAR {
				for _, spec := range x.Specs {
					if vspec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range vspec.Names {
							variables = append(variables, name.Name)
						}
					}
				}
			}
		}
		return true
	})
	return variables
}

func TestParseGoFile(t *testing.T) {
	fset := token.NewFileSet()
	// 包含错误的代码示例
	code := `
package main

func main() {
    var a int
    var b, c string
    d := 10
    // 此处存在语法错误，缺少右括号
    fmt.log.Info().Msg("Hello, World!
}
`
	file, err := parser.ParseFile(fset, "", code, parser.AllErrors)
	if err != nil {
		fmt.Println("解析代码时出错:", err)
	}
	variables := findVariableDefinitions(file)
	fmt.Println("定义的变量有:", variables)
}

func TestParseGoFile2(t *testing.T) {
}
