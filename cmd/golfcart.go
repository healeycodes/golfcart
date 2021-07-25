package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	debug := flag.Bool("debug", false, "Dump state after execution")
	astOnly := flag.Bool("ast", false, "Print AST and quit.")
	ebnf := flag.Bool("ebnf", false, "Print EBNF grammar of the parser and quit")
	version := flag.Bool("version", false, "Print version and quit")
	flag.Parse()

	config := golfcart.Config{}
	if *version {
		fmt.Printf("Golfcart v%v\n", golfcart.VERSION)
		return
	}
	if *ebnf {
		fmt.Println(golfcart.GetGrammer())
		return
	}
	if *debug {
		config.Debug = true
	}
	if *astOnly {
		config.AST = true
	}

	// If no file path, user probably wants to run the REPL
	file := flag.Arg(0)
	if file == "" {
		golfcart.REPL()
		return
	}

	result, err := golfcart.RunProgram(readFile(file), config)
	if err != nil {
		fmt.Printf("while running %v: %v\n", file, err)
		os.Exit(1)
	}

	fmt.Println(*result)
}

func readFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("while reading %v: %v\n", file, err)
		os.Exit(1)
	}
	return string(b)
}
