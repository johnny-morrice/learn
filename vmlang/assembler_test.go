package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestAssembler(t *testing.T) {
	type testCase struct {
		ast              asmScript
		expectedBytecode []uint64
		expectedHeap     []uint64
		expectedError    error
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
			expectedBytecode: []uint64{uint64(Push), 4, 0},
		},
		"write to heap var": {
			ast: asmScript{
				stmts: []asmStmt{
					{
						varStmt: &varStmt{
							varNames: []string{"TestVar"},
						},
					},
					{
						opStmt: &opStmt{
							op: WriteMemory,
							parameters: []param{
								{
									variable: "TestVar",
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(WriteMemory), 2 + gapSize + stackSize + gapSize, 0},
		},
		"var can be defined anywhere": {
			ast: asmScript{
				stmts: []asmStmt{
					{
						opStmt: &opStmt{
							op: WriteMemory,
							parameters: []param{
								{
									variable: "TestVar",
								},
							},
						},
					},
					{
						varStmt: &varStmt{
							varNames: []string{"TestVar"},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(WriteMemory), 2 + gapSize + stackSize + gapSize, 0},
		},
		"go to label": {
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
					{
						labelStmt: &labelStmt{
							labelName: "TestLabel",
						},
					},
					{
						opStmt: &opStmt{
							op: Push,
							parameters: []param{
								{
									literal: 5,
								},
							},
						},
					},
					{
						opStmt: &opStmt{
							op: Pop,
						},
					},
					{
						opStmt: &opStmt{
							op: Goto,
							parameters: []param{
								{
									variable: "TestLabel",
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(Push), 4, uint64(Push), 5, uint64(Pop), uint64(Goto), 2, 0},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			vm, err := assemble(tc.ast)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected err: %s\nactual: %s", tc.expectedError, err)
			}

			actualBytecode := vm.Memory[:len(tc.expectedBytecode)]
			if !reflect.DeepEqual(tc.expectedBytecode, actualBytecode) {
				t.Errorf("expected bytecode: %v\nactual: %v", tc.expectedBytecode, actualBytecode)
				for i, exp := range tc.expectedBytecode {
					act := actualBytecode[i]
					if exp != act {
						t.Errorf("first difference at %v\nexpected: %v\nactual: %v", i, exp, act)
						return
					}
				}
			}

			if len(tc.expectedHeap) > 0 {
				actualHeap := vm.Memory[vm.HeapStart:len(tc.expectedHeap)]
				if !reflect.DeepEqual(tc.expectedHeap, actualHeap) {
					t.Errorf("expected heap: %v\nactual: %v", tc.expectedHeap, actualHeap)
					for i, exp := range tc.expectedHeap {
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
