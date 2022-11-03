package parser

import "github.com/johnny-morrice/learn/vmlang/assembler"

type ParseCombinator func(ParseContext) ParseContext

func PAnd(combs ...ParseCombinator) ParseCombinator {
	return func(pCtx ParseContext) ParseContext {
		loopCtx := pCtx
		for _, combinator := range combs {
			loopCtx = combinator(loopCtx)
			if loopCtx.Failed {
				return loopCtx
			}
		}
		return loopCtx
	}
}

func POr(combs ...ParseCombinator) ParseCombinator {
	return func(pCtx ParseContext) ParseContext {
		loopCtx := pCtx
		for _, combinator := range combs {
			loopCtx = combinator(loopCtx)
			if !loopCtx.Failed {
				return loopCtx
			}

		}
		pCtx.Failed = true
		return pCtx
	}
}

func ParseFile(fileName string) (*assembler.AST, error) {
	panic("not implemented")
}

func PEof() ParseCombinator {
	return func(pCtx ParseContext) ParseContext {
		pCtx.Failed = len(pCtx.RemainingInput) != 0
		return pCtx
	}
}

func PStmt() ParseCombinator {
	panic("not implemented")
}

func PAst() ParseCombinator {

	return func(pCtx ParseContext) ParseContext {
		f := POr(PEof(), PAnd(PStmt(), PAst()))
		return f(pCtx)
	}
}

type ParseContext struct {
	FileName       string
	Line           string
	Char           string
	RemainingInput string
	Failed         bool
}

func Parse(pCtx ParseContext) (*assembler.AST, error) {
	return &assembler.AST{}, nil
}
