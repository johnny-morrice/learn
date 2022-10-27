package main

import (
	"io"
)

func ParseFile(fileName string) (*asmScript, error) {
	panic("not implemented")
}

type ParseContext struct {
	FileName string
	Line     string
	Char     string
}

func Parse(pCtx ParseContext, reader io.Reader) (*asmScript, error) {
	return &asmScript{}, nil
}
