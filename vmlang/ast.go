package main

type asmStmt struct {
	varStmt   *varStmt
	opStmt    *opStmt
	labelStmt *labelStmt
}
type varStmt struct {
	varNames []string
}
type opStmt struct {
	op         Bytecode
	parameters []param
}
type labelStmt struct {
	labelName string
}
type param struct {
	literal  uint64
	variable string
}
type asmScript struct {
	stmts []asmStmt
}
