package file

import (
	"context"
	"fmt"
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
	Variables map[string]map[string]TypeInfoSpec
	Functions map[string]map[string]FuncSpec
	Types     map[string]map[string]TypeSpec
	Imports   map[string]ImportSpec
}

func ParseGoFile(file *os.File) (*GoFile, error) {
	var gof = GoFile{
		Variables: make(map[string]map[string]TypeInfoSpec),
		Functions: make(map[string]map[string]FuncSpec),
		Types:     make(map[string]map[string]TypeSpec),
		Imports:   make(map[string]ImportSpec),
	}

	fest := token.NewFileSet()
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
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
	return &gof, nil
}

func (g *GoFile) parse(ctx context.Context, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		g.parseHandle(ctx, n)
		return false
	})
}

func (g *GoFile) parseHandle(ctx context.Context, n ast.Node) {
	println("parseHandle")
	if n == nil {
		return
	}
	ast.Inspect(n, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			g.funcDeclHandle(ctx, x)
		case *ast.GenDecl:
			g.genDeclHandle(ctx, x)
		case *ast.AssignStmt:
			g.assignDeclStmtHandle(ctx, x)
		case *ast.TypeSpec:
			g.typeDeclHandle(ctx, x)
		}
		return true
	})
}

func (g *GoFile) typeDeclHandle(ctx context.Context, x *ast.TypeSpec) {
	println("typeDeclHandle")
	fields := []TypeInfoSpec{}
	if structType, ok := x.Type.(*ast.StructType); ok {
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

	g.withTypeInfo(ctx, TypeSpec{
		Name:    x.Name.Name,
		Fields:  fields,
		Scope:   Scope{x.Pos(), x.End()},
		Comment: x.Doc.Text(),
		Type:    getTypeString(x.Type),
	})
}

func (g *GoFile) funcDeclHandle(ctx context.Context, n *ast.FuncDecl) {
	println("funcDeclHandle", n.Name.Name)
	ctx = context.WithValue(ctx, BlockName, n.Name.Name)
	params := []TypeInfoSpec{}
	returns := []TypeInfoSpec{}

	if n.Type.Params != nil {
		for _, param := range n.Type.Params.List {
			for _, ident := range param.Names {
				params = append(params, TypeInfoSpec{
					Name:    &ident.Name,
					Scope:   Scope{ident.Pos(), ident.End()},
					Type:    ident.Obj.Kind.String(),
					Comment: "",
				})
			}
		}
	}

	if n.Type.Results != nil {
		for _, parm := range n.Type.Results.List {
			for _, ident := range parm.Names {
				returns = append(returns, TypeInfoSpec{
					Name:    &ident.Name,
					Scope:   Scope{ident.Pos(), ident.End()},
					Type:    getTypeString(ident),
					Comment: "",
				})
			}
		}
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
	println("assignDeclStmtHandle")
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
	println("genDeclHandle")
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

// from: fileName dest: global|funcName name: keyword
func (g *GoFile) withVariable(ctx context.Context, n TypeInfoSpec) {
	blackName := getBlockName(ctx)
	if _, ok := g.Variables[blackName]; !ok {
		g.Variables[blackName] = make(map[string]TypeInfoSpec)
	}

	g.Variables[blackName][*n.Name] = n
}

const BlockName = "block-name"

func getBlockName(ctx context.Context) string {
	if v, ok := ctx.Value(BlockName).(string); ok {
		return v
	}
	return "global"
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
