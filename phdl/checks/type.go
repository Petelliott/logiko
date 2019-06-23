package checks

import (
	"github.com/petelliott/logiko/phdl"
	"fmt"
	"math/bits"
)

// TypeCheckFile checks that all types match, and assigns types to unknown
// typed conns
func TypeCheckFile(file *phdl.AstFile) error {
	for _, block := range file.Blocks {
		err := TypeCheckBlock(block)
		if err != nil {
			return err
		}
	}
	return nil
}

func TypeCheckBlock(block *phdl.AstBlock) error {
	for _, stmt := range block.Stmts {
		err := TypeCheckStmt(stmt)
		if err != nil {
			return fmt.Errorf("block '%s': %s", block.Name, err)
		}
	}

	// resolve twice for unresolved conns with indexes
	for _, stmt := range block.Stmts {
		err := TypeCheckStmt(stmt)
		if err != nil {
			return fmt.Errorf("block '%s': %s", block.Name, err)
		}
	}

	for _, v := range block.Vars {
		if !v.HasType() {
			fmt.Errorf(
				"block '%s': type of conn '%s' cannot be determined",
				block.Name, v.Name)
		}
	}

	return nil
}

func TypeCheckStmt(stmt *phdl.AstStmt) error {
	for idx, arg := range stmt.Args {
		err := TypeCheckExpr(stmt.Op.Args[idx].Width, arg)
		if err != nil {
			return err
		}
	}
	return nil

	for idx, ret := range stmt.Rets {
		err := TypeCheckExpr(stmt.Op.Rets[idx].Width, ret)
		if err != nil {
			return err
		}
	}
	return nil
}

func TypeCheckExpr(expected int, expr *phdl.AstExpr) error {
	if expr.HasIndex() {
		width := (expr.Hi-expr.Lo)+1
		if width != expected {
			return fmt.Errorf(
				"expected d%v, got range [%v..%v] (d%v)",
				expected, expr.Lo, expr.Hi, width)
		}
	}

	if expr.Conn != nil {
		if expr.HasIndex() {
			if expr.Conn.HasType() && expr.Hi >= expr.Conn.Width {
				return fmt.Errorf(
					"attempting to get range [%v..%v] (d%v) of '%v' (d%v)",
					expr.Lo, expr.Hi, (expr.Hi-expr.Lo+1),
					expr.Conn.Name, expr.Conn.Width)
			}
		} else {
			if expr.Conn.HasType() {
				if expr.Conn.Width != expected {
					return fmt.Errorf(
						"expected d%v, got '%v' (d%v)",
						expected, expr.Conn.Name, expr.Conn.Width)
				} else {
					// type resolution
					expr.Conn.Width = expected
				}
			}
		}
	} else {
		if bits.Len64(uint64(expr.Literal)) > expected {
			return fmt.Errorf(
				"Literal '%v' does not fit in d%v",
				expr.Literal, expected)
		}
	}
	return nil
}
