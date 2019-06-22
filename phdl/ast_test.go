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

func TestCompileStmt(t *testing.T) {
	Parser := participle.MustBuild(
		&Statement{},
		participle.Lexer(Lexer),
		participle.Elide("Whitespace", "OneLineComment", "MultiLineComment"),
	)

	prog1 := "(a)f -> a;"
	prog2 := "(4)f -> a;"

	ptree := &Statement{}
	err := Parser.ParseString(prog1, ptree)
	if err != nil {
		t.Fatal(err)
	}

	astfile := &AstFile{Blocks: map[string]*AstBlock{
		"f": &AstBlock{Name: "f"},
		"tb": &AstBlock{Name: "tb", Vars: make(map[string]*AstConn)},
	}}
	ast, err := CompileStmt(astfile, astfile.Blocks["tb"], ptree)
	if err != nil {
		t.Fatal(err)
	}

	expected := "(f [((a 0) 0 0)] [((a 0) 0 0)])"
	if ast.String() != expected {
		t.Errorf("expected/got:\n%s\n%s\n", expected, ast.String())
	}

	if ast.Op != astfile.Blocks["f"] {
		t.Error("stmt block is not file block")
	}

	if ast.Args[0].Conn != ast.Rets[0].Conn || ast.Rets[0].Conn != astfile.Blocks["tb"].Vars["a"] {
		t.Error("different sybmols for 'a'")
	}

	err = Parser.ParseString(prog2, ptree)
	if err != nil {
		t.Fatal(err)
	}

	ast, err = CompileStmt(astfile, astfile.Blocks["tb"], ptree)
	if err != nil {
		t.Fatal(err)
	}

	expected = "(f [(4 0 0)] [((a 0) 0 0)])"
	if ast.String() != expected {
		t.Errorf("expected/got:\n%s\n%s\n", expected, ast.String())
	}

	if ast.Rets[0].Conn != astfile.Blocks["tb"].Vars["a"] {
		t.Error("different sybmols for 'a'")
	}
}
