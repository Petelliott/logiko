package main

import (
	"github.com/petelliott/logiko/phdl"
	"github.com/petelliott/logiko/phdl/checks"
	"os"
	"fmt"
)

func main() {
	ptree := &phdl.File{}
	err := phdl.Parser.Parse(os.Stdin, ptree)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast, err := phdl.CompileFile(ptree)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = checks.TypeCheckFile(ast)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ast)
}
