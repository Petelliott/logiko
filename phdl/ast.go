package phdl

import (
	"strconv"
	"fmt"
	"errors"
)

type AstFile struct {
	Blocks map[string]*AstBlock
	Tests  map[string]*AstTest
}

func (af AstFile) String() string {
	return fmt.Sprintf(
		"(%v %v)",
		af.Blocks,
		af.Tests,
	)
}

func CompileFile(file *File) (*AstFile, error) {
	astfile := &AstFile{make(map[string]*AstBlock), make(map[string]*AstTest)}
	var err error

	for _, ablock := range file.Blocks {
		if ablock.Block != nil {
			astfile.Blocks[ablock.Block.Ident.Value], err = CompileBlock(astfile,
				ablock.Block)
		} else {
			astfile.Tests[ablock.TestBlock.Ident.Value], err = CompileTestBlock(astfile,
				ablock.TestBlock)
		}

		if err != nil {
			break
		}
	}

	return astfile, err
}

type AstBlock struct {
	Name   string
	Args   []*AstConn
	Rets   []*AstConn
	Vars map[string]*AstConn
	Stmts []*AstStmt
}

func (ab AstBlock) String() string {
	return fmt.Sprintf(
		"(%v %v %v %v %v)",
		ab.Name,
		ab.Args,
		ab.Rets,
		ab.Vars,
		ab.Stmts,
	)
}

func declToConn(d *Declaration) (*AstConn, error) {
	width, err := strconv.Atoi(d.Type[1:])
	conn := &AstConn{d.Ident.Value, width}
	return conn, err
}

func CompileBlock(astfile *AstFile, block *Block) (*AstBlock, error) {
	ablock := &AstBlock{
		Name: block.Ident.Value,
		Args: make([]*AstConn, 0),
		Rets: make([]*AstConn, 0),
		Vars: make(map[string]*AstConn, 0),
		Stmts: make([]*AstStmt, 0),
	}

	for _, arg := range block.Args {
		conn, err := declToConn(arg)
		if err != nil {
			return nil, err
		}
		ablock.Vars[arg.Ident.Value] = conn
		ablock.Args = append(ablock.Args, conn)
	}

	for _, ret := range block.Rets {
		conn, err := declToConn(ret)
		if err != nil {
			return nil, err
		}
		ablock.Vars[ret.Ident.Value] = conn
		ablock.Rets = append(ablock.Rets, conn)
	}

	for _, stmt := range block.Stmts {
		astmt, err := CompileStmt(astfile, ablock, stmt)
		if err != nil {
			return nil, err
		}
		ablock.Stmts = append(ablock.Stmts, astmt)
	}

	return ablock, nil
}

type AstConn struct {
	Name  string
	Width int
}

func (ac AstConn) HasType() bool {
	return ac.Width != 0
}

func (ac AstConn) String() string {
	return fmt.Sprintf(
		"(%v %v)",
		ac.Name,
		ac.Width,
	)
}

type AstStmt struct {
	Args []*AstExpr
	Op   *AstBlock
	Rets []*AstExpr
}

func (as AstStmt) String() string {
	return fmt.Sprintf(
		"(%v %v %v)",
		as.Op.Name,
		as.Args,
		as.Rets,
	)
}

func CompileStmt(astfile *AstFile, block *AstBlock, stmt *Statement) (*AstStmt, error) {
	op, ok := astfile.Blocks[stmt.Ident.Value]
	if !ok {
		return nil, fmt.Errorf("block '%s' not defined", stmt.Ident.Value)
	}

	astmt := &AstStmt{
		Args: make([]*AstExpr, 0),
		Op: op,
		Rets: make([]*AstExpr, 0),
	}

	for _, arg := range stmt.Args {
		expr, err := CompileExpr(block, arg)
		if err != nil {
			return nil, fmt.Errorf("CompileStmt: %s", err)
		}
		astmt.Args = append(astmt.Args, expr)
	}

	for _, ret := range stmt.Rets {
		expr, err := CompileExpr(block, ret)
		if err != nil {
			return nil, fmt.Errorf("CompileStmt: %s", err)
		}
		astmt.Rets = append(astmt.Rets, expr)
	}

	return astmt, nil
}

type AstExpr struct {
	Literal int64
	Conn    *AstConn // nil indicates Literal expr
	Lo      int      // -1 indicates no index
	Hi      int
}

