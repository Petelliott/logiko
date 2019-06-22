package main

import (
	"github.com/petelliott/logiko/phdl"
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
	fmt.Println(ast)
}
