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
	t.Skip()
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
				t.Errorf("expected: %v\n\nbut was: %v", tc.expectedAst, actualAst)
			}
		})
	}
}
