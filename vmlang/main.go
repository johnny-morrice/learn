package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/johnny-morrice/learn/vmlang/assembler"
	"github.com/johnny-morrice/learn/vmlang/assembler/parser"
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
	ast, err := parser.ParseFile(*scriptInput)
	if err != nil {
		return err
	}
	vm, err := assembler.Assemble(ast)
	if err != nil {
		return err
	}
	return vm.Execute()
}
