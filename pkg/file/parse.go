package file

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

type GoFile struct {
	FileInfo os.FileInfo
	*ast.File
	*token.FileSet
	Variables    map[string]map[string]TypeInfoSpec
	Functions    map[string]map[string]FuncSpec
	Types        map[string]map[string]TypeSpec
	Imports      map[string]ImportSpec
	buffer       map[Position]byte
	BlockType    []string
	scopeIsParse map[Scope]struct{} //保存已经解析过的范围
}

func (g *GoFile) GetByteByPosition(position Position) (byte, error) {
	if b, ok := g.buffer[position]; ok {
		return b, nil
	}
	return byte(0), fmt.Errorf("no byte found for position %v", position)
}

func (g *GoFile) GetCursorNode(pos Position) ast.Node {
	var respNode ast.Node
	ast.Inspect(g, func(n ast.Node) bool {
		startLine := g.Position(n.Pos())
		endLine := g.Position(n.End())
		// 说明还没找到
		if pos.Line < startLine.Line {
			return true
		}

		// 说明已经越过了最可能的的节点了
		if pos.Line > endLine.Line {
			return false
		}
		respNode = n
		return true
	})
	return respNode
}

func (g *GoFile) SetScopeDest(ctx context.Context, start, end token.Pos) {
	startLine := g.Position(start).Line
	endLine := g.Position(end).Line
	for i := startLine - 1; i < endLine; i++ {
		g.BlockType[i] = getBlockName(ctx)
	}
}

// 获取当前光标前的单词
func (g *GoFile) GetCursorWord(pos Position) string {
	var word []byte
	for ; pos.Column >= 0; pos.Column-- {
		b := g.buffer[pos]
		if b == ' ' {
			return string(word)
		}
		word = append(word, b)
	}
	return string(word)
}

// 返回当前块节点的名字
func (g *GoFile) GetCursorBlock(position Position) string {
	if position.Line < 0 || position.Line >= len(g.BlockType) || g.BlockType[position.Line] == "" {
		return "global"
	}
	return g.BlockType[position.Line-1]
}

type Position struct {
	Filename string
	Line     int
	Column   int
}

func ParseGoFile(file *os.File) (*GoFile, error) {
	var gof = GoFile{
		Variables:    make(map[string]map[string]TypeInfoSpec),
		Functions:    make(map[string]map[string]FuncSpec),
		Types:        make(map[string]map[string]TypeSpec),
		Imports:      make(map[string]ImportSpec),
		buffer:       make(map[Position]byte),
		scopeIsParse: make(map[Scope]struct{}),
		BlockType:    make([]string, 0),
	}

	fest := token.NewFileSet()
	gof.FileSet = fest
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	x := 0
	y := 0
	for _, b := range code {
		gof.buffer[Position{Filename: file.Name(), Line: y, Column: x}] = b
		x++
		if b == '\r' || b == '\n' {
			y++
			x = 0
		}
	}
	gof.BlockType = make([]string, y+1)

	gof.FileInfo, err = file.Stat()
	if err != nil {
		return nil, err
	}
	astFile, err := parser.ParseFile(fest, file.Name(), code, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	gof.File = astFile
	gof.parse(context.Background(), gof.File)
	gof.scopeIsParse = nil
	return &gof, nil
}

func (g *GoFile) parse(ctx context.Context, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		g.parseHandle(ctx, n)
		return false
	})
}

func (g *GoFile) parseHandle(ctx context.Context, node ast.Node) {
	log.Info().Msg("parseHandle")
	if node == nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if g.InParse(n.Pos(), n.End()) {
			return true
		}

		switch x := n.(type) {
		case *ast.FuncDecl:
			g.funcDeclHandle(ctx, x)
		case *ast.GenDecl:
			g.genDeclHandle(ctx, x)
		case *ast.AssignStmt:
			g.assignDeclStmtHandle(ctx, x)
		case *ast.TypeSpec:
			g.typeDeclHandle(ctx, x)
		case *ast.FuncLit:
			g.funcLitHandle(ctx, x)
		case *ast.BlockStmt:
			g.blockDeclHandle(ctx, x)
		default:
			return true
		}
		return true
	})
}

