package asm

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/johnny-morrice/learn/vmlang/vm"
)

func TestAssembleAndRunProgram(t *testing.T) {
	type testCase struct {
		ast            AST
		expectedOutput []byte
		expectedError  error
	}

	testCases := map[string]testCase{
		"goto label": {
			ast: AST{
				Stmts: []Stmt{
					{
						Op: &OpStmt{
							vm.Push,
							[]Param{
								{
									Literal: 5,
								},
							},
						},
					},
					{
						Op: &OpStmt{
							vm.Goto,
							[]Param{
								{
									Variable: "TestLabel",
								},
							},
						},
					},
					{
						Op: &OpStmt{
							Op: vm.Pop,
						},
					},
					{
						Label: &LabelStmt{
							"TestLabel",
						},
					},
					{
						Op: &OpStmt{
							Op: vm.OutputByte,
						},
					},
					{
						Op: &OpStmt{
							Op: vm.Exit,
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

func factorialAst() AST {
	return AST{
		Stmts: []Stmt{
			{
				Var: &VarStmt{
					[]string{"acc"},
				},
			},
			{
				Op: &OpStmt{
					vm.Push, []Param{{Literal: 4}},
				},
			},
			{
				Op: &OpStmt{
					vm.Push, []Param{{Variable: "acc"}},
				},
			},
			{
				Op: &OpStmt{
					Op: vm.WriteMemory,
				},
			},
			{
				Label: &LabelStmt{"fac"},
			},
			{
				Op: &OpStmt{
					Op: vm.Decrement,
				},
			},
			{
				Op: &OpStmt{
					vm.JumpNotZero, []Param{{Variable: "body"}},
				},
			},
			{
				Op: &OpStmt{
					vm.Goto, []Param{{Variable: "output"}},
				},
			},
			{
				Label: &LabelStmt{"body"},
			},
			{
				Op: &OpStmt{
					Op: vm.Duplicate,
				},
			},
			{
				Op: &OpStmt{
					vm.Push, []Param{{Variable: "acc"}},
				},
			},
			{
				Op: &OpStmt{
					Op: vm.ReadMemory,
				},
			},
			{
				Op: &OpStmt{
					Op: vm.Multiply,
				},
			},
			{
				Op: &OpStmt{
					vm.Push, []Param{{Variable: "acc"}},
				},
			},
			{
				Op: &OpStmt{
					Op: vm.WriteMemory,
				},
			},
			{
				Op: &OpStmt{
					Op: vm.Pop,
				},
			},
			{
				Op: &OpStmt{
					vm.Goto, []Param{{Variable: "fac"}},
				},
			},
			{
				Label: &LabelStmt{"output"},
			},
			{
				Op: &OpStmt{
					vm.Push, []Param{{Variable: "acc"}},
				},
			},
			{
				Op: &OpStmt{
					Op: vm.ReadMemory,
				},
			},
			{
				Op: &OpStmt{
					Op: vm.OutputByte,
				},
			},
		},
	}
}

func TestAssembleAsmScript(t *testing.T) {
	type testCase struct {
		ast              AST
		expectedBytecode []uint64
		expectedHeap     []uint64
		expectedError    error
	}

	testCases := map[string]testCase{
		"simple push": {
			ast: AST{
				Stmts: []Stmt{
					{
						Op: &OpStmt{
							Op: vm.Push,
							Params: []Param{
								{
									Literal: 4,
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.Push), 4, uint64(vm.Exit), 0},
		},
		"write to heap var": {
			ast: AST{
				Stmts: []Stmt{
					{
						Var: &VarStmt{
							VarNames: []string{"TestVar"},
						},
					},
					{
						Op: &OpStmt{
							Op: vm.WriteMemory,
							Params: []Param{
								{
									Variable: "TestVar",
								},
							},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.WriteMemory), 3 + gapSize + stackSize + gapSize, uint64(vm.Exit), 0},
		},
		"var can be defined anywhere": {
			ast: AST{
				Stmts: []Stmt{
					{
						Op: &OpStmt{
							Op: vm.WriteMemory,
							Params: []Param{
								{
									Variable: "TestVar",
								},
							},
						},
					},
					{
						Var: &VarStmt{
							VarNames: []string{"TestVar"},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.WriteMemory), 3 + gapSize + stackSize + gapSize, uint64(vm.Exit), 0},
		},
		"go to missing label": {
			ast: AST{
				Stmts: []Stmt{
					{
						Op: &OpStmt{
							vm.Goto, []Param{{Variable: "foo"}},
						},
					},
				},
			},
			expectedError: ErrAssembler,
		},
		"go to label": {
			ast: AST{
				Stmts: []Stmt{
					{
						Op: &OpStmt{
							Op: vm.Push,
							Params: []Param{
								{
									Literal: 4,
								},
							},
						},
					},
					{
						Label: &LabelStmt{
							Label: "TestLabel",
						},
					},
					{
						Op: &OpStmt{
							Op: vm.Push,
							Params: []Param{
								{
									Literal: 5,
								},
							},
						},
					},
					{
						Op: &OpStmt{
							Op: vm.Pop,
						},
					},
					{
						Op: &OpStmt{
							Op: vm.Goto,
							Params: []Param{
								{
									Variable: "TestLabel",
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
