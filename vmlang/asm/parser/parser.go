package parser

import (
	"errors"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
)

func ParseFile(fileName string) (*ast.AST, error) {
	panic("not implemented")
}

type ProgressFunc func(bldr ast.Builder) (ast.Builder, error)

type ParseContext struct {
	FileName       string
	Line           string
	Char           string
	RemainingInput string
	Failed         bool
	Error          error
	Bldr           *ast.Builder
	Progress       []ProgressFunc
}

func Parse(pc ParseContext) (*ast.AST, error) {
	pc = AST()(pc)
	if pc.Failed {
		return nil, errors.New("parse error")
	}
	return pc.Bldr.Build(), nil
}
