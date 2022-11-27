package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/johnny-morrice/learn/vmlang/asm"
	"github.com/johnny-morrice/learn/vmlang/asm/parser"
)

var asmInput = flag.String("run-asm", "", "run asm file")
var byteToDec = flag.Bool("byte2dec", false, "make output bytes human readable")

func main() {
	flag.Parse()
	if *asmInput != "" {
		err := runAsm()
		if err != nil {
			fmt.Printf("error running asm: %s", err)
			os.Exit(1)
		}
	} else {
		flag.Usage()
	}
}

func runAsm() error {
	ast, err := parser.ParseFile(*asmInput)
	if err != nil {
		return err
	}
	vm, err := asm.Assemble(ast)
	if err != nil {
		return err
	}
	if *byteToDec {
		vm.Output = byte2dec{out: vm.Output}
	}
	return vm.Execute()
}

type byte2dec struct {
	out io.Writer
}

func (w byte2dec) Write(bs []byte) (int, error) {
	written := 0
	for _, b := range bs {
		_, err := fmt.Fprintf(w.out, "%d", b)
		if err == nil {
			written++
		} else {
			return written, err
		}
		fmt.Println()
	}
	return written, nil
}
