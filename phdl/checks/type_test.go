package checks

import (
	"testing"
	"github.com/petelliott/logiko/phdl"
)

func TestTypeCheckFile(t *testing.T) {
	err := TypeCheckFile(&phdl.AstFile{})
	if err != nil {
		t.Error(err)
	}

	err = TypeCheckFile(&phdl.AstFile{Blocks: map[string]*phdl.AstBlock{
		"blk": &phdl.AstBlock{Vars: map[string]*phdl.AstConn{
			"a": &phdl.AstConn{Name: "a"},
		}},
	}})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTypeCheckBlock(t *testing.T) {
	err := TypeCheckBlock(&phdl.AstBlock{})
	if err != nil {
		t.Error(err)
	}

	Comp := func(t *testing.T, prog string) *phdl.AstFile {
		t.Helper()

		ptree := &phdl.File{}
		err := phdl.Parser.ParseString(prog, ptree)
		if err != nil {
			t.Fatal(err)
		}

		ast, err := phdl.CompileFile(ptree)
		if err != nil {
			t.Fatal(err)
		}

		return ast
	}

	ast := Comp(t, `
		block a (a d3) -> (b d3) {}

		block b (a d3) -> (b d3) {
			(a)a -> b;
		}
	`)

	err = TypeCheckBlock(ast.Blocks["b"])
	if err != nil {
		t.Error(err)
	}

	ast = Comp(t, `
		block a (a d3) -> (b d3) {}

		block b (a d4) -> (b d3) {
			(a)a -> b;
		}
	`)

	err = TypeCheckBlock(ast.Blocks["b"])
	if err == nil {
		t.Error("expected error")
	}

	ast = Comp(t, `
		block a (a d5) -> (b d3) {}
		block c (a d2) -> (b d3) {}

		block b (a d3) -> (b d3) {
			(c[0..4])a -> b;
			(c)c -> b;
		}
	`)

	err = TypeCheckBlock(ast.Blocks["b"])
	if err == nil {
		t.Error("expected error")
	}

	ast = Comp(t, `
		block a (a d5) -> (b d3) {}
		block c (a d2) -> (b d3) {}

		block b (a d3) -> (b d3) {
			(1)a -> c[0..4];
			(1)c -> c[0..2];
		}
	`)

	err = TypeCheckBlock(ast.Blocks["b"])
	if err == nil {
		t.Error("expected error")
	}
}

func TestTypeCheckExpr(t *testing.T) {
	err := TypeCheckExpr(1, &phdl.AstExpr{Literal: 0})
	if err != nil {
		t.Error(err)
	}

	err = TypeCheckExpr(1, &phdl.AstExpr{Hi: 2, Lo: 0})
	if err == nil {
		t.Error("expected error")
	}

	err = TypeCheckExpr(5, &phdl.AstExpr{
		Conn: &phdl.AstConn{Name: "a", Width: 3},
		Lo: 0,
		Hi: 4,
	})
	if err == nil {
		t.Error("expected error")
	}

	err = TypeCheckExpr(5, &phdl.AstExpr{
		Conn: &phdl.AstConn{Name: "a", Width: 3},
	})
	if err == nil {
		t.Error("expected error")
	}

	err = TypeCheckExpr(5, &phdl.AstExpr{
		Conn: &phdl.AstConn{Name: "a", Width: 0},
		Lo: -1,
	})
	if err != nil {
		t.Error(err)
	}

	err = TypeCheckExpr(1, &phdl.AstExpr{Literal: 5})
	if err == nil {
		t.Error("expected error")
	}
}
