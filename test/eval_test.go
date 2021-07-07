package golfcart

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/healeycodes/golfcart/pkg/golfcart"
)

func TestPrograms(t *testing.T) {
	path := "./programs/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		cur := path + file.Name()
		b, err := ioutil.ReadFile(cur)
		if err != nil {
			log.Fatal(err)
		}
		source := string(b)
		_, err = golfcart.RunProgram(source)
		if err != nil {
			t.Errorf("RunProgram(%s): %v", cur, err)
		}
	}
}
