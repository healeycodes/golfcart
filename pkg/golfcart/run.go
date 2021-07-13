package golfcart

import (
	"bufio"
	"fmt"
	"os"
)

func RunProgram(source string, debug bool) (*string, error) {
	ast, err := GenerateAST(source)
	if err != nil {
		return nil, err
	}

	context := Context{}
	context.Init()
	InjectRuntime(&context)

	result, err := ast.Eval(&context)
	if err != nil {
		return nil, err
	}

	if debug {
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
      '-..::-' golfcart v%v`+"\n", VERSION)
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
