package parser

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
)

func ParseFile(fileName string) (ast.AST, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return ast.AST{}, fmt.Errorf("failed to open asm file: %s", err)
	}
	defer file.Close()
	bs, err := io.ReadAll(file)
	if err != nil {
		return ast.AST{}, fmt.Errorf("failed to read asm file: %s", err)
	}
	pc := ParseContext{
		FileName:       fileName,
		RemainingInput: string(bs),
	}
	return Parse(pc)
}

type ParseContext struct {
	FileName       string
	Line           string
	Char           string
	RemainingInput string
	Failed         bool
	IsCapturing    bool
	CapturedText   string
	ErrorMessage   string
	Bldr           ast.Builder
}

func Parse(pc ParseContext) (ast.AST, error) {
	pc = AST()(pc)
	if pc.Failed {
		err := errors.New("parse error")
		if pc.ErrorMessage != "" {
			err = errors.New(pc.ErrorMessage)
		}
		return ast.AST{}, err
	}
	return pc.Bldr.Build(), nil
}
