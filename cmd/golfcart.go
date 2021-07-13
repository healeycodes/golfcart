package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	repl := flag.Bool("repl", false, "Enter language shell")
	debug := flag.Bool("debug", false, "Dump state after execution")
	ebnf := flag.Bool("ebnf", false, "Print EBNF grammar of the parser and quit")
	version := flag.Bool("version", false, "Print version and quit")
	flag.Parse()

	if *version {
		fmt.Printf("Golfcart v%v\n", golfcart.VERSION)
		return
	}
	if *repl {
		golfcart.REPL()
		return
	}
	if *ebnf {
		fmt.Println(golfcart.GetGrammer())
		os.Exit(0)
	}

	// If no file path, user probably wants to run the REPL
	file := flag.Arg(0)
	if file == "" {
		golfcart.REPL()
		return
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("while parsing %v: %v\n", file, err)
		os.Exit(1)
	}
	source := string(b)

	result, err := golfcart.RunProgram(source, *debug)
	if err != nil {
		fmt.Printf("while running %v: %v\n", file, err)
		os.Exit(1)
	}

	fmt.Println(*result)
}
