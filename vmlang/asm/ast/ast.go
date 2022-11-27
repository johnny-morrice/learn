package ast

import (
	"fmt"
	"strings"

	"github.com/johnny-morrice/learn/vmlang/vm"
)

type Stmt struct {
	Var   *VarStmt
	Op    *OpStmt
	Label *LabelStmt
}

func (stmt Stmt) String() string {
	isVar := stmt.Var != nil
	isOp := stmt.Op != nil
	isLabel := stmt.Label != nil

	if countTrue(isVar, isOp, isLabel) > 1 {
		return "[invalid AsmStmt]"
	}

	if isVar {
		return stmt.Var.String()
	}
	if isOp {
		return stmt.Op.String()
	}
	if isLabel {
		return stmt.Label.String()
	}
	return "[empty AsmStmt]"
}

func countTrue(bools ...bool) int {
	x := 0
	for _, b := range bools {
		if b {
			x++
		}
	}
	return x
}

type VarStmt struct {
	VarNames []string
}

func (stmt VarStmt) String() string {
	return "var " + strings.Join(stmt.VarNames, " ")
}

type OpStmt struct {
	Op     vm.Bytecode
	Params []Param
}

func (stmt OpStmt) String() string {
	builder := strings.Builder{}
	builder.WriteString(stmt.Op.String())
	for _, param := range stmt.Params {
		builder.WriteString(" ")
		builder.WriteString(param.String())
	}
	return builder.String()
}

type LabelStmt struct {
	Label string
}

func (stmt LabelStmt) String() string {
	return stmt.Label + ":"
}

type Param struct {
	Literal  uint64
	Variable string
}

func (p Param) String() string {
	if p.Variable != "" {
		return p.Variable
	}
	return fmt.Sprint(p.Literal)
}

type AST struct {
	Stmts []Stmt
}

func (tree AST) String() string {
	builder := strings.Builder{}
	for i, stmt := range tree.Stmts {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(stmt.String())
	}
	return builder.String()
}