func (ae AstExpr) HasIndex() bool {
	return ae.Lo != -1
}

func (ae AstExpr) String() string {
	if ae.Conn == nil {
		return fmt.Sprintf(
			"(%v %v %v)",
			ae.Literal,
			ae.Lo,
			ae.Hi,
		)
	} else {
		return fmt.Sprintf(
			"(%v %v %v)",
			ae.Conn,
			ae.Lo,
			ae.Hi,
		)
	}
}

func CompileExpr(block *AstBlock, expr *Expr) (*AstExpr, error) {
	var conn *AstConn
	if expr.Ident == nil {
		conn = nil
	} else if c, ok := block.Vars[expr.Ident.Value]; ok {
		conn = c
	} else {
		// 0 represents unknown width
		conn = &AstConn{expr.Ident.Value, 0}
		block.Vars[conn.Name] = conn
	}

	var lit int64
	if conn == nil {
		var err error
		// TODO: custom string literal parsing
		lit, err = strconv.ParseInt(expr.Literal, 0, 64)
		if err != nil {
			return nil, err
		}
	}

	lo := -1
	hi := 0
	if expr.Index != nil {
		var err error
		lo, err = strconv.Atoi(expr.Index.Lo)
		if err != nil {
			return nil, fmt.Errorf("CompileExpr: %s", err)
		}

		hi = lo
		if expr.Index.Hi != "" {
			var err error
			hi, err = strconv.Atoi(expr.Index.Hi)
			if err != nil {
				return nil, fmt.Errorf("CompileExpr: %s", err)
			} else if hi < lo {
				return nil, fmt.Errorf("CompileExpr: low index (%v)" +
									   " is greater than high index (%v)",
									   lo, hi)
			}
		}
	}

	return &AstExpr{
		Literal: lit,
		Conn: conn,
		Lo: lo,
		Hi: hi,
	}, nil

}

type AstTest struct {
	Name  string
	Block *AstBlock
	Stmts []*AstTestStmt
}

func (at AstTest) String() string {
	return fmt.Sprintf(
		"(%v %v %v)",
		at.Name,
		at.Block.Name,
		at.Stmts,
	)
}

func CompileTestBlock(file *AstFile, test *TestBlock) (*AstTest, error) {
	block, ok := file.Blocks[test.Block.Value]
	if !ok {
		return nil, fmt.Errorf("CompileTest: block '%s' is not defined",
			test.Block.Value)
	}
	atest := &AstTest{
		Name: test.Ident.Value,
		Block: block,
		Stmts: make([]*AstTestStmt, 0),
	}

	for _, stmt := range test.Stmts {
		astmt, err := CompileTestStmt(atest, stmt)
		if err != nil {
			return nil, fmt.Errorf("CompileTest: %s", err)
		}
		atest.Stmts = append(atest.Stmts, astmt)
	}

	return atest, nil
}

type AstTestStmt struct {
	Args []*AstExpr
	Rets []*AstExpr
}

func (ats AstTestStmt) String() string {
	return fmt.Sprintf(
		"(%v %v)",
		ats.Args,
		ats.Rets,
	)
}

func CompileTestStmt(test *AstTest, stmt *TestStmt) (*AstTestStmt, error) {
	fakeblock := &AstBlock{Vars: make(map[string]*AstConn)}

	atstmt := &AstTestStmt{
		Args: make([]*AstExpr, 0),
		Rets: make([]*AstExpr, 0),
	}
	for _, arg := range stmt.Args {
		expr, err := CompileExpr(fakeblock, arg)
		if err != nil {
			return nil, fmt.Errorf("CompileTestStmt: %s", err)
		} else if expr.Conn != nil {
			return nil, errors.New(
				"CompileTestStmt: connections are not allowed in tests")
		}
		atstmt.Args = append(atstmt.Args, expr)
	}

	for _, ret := range stmt.Rets {
		expr, err := CompileExpr(fakeblock, ret)
		if err != nil {
			return nil, fmt.Errorf("CompileTestStmt: %s", err)
		} else if expr.Conn != nil {
			return nil, errors.New(
				"CompileTestStmt: connections are not allowed in tests")
		}
		atstmt.Rets = append(atstmt.Rets, expr)
	}

	return atstmt, nil
}
