package file

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

type GoFile struct {
	FileInfo os.FileInfo
	*ast.File
}

func ParseGoFile(file *os.File) (*GoFile, error) {
	var gof GoFile
	fest := token.NewFileSet()
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	gof.FileInfo, err = file.Stat()
	if err != nil {
		return nil, err
	}
	ast, err := parser.ParseFile(fest, file.Name(), code, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	gof.File = ast
	return &gof, nil
}
