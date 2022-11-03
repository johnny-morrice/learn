package assembler

import (
	"errors"
	"fmt"

	"github.com/johnny-morrice/learn/vmlang/vm"
)

var ErrAssembler = errors.New("assembly error")

type intrParam struct {
	value     *uint64
	varName   string
	labelName string
}

func (param intrParam) getParamName() string {
	if param.varName != "" {
		return param.varName
	}
	return param.labelName
}

func (param intrParam) missingValueError() error {
	name := param.getParamName()
	if name != "" {
		return fmt.Errorf("variable not defined: %v; %w", name, ErrAssembler)
	}
	return fmt.Errorf("missing variable definition: %w", ErrAssembler)
}

func (param intrParam) valueString() string {
	if param.value == nil {
		return "nil"
	}
	return fmt.Sprint(*param.value)
}

func (param intrParam) String() string {
	if param.varName != "" {
		return fmt.Sprintf("[varName: %v, value: %v]", param.varName, param.valueString())
	}
	if param.labelName != "" {
		return fmt.Sprintf("[labelName: %v, value: %v]", param.labelName, param.valueString())
	}
	return fmt.Sprintf("[value: %v]", param.valueString())
}

type intrOp struct {
	op         vm.Bytecode
	size       int
	parameters []intrParam
	label      string
}

func (op intrOp) String() string {
	if op.label != "" {
		return fmt.Sprintf("[label: %v]", op.label)
	}
	return fmt.Sprintf("[op: %v, size: %v, params: %v]", op.op, op.size, op.parameters)
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
		return fmt.Errorf("duplicate variable definition: %s; %w", varName, ErrAssembler)
	}
	asm.varTable[varName] = len(asm.varTable)
	val := uint64(0)
	asm.nameTable[varName] = &val

	return nil
}

func (asm *assembler) defineLabel(labelName string) error {
	_, exists := asm.nameTable[labelName]
	if exists {
		return fmt.Errorf("duplicate variable definition: %s; %w", labelName, ErrAssembler)
	}
	asm.labelTable[labelName] = struct{}{}
	val := uint64(0)
	asm.nameTable[labelName] = &val

	return nil
}

func (asm *assembler) addOpStmt(stmt opStmt) {
	iOp := intrOp{}
	iOp.size = 1 + len(stmt.parameters)
	iOp.op = stmt.op

	// fmt.Printf("add op stmt: %v\n", stmt)

	for _, param := range stmt.parameters {
		iParam := intrParam{}

		if param.variable == "" {
			iParam.value = &param.literal
			iOp.parameters = append(iOp.parameters, iParam)
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
		iOp.parameters = append(iOp.parameters, iParam)
	}

	asm.stmts = append(asm.stmts, iOp)
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

const stackSize = 2_000_000
const gapSize = 100

func Assemble(tree *AsmScript) (*vm.VirtualMachine, error) {
	asm := assembler{
		varTable:   map[string]int{},
		nameTable:  map[string]*uint64{},
		labelTable: map[string]struct{}{},
	}

	machine := &vm.VirtualMachine{}

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
			asm.addOpStmt(*stmt.opStmt)
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

	bytecodeSize++

	stackStart := bytecodeSize + gapSize
	stackEnd := stackStart + stackSize
	heapStart := stackStart + stackSize + gapSize

	for varName, offset := range asm.varTable {
		asm.setNameAddress(varName, heapStart+uint64(offset))
	}

	machine.Memory = make([]uint64, heapStart)
	index := 0
	for _, iStmt := range asm.stmts {
		if iStmt.label != "" {
			continue
		}
		machine.Memory[index] = uint64(iStmt.op)
		index++
		for _, iParam := range iStmt.parameters {
			if iParam.value == nil {
				return nil, iParam.missingValueError()
			}
			machine.Memory[index] = *iParam.value
			index++
		}
	}
	machine.Memory[index] = uint64(vm.Exit)
	machine.StackEnd = stackEnd
	machine.SP = stackStart
	machine.HeapStart = heapStart

	return machine, nil
}
