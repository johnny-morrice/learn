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
		expectedHeap   []uint64
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

			actualBytecode := vm.Memory[:len(tc.expectedMemory)]
			if !reflect.DeepEqual(tc.expectedMemory, actualBytecode) {
				t.Errorf("expected bytecode: %v\nactual: %v", tc.expectedMemory, actualBytecode)
				for i, exp := range tc.expectedMemory {
					act := actualBytecode[i]
					if exp != act {
						t.Errorf("first difference at %v\nexpected: %v\nactual: %v", i, exp, act)
						return
					}
				}
			}

			if len(tc.expectedHeap) > 0 {
				actualHeap := vm.Memory[vm.HeapStart:len(tc.expectedHeap)]
				if !reflect.DeepEqual(tc.expected heap, actualHeap) {
					t.Errorf("expected heap: %v\nactual: %v", tc.expected heap, actualHeap)
					for i, exp := range tc.expected heap {
						act := actualHeap[i]
						if exp != act {
							t.Errorf("first difference at %v\nexpected: %v\nactual: %v", i, exp, act)
							return
						}
					}
				}
			}

		})
	}
}
