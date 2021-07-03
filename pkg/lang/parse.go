package lang

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/v2"
)

type ExpressionList struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Expression []*Expression `@@*`
}

type Expression struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Break      *Break      `  @@`
	Continue   *Continue   `| @@`
	For        *For        `| @@`
	While      *While      `| @@`
	If         *If         `| @@`
	Assignment *Assignment `| @@`
	Function   *Function   `| @@`
	Binary     *Binary     `| @@`
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

	Variable   string      `@Ident`
	Op         string      `"="`
	Expression *Expression `@@`
}

type Binary struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Arithmetic *Arithmetic `@@`
	Op         string      `[ @( "!" "=" | "=" "=" | ">" | ">" "=" | "<" | "<" "=" | "or" | "and" )`
	Next       *Binary     `  @@ ]`
}

type Arithmetic struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Unary *Unary      `@@`
	Op    string      `[ @( "-" "="* | "+" "="* | "/" "="* | "*" "="* | "^" | "%" | "&" | "|" )`
	Next  *Arithmetic `  @@ ]`
}

type Unary struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Op       string    `( @( "!" | "-" )`
	Unary    *Unary    `  @@ )`
	Primary  *Primary  `| @@`
	Function *Function `| @@`
}

type Primary struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Call          *Call          `  @@`
	Access        *Access        `| @@`
	Ident         string         `| @Ident`
	Number        *float64       `| @Float | @Int`
	String        *string        `| @String`
	Bool          *bool          `| ( @"true" | "false" )`
	Nil           bool           `| @"nil"`
	SubExpression *Expression    `| "(" @@ ")" `
	ListLiteral   *[]Expression  `| "[" ( @@ ("," @@)* ","? )? "]"`
	ObjectLiteral *[]ObjectEntry `| "{" ( @@ ("," @@)* ","? )? "}"`
}

type Call struct {
	Ident      string        `@Ident`
	Parameters *[]Expression `"(" ( @@ ("," @@)* )? ")"`
}

type Access struct {
	Ident    string     `@Ident`
	Dot      string     `( "." @Ident`
	Brackets Expression `| "[" @@ "]" )`
}

type ObjectEntry struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Key   *Expression `@@ ":"`
	Value *Expression `@@`
}

type Function struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Parameters []*string     `( ( "(" ")" | "(" (@Ident ("," @Ident)*) ")" | @Ident ) "=" ">" )`
	Expression []*Expression `( "{" @@* "}" | @@ )`
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
