package main

import (
	"fmt"
	"io"
)

type VmPackage struct {
	Memory []uint64
	Output io.Writer
	SP     uint64
	IP     uint64
}

func (vm *VmPackage) Execute() error {
	for {
		op := Bytecode(vm.Memory[vm.IP])
		switch op {
		case Push:
			vm.SP++
			x := vm.Memory[vm.IP+1]
			vm.Memory[vm.SP] = x
			vm.IP += 2
		case Pop:
			vm.Memory[vm.SP] = 0
			vm.SP--
			vm.IP++
		case Increment:
			vm.Memory[vm.SP]++
			vm.IP++
		case Decrement:
			vm.Memory[vm.SP]--
			vm.IP++
		case Duplicate:
			x := vm.Memory[vm.SP]
			vm.SP++
			vm.Memory[vm.SP] = x
			vm.IP++
		case ReadMemory:
		case WriteMemory:
		case OutputByte:
			x := vm.Memory[vm.SP]
			bs := []byte{byte(x)}
			_, err := vm.Output.Write(bs)
			if err != nil {
				return err
			}
			vm.IP++
		case Goto:
			x := vm.Memory[vm.IP+1]
			vm.IP = x
		case JumpNotZero:
			x := vm.Memory[vm.SP]
			if x == 0 {
				vm.IP++
				continue
			}
			y := vm.Memory[vm.IP+1]
			vm.IP = y
		case Call:
		case Return:
		case Exit:
			return nil
		default:
			return fmt.Errorf("Unknown bytecode: %v", op)
		}
	}
}

func LoadBytecodeFile(filePath string) (*VmPackage, error) {
	panic("not implemented")
}

type Bytecode uint64

const (
	Push = Bytecode(iota + 1)
	Pop
	Increment
	Decrement
	Duplicate
	ReadMemory
	WriteMemory
	OutputByte
	Goto
	JumpNotZero
	Call
	Return
	Exit
)
