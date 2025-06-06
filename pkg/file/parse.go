package file

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const Global = "Global"

type FileInfo struct {
	FileName  string
	FilePath  string
	LineCount int
	ModTime   time.Time
}

type GoFile struct {
	FileInfo       FileInfo `json:"file_info"`
	*ast.File      `json:"ast_file"`
	*token.FileSet `json:"file_set"`
	Variables      map[string]map[string]TypeInfoSpec `json:"variables"`
	Functions      map[string]map[string]FuncSpec     `json:"functions"`
	Types          map[string]map[string]TypeSpec     `json:"types"`
	Imports        map[string][]ImportSpec            `json:"imports"`
	Buffer         map[Position]byte                  `json:"buffer"`
	BlockType      []string                           `json:"blockType"`
	ScopeIsParse   map[Scope]struct{}                 //保存已经解析过的范围
}

func (g *GoFile) GetByteByPosition(position Position) (byte, error) {
	if b, ok := g.Buffer[position]; ok {
		return b, nil
	}
	return byte(0), fmt.Errorf("no byte found for position %v", position)
}

func (g *GoFile) GetCursorNode(pos Position) ast.Node {
	var respNode ast.Node
	for _, n := range g.Scope.Objects {
		x := n.Decl.(ast.Node)
		startLine := g.Position(x.Pos())
		endLine := g.Position(x.End())
		// 说明还没找到
		if pos.Line < startLine.Line {
			continue
		}

		// 说明已经越过了最可能的的节点了
		if pos.Line > endLine.Line {
			break
		}
		respNode = x
	}
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
		b := g.Buffer[pos]
		if unicode.IsSpace(rune(b)) {
			return string(word)
		}

		if unicode.IsLetter(rune(b)) || b == '.' {
			word = append([]byte{b}, word...)
		}
	}
	return string(word)
}

// 返回当前块节点的名字
func (g *GoFile) GetCursorBlock(position Position) string {
	if position.Line < 0 || position.Line >= len(g.BlockType) || g.BlockType[position.Line] == "" {
		return Global
	}
	return g.BlockType[position.Line-1]
}

type Position struct {
	Line   int
	Column int
}

func newGoFile() GoFile {
	return GoFile{
		Variables:    make(map[string]map[string]TypeInfoSpec),
		Functions:    make(map[string]map[string]FuncSpec),
		Types:        make(map[string]map[string]TypeSpec),
		Imports:      make(map[string][]ImportSpec),
		Buffer:       make(map[Position]byte),
		ScopeIsParse: make(map[Scope]struct{}),
		BlockType:    make([]string, 0),
		FileInfo:     FileInfo{},
	}
}

func (g *GoFile) UpdateBuffer(code []byte) {
	x := 0
	y := 0
	for _, b := range code {
		g.Buffer[Position{Line: y, Column: x}] = b
		x++
		if b == '\r' || b == '\n' {
			y++
			x = 0
		}
	}
	g.FileInfo.LineCount = y
}

func ParseGoFile(file *os.File) (*GoFile, error) {
	var gof = newGoFile()

	fest := token.NewFileSet()
	gof.FileSet = fest
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	gof.UpdateBuffer(code)
	gof.BlockType = make([]string, gof.FileInfo.LineCount+1)

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	gof.FileInfo = FileInfo{
		FileName: filepath.Ext(fileInfo.Name()),
		FilePath: fileInfo.Name(),
		ModTime:  fileInfo.ModTime(),
	}
	astFile, err := parser.ParseFile(fest, file.Name(), code, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	gof.File = astFile
	gof.parse(context.Background(), gof.File)
	gof.ScopeIsParse = nil
	return &gof, nil
}

func ParseGoBuffer(fileName string, code []byte) (*GoFile, error) {
	gof := newGoFile()
	fest := token.NewFileSet()
	gof.FileSet = fest
	var err error

	gof.UpdateBuffer(code)
	gof.BlockType = make([]string, gof.FileInfo.LineCount+1)

	astFile, err := parser.ParseFile(fest, fileName, code, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	gof.File = astFile
	gof.FileInfo = FileInfo{
		FileName:  filepath.Ext(fileName),
		FilePath:  fileName,
		ModTime:   time.Now(),
		LineCount: gof.FileInfo.LineCount,
	}

	gof.parse(context.Background(), gof.File)
	gof.ScopeIsParse = nil
	return &gof, nil
}

func (g *GoFile) PackageName() string {
	if g.File.Name != nil {
		return g.File.Name.Name
	}
	return ""
}

func (g *GoFile) parse(ctx context.Context, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		g.parseHandle(ctx, n)
		return false
	})
}

