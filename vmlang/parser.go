package main

import "fmt"

func ParseFile(fileName string) (*VirtualMachine, error) {
	panic("not implemented")
}

type asmStmt struct {
	varStmt         *varStmt
	instructionStmt *instructionStmt
	labelStmt       *labelStmt
}
type varStmt struct {
	varNames []string
}
type instructionStmt struct {
	instruction Bytecode
	parameters  []param
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
type intrParam struct {
	value     *uint64
	varName   string
	labelName string
}
type intrStmt struct {
	size        int
	instruction Bytecode
	parameters  []intrParam
	label       string
}
type asmState struct {
	varTable   map[string]int
	nameTable  map[string]*uint64
	labelTable map[string]struct{}
	stmts      []intrStmt
}

func (state *asmState) defineVar(varName string) error {
	_, exists := state.nameTable[varName]
	if exists {
		return fmt.Errorf("duplicate variable definition: %s", varName)
	}
	state.varTable[varName] = len(state.varTable)
	val := uint64(0)
	state.nameTable[varName] = &val

	return nil
}

func (state *asmState) defineLabel(labelName string) error {
	_, exists := state.nameTable[labelName]
	if exists {
		return fmt.Errorf("duplicate variable definition: %s", labelName)
	}
	state.labelTable[labelName] = struct{}{}
	val := uint64(0)
	state.nameTable[labelName] = &val

	return nil
}

func (state *asmState) addIntructionStmt(stmt instructionStmt) {
	intr := intrStmt{}
	intr.size = 1 + len(stmt.parameters)

	for _, param := range stmt.parameters {
		iParam := intrParam{}

		if param.variable == "" {
			iParam.value = &param.literal
			continue
		}

		_, varExists := state.varTable[param.variable]
		_, labelExists := state.labelTable[param.variable]

		addr := state.nameTable[param.variable]

		if varExists {
			iParam.varName = param.variable
		}
		if labelExists {
			iParam.labelName = param.variable
		}
		iParam.value = addr
		intr.parameters = append(intr.parameters, iParam)
	}

	state.stmts = append(state.stmts, intr)
}

func (state *asmState) addLabelStmt(stmt labelStmt) {
	intr := intrStmt{}

	intr.label = stmt.labelName

	state.stmts = append(state.stmts, intr)
}

func (state *asmState) setNameAddress(name string, addr uint64) {
	ptr := state.nameTable[name]

	*ptr = addr
}

func parseTree(tree asmScript) (*VirtualMachine, error) {
	state := asmState{
		varTable:   map[string]int{},
		nameTable:  map[string]*uint64{},
		labelTable: map[string]struct{}{},
	}

	vm := &VirtualMachine{}

	for _, stmt := range tree.stmts {
		if stmt.varStmt != nil {
			for _, varName := range stmt.varStmt.varNames {
				state.defineVar(varName)
			}
		}
		if stmt.labelStmt != nil {
			state.defineLabel(stmt.labelStmt.labelName)
		}
	}

	for _, stmt := range tree.stmts {
		if stmt.instructionStmt != nil {
			state.addIntructionStmt(*stmt.instructionStmt)
		}
		if stmt.labelStmt != nil {
			state.addLabelStmt(*stmt.labelStmt)
		}
	}

	bytecodeSize := uint64(0)

	for _, iStmt := range state.stmts {
		if iStmt.label != "" {
			state.setNameAddress(iStmt.label, bytecodeSize)
			continue
		}
		bytecodeSize += uint64(iStmt.size)
	}

	return vm, nil

}
