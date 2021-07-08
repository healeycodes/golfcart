package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", true, "Debug mode")
	flag.Parse()

	file := flag.Arg(0)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("while parsing %v: %v\n", file, err)
		os.Exit(1)
	}
	source := string(b)

	result, err := golfcart.RunProgram(source, debug)
	if err != nil {
		fmt.Printf("while running %v: %v\n", file, err)
		os.Exit(1)
	}

	println(*result)
}
