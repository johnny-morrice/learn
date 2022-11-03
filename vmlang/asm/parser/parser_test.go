package parser

import (
	"errors"
	"reflect"
	"testing"

	_ "embed"

	"github.com/johnny-morrice/learn/vmlang/asm"
	"github.com/johnny-morrice/learn/vmlang/example"
)

//go:embed fac.vmsm
var factorialSourceCode string

func TestParserCreatesAST(t *testing.T) {
	type testCase struct {
		pCtx          ParseContext
		expectedAst   asm.AST
		expectedError error
	}

	testCases := map[string]testCase{
		"factorial": {
			pCtx: ParseContext{
				FileName:       "fac.vmsm",
				RemainingInput: factorialSourceCode,
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

func TestParserCombinators(t *testing.T) {
	type testCase struct {
		input    ParseContext
		comb     ParseCombinator
		expected ParseContext
	}

	testCases := map[string]testCase{
		"eof WhenIsEof": {
			comb: PEof(),
		},
		"eof WhenNotEof": {
			input: ParseContext{
				RemainingInput: "  ",
			},
			comb: PEof(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "  ",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := tc.comb(tc.input)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected: %v\n\nbut was: %v", tc.expected, actual)
			}
		})
	}
}
