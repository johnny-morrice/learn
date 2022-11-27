package ast

import (
	"errors"

	"github.com/johnny-morrice/learn/vmlang/collections"
	"github.com/johnny-morrice/learn/vmlang/vm"
)

type Builder struct {
	Stmts       collections.List[Stmt]
	CurrentStmt Stmt
	Vars        collections.List[string]
	Params      collections.List[Param]
}

func (bldr Builder) AddVarStmt() Builder {
	bldr.CurrentStmt = Stmt{
		Var: &VarStmt{},
	}
	return bldr
}

func (bldr Builder) AddLabelStmt(label string) Builder {
	bldr.CurrentStmt = Stmt{
		Label: &LabelStmt{label},
	}
	return bldr
}

func (bldr Builder) AddOpStmt(op vm.Bytecode) Builder {
	bldr.CurrentStmt = Stmt{
		Op: &OpStmt{Op: op},
	}
	return bldr
}

func (bldr Builder) AddVar(varName string) (Builder, error) {
	var nope Builder

	if bldr.CurrentStmt.Var == nil {
		return nope, errors.New("expected var statement")
	}

	bldr.Vars = bldr.Vars.Append(varName)
	return bldr, nil
}

func (bldr Builder) AddParam(param Param) (Builder, error) {
	var nope Builder

	if bldr.CurrentStmt.Op == nil {
		return nope, errors.New("expected op statement")
	}

	bldr.Params = bldr.Params.Append(param)

	return bldr, nil
}

func (bldr Builder) CompleteStmt() (Builder, error) {
	var nope Builder
	if bldr.CurrentStmt.Label == nil && bldr.CurrentStmt.Var == nil && bldr.CurrentStmt.Op == nil {
		return nope, errors.New("expected initialised statement")
	}
	if bldr.CurrentStmt.Var != nil {
		bldr.CurrentStmt.Var.VarNames = bldr.Vars.Slice()
	}
	if bldr.CurrentStmt.Op != nil {
		bldr.CurrentStmt.Op.Params = bldr.Params.Slice()
	}
	bldr.Stmts = bldr.Stmts.Append(bldr.CurrentStmt)
	bldr.CurrentStmt = Stmt{}
	bldr.Params = collections.List[Param]{}
	bldr.Vars = collections.List[string]{}
	return bldr, nil
}

func (bldr Builder) Build() AST {
	return AST{
		Stmts: bldr.Stmts.Slice(),
	}
}
