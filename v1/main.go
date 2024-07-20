package main

import (
	"fmt"
	"os"

	"github.com/bjatkin/silabex/font"
)

func main() {
	f, err := font.NewFont("reference/font2.svg")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	char := f.NewCharacter("0459", "0123", "")
	err = os.WriteFile("reference/test.svg", []byte(char.SVG()), 0o0655)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}
