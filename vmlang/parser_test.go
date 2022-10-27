package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	_ "embed"
)

//go:embed asm/fac.vmsm
var factorialSourceCode string

func TestParserCreatesAST(t *testing.T) {
	type testCase struct {
		pCtx          ParseContext
		sourceCode    string
		expectedAst   asmScript
		expectedError error
	}

	testCases := map[string]testCase{
		"factorial": {
			pCtx: ParseContext{
				FileName: "fac.vmsm",
			},
			sourceCode:  factorialSourceCode,
			expectedAst: factorialAst(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			reader := strings.NewReader(tc.sourceCode)
			actualAst, err := Parse(tc.pCtx, reader)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected err:%s\nactual: %s", tc.expectedError, err)
			}
			if !reflect.DeepEqual(tc.expectedAst, actualAst) {
				t.Errorf("expected: %v\n\nbut was: %v", tc.expectedAst, actualAst)
			}
		})
	}
}
