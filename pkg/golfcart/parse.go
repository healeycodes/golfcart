package golfcart

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

type ExpressionList struct {
	Pos lexer.Position

	Expressions []*Expression `@@*`
}

type Expression struct {
	Pos lexer.Position

	Assignment          *Assignment `@@`
	NativeFunctionValue *NativeFunctionValue
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
	Op      string    `( @( "and" )`
	Next    *LogicAnd `  @@ )?`
}

type LogicOr struct {
	Pos lexer.Position

	Equality *Equality `@@`
	Op       string    `( @( "or" )`
	Next     *LogicOr  `  @@ )?`
}

type Equality struct {
	Pos lexer.Position

	Comparison *Comparison `@@`
	Op         string      `( @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ )?`
}

type Comparison struct {
	Pos lexer.Position

	Addition *Addition   `@@`
	Op       string      `( @( ">" "=" | ">" | "<" "=" | "<" )`
	Next     *Comparison `  @@ )?`
}

type Addition struct {
	Pos lexer.Position

	Multiplication *Multiplication `@@`
	Op             string          `( @( "-" | "+" )`
	Next           *Addition       `  @@ )?`
}

type Multiplication struct {
	Pos lexer.Position

	Unary *Unary          `@@`
	Op    string          `( @( "/" | "*" | "%")`
	Next  *Multiplication `  @@ )?`
}

type Unary struct {
	Pos lexer.Position

	Op      string   `( @( "!" | "-" )`
	Unary   *Unary   `  @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	Pos lexer.Position

	If            *If          `@@`
	DataLiteral   *DataLiteral `| @@`
	SubExpression *Expression  `| "(" @@ ")"`
	Call          *Call        `| @@`
	// TODO: `for {}` is a parser error
	ForKeyValue *ForKeyValue `| @@`
	ForValue    *ForValue    `| @@`
	For         *For         `| @@`
	ForWhile    *ForWhile    `| @@`
	Return      *Return      `| @@`
	Break       *Break       `| @@`
	Continue    *Continue    `| @@`
	Number      *float64     `| @Float | @Int`
	Str         *string      `| @String`
	True        *bool        `| @"true"`
	False       *bool        `| @"false"`
	Nil         *bool        `| @"nil"`
	Ident       *string      `| @Ident`
}

type DataLiteral struct {
	FunctionLiteral *FunctionLiteral `@@`
	ListLiteral     *ListLiteral     `| @@`
	DictLiteral     *DictLiteral     `| @@`
}

type If struct {
	Pos lexer.Position

	Condition *Expression   `"if" @@`
	IfBody    []*Expression `"{" @@* "}"`
	ElseIf    *ElseIf       `@@*`
	ElseBody  []*Expression `( "else" "{" @@* "}" )?`
}

type ElseIf struct {
	Condition *Expression   `"else" "if" @@`
	IfBody    []*Expression `"{" @@* "}"`
	Next      *ElseIf       `@@*`
}

type FunctionLiteral struct {
	Pos lexer.Position

	Parameters []string      `( "(" ( @Ident ( "," @Ident )* )? ")" | @Ident )`
	Body       []*Expression `"=" ">" ( "{" @@* "}" | @@ )`
}

type ListLiteral struct {
	Pos lexer.Position

	Expressions *[]Expression `"[" ( @@ ( "," @@ )* )? "]"`
}

type DictLiteral struct {
	Pos lexer.Position

	DictEntry *[]DictEntry `"{" ( @@ ("," @@)* ","? )? "}"`
}

type DictEntry struct {
	Pos lexer.Position

	Ident *string     `( @Ident`
	Key   *Expression `| @@ ) ":" `
	Value *Expression `@@`
}

type Call struct {
	Pos lexer.Position

	Ident         *string     `( @Ident`
	SubExpression *Expression `| "(" @@ ")" )`
	CallChain     *CallChain  `@@`
}

type CallChain struct {
	Parameters     *[]Expression `( "(" ( @@ ( "," @@ )* )? ")" `
	Access         *string       `    | "." @Ident`
	ComputedAccess *Expression   `    | "[" @@ "]" )`
	Next           *CallChain    `@@?`
}

type Break struct {
	Pos lexer.Position

	Break *string `"break"`
}

type Continue struct {
	Pos lexer.Position

	Continue *string `"continue"`
}

type Return struct {
	Pos lexer.Position

	Return     *string     `( "return" `
	Expression *Expression `@@ )`
}

type For struct {
	Pos lexer.Position

	Init      []*Assignment `"for" ( @@ ";"`
	Condition *Expression   `@@ ";"`
	Post      *Expression   `@@`
	Body      []*Expression `"{" @@* "}" )`
}

type ForValue struct {
	Pos lexer.Position

	Value                *string       `( "for" @Ident "in"`
	Collection           *string       `( @Ident`
	CollectionExpression *Expression   `| @@) `
	Body                 []*Expression `"{" @@* "}" )`
}

type ForKeyValue struct {
	Pos lexer.Position

	Key                  *string       `( "for" @Ident ","`
	Value                *string       `@Ident "in"`
	Collection           *string       `( @Ident`
	CollectionExpression *Expression   `| @@) `
	Body                 []*Expression `"{" @@* "}" )`
}

type ForWhile struct {
	Pos lexer.Position

	Init      []*Assignment `"for" (`
	Condition *Expression   `@@?`
	Post      *Expression   ``
	Body      []*Expression `"{" @@* "}" )`
}

var (
	_lexer = lexer.Must(stateful.New(stateful.Rules{
		"Root": {
			{"comment", `//.*|/\*.*?\*/`, nil},
			{"whitespace", `[\n\r\t ]+`, nil},
			{"Float", `[+-]?([0-9]*[.])?[0-9]+`, nil},
			{"Int", `[\d]+`, nil},
			{"String", `"([^"]*)"`, nil},
			{"Ident", `[\w]+`, nil},
			{"Punct", `[-[!*%()+_={}\|:;"<,>./]|]`, nil},
		},
	}))
	parser = participle.MustBuild(&ExpressionList{}, participle.Lexer(_lexer),
		participle.Elide("whitespace", "comment"), participle.UseLookahead(2))
)

func GetGrammer() string {
	return parser.String()
}

func GenerateAST(source string) (*ExpressionList, error) {
	expressionList := &ExpressionList{}

	err := parser.ParseString("", source, expressionList)
	if err != nil {
		return nil, err
	}

	return expressionList, nil
}
