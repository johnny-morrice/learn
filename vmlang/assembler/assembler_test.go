package assembler

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/johnny-morrice/learn/vmlang/vm"
)

func TestAssembleAndRunProgram(t *testing.T) {
	type testCase struct {
		ast            AsmScript
		expectedOutput []byte
		expectedError  error
	}

	testCases := map[string]testCase{
		"goto label": {
			ast: AsmScript{
				stmts: []AsmStmt{
					{
						opStmt: &opStmt{
							vm.Push,
							[]param{
								{
									literal: 5,
								},
							},
						},
					},
					{
						opStmt: &opStmt{
							vm.Goto,
							[]param{
								{
									variable: "TestLabel",
								},
							},
						},
					},
					{
						opStmt: &opStmt{
							op: vm.Pop,
						},
					},
					{
						labelStmt: &labelStmt{
							"TestLabel",
						},
					},
					{
						opStmt: &opStmt{
							op: vm.OutputByte,
						},
					},
					{
						opStmt: &opStmt{
							op: vm.Exit,
						},
					},
				},
			},
			expectedOutput: []byte{5},
		},
		"factorial": {
			ast:            factorialAst(),
			expectedOutput: []byte{24},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			vm, err := Assemble(&tc.ast)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected assemble err: %s\nactual: %s", tc.expectedError, err)
			}
			buf := &bytes.Buffer{}
			vm.Output = buf
			err = vm.Execute()
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected vm err: %s\nactual: %s", tc.expectedError, err)
			}
			actual := buf.Bytes()
			if !reflect.DeepEqual(tc.expectedOutput, actual) {
				t.Fatalf("expected output: %v but received: %v", tc.expectedOutput, actual)
				return
			}
		})
	}
}

func factorialAst() AsmScript {
	return AsmScript{
		stmts: []AsmStmt{
			{
				varStmt: &varStmt{
					[]string{"acc"},
				},
			},
			{
				opStmt: &opStmt{
					vm.Push, []param{{literal: 4}},
				},
			},
			{
				opStmt: &opStmt{
					vm.Push, []param{{variable: "acc"}},
				},
			},
			{
				opStmt: &opStmt{
					op: vm.WriteMemory,
				},
			},
			{
				labelStmt: &labelStmt{"fac"},
			},
			{
				opStmt: &opStmt{
					op: vm.Decrement,
				},
			},
			{
				opStmt: &opStmt{
					vm.JumpNotZero, []param{{variable: "body"}},
				},
			},
			{
				opStmt: &opStmt{
					vm.Goto, []param{{variable: "output"}},
				},
			},
			{
				labelStmt: &labelStmt{"body"},
			},
			{
				opStmt: &opStmt{
					op: vm.Duplicate,
				},
			},
			{
				opStmt: &opStmt{
					vm.Push, []param{{variable: "acc"}},
				},
			},
			{
				opStmt: &opStmt{
					op: vm.ReadMemory,
				},
			},
			{
				opStmt: &opStmt{
					op: vm.Multiply,
				},
			},
			{
				opStmt: &opStmt{
					vm.Push, []param{{variable: "acc"}},
				},
			},
			{
				opStmt: &opStmt{
					op: vm.WriteMemory,
				},
			},
			{
				opStmt: &opStmt{
					op: vm.Pop,
				},
			},
			{
				opStmt: &opStmt{
					vm.Goto, []param{{variable: "fac"}},
				},
			},
			{
				labelStmt: &labelStmt{"output"},
			},
			{
				opStmt: &opStmt{
					vm.Push, []param{{variable: "acc"}},
				},
			},
			{
				opStmt: &opStmt{
					op: vm.ReadMemory,
				},
			},
			{
				opStmt: &opStmt{
					op: vm.OutputByte,
				},
			},
		},
	}
}

func TestAssembleAsmScript(t *testing.T) {
	type testCase struct {
		ast              AsmScript
		expectedBytecode []uint64
		expectedHeap     []uint64
		expectedError    error
	}

	testCases := map[string]testCase{
		"simple push": {
			ast: AsmScript{
				stmts: []AsmStmt{
					{
						opStmt: &opStmt{
							op: vm.Push,
							parameters: []param{
								{
									literal: 4,
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.Push), 4, uint64(vm.Exit), 0},
		},
		"write to heap var": {
			ast: AsmScript{
				stmts: []AsmStmt{
					{
						varStmt: &varStmt{
							varNames: []string{"TestVar"},
						},
					},
					{
						opStmt: &opStmt{
							op: vm.WriteMemory,
							parameters: []param{
								{
									variable: "TestVar",
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.WriteMemory), 3 + gapSize + stackSize + gapSize, uint64(vm.Exit), 0},
		},
		"var can be defined anywhere": {
			ast: AsmScript{
				stmts: []AsmStmt{
					{
						opStmt: &opStmt{
							op: vm.WriteMemory,
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
			expectedBytecode: []uint64{uint64(vm.WriteMemory), 3 + gapSize + stackSize + gapSize, uint64(vm.Exit), 0},
		},
		"go to missing label": {
			ast: AsmScript{
				stmts: []AsmStmt{
					{
						opStmt: &opStmt{
							vm.Goto, []param{{variable: "foo"}},
						},
					},
				},
			},
			expectedError: ErrAssembler,
		},
		"go to label": {
			ast: AsmScript{
				stmts: []AsmStmt{
					{
						opStmt: &opStmt{
							op: vm.Push,
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
							op: vm.Push,
							parameters: []param{
								{
									literal: 5,
								},
							},
						},
					},
					{
						opStmt: &opStmt{
							op: vm.Pop,
						},
					},
					{
						opStmt: &opStmt{
							op: vm.Goto,
							parameters: []param{
								{
									variable: "TestLabel",
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.Push), 4, uint64(vm.Push), 5, uint64(vm.Pop), uint64(vm.Goto), 2, uint64(vm.Exit), 0},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			vm, err := Assemble(&tc.ast)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected err: %s\nactual: %s", tc.expectedError, err)
			}

			if len(tc.expectedBytecode) > 0 {
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
