package asm

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/johnny-morrice/learn/vmlang/asm/ast"
	"github.com/johnny-morrice/learn/vmlang/example"
	"github.com/johnny-morrice/learn/vmlang/vm"
)

func TestAssembleAndRunProgram(t *testing.T) {
	type testCase struct {
		ast            ast.AST
		expectedOutput []byte
		expectedError  error
	}

	testCases := map[string]testCase{
		"goto label": {
			ast: ast.AST{
				Stmts: []ast.Stmt{
					{
						Op: &ast.OpStmt{
							vm.Push,
							[]ast.Param{
								{
									Literal: 5,
								},
							},
						},
					},
					{
						Op: &ast.OpStmt{
							vm.Goto,
							[]ast.Param{
								{
									Variable: "TestLabel",
								},
							},
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.Pop,
						},
					},
					{
						Label: &ast.LabelStmt{
							"TestLabel",
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.OutputByte,
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.Exit,
						},
					},
				},
			},
			expectedOutput: []byte{5},
		},
		"factorial": {
			ast:            example.FactorialAst(),
			expectedOutput: []byte{24},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			vm, err := Assemble(tc.ast)
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

func TestAssembleAsmScript(t *testing.T) {
	type testCase struct {
		ast              ast.AST
		expectedBytecode []uint64
		expectedHeap     []uint64
		expectedError    error
	}

	testCases := map[string]testCase{
		"simple push": {
			ast: ast.AST{
				Stmts: []ast.Stmt{
					{
						Op: &ast.OpStmt{
							Op: vm.Push,
							Params: []ast.Param{
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
			ast: ast.AST{
				Stmts: []ast.Stmt{
					{
						Var: &ast.VarStmt{
							VarNames: []string{"TestVar"},
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.WriteMemory,
							Params: []ast.Param{
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
			ast: ast.AST{
				Stmts: []ast.Stmt{
					{
						Op: &ast.OpStmt{
							Op: vm.WriteMemory,
							Params: []ast.Param{
								{
									Variable: "TestVar",
								},
							},
						},
					},
					{
						Var: &ast.VarStmt{
							VarNames: []string{"TestVar"},
						},
					},
				},
			},
			expectedBytecode: []uint64{uint64(vm.WriteMemory), 3 + gapSize + stackSize + gapSize, uint64(vm.Exit), 0},
		},
		"go to missing label": {
			ast: ast.AST{
				Stmts: []ast.Stmt{
					{
						Op: &ast.OpStmt{
							vm.Goto, []ast.Param{{Variable: "foo"}},
						},
					},
				},
			},
			expectedError: ErrAssembler,
		},
		"go to label": {
			ast: ast.AST{
				Stmts: []ast.Stmt{
					{
						Op: &ast.OpStmt{
							Op: vm.Push,
							Params: []ast.Param{
								{
									Literal: 4,
								},
							},
						},
					},
					{
						Label: &ast.LabelStmt{
							Label: "TestLabel",
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.Push,
							Params: []ast.Param{
								{
									Literal: 5,
								},
							},
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.Pop,
						},
					},
					{
						Op: &ast.OpStmt{
							Op: vm.Goto,
							Params: []ast.Param{
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
			vm, err := Assemble(tc.ast)
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
