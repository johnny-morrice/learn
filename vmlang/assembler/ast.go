package assembler

import (
	"fmt"
	"strings"

	"github.com/johnny-morrice/learn/vmlang/vm"
)

type AsmStmt struct {
	varStmt   *varStmt
	opStmt    *opStmt
	labelStmt *labelStmt
}

func (stmt AsmStmt) String() string {
	if stmt.varStmt != nil {
		return stmt.varStmt.String()
	}
	if stmt.opStmt != nil {
		return stmt.opStmt.String()
	}
	if stmt.labelStmt != nil {
		return stmt.labelStmt.String()
	}
	return "[invalid AsmStmt]"
}

type varStmt struct {
	varNames []string
}

func (stmt varStmt) String() string {
	return "var " + strings.Join(stmt.varNames, " ")
}

type opStmt struct {
	op         vm.Bytecode
	parameters []param
}

func (stmt opStmt) String() string {
	builder := strings.Builder{}
	builder.WriteString(stmt.op.String())
	for _, param := range stmt.parameters {
		builder.WriteString(" ")
		builder.WriteString(param.String())
	}
	return builder.String()
}

type labelStmt struct {
	labelName string
}

func (stmt labelStmt) String() string {
	return stmt.labelName + ":"
}

type param struct {
	literal  uint64
	variable string
}

func (p param) String() string {
	if p.variable != "" {
		return p.variable
	}
	return fmt.Sprint(p.literal)
}

type AsmScript struct {
	stmts []AsmStmt
}

func (tree AsmScript) String() string {
	builder := strings.Builder{}
	for i, stmt := range tree.stmts {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(stmt.String())
	}
	return builder.String()
}
