package main

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestVM(t *testing.T) {
	testCases := map[string]struct {
		vm            *VmPackage
		expected      []byte
		expectedError error
	}{
		"push": {
			expected: []byte("z"),
			vm: &VmPackage{
				Memory: []uint64{1, 122, 8, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				IP:     0,
				SP:     10,
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
