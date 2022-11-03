package example

import (
	"github.com/johnny-morrice/learn/vmlang/asm/ast"
	"github.com/johnny-morrice/learn/vmlang/vm"

	_ "embed"
)

//go:embed asm/fac.vmsm
var FactorialSourceCode string

func FactorialAst() ast.AST {
	return ast.AST{
		Stmts: []ast.Stmt{
			{
				Var: &ast.VarStmt{
					[]string{"acc"},
				},
			},
			{
				Op: &ast.OpStmt{
					vm.Push, []ast.Param{{Literal: 4}},
				},
			},
			{
				Op: &ast.OpStmt{
					vm.Push, []ast.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.WriteMemory,
				},
			},
			{
				Label: &ast.LabelStmt{"fac"},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.Decrement,
				},
			},
			{
				Op: &ast.OpStmt{
					vm.JumpNotZero, []ast.Param{{Variable: "body"}},
				},
			},
			{
				Op: &ast.OpStmt{
					vm.Goto, []ast.Param{{Variable: "output"}},
				},
			},
			{
				Label: &ast.LabelStmt{"body"},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.Duplicate,
				},
			},
			{
				Op: &ast.OpStmt{
					vm.Push, []ast.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.ReadMemory,
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.Multiply,
				},
			},
			{
				Op: &ast.OpStmt{
					vm.Push, []ast.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.WriteMemory,
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.Pop,
				},
			},
			{
				Op: &ast.OpStmt{
					vm.Goto, []ast.Param{{Variable: "fac"}},
				},
			},
			{
				Label: &ast.LabelStmt{"output"},
			},
			{
				Op: &ast.OpStmt{
					vm.Push, []ast.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.ReadMemory,
				},
			},
			{
				Op: &ast.OpStmt{
					Op: vm.OutputByte,
				},
			},
		},
	}
}
