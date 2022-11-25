package parser

import (
	"unicode"
	"unicode/utf8"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
	"github.com/johnny-morrice/learn/vmlang/vm"
)

type ParseCombinator func(pc ParseContext) ParseContext

func Seq(combs ...ParseCombinator) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		initCtx := pc
		loopCtx := pc
		for _, combinator := range combs {
			loopCtx = combinator(loopCtx)
			if loopCtx.Failed {
				initCtx.Failed = true
				return initCtx
			}
		}
		return loopCtx
	}
}

func Alt(combs ...ParseCombinator) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		loopCtx := pc
		for _, combinator := range combs {
			loopCtx = combinator(loopCtx)
			if !loopCtx.Failed {
				return loopCtx
			}
		}
		pc.Failed = true
		return pc
	}
}

func EOF() ParseCombinator {
	return func(pc ParseContext) ParseContext {
		pc.Failed = len(pc.RemainingInput) != 0
		return pc
	}
}

func TextEq(text string) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		if len(pc.RemainingInput) < len(text) {
			pc.Failed = true
			return pc
		}
		input := pc.RemainingInput[:len(text)]
		pc.Failed = text != input
		if pc.Failed {
			return pc
		}
		pc.RemainingInput = pc.RemainingInput[len(text):]
		return pc
	}
}

func WhiteChar() ParseCombinator {
	return func(pc ParseContext) ParseContext {
		return Alt(TextEq(" "), TextEq("\t"))(pc)
	}
}

func Repeat(comb ParseCombinator) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		loopCtx := pc
		for {
			nextCtx := comb(loopCtx)
			if nextCtx.Failed {
				return loopCtx
			}
			loopCtx = nextCtx
		}
	}
}

func OpName() ParseCombinator {
	opCombs := []ParseCombinator{}
	for _, op := range vm.Bytecodes() {
		opCombs = append(opCombs, TextEq(op.String()))
	}
	return Alt(opCombs...)
}

func Letter() ParseCombinator {
	return MatchRune(unicode.IsLetter)
}

func Digit() ParseCombinator {
	return MatchRune(unicode.IsDigit)
}

func VarName() ParseCombinator {
	return Seq(Letter(), Repeat(Alt(Letter(), Digit())))
}

func Number() ParseCombinator {
	return Seq(Digit(), Repeat(Digit()))
}

func MatchRune(matcher func(r rune) bool) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		r, size := utf8.DecodeRuneInString(pc.RemainingInput)
		pc.Failed = utf8.RuneError == r || size == 0 || !matcher(r)
		if pc.Failed {
			return pc
		}
		pc.RemainingInput = pc.RemainingInput[size:]
		return pc
	}
}

func Whitespace() ParseCombinator {
	return Seq(WhiteChar(), Repeat(WhiteChar()))
}

func OptionalWhitespace() ParseCombinator {
	return Repeat(WhiteChar())
}

func Newline() ParseCombinator {
	return Alt(TextEq("\n"), TextEq("\r\n"))
}

func StmtEnd() ParseCombinator {
	return Seq(OptionalWhitespace(), Alt(Newline(), EOF()))
}

func VarStmt() ParseCombinator {
	return Seq(
		TextEq("var"),
		WithBuilder(func(bldr ast.Builder) (ast.Builder, error) {
			return bldr.AddVarStmt(), nil
		}),
		Whitespace(),
		VarName(),
		Repeat(Seq(Whitespace(), VarName())),
		CompleteStmt(),
	)
}

func OpStmt() ParseCombinator {
	return Seq(
		OpName(),
		Repeat(
			Seq(
				Whitespace(),
				Alt(VarName(), Number()),
			),
		),
	)
}

func LabelStmt() ParseCombinator {
	return Seq(
		VarName(), TextEq(":"),
	)
}

func Stmt() ParseCombinator {
	return Seq(
		OptionalWhitespace(),
		Alt(LabelStmt(), VarStmt(), OpStmt()),
		StmtEnd(),
	)
}

func AST() ParseCombinator {

	return func(pc ParseContext) ParseContext {
		f := Seq(Repeat(Stmt()), EOF())
		return f(pc)
	}
}

type BuilderFunc func(bldr ast.Builder) (ast.Builder, error)

func WithBuilder(f BuilderFunc) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		bldr, err := f(pc.Bldr)
		if err != nil {
			pc.Failed = true
			pc.Error = err
		}
		pc.Bldr = bldr
		return pc
	}
}

func CompleteStmt() ParseCombinator {
	return WithBuilder(func(bldr ast.Builder) (ast.Builder, error) {
		return bldr.CompleteStmt()
	})
}
