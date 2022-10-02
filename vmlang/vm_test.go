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
		expected      string
		expectedError error
	}{}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			output := &bytes.Buffer{}
			testCase.vm.Output = output
			err := testCase.vm.Execute()
			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("expected error: %s but received: %s", testCase.expectedError, err)
				return
			}
			actual := output.String()
			if !reflect.DeepEqual(testCase.expected, actual) {
				t.Fatalf("expected output: %s but received: %s", testCase.expected, actual)
				return
			}
		})
	}
}
