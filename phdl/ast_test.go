package phdl

import (
	"github.com/alecthomas/participle"
	"testing"
)


func TestCompileFile(t *testing.T) {
	prog := `
		block f (a d1) -> (b d2) {}
		test ftest(f) {}
	`
	ptree := &File{}
	err := Parser.ParseString(prog, ptree)
	if err != nil {
		t.Fatal(err)
	}

	ast, err := CompileFile(ptree)
	if err != nil {
		t.Fatal(err)
	}

	expected := "(map[f:(f [(a 1)] [(b 2)] map[a:(a 1) b:(b 2)] [])] map[ftest:(ftest f [])])"
	if ast.String() != expected {
		t.Errorf("expected/got:\n%s\n%s\n", expected, ast.String())
	}

	if ast.Blocks["f"] != ast.Tests["ftest"].Block {
		t.Error("test does not point to block")
	}
}

func TestCompileBlock(t *testing.T) {
	Parser := participle.MustBuild(
		&Block{},
		participle.Lexer(Lexer),
		participle.Elide("Whitespace", "OneLineComment", "MultiLineComment"),
	)

	prog := `
		block f (a d1) -> (b d2) {
			(a)f2 -> b;
		}
	`
	ptree := &Block{}
	err := Parser.ParseString(prog, ptree)
	if err != nil {
		t.Fatal(err)
	}

	astfile := &AstFile{Blocks: map[string]*AstBlock{
		"f2": &AstBlock{Name: "f2"},
	}}
	ast, err := CompileBlock(astfile, ptree)
	if err != nil {
		t.Fatal(err)
	}

	expected := "(f [(a 1)] [(b 2)] map[a:(a 1) b:(b 2)] [(f2 [((a 1) 0 0)] [((b 2) 0 0)])])"
	if ast.String() != expected {
		t.Errorf("expected/got:\n%s\n%s\n", expected, ast.String())
	}

	if ast.Args[0] != ast.Vars["a"] {
		t.Error("block arg0 does not match var")
	}

	if ast.Rets[0] != ast.Stmts[0].Rets[0].Conn {
		t.Error("stmt ret0 does not match block ret0")
	}
}
