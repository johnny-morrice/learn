package main

import "fmt"

type intrParam struct {
	value     *uint64
	varName   string
	labelName string
}
type intrOp struct {
	size        int
	instruction Bytecode
	parameters  []intrParam
	label       string
}

type assembler struct {
	varTable   map[string]int
	nameTable  map[string]*uint64
	labelTable map[string]struct{}
	stmts      []intrOp
}

func (asm *assembler) defineVar(varName string) error {
	_, exists := asm.nameTable[varName]
	if exists {
		return fmt.Errorf("duplicate variable definition: %s", varName)
	}
	asm.varTable[varName] = len(asm.varTable)
	val := uint64(0)
	asm.nameTable[varName] = &val

	return nil
}

func (asm *assembler) defineLabel(labelName string) error {
	_, exists := asm.nameTable[labelName]
	if exists {
		return fmt.Errorf("duplicate variable definition: %s", labelName)
	}
	asm.labelTable[labelName] = struct{}{}
	val := uint64(0)
	asm.nameTable[labelName] = &val

	return nil
}

func (asm *assembler) addIntructionStmt(stmt opStmt) {
	intr := intrOp{}
	intr.size = 1 + len(stmt.parameters)

	for _, param := range stmt.parameters {
		iParam := intrParam{}

		if param.variable == "" {
			iParam.value = &param.literal
			continue
		}

		_, varExists := asm.varTable[param.variable]
		_, labelExists := asm.labelTable[param.variable]

		addr := asm.nameTable[param.variable]

		if varExists {
			iParam.varName = param.variable
		}
		if labelExists {
			iParam.labelName = param.variable
		}
		iParam.value = addr
		intr.parameters = append(intr.parameters, iParam)
	}

	asm.stmts = append(asm.stmts, intr)
}

func (asm *assembler) addLabelStmt(stmt labelStmt) {
	intr := intrOp{}

	intr.label = stmt.labelName

	asm.stmts = append(asm.stmts, intr)
}

func (asm *assembler) setNameAddress(name string, addr uint64) {
	ptr := asm.nameTable[name]
	*ptr = addr
}

func assemble(tree asmScript) (*VirtualMachine, error) {
	asm := assembler{
		varTable:   map[string]int{},
		nameTable:  map[string]*uint64{},
		labelTable: map[string]struct{}{},
	}

	vm := &VirtualMachine{}

	for _, stmt := range tree.stmts {
		if stmt.varStmt != nil {
			for _, varName := range stmt.varStmt.varNames {
				asm.defineVar(varName)
			}
		}
		if stmt.labelStmt != nil {
			asm.defineLabel(stmt.labelStmt.labelName)
		}
	}

	for _, stmt := range tree.stmts {
		if stmt.opStmt != nil {
			asm.addIntructionStmt(*stmt.opStmt)
		}
		if stmt.labelStmt != nil {
			asm.addLabelStmt(*stmt.labelStmt)
		}
	}

	bytecodeSize := uint64(0)

	for _, iStmt := range asm.stmts {
		if iStmt.label != "" {
			asm.setNameAddress(iStmt.label, bytecodeSize)
			continue
		}
		bytecodeSize += uint64(iStmt.size)
	}

	const stackSize = 2_000_000
	const gapSize = 100
	stackStart := bytecodeSize + gapSize
	stackEnd := stackStart + stackSize
	heapStart := stackStart + stackSize + gapSize

	for varName, offset := range asm.varTable {
		asm.setNameAddress(varName, heapStart+uint64(offset))
	}

	memory := make([]uint64, heapStart)
	for _, iStmt := range asm.stmts {
		if iStmt.label != "" {
			continue
		}
		memory = append(memory, uint64(iStmt.instruction))
		for _, iParam := range iStmt.parameters {
			memory = append(memory, *iParam.value)
		}
	}
	vm.StackEnd = stackEnd
	vm.SP = stackStart

	return vm, nil
}
