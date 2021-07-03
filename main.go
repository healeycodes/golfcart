package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"
	"github.com/healeycodes/lang"
)

func main() {

	var cli struct {
		ExpressionList []string `arg required help:"ExpressionList to parse."`
	}

	ctx := kong.Parse(&cli)
	ast, err := lang.GenerateAST(strings.Join(cli.ExpressionList, " "))

	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	repr.Println(ast)
}
