package main

import "io"

type VmPackage struct {
	Memory []uint64
	Output io.Writer
	SP     uint64
	IP     uint64
}

func (vm *VmPackage) Execute() error {
	panic("not implemented")
}

func LoadBytecodeFile(filePath string) (*VmPackage, error) {
	panic("not implemented")
}

type Bytecode uint64

const (
	Push = Bytecode(iota + 1)
	Pop
	ReadMemory
	WriteMemory
	OutputByte
	Goto
	JumpNotZero
	Call
	Return
)
