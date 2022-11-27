package parser

import (
	"errors"
	"reflect"
	"testing"

	_ "embed"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
	"github.com/johnny-morrice/learn/vmlang/example"
)

func TestParserCreatesAST(t *testing.T) {
	type testCase struct {
		pCtx          ParseContext
		expectedAst   ast.AST
		expectedError error
	}

	testCases := map[string]testCase{
		"factorial": {
			pCtx: ParseContext{
				FileName:       "fac.vmsm",
				RemainingInput: example.FactorialSourceCode,
			},
			expectedAst: example.FactorialAst(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualAst, err := Parse(tc.pCtx)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected err:%s\nactual: %s", tc.expectedError, err)
			}
			if !reflect.DeepEqual(tc.expectedAst, actualAst) {
				t.Errorf("expected:\n%v\n\nactual:\n%v", tc.expectedAst, actualAst)

				lenExpect := len(tc.expectedAst.Stmts)
				lenActual := len(actualAst.Stmts)

				if lenExpect != lenActual {
					t.Errorf("expected %d statements but got %d", lenExpect, lenActual)
				}
				minLen := min(lenExpect, lenActual)
				for i := 0; i < minLen; i++ {
					expectStmt := tc.expectedAst.Stmts[i]
					actualStmt := actualAst.Stmts[i]
					if !reflect.DeepEqual(expectStmt, actualStmt) {
						t.Errorf("unexpected statement at %v\nexpected %v\nactual: %v", i, expectStmt, actualStmt)
					}
				}

			}
		})
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
