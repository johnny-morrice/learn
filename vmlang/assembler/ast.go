package assembler

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
	if stmt.Var != nil {
		return stmt.Var.String()
	}
	if stmt.Op != nil {
		return stmt.Op.String()
	}
	if stmt.Label != nil {
		return stmt.Label.String()
	}
	return "[invalid AsmStmt]"
}

type VarStmt struct {
	varNames []string
}

func (stmt VarStmt) String() string {
	return "var " + strings.Join(stmt.varNames, " ")
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
