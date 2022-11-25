package ast

import (
	"errors"

	"github.com/johnny-morrice/learn/vmlang/vm"
)

type Builder struct {
	ast *AST
}

func (bldr Builder) AddVarStmt() Builder {
	stmt := Stmt{
		Var: &VarStmt{},
	}
	bldr.ast.Stmts = append(bldr.ast.Stmts, stmt)

	return bldr
}

func (bldr Builder) AddLabelStmt(label string) Builder {
	stmt := Stmt{
		Label: &LabelStmt{label},
	}
	bldr.ast.Stmts = append(bldr.ast.Stmts, stmt)

	return bldr
}

func (bldr Builder) AddOpStmt(op vm.Bytecode) Builder {
	stmt := Stmt{
		Op: &OpStmt{Op: op},
	}
	bldr.ast.Stmts = append(bldr.ast.Stmts, stmt)

	return bldr
}

func (bldr Builder) AddVar(varName string) (Builder, error) {
	var nope Builder
	if len(bldr.ast.Stmts) == 0 {
		return nope, errors.New("expected statement")
	}
	lastIndex := len(bldr.ast.Stmts) - 1
	lastItem := bldr.ast.Stmts[lastIndex]

	if lastItem.Var == nil {
		return nope, errors.New("expected var statement")
	}

	lastItem.Var.VarNames = append(lastItem.Var.VarNames, varName)

	bldr.ast.Stmts[lastIndex] = lastItem

	return bldr, nil
}

func (bldr Builder) AddParam(param Param) (Builder, error) {
	var nope Builder
	if len(bldr.ast.Stmts) == 0 {
		return nope, errors.New("expected statement")
	}
	lastIndex := len(bldr.ast.Stmts) - 1
	lastItem := bldr.ast.Stmts[lastIndex]

	if lastItem.Op == nil {
		return nope, errors.New("expected op statement")
	}

	lastItem.Op.Params = append(lastItem.Op.Params, param)

	bldr.ast.Stmts[lastIndex] = lastItem

	return bldr, nil
}

func (bldr Builder) Build() *AST {
	return bldr.ast
}

func (bldr Builder) RemoveLastStmt() Builder {
	if len(bldr.ast.Stmts) > 0 {
		bldr.ast.Stmts = bldr.ast.Stmts[:len(bldr.ast.Stmts)-1]
	}
	return bldr
}

func (bldr Builder) ensureAST() Builder {
	if bldr.ast == nil {
		bldr.ast = &AST{}
	}
	return bldr
}
