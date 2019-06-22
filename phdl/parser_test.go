package phdl

import (
	"testing"
	"strings"
	_ "github.com/alecthomas/participle/lexer"
)

type testToken struct {
	Type string
	Value string
}

func lexerExpect(t *testing.T, input string, expected []testToken) {
	t.Helper()
	l, err := Lexer.Lex(strings.NewReader(input))
	if err != nil {
		t.Error("Lexer.Lex():", err)
		return
	}

	m := Lexer.Symbols()

	for i, exp := range expected {
		tok, err := l.Next()
		if err != nil {
			t.Errorf("token %d: l.Next(): %v", i, err)
			return
		}
		if tok.Type != m[exp.Type] || tok.Value != exp.Value {
			t.Errorf("token %d: expected: '%s', got: '%s'", i, exp.Value, tok.Value)
			return
		}
	}

	n, err := l.Next()
	if err != nil {
		t.Error("l.Next():", err)
		return
	}
	if !n.EOF() {
		t.Errorf("expected <EOF>, got: '%s'", n.Value)
	}
}

func TestLexerKeyword(t *testing.T) {
	lexerExpect(
		t,
		"ftest test testf testf",
		[]testToken{
			{"Ident3", "ftest"},
			{"Whitespace", " "},
			{"Test", "test"},
			{"Whitespace", " "},
			{"Ident1", "testf"},
			{"Whitespace", " "},
			{"Ident1", "testf"},
		},
	)

	lexerExpect(
		t,
		"fblock block blockf blockf",
		[]testToken{
			{"Ident3", "fblock"},
			{"Whitespace", " "},
			{"Block", "block"},
			{"Whitespace", " "},
			{"Ident1", "blockf"},
			{"Whitespace", " "},
			{"Ident1", "blockf"},
		},
	)
}

func TestLexerLiteral(t *testing.T) {
	lexerExpect(
		t,
		"0 0o0 0x0 0b0 -55 -0o55 -0xa5 -0b1010",
		[]testToken{
			{"Number", "0"}, {"Whitespace", " "},
			{"Number", "0o0"}, {"Whitespace", " "},
			{"Number", "0x0"}, {"Whitespace", " "},
			{"Number", "0b0"}, {"Whitespace", " "},
			{"Number", "-55"}, {"Whitespace", " "},
			{"Number", "-0o55"}, {"Whitespace", " "},
			{"Number", "-0xa5"}, {"Whitespace", " "},
			{"Number", "-0b1010"},
		},
	)
}

func TestLexerComment(t *testing.T) {
	lexerExpect(
		t,
		"// blah // b\n/*h\nb a*/\nhello /* */",
		[]testToken{
			{"OneLineComment", "// blah // b"},
			{"Whitespace", "\n"},
			{"MultiLineComment", "/*h\nb a*/"},
			{"Whitespace", "\n"},
			{"Ident3", "hello"},
			{"Whitespace", " "},
			{"MultiLineComment", "/* */"},
		},
	)
}

func TestLexerAsterixComment(t *testing.T) {
	t.Skip("see bug #1")
	lexerExpect(
		t,
		"/* * */",
		[]testToken{
			{"MultiLineComment", "/* * */"},
		},
	)
}

func TestLexerSeperators(t *testing.T) {
	lexerExpect(
		t,
		"(){}[]->==>;..",
		[]testToken{
			{"Lparen", "("},
			{"Rparen", ")"},
			{"Lbrace", "{"},
			{"Rbrace", "}"},
			{"Lbrak", "["},
			{"Rbrak", "]"},
			{"Arrow", "->"},
			{"TestArrow", "==>"},
			{"Semicolon", ";"},
			{"Ellipsis", ".."},
		},
	)
}

func TestLexerType(t *testing.T) {
	lexerExpect(
		t,
		"d0 d32 dd5 d54d",
		[]testToken{
			{"Type", "d0"}, {"Whitespace", " "},
			{"Type", "d32"}, {"Whitespace", " "},
			{"Ident2", "dd5"}, {"Whitespace", " "},
			{"Ident2", "d54d"},
		},
	)
}
