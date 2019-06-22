package phdl

type File struct {
	Blocks []*AnyBlock `@@*`
}

type AnyBlock struct {
	Block *Block         `  @@`
	TestBlock *TestBlock `| @@`
}

type Ident struct {
	Value string ` @Ident1 | @Ident2 | @Ident3 `
}

type Block struct {
	Ident *Ident         ` Block @@ `
	Args  []*Declaration ` Lparen (@@ (Comma @@)* )? Rparen `
	Rets  []*Declaration ` (Arrow Lparen (@@ (Comma @@)* )? Rparen)? `
	Stmts []*Statement   ` Lbrace @@* Rbrace `
}

type Declaration struct {
	Ident *Ident ` @@ `
	Type  string ` @Type `
}

type Statement struct {
	Args  []*Expr ` Lparen ( @@ (Comma @@)* )? Rparen `
	Ident *Ident  ` ( @@ Arrow `
	Rets  []*Expr ` @@ (Comma @@)* )? Semicolon `
}

type Expr struct {
	Literal string `@Number `
	Ident   *Ident `| ( @@ `
	Index   *Index `( Lbrak @@ Rbrak )? )? `
}

type Index struct {
	Lo string ` @Number `
	Hi string ` ( Ellipsis @Number )? `
}

type TestBlock struct {
	Ident *Ident      ` Test @@ `
	Block *Ident      ` Lparen @@ Rparen `
	Stmts []*TestStmt ` Lbrace @@* Rbrace `
}

type TestStmt struct {
	Args []*Expr ` ( @@ (Comma @@)* )? TestArrow `
	Rets []*Expr ` ( @@ (Comma @@)* )? Semicolon`
}
