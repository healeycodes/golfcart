package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	file := os.Args[1]
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("while parsing %v: %v\n", file, err)
		os.Exit(1)
	}
	source := string(b)

	result, err := golfcart.RunProgram(source)
	if err != nil {
		fmt.Printf("while parsing %v: %v\n", file, err)
		os.Exit(1)
	}

	println(*result)
}
