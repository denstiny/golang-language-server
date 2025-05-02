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
