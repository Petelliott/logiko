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

func TestCompileExpr(t *testing.T) {
	Parser := participle.MustBuild(
		&Expr{},
		participle.Lexer(Lexer),
		participle.Elide("Whitespace", "OneLineComment", "MultiLineComment"),
	)
	testblock := &AstBlock{Name: "tb", Vars: map[string]*AstConn{
		"a": &AstConn{"a", 32},
		"b": &AstConn{"b", 32},
	}}

	Comp := func(prog string) (*AstExpr, error) {
		t.Helper()

		ptree := &Expr{}
		err := Parser.ParseString(prog, ptree)
		if err != nil {
			t.Fatal(err)
		}

		return CompileExpr(testblock, ptree)
	}

	ast, err := Comp("a")
	if err != nil {
		t.Error(err)
	}

	if ast.Conn != testblock.Vars["a"] {
		t.Error("conn is not 'a'")
	}

	ast, err = Comp("a[7]")
	if err != nil {
		t.Error(err)
	}

	if ast.Lo != 7 || ast.Hi != 7 {
		t.Errorf("incorrect index bounds (l=%v, h=%v)", ast.Lo, ast.Hi)
	}

	ast, err = Comp("a[7..6]")
	if err == nil {
		t.Error("hi index to be higher than lo index")
	}

	Lcheck := func(lit string, exp int64) {
		t.Helper()

		ast, err = Comp(lit)
		if err != nil {
			t.Error(err)
		}

		if ast.Conn != nil {
			t.Error("non-nil conn on literal")
		}

		if ast.Literal != exp {
			t.Errorf("expected: %v, got: %v", exp, ast.Literal)
		}
	}

	Lfail := func(lit string) {
		t.Helper()

		_, err = Comp(lit)
		if err == nil {
			t.Error("expected to fail, didn't")
		}
	}

	Lcheck("32", 32)
	Lcheck("-32", -32)
	Lcheck("0", 0)
	Lcheck("-0", 0)
	Lcheck("0x8f", 0x8f)
	Lcheck("-0x8f", -0x8f)
	Lcheck("0x11", 0x11)

	Lfail("5a")

	t.Skip("skipping rest of test: see bug #3")

	Lcheck("0o11", 011)
	Lcheck("-0o11", -011)
	Lfail("0o8")

	Lcheck("0b1010", 10)
	Lcheck("-0b1010", -10)
	Lfail("0b2")

}

func TestCompileTestBlock(t *testing.T) {
	Parser := participle.MustBuild(
		&TestBlock{},
		participle.Lexer(Lexer),
		participle.Elide("Whitespace", "OneLineComment", "MultiLineComment"),
	)

	astfile := &AstFile{Blocks: map[string]*AstBlock{
		"f": &AstBlock{Name: "f"},
	}}

	Comp := func(prog string) (*AstTest, error) {
		t.Helper()

		ptree := &TestBlock{}
		err := Parser.ParseString(prog, ptree)
		if err != nil {
			t.Fatal(err)
		}

		return CompileTestBlock(astfile, ptree)
	}

	ast, err := Comp(`
		test test1(f) {
			1, 2, 3 ==> 1, 2;
			3, 2, 1 ==> 2, 1;
		}
	`)

	if err != nil {
		t.Error(err)
	}

	expected := "(test1 f [([(1 0 0) (2 0 0) (3 0 0)] [(1 0 0) (2 0 0)]) ([(3 0 0) (2 0 0) (1 0 0)] [(2 0 0) (1 0 0)])])"
	if ast.String() != expected {
		t.Errorf("expected/got:\n%s\n%s\n", expected, ast.String())
	}

	if ast.Block != astfile.Blocks["f"] {
		t.Error("expected test's block to be same as global block")
	}

	_, err = Comp("test test2(g) {}")
	if err == nil {
		t.Error("expected block not found error")
	}

	_, err = Comp("test test3(f) { a ==> 8; }")
	if err == nil {
		t.Error("conn not found error")
	}
}

func TestCompileTestStmt(t *testing.T) {
	Parser := participle.MustBuild(
		&TestStmt{},
		participle.Lexer(Lexer),
		participle.Elide("Whitespace", "OneLineComment", "MultiLineComment"),
	)

	asttest := &AstTest{}

	Comp := func(prog string) (*AstTestStmt, error) {
		t.Helper()

		ptree := &TestStmt{}
		err := Parser.ParseString(prog, ptree)
		if err != nil {
			t.Fatal(err)
		}

		return CompileTestStmt(asttest, ptree)
	}

	ast, err := Comp("1,2 ==> 4;")
	if err != nil {
		t.Error(err)
	}

	expected := "([(1 0 0) (2 0 0)] [(4 0 0)])"
	if ast.String() != expected {
		t.Errorf("expected/got:\n%s\n%s\n", expected, ast.String())
	}


	_, err = Comp("1,2 ==> a;")
	if err == nil {
		t.Error("expected that conns are not allowed in tests")
	}

	_, err = Comp("1,a ==> 4;")
	if err == nil {
		t.Error("expected that conns are not allowed in tests")
	}
}