func (g *GoFile) InParse(start, end token.Pos) bool {
	_, ok := g.scopeIsParse[Scope{Start: start, End: end}]
	return ok
}

func (g *GoFile) registerParse(start, end token.Pos) {
	g.scopeIsParse[Scope{Start: start, End: end}] = struct{}{}
}

func (g *GoFile) funcLitHandle(ctx context.Context, n *ast.FuncLit) {
	g.registerParse(n.Pos(), n.End())
	g.SetScopeDest(ctx, n.Pos(), n.End())
	log.Info().Msg("funcLitHandle")
	ctx = withBlockName(ctx, BlockSpecTypeLambda.String())
	g.parseHandle(ctx, n.Body)
}

func (g *GoFile) blockDeclHandle(ctx context.Context, block *ast.BlockStmt) {
	g.registerParse(block.Pos(), block.End())
	g.SetScopeDest(ctx, block.Pos(), block.End())
	log.Info().Msg("blockDeclHandle")
	ctx = withBlockName(ctx, BlockSpecTypeBlock.String())
	for _, stmt := range block.List {
		g.parseHandle(ctx, stmt)
	}
}

func (g *GoFile) typeDeclHandle(ctx context.Context, x *ast.TypeSpec) {
	log.Info().Msg("typeDeclHandle")
	fields := []TypeInfoSpec{}
	if structType, ok := x.Type.(*ast.StructType); ok {
		ctx = withBlockName(ctx, x.Name.Name)
		for _, field := range structType.Fields.List {
			for _, ident := range field.Names {
				fieldInfo := TypeInfoSpec{
					Name:  &ident.Name,
					Type:  getTypeString(field.Type),
					Scope: Scope{ident.Pos(), ident.End()},
				}
				fields = append(fields, fieldInfo)
			}
		}
	}

	g.registerParse(x.Pos(), x.End())
	g.SetScopeDest(ctx, x.Pos(), x.End())
	g.withTypeInfo(ctx, TypeSpec{
		Name:    x.Name.Name,
		Fields:  fields,
		Scope:   Scope{x.Pos(), x.End()},
		Comment: x.Doc.Text(),
		Type:    getTypeString(x.Type),
	})
}

func (g *GoFile) parseFieldList(ctx context.Context, fields *ast.FieldList) []TypeInfoSpec {
	var fieldlist []TypeInfoSpec
	for _, param := range fields.List {
		for _, ident := range param.Names {
			fieldlist = append(fieldlist, TypeInfoSpec{
				Name:    &ident.Name,
				Scope:   Scope{ident.Pos(), ident.End()},
				Type:    ident.Obj.Kind.String(),
				Comment: "",
			})
		}
	}
	return fieldlist
}

func (g *GoFile) funcDeclHandle(ctx context.Context, n *ast.FuncDecl) {
	g.registerParse(n.Pos(), n.End())
	g.SetScopeDest(ctx, n.Pos(), n.End())
	log.Info().Msg(fmt.Sprintf("funcDeclHandle: %v", n.Name.Name))
	ctx = withBlockName(ctx, n.Name.Name)
	params := []TypeInfoSpec{}
	returns := []TypeInfoSpec{}

	if n.Type.Params != nil {
		params = g.parseFieldList(ctx, n.Type.Params)
	}

	if n.Type.Results != nil {
		returns = g.parseFieldList(ctx, n.Type.Results)
	}
	g.withFunction(ctx, FuncSpec{
		Name:    n.Name.Name,
		Scope:   Scope{n.Pos(), n.End()},
		Params:  params,
		Returns: returns,
		Comment: n.Doc.Text(),
	})

	g.parseHandle(ctx, n.Body)
}

