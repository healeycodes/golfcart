package funlang

import (

	"github.com/alecthomas/participle/v2"
)

parser, err := participle.Build(&INI{})

ast := &INI{}
err := parser.ParseString("", "size = 10", ast)