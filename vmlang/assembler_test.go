package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestAssembler(t *testing.T) {
	type testCase struct {
		ast            asmScript
		expectedMemory []uint64
		expectedError  error
	}

	testCases := map[string]testCase{
		"simple push": {
			ast: asmScript{
				stmts: []asmStmt{
					{
						opStmt: &opStmt{
							op: Push,
							parameters: []param{
								{
									literal: 4,
								},
							},
						},
					},
				},
			},
			expectedMemory: []uint64{1, 4, 0},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			vm, err := assemble(tc.ast)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected err: %s\nactual: %s", tc.expectedError, err)
			}
			actualMem := vm.Memory[:len(tc.expectedMemory)]
			if !reflect.DeepEqual(tc.expectedMemory, actualMem) {
				t.Errorf("expected: %v\nactual: %v", tc.expectedMemory, actualMem)
				for i, exp := range tc.expectedMemory {
					act := actualMem[i]
					if exp != act {
						t.Errorf("first difference at %v\nexpected: %v\nactual: %v", i, exp, act)
						return
					}
				}
			}
		})
	}
}
