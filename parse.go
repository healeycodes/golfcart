package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

type ExpressionList struct {
	Pos lexer.Position

	Expression []*Expression `@@*`
}

type Expression struct {
	Pos lexer.Position

	For        *For        `@@`
	Assignment *Assignment `| @@`
	Function   *Function   `| @@`
	Binary     *Binary     `| @@`
}

type For struct {
	Assignment []*Assignment `"for" @@* ";"`
}

type Assignment struct {
	Pos lexer.Position

	Variable   string      `@Ident`
	Op         string      `"="`
	Expression *Expression `@@`
}

type Binary struct {
	Arithmetic *Arithmetic `@@`
	Op         string      `[ @( "!" "=" | "=" "=" | ">" | ">" "=" | "<" | "<" "=" | "or" | "and" )`
	Next       *Binary     `  @@ ]`
}

type Arithmetic struct {
	Unary *Unary      `@@`
	Op    string      `[ @( "-" | "+" | "/" | "*" | "^" | "%" | "&" | "|")`
	Next  *Arithmetic `  @@ ]`
}

type Unary struct {
	Op       string    `  ( @( "!" | "-" )`
	Unary    *Unary    `    @@ )`
	Primary  *Primary  `| @@`
	Function *Function `| @@`
}

type Primary struct {
	Ident         string      `@Ident`
	Number        *float64    `| @Float | @Int`
	String        *string     `| @String`
	Bool          *bool       `| ( @"true" | "false" )`
	Nil           bool        `| @"nil"`
	SubExpression *Expression `| "(" @@ ")" `
}

type Function struct {
	Parameters []*string     `( ( "(" (@Ident ("," @Ident)*) ")" | @Ident ) "=" ">" )`
	Expression []*Expression `( "{" @@* "}" | @@ )`
}

var parser = participle.MustBuild(&ExpressionList{}, participle.UseLookahead(2))

func main() {
	var cli struct {
		ExpressionList []string `arg required help:"ExpressionList to parse."`
	}
	ctx := kong.Parse(&cli)

	exprList := &ExpressionList{}
	err := parser.ParseString("", strings.Join(cli.ExpressionList, " "), exprList)
	ctx.FatalIfErrorf(err)

	repr.Println(exprList)
	// println(parser.String())
}
