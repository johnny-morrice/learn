package parser

import (
	"errors"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
)

func ParseFile(fileName string) (*ast.AST, error) {
	panic("not implemented")
}

type ProgressFunc func(bldr *ast.Builder) (*ast.Builder, error)

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

func (pc ParseContext) CompleteStmt() error {

	bldr := pc.Bldr

	for _, progress := range pc.Progress {
		var err error
		bldr, err = progress(bldr)
		if err != nil {
			return err
		}
	}
	pc.Progress = []ProgressFunc{}

	return nil
}

func (pc ParseContext) AddProgress(f ProgressFunc) {
	pc.Progress = append(pc.Progress, f)
}

func (pc ParseContext) Copy() ParseContext {
	copy := pc

	copy.Progress = []ProgressFunc{}
	copy.Progress = append(copy.Progress, pc.Progress...)

	return copy
}

func Parse(pc ParseContext) (*ast.AST, error) {
	pc = AST()(pc)
	if pc.Failed {
		return nil, errors.New("parse error")
	}
	return pc.Bldr.Build(), nil
}