func (g *GoFile) assignDeclStmtHandle(ctx context.Context, stmt *ast.AssignStmt) {
	g.registerParse(stmt.Pos(), stmt.End())
	g.SetScopeDest(ctx, stmt.Pos(), stmt.End())
	log.Info().Msg("assignDeclStmtHandle")
	if stmt.Tok != token.DEFINE {
		return
	}

	for _, lhs := range stmt.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok {
			g.withVariable(ctx, TypeInfoSpec{
				Name: &ident.Name,
				Type: getTypeString(ident),
				Scope: Scope{
					ident.Pos(),
					ident.End(),
				},
			})
		}
	}
}

func (g *GoFile) genDeclHandle(ctx context.Context, x *ast.GenDecl) {
	g.registerParse(x.Pos(), x.End())
	g.SetScopeDest(ctx, x.Pos(), x.End())
	log.Info().Msg("genDeclHandle")
	if x.Tok == token.VAR || x.Tok == token.CONST {
		for _, spec := range x.Specs {
			if vspec, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range vspec.Names {
					g.withVariable(ctx, TypeInfoSpec{
						Name: &name.Name,
						Scope: Scope{
							vspec.Pos(),
							vspec.End(),
						},
						Type:    getTypeString(vspec.Type),
						Comment: vspec.Comment.Text(),
					})
				}
			}
		}
	}

	if x.Tok == token.IMPORT {
		for _, spec := range x.Specs {
			if vspec, ok := spec.(*ast.ImportSpec); ok {
				path := strings.Trim(vspec.Path.Value, "\"")
				names := strings.Split(path, "/")
				name := names[len(names)-1]
				g.withImports(ctx, ImportSpec{
					Name: name,
					Path: path,
					Scope: Scope{
						Start: vspec.Pos(),
						End:   vspec.End(),
					},
				})
			}
		}
	}
}

// from: fileName dest: global|funcName|block name: keyword
func (g *GoFile) withVariable(ctx context.Context, n TypeInfoSpec) {
	blackName := getBlockName(ctx)
	if _, ok := g.Variables[blackName]; !ok {
		g.Variables[blackName] = make(map[string]TypeInfoSpec)
	}

	g.Variables[blackName][*n.Name] = n
}

const BlockName = "block-name"

func getBlockName(ctx context.Context) string {
	if v := ctx.Value(BlockName); v != nil {
		return v.(string)
	}
	return "global"
}

const blockNameJoinSep = "/"

func withBlockName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, BlockName, getBlockName(ctx)+blockNameJoinSep+name)
}

// dest: global|funcName name: keyword
func (g *GoFile) withFunction(ctx context.Context, n FuncSpec) {
	dest := getBlockName(ctx)
	if _, ok := g.Functions[dest]; !ok {
		g.Functions[dest] = make(map[string]FuncSpec)
	}
	g.Functions[dest][n.Name] = n
}

func (g *GoFile) withImports(ctx context.Context, n ImportSpec) {
	dest := getBlockName(ctx)
	if _, ok := g.Imports[dest]; !ok {
		g.Imports = make(map[string]ImportSpec)
	}
	g.Imports[dest] = n
}

func (g *GoFile) withTypeInfo(ctx context.Context, n TypeSpec) {
	dest := getBlockName(ctx)
	if _, ok := g.Types[dest]; !ok {
		g.Types[dest] = make(map[string]TypeSpec)
	}
	g.Types[dest][n.Name] = n
}

// getTypeString 辅助函数，将 ast 类型节点转换为字符串
func getTypeString(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.ArrayType:
		return "[]" + getTypeString(n.Elt)
	case *ast.StructType:
		return "struct{}"
	case *ast.FuncType:
		return "func()"
	default:
		return fmt.Sprintf("%T", n)
	}
}