func (g *GoFile) parseHandle(ctx context.Context, node ast.Node) {
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
	_, ok := g.ScopeIsParse[Scope{Start: start, End: end}]
	return ok
}

func (g *GoFile) registerParse(start, end token.Pos) {
	g.ScopeIsParse[Scope{Start: start, End: end}] = struct{}{}
}

func (g *GoFile) funcLitHandle(ctx context.Context, n *ast.FuncLit) {
	g.registerParse(n.Pos(), n.End())
	g.SetScopeDest(ctx, n.Pos(), n.End())
	ctx = withBlockName(ctx, BlockSpecTypeLambda.String())
	g.parseHandle(ctx, n.Body)
}

func (g *GoFile) blockDeclHandle(ctx context.Context, block *ast.BlockStmt) {
	g.registerParse(block.Pos(), block.End())
	g.SetScopeDest(ctx, block.Pos(), block.End())
	ctx = withBlockName(ctx, BlockSpecTypeBlock.String())
	for _, stmt := range block.List {
		g.parseHandle(ctx, stmt)
	}
}

func (g *GoFile) typeDeclHandle(ctx context.Context, x *ast.TypeSpec) {
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

	ctx = withBlockName(ctx, n.Name.Name)
	g.parseHandle(ctx, n.Body)
}

func (g *GoFile) assignDeclStmtHandle(ctx context.Context, stmt *ast.AssignStmt) {
	g.registerParse(stmt.Pos(), stmt.End())
	g.SetScopeDest(ctx, stmt.Pos(), stmt.End())
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

// from: fileName dest: Global|funcName|block name: keyword
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
	return Global
}

const blockNameJoinSep = "/"

func withBlockName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, BlockName, getBlockName(ctx)+blockNameJoinSep+name)
}

// dest: Global|funcName name: keyword
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
		g.Imports = make(map[string][]ImportSpec)
	}
	g.Imports[dest] = append(g.Imports[dest], n)
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

func FindPackageName(file string, packageName string) (string, error) {
	println("package:", packageName)
	var err error
	if IsDir(file) {
		if Exists(filepath.Join(file, "go.mod")) {
			f, err := Open(filepath.Join(file, "go.mod"))
			if err != nil {
				return "", err
			}
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				return "", err
			}
			mod, err := modfile.Parse(filepath.Join(file, "go.mod"), b, nil)
			if err != nil {
				return "", err
			}
			packages := strings.Split(packageName, "/")
			return fmt.Sprintf("%v/%v", mod.Module.Mod.String(), strings.Join(packages[1:], "/")), nil
		} else {
			fold := getLastFolder(file)
			return FindPackageName(filepath.Dir(file), fmt.Sprintf("%v/%v", fold, packageName))
		}
	}

	fold := getLastFolder(file)
	if Exists(file) {
		packageName, err = FindPackageName(filepath.Dir(file), filepath.Join(fold, packageName))
		if err != nil {
			return "", err
		}
	}
	return packageName, nil
}

func getLastFolder(path string) string {
	// 获取路径的目录部分
	dir := filepath.Dir(path)
	// 获取目录部分的最后一个元素
	return filepath.Base(dir)
}

func (p *GoFile) FindPackage(packageName string) (ImportSpec, error) {
	for _, pkgs := range p.Imports {
		for _, pkg := range pkgs {
			if pkg.Name == packageName {
				return pkg, nil
			}
		}
	}
	return ImportSpec{}, errors.New("package not found")
}
