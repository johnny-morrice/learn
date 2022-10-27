package main

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func factorialMemory() []uint64 {
	facM := make([]uint64, 200)

	facM[0] = uint64(Push)
	facM[1] = 4
	facM[2] = uint64(Push)
	facM[3] = 100
	facM[4] = uint64(WriteMemory)
	facM[5] = uint64(Decrement)
	facM[6] = uint64(JumpNotZero)
	facM[7] = 10
	facM[8] = uint64(Goto)
	facM[9] = 21
	facM[10] = uint64(Duplicate)
	facM[11] = uint64(Push)
	facM[12] = 100
	facM[13] = uint64(ReadMemory)
	facM[14] = uint64(Multiply)
	facM[15] = uint64(Push)
	facM[16] = 100
	facM[17] = uint64(WriteMemory)
	facM[18] = uint64(Pop)
	facM[19] = uint64(Goto)
	facM[20] = 5
	facM[21] = uint64(Push)
	facM[22] = 100
	facM[23] = uint64(ReadMemory)
	facM[24] = uint64(OutputByte)
	facM[25] = uint64(Exit)
	return facM
}

func TestVM(t *testing.T) {
	testCases := map[string]struct {
		vm            *VirtualMachine
		expected      []byte
		expectedError error
	}{
		"push": {
			expected: []byte("z"),
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 122, uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       10,
				StackEnd: 100,
			},
		},
		"pop": {
			expected: []byte("j"),
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 106, uint64(Push), 122, uint64(Pop), uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       10,
				StackEnd: 100,
			},
		},
		"increment": {
			expected: []byte{10},
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 9, uint64(Increment), uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       10,
				StackEnd: 100,
			},
		},
		"decrement": {
			expected: []byte{8},
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 9, uint64(Decrement), uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       10,
				StackEnd: 100,
			},
		},
		"duplicate": {
			expected: []byte{9, 9},
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 9, uint64(Duplicate), uint64(OutputByte), uint64(Pop), uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       10,
				StackEnd: 100,
			},
		},
		"multiply": {
			expected: []byte{18},
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 9, uint64(Push), 2, uint64(Multiply), uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       10,
				StackEnd: 100,
			},
		},
		"factorial": {
			expected: []byte{24},
			vm: &VirtualMachine{
				Memory:   factorialMemory(),
				IP:       0,
				SP:       50,
				StackEnd: 100,
			},
		},
		"dynamic memory": {
			expected: []byte{64},
			vm: &VirtualMachine{
				Memory:   []uint64{uint64(Push), 64, uint64(Push), 2048, uint64(WriteMemory), uint64(Pop), uint64(Push), 2048, uint64(ReadMemory), uint64(OutputByte), uint64(Exit), 0, 0, 0, 0, 0, 0, 0, 0},
				IP:       0,
				SP:       12,
				StackEnd: 100,
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			output := &bytes.Buffer{}
			testCase.vm.Output = output
			err := testCase.vm.Execute()
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected error: %s but received: %s", testCase.expectedError, err)
				return
			}
			actual := output.Bytes()
			if !reflect.DeepEqual(testCase.expected, actual) {
				t.Fatalf("expected output: %v but received: %v", testCase.expected, actual)
				return
			}
		})
	}
}

func TestBytecodeRepresentation(t *testing.T) {
	testCases := map[uint64]Bytecode{
		1: Push,
		2: Pop,
	}
	for rep, bc := range testCases {
		if rep != uint64(bc) {
			t.Errorf("expected representation %v to match bytecode %v", rep, bc)
		}
	}
}
