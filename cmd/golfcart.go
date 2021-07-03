package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	var cli struct {
		ExpressionList []string `arg required help:"ExpressionList to parse."`
	}

	ctx := kong.Parse(&cli)
	ast, err := golfcart.GenerateAST(strings.Join(cli.ExpressionList, " "))

	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	// repr.Println(ast)
	frame := golfcart.StackFrame{Values: make(map[string]golfcart.Value)}
	result, err := ast.Eval(&frame)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%T\n", result)
	println(result.String())
}
