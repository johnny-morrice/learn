package example

import (
	"github.com/johnny-morrice/learn/vmlang/asm"
	"github.com/johnny-morrice/learn/vmlang/vm"
)

func FactorialAst() asm.AST {
	return asm.AST{
		Stmts: []asm.Stmt{
			{
				Var: &asm.VarStmt{
					[]string{"acc"},
				},
			},
			{
				Op: &asm.OpStmt{
					vm.Push, []asm.Param{{Literal: 4}},
				},
			},
			{
				Op: &asm.OpStmt{
					vm.Push, []asm.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.WriteMemory,
				},
			},
			{
				Label: &asm.LabelStmt{"fac"},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.Decrement,
				},
			},
			{
				Op: &asm.OpStmt{
					vm.JumpNotZero, []asm.Param{{Variable: "body"}},
				},
			},
			{
				Op: &asm.OpStmt{
					vm.Goto, []asm.Param{{Variable: "output"}},
				},
			},
			{
				Label: &asm.LabelStmt{"body"},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.Duplicate,
				},
			},
			{
				Op: &asm.OpStmt{
					vm.Push, []asm.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.ReadMemory,
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.Multiply,
				},
			},
			{
				Op: &asm.OpStmt{
					vm.Push, []asm.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.WriteMemory,
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.Pop,
				},
			},
			{
				Op: &asm.OpStmt{
					vm.Goto, []asm.Param{{Variable: "fac"}},
				},
			},
			{
				Label: &asm.LabelStmt{"output"},
			},
			{
				Op: &asm.OpStmt{
					vm.Push, []asm.Param{{Variable: "acc"}},
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.ReadMemory,
				},
			},
			{
				Op: &asm.OpStmt{
					Op: vm.OutputByte,
				},
			},
		},
	}
}
