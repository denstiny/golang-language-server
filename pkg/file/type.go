package file

import "go/token"

type Scope struct {
	Start token.Pos
	End   token.Pos
}

type TypeInfoSpec struct {
	Name    *string // 名称可以为空比如函数的返回参数可以定义变量类型的名称
	Scope   Scope
	Type    string
	Comment string
}

type BlockSpec struct {
	Type  BlockSpecType // 函数块，匿名函数块，普通块
	Scope Scope
}

type BlockSpecType int

const (
	BlockSpecTypeFunc   BlockSpecType = 0
	BlockSpecTypeLambda BlockSpecType = 1
	BlockSpecTypeBlock  BlockSpecType = 2
)

func (t BlockSpecType) String() string {
	switch t {
	case BlockSpecTypeFunc:
		return "func"
	case BlockSpecTypeLambda:
		return "lambda"
	case BlockSpecTypeBlock:
		return "block"
	}
	return "unknown"
}

type FuncSpec struct {
	Name    string
	Params  []TypeInfoSpec
	Returns []TypeInfoSpec
	Scope   Scope
	Comment string
}

type ImportSpec struct {
	Name  string
	Path  string
	Scope Scope
}

type TypeSpec struct {
	Name    string
	Type    string
	Fields  []TypeInfoSpec
	Scope   Scope
	Comment string
}
