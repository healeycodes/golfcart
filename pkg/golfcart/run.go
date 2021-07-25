package golfcart

import (
	"bufio"
	"fmt"
	"os"
)

const VERSION = 0.1

func RunProgram(source string, config Config) (*string, error) {
	ast, err := GenerateAST(source)
	if err != nil {
		return nil, err
	}

	if config.AST {
		fmt.Println(ast)
		return nil, nil
	}

	context := Context{}
	context.Init()
	InjectRuntime(&context)

	result, err := ast.Eval(&context)
	if err != nil {
		return nil, err
	}

	if config.Debug {
		fmt.Println(context.stackFrame.String())
	}

	ret := result.String()
	return &ret, nil
}

func REPL() {
	fmt.Printf(`
      .-::":-.
    .'''..''..'.
   /..''..''..''\
  ;'..''..''..''.;
  ;'..''..''..'..;
   \..''..''..''/
    '.''..''...'
      '-..::-' Golfcart v%v`+"\n", VERSION)
	context := Context{}
	context.Init()
	InjectRuntime(&context)
	for {
		endlessREPL(&context)
	}
}

func endlessREPL(context *Context) {
	defer func() {
		recover()
	}()
	for {
		fmt.Print("λ ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		ast, err := GenerateAST(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		result, err := ast.Eval(context)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(result)
	}
}

type Config struct {
	Debug bool
	AST   bool
}
