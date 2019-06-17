package main

import (
	"github.com/petelliott/logiko/phdl"
	"os"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	ast := &phdl.File{}
	err := phdl.Parser.Parse(os.Stdin, ast)
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(ast)
}
