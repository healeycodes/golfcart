package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	// var cli struct {
	// 	ExpressionList []string `arg required help:"ExpressionList to parse."`
	// }

	// ctx := kong.Parse(&cli)
	// ast, err := golfcart.GenerateAST(strings.Join(cli.ExpressionList, " "))

	// if err != nil {
	// 	ctx.FatalIfErrorf(err)
	// }

	// repr.Println(ast)

	file := os.Args[1]
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	ast, err := golfcart.GenerateAST(string(b))
	if err != nil {
		println(fmt.Sprintf("%v", err))
		return
	}

	context := golfcart.Context{}
	context.Init()
	golfcart.InjectRuntimeFunctions(&context)

	result, err := ast.Eval(&context)
	if err != nil {
		println(fmt.Sprintf("%v", err))
		return
	}

	println(result.String())
}
