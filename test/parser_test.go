package golfcart

import (
	"testing"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func TestGenerateAST(t *testing.T) {
	program := "a = 1"
	_, err := golfcart.GenerateAST(program)
	if err != nil {
		t.Errorf("GenerateAST: %v", err)
	}
}

func TestEval(t *testing.T) {
	program := "a = 1"

	ast, _ := golfcart.GenerateAST(program)

	context := golfcart.Context{}
	context.Init()
	golfcart.InjectRuntime(&context)

	_, err := ast.Eval(&context)
	if err != nil {
		t.Errorf("Eval: %v", err)
	}
}
