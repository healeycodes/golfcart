package golfcart

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ExpressionList struct {
	Pos lexer.Position

	Expressions []*Expression `@@*`
}

type Expression struct {
	Pos lexer.Position

	If         *If         `@@`
	For        *For        `| @@`
	While      *While      `| @@`
	Break      *Break      `| @@`
	Continue   *Continue   `| @@`
	Assignment *Assignment `| @@`
}

type Break struct {
	Pos lexer.Position

	Break string `"break"`
}

type Continue struct {
	Pos lexer.Position

	Continue string `"continue"`
}

type For struct {
	Pos lexer.Position

	Init      []*Assignment `"for" ( @@* ("," @@ )* )* ";"`
	Condition *Expression   `@@* ";"`
	Post      *Expression   `@@*`
	Body      []*Expression `"{" @@* "}"`
}

type While struct {
	Pos lexer.Position

	Condition *Expression   `"while" @@`
	Body      []*Expression `"{" @@* "}"`
}

type If struct {
	Pos lexer.Position

	Init      []*Assignment `"if" ( @@ ";" )?`
	Condition *Expression   `@@`
	IfBody    []*Expression `"{" @@* "}"`
	ElseBody  []*Expression `( "else" "{" @@* "}" )?`
}

type Assignment struct {
	Pos lexer.Position

	LogicAnd *LogicAnd `@@`
	Op       string    `( @"="`
	Next     *LogicAnd `  @@ )?`
}

type LogicAnd struct {
	Pos lexer.Position

	LogicOr *LogicOr  `@@`
	Op      string    `[ @( "and" )`
	Next    *LogicAnd `  @@ ]`
}

type LogicOr struct {
	Pos lexer.Position

	Equality *Equality `@@`
	Op       string    `[ @( "or" )`
	Next     *LogicOr  `  @@ ]`
}

type Equality struct {
	Pos lexer.Position

	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

type Comparison struct {
	Pos lexer.Position

	Addition *Addition   `@@`
	Op       string      `[ @( ">" "=" | ">" | "<" "=" | "<" )`
	Next     *Comparison `  @@ ]`
}

type Addition struct {
	Pos lexer.Position

	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

type Multiplication struct {
	Pos lexer.Position

	Unary *Unary          `@@`
	Op    string          `[ @( "/" | "*" | "%")`
	Next  *Multiplication `  @@ ]`
}

type Unary struct {
	Pos lexer.Position

	Op      string   `( @( "!" | "-" )`
	Unary   *Unary   `  @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	Pos lexer.Position

	FunctionLiteral *FunctionLiteral `@@`
	ListLiteral     *[]Expression    `| "[" ( @@ ("," @@)* ","? )? "]"`
	ObjectLiteral   *[]ObjectEntry   `| "{" ( @@ ("," @@)* ","? )? "}"`
	Call            *Call            `| @@`
	SubExpression   *Expression      `| "(" @@ ")"`
	Ident           string           `| @Ident`
	Number          *float64         `| @Float | @Int`
	Str             *string          `| @String`
	Bool            *bool            `| ( @"true" | "false" )`
	Nil             *bool            `| @"nil"`
}

type FunctionLiteral struct {
	Pos lexer.Position

	Parameters []string      `"(" ( @Ident ( "," @Ident )* )? ")"`
	Body       []*Expression `"=" ">" ( "{" @@* "}" | @@ )`
}

type ObjectEntry struct {
	Pos lexer.Position

	Key   *Expression `@@ ":"`
	Value *Expression `@@`
}

type Call struct {
	Pos lexer.Position

	Ident          *string       `( @Ident`
	SubExpression  *Expression   `| "(" @@ ")" )`
	Parameters     *[]Expression `( "(" ( @@ ( "," @@ )* )? ")" `
	Access         string        `    | "." @Ident`
	ComputedAccess *Expression   `    | "[" @@ "]" )`
}

var (
	parser = participle.MustBuild(&ExpressionList{}, participle.UseLookahead(2))
)

func GenerateAST(source string) (*ExpressionList, error) {
	expressionList := &ExpressionList{}

	// Print grammar
	// println(parser.String())

	err := parser.ParseString("", source, expressionList)
	if err != nil {
		return nil, err
	}
	return expressionList, nil
}
