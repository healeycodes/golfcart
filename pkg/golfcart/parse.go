package golfcart

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/v2"
)

type ExpressionList struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Expressions []*Expression `@@*`
}

type Expression struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Break           *Break           `  @@`
	Continue        *Continue        `| @@`
	For             *For             `| @@`
	While           *While           `| @@`
	If              *If              `| @@`
	Assignment      *Assignment      `| @@`
	LogicAnd        *LogicAnd        `| @@`
	FunctionLiteral *FunctionLiteral `| @@`
	ListLiteral     *[]Expression    `| "[" ( @@ ("," @@)* ","? )? "]"`
	ObjectLiteral   *[]ObjectEntry   `| "{" ( @@ ("," @@)* ","? )? "}"`
}

type Break struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Break string `"break"`
}

type Continue struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Continue string `"continue"`
}

type For struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Init      []*Assignment `"for" ( @@* ("," @@ )* )* ";"`
	Condition *Expression   `@@ ";"`
	Post      *Expression   `@@`
	Body      []*Expression `"{" @@* "}"`
}

type While struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Condition *Expression   `"while" @@`
	Body      []*Expression `"{" @@* "}"`
}

type If struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Init      []*Assignment `"if" ( @@ ("," @@ )* ";" )?`
	Condition *Expression   `@@`
	Body      []*Expression `"{" @@* "}"`
	ElseBody  []*Expression `( "else" "{" @@* "}" )?`
}

type Assignment struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Ident      *string     `@Ident`
	Op         string      `"="`
	Expression *Expression `@@`
}

type LogicAnd struct {
	Pos    lexer.Position
	EndPos lexer.Position

	LogicOr *LogicOr  `@@`
	Op      string    `[ @( "and" )`
	Next    *LogicAnd `  @@ ]`
}

type LogicOr struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Equality *Equality `@@`
	Op       string    `[ @( "or" )`
	Next     *LogicOr  `  @@ ]`
}

type Equality struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

type Comparison struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Addition *Addition   `@@`
	Op       string      `[ @( ">" "=" | ">" | "<" "=" | "<" )`
	Next     *Comparison `  @@ ]`
}

type Addition struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

type Multiplication struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Logical *Unary          `@@`
	Op      string          `[ @( "/" | "*" | "%")`
	Next    *Multiplication `  @@ ]`
}

type Unary struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Op      string   `[ ( @( "!" | "-" )`
	Unary   *Unary   `  @@ ) ]`
	Primary *Primary `| @@`
	Call    *Call    `| @@`
}

type Primary struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Ident         *string     `@Ident`
	Number        *float64    `| @Float | @Int`
	Str           *string     `| @String`
	Bool          *bool       `| ( @"true" | "false" )`
	Nil           *bool       `| @"nil"`
	SubExpression *Expression `| "(" @@ ")" `
}

type Call struct {
	Expression *Expression   `@@`
	Parameters *[]Expression `[ ( "(" ( @@ ("," @@)* )? ")"`
	Ident      *string       `	| ( "." @Ident `
	Brackets   *Expression   `    | "[" @@ "]" ) ]`
}

type ObjectEntry struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Key   *Expression `@@ ":"`
	Value *Expression `@@`
}

type FunctionLiteral struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Parameters []*string     `( ( "(" ")" | "(" (@Ident ("," @Ident)*) ")" | @Ident ) "=" ">" )`
	Body       []*Expression `( "{" @@* "}" | @@ )`
}

func GenerateAST(source string) (*ExpressionList, error) {
	var parser = (participle.MustBuild(&ExpressionList{}, participle.UseLookahead(2)))
	expressionList := &ExpressionList{}

	err := parser.ParseString("", source, expressionList)
	if err != nil {
		return nil, err
	}
	return expressionList, nil

	// Print grammar
	// println(parser.String())
}
