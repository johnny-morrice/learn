package parser

import (
	"errors"
	"log"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
	"github.com/johnny-morrice/learn/vmlang/vm"
)

const logBacktrack = true

type ParseCombinator func(pc ParseContext) ParseContext

func Seq(name string, combs ...ParseCombinator) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		loopCtx := pc
		for _, combinator := range combs {
			loopCtx = combinator(loopCtx)
			if loopCtx.Failed {
				if logBacktrack {
					log.Printf("%s Seq backtrack: %s", name, pc.Error)
				}
				pc.Failed = true
				return pc
			}
		}
		return loopCtx
	}
}

func Alt(name string, combs ...ParseCombinator) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		loopCtx := pc
		for _, combinator := range combs {
			loopCtx = combinator(loopCtx)
			if !loopCtx.Failed {
				if logBacktrack {
					log.Printf("%s Alt backtrack: %s", name, pc.Error)
				}
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
			pc.Error = errors.New("expected: " + text)
			if logBacktrack {
				log.Printf("TextEq backtrack: %s", pc.Error)
			}
			return pc
		}
		input := pc.RemainingInput[:len(text)]
		pc.Failed = text != input
		if pc.Failed {
			pc.Error = errors.New("expected: " + text)
			if logBacktrack {
				log.Printf("TextEq backtrack: %s", pc.Error)
			}
			return pc
		}
		if pc.IsCapturing {
			pc.CapturedText = pc.CapturedText + text
		}
		pc.RemainingInput = pc.RemainingInput[len(text):]
		return pc
	}
}

func WhiteChar() ParseCombinator {
	return func(pc ParseContext) ParseContext {
		return Alt("WhiteChar", TextEq(" "), TextEq("\t"))(pc)
	}
}

func Repeat(comb ParseCombinator) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		loopCtx := pc
		for {
			nextCtx := comb(loopCtx)
			if nextCtx.Failed {
				if logBacktrack {
					log.Printf("Repeat backtrack: %s", nextCtx.Error)
				}
				return loopCtx
			}
			loopCtx = nextCtx
		}
	}
}

func OpName() ParseCombinator {
	opCombs := []ParseCombinator{}
	for _, iterOp := range vm.Bytecodes() {
		op := iterOp
		opCombs = append(opCombs,
			Seq("Op"+op.String(),
				TextEq(op.String()),
				WithBuilder(func(bldr ast.Builder) (ast.Builder, error) {
					return bldr.AddOpStmt(op), nil
				}),
			),
		)
	}
	return Alt("OpName", opCombs...)
}

func Letter() ParseCombinator {
	return MatchRune("IsLetter", unicode.IsLetter)
}

func Digit() ParseCombinator {
	return MatchRune("IsDigit", unicode.IsDigit)
}

func VarName() ParseCombinator {
	return Seq("VarName", Letter(), Repeat(Alt("VarNameContinue", Letter(), Digit())))
}

func Number() ParseCombinator {
	return Seq("Number", Digit(), Repeat(Digit()))
}

func MatchRune(name string, matcher func(r rune) bool) ParseCombinator {
	return func(pc ParseContext) ParseContext {
		r, size := utf8.DecodeRuneInString(pc.RemainingInput)
		pc.Failed = utf8.RuneError == r || size == 0 || !matcher(r)
		if pc.Failed {
			pc.Error = errors.New("unexpected rune: " + name)
			if logBacktrack {
				log.Printf("MatchRune backtrack: %s", pc.Error)
			}
			return pc
		}
		if pc.IsCapturing {
			pc.CapturedText = pc.CapturedText + string(r)
		}
		pc.RemainingInput = pc.RemainingInput[size:]
		return pc
	}
}

func Whitespace() ParseCombinator {
	return Seq("Whitespace", WhiteChar(), Repeat(WhiteChar()))
}

func OptionalWhitespace() ParseCombinator {
	return Repeat(WhiteChar())
}

func Newline() ParseCombinator {
	return Alt("Newline", TextEq("\n"), TextEq("\r\n"))
}

func StmtEnd() ParseCombinator {
	return Seq("StmtEnd", OptionalWhitespace(), Alt("Newline", Newline(), EOF()))
}

func VarStmt() ParseCombinator {
	return Seq(
		"VarStmt",
		TextEq("var"),
		WithBuilder(func(bldr ast.Builder) (ast.Builder, error) {
			return bldr.AddVarStmt(), nil
		}),
		Whitespace(),
		VarDecl(),
		Repeat(Seq("VarDecls", Whitespace(), VarDecl())),
		CompleteStmt(),
	)
}

func VarDecl() ParseCombinator {
	return Seq(
		"VarDecl",
		StartCapture(),
		VarName(),
		StopCapture(),
		func(pc ParseContext) ParseContext {
			bldr, err := pc.Bldr.AddVar(pc.CapturedText)
			if err != nil {
				pc.Failed = true
				pc.Error = err
			} else {
				pc.CapturedText = ""
				pc.Bldr = bldr
			}
			return pc
		},
	)
}

func OpStmt() ParseCombinator {
	return Seq(
		"OpStmt",
		OpName(),
		Repeat(
			Seq(
				"OpParams",
				Whitespace(),
				Alt(
					"OpParamAlt",
					Seq(
						"OpVarParam",
						StartCapture(),
						VarName(),
						StopCapture(),
						func(pc ParseContext) ParseContext {
							bldr, err := pc.Bldr.AddParam(ast.Param{Variable: pc.CapturedText})
							if err != nil {
								pc.Failed = true
								pc.Error = err
							} else {
								pc.CapturedText = ""
								pc.Bldr = bldr
							}
							return pc
						},
					),
					Seq(
						"OpLitParam",
						StartCapture(),
						Number(),
						StopCapture(),
						func(pc ParseContext) ParseContext {
							num, err := strconv.ParseUint(pc.CapturedText, 10, 64)
							if err != nil {
								pc.Failed = true
								pc.Error = err
							}
							bldr, err := pc.Bldr.AddParam(ast.Param{Literal: num})
							if err != nil {
								pc.Failed = true
								pc.Error = err
							} else {
								pc.CapturedText = ""
								pc.Bldr = bldr
							}
							return pc
						},
					),
				),
			),
		),
		CompleteStmt(),
	)
}

func LabelStmt() ParseCombinator {
	return Seq(
		"LabelStmt",
		StartCapture(),
		VarName(),
		StopCapture(),
		TextEq(":"),
		func(pc ParseContext) ParseContext {
			pc.Bldr = pc.Bldr.AddLabelStmt(pc.CapturedText)
			pc.CapturedText = ""
			return pc
		},
		CompleteStmt(),
	)
}

func Stmt() ParseCombinator {
	return Seq(
		"Stmt",
		OptionalWhitespace(),
		Alt(
			"StmtAlt",
			LabelStmt(), VarStmt(), OpStmt()),
		StmtEnd(),
	)
}

func AST() ParseCombinator {
	return func(pc ParseContext) ParseContext {
		f := Seq("AST", Repeat(Stmt()), EOF())
		return f(pc)
	}
}

func StartCapture() ParseCombinator {
	return func(pc ParseContext) ParseContext {
		pc.IsCapturing = true
		return pc
	}
}

func StopCapture() ParseCombinator {
	return func(pc ParseContext) ParseContext {
		pc.IsCapturing = false
		return pc
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
