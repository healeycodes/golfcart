package main

import (
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	file := os.Args[1]
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		println(err)
	}
	source := string(b)

	result, err := RunProgram(source)
	if err != nil {
		panic(err)
	}

	println(result)
}

func RunProgram(source string) (*string, error) {
	ast, err := golfcart.GenerateAST(source)
	if err != nil {
		return nil, err
	}

	context := golfcart.Context{}
	context.Init()
	golfcart.InjectRuntime(&context)

	result, err := ast.Eval(&context)
	if err != nil {
		return nil, err
	}

	ret := result.String()
	return &ret, nil
}
