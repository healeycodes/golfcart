package golfcart

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func runAllProgramsInDir(t *testing.T, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		cur := path + file.Name()
		b, err := ioutil.ReadFile(cur)
		if err != nil {
			log.Fatal(err)
		}
		source := string(b)
		_, err = golfcart.RunProgram(source, false)
		if err != nil {
			t.Errorf("RunProgram(%s): %v", cur, err)
		}
	}
}

func TestPrograms(t *testing.T) {
	path := "../example programs/spec programs/"
	runAllProgramsInDir(t, path)

	path = "../example programs/"
	runAllProgramsInDir(t, path)
}

func TestBadPrograms(t *testing.T) {
	path := "../example programs/error programs/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		cur := path + file.Name()
		b, err := ioutil.ReadFile(cur)
		if err != nil {
			log.Fatal(err)
		}
		source := string(b)
		_, err = golfcart.RunProgram(source, false)
		if err == nil {
			t.Errorf("RunProgram(%s): didn't throw an error", cur)
		}
	}
}
