package vm

import (
	"errors"
	"fmt"
	"io"
)

type VirtualMachine struct {
	Memory    []uint64
	Output    io.Writer
	SP        uint64
	StackEnd  uint64
	HeapStart uint64
	IP        uint64
}

func (vm *VirtualMachine) Execute() error {
	const debug = false
	for {
		op := Bytecode(vm.Memory[vm.IP])
		if debug {
			fmt.Printf("vm debug; op: %v; sp: %v; ip: %v\n", op, vm.SP, vm.IP)
		}
		switch op {
		case Push:
			err := vm.incrementSP()
			if err != nil {
				return err
			}
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
			err := vm.incrementSP()
			if err != nil {
				return err
			}
			vm.Memory[vm.SP] = x
			vm.IP++
		case ReadMemory:
			i := vm.Memory[vm.SP]
			vm.growMemory(i)
			x := vm.Memory[i]
			vm.Memory[vm.SP] = x
			vm.IP++
		case WriteMemory:
			i := vm.Memory[vm.SP]
			vm.growMemory(i)
			x := vm.Memory[vm.SP-1]
			vm.Memory[i] = x
			vm.Memory[vm.SP] = 0
			vm.SP--
			vm.IP++
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
		case Multiply:
			first, second := vm.Memory[vm.SP], vm.Memory[vm.SP-1]
			x := first * second
			vm.Memory[vm.SP] = 0
			vm.SP--
			vm.Memory[vm.SP] = x
			vm.IP++
		default:
			return fmt.Errorf("unknown bytecode: %v", op)
		}
	}
}

func (vm *VirtualMachine) incrementSP() error {
	vm.SP++
	if vm.SP >= vm.StackEnd {
		return errors.New("stack overflow")
	}

	return nil
}

func (vm *VirtualMachine) growMemory(i uint64) {
	memSize := uint64(len(vm.Memory))
	if memSize-1 < i {
		expand := (i - memSize) + 1
		if expand < memSize {
			expand = memSize
		}
		extra := make([]uint64, expand)
		vm.Memory = append(vm.Memory, extra...)
	}
}

func LoadBytecodeFile(filePath string) (*VirtualMachine, error) {
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
	Multiply
	// Make sure you update the Bytecodes array below.
)

func Bytecodes() []Bytecode {
	const max = Multiply
	bc := []Bytecode{}
	for i := Push; i <= max; i++ {
		bc = append(bc, i)
	}
	return bc
}

func (code Bytecode) String() string {
	switch code {
	case Push:
		return "push"
	case Pop:
		return "pop"
	case Increment:
		return "incr"
	case Decrement:
		return "decr"
	case Duplicate:
		return "dupl"
	case ReadMemory:
		return "rmem"
	case WriteMemory:
		return "wmem"
	case OutputByte:
		return "outb"
	case Goto:
		return "goto"
	case JumpNotZero:
		return "jnz"
	case Call:
		return "call"
	case Return:
		return "rtn"
	case Exit:
		return "exit"
	case Multiply:
		return "mult"
	default:
		return fmt.Sprint(uint64(code))
	}
}
