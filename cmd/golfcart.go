package main

import (
	"io/ioutil"
	"os"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func main() {
	file := os.Args[1]
	b, err := ioutil.ReadFile(file)
	if err != nil {
		println(err)
	}
	source := string(b)

	result, err := golfcart.RunProgram(source)
	if err != nil {
		panic(err)
	}

	println(result)
}
