package main

import "io"

type VmPackage struct {
	Bytecode []Bytecode
	Memory   []uint64
	Output   io.Writer
	SP       uint64
	IP       uint64
}

func (vm *VmPackage) Execute() error {
	panic("not implemented")
}

func LoadBytecodeFile(filePath string) (*VmPackage, error) {
	panic("not implemented")
}

type Bytecode uint64
