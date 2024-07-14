package main

import (
	"fmt"
	"strings"
)

func main() {
	dict, err := NewDict("cmd/v2/dict.json")
	if err != nil {
		fmt.Println("dict err:", err)
		return
	}

	phrase := "hello tiny moon"
	words := strings.Split(phrase, " ")
	for _, word := range words {
		symbols, err := dict.Lookup(word)
		if err != nil {
			fmt.Println("symbol err:", err)
			return
		}

		fmt.Println(word, symbols)
	}
}
