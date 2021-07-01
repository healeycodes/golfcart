package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
)

type Program struct {
	Expression []*Expression `@@*`
}

// type Declaration struct {
// 	Pos lexer.Position

// 	Type string `@Type`
// 	Name string `@Ident`
// }

// type FunctionBody struct {
// 	Pos lexer.Position

// 	Expression []*Expression `@@*`
// }

type Expression struct {
	Equality *Equality `@@`
}

type Equality struct {
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

type Comparison struct {
	Addition *Addition   `@@`
	Op       string      `[ @( ">" | ">" "=" | "<" | "<" "=" )`
	Next     *Comparison `  @@ ]`
}

type Addition struct {
	Multiplication *Factor   `@@`
	Op             string    `[ @( "-" | "+" )`
	Next           *Addition `  @@ ]`
}

type Factor struct {
	Unary *Unary  `@@`
	Op    string  `[ @( "/" | "*" | "^" | "%" | "&" | "|")`
	Next  *Factor `  @@ ]`
}

type Unary struct {
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	Number        *float64    `  @Float | @Int`
	String        *string     `| @String`
	Bool          *bool       `| ( @"true" | "false" )`
	Nil           bool        `| @"nil"`
	SubExpression *Expression `| "(" @@ ")" `
}

var parser = participle.MustBuild(&Program{}, participle.UseLookahead(2))

func main() {
	var cli struct {
		Program []string `arg required help:"Expression to parse."`
	}
	ctx := kong.Parse(&cli)

	program := &Program{}
	err := parser.ParseString("", strings.Join(cli.Program, " "), program)
	ctx.FatalIfErrorf(err)

	repr.Println(program)
}
