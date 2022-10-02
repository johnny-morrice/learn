package main

import (
	"flag"
	"fmt"
	"os"
)

var scriptInput = flag.String("input", "", "input script")

func main() {
	flag.Parse()
	if *scriptInput != "" {
		err := runScript()
		if err != nil {
			fmt.Printf("error running script: %s", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("expected `input`")
	}
}

func runScript() error {
	vm, err := ParseFile(*scriptInput)
	if err != nil {
		return err
	}
	return vm.Execute()
}
