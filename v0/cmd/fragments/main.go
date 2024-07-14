package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/JoshVarga/svgparser"
)

func main() {
	raw, err := os.ReadFile("cmd/universal_char.svg")
	if err != nil {
		fmt.Println("read err:", err)
		return
	}

	reader := strings.NewReader(string(raw))

	element, err := svgparser.Parse(reader, false)
	if err != nil {
		fmt.Println("parse err:", err)
		return
	}

	for _, e := range element.Children {
		if e.Name != "g" {
			continue
		}
		for _, e := range e.Children {
			switch e.Name {
			case "path":
				fmt.Printf("%s Fragment = Path{d: \"%s\"}\n",
					snakeToPascal(e.Attributes["label"]),
					e.Attributes["d"],
				)
			case "circle":
				fmt.Printf("%s Fragment = Circle{cx: %s, cy: %s, r: %s}\n",
					snakeToPascal(e.Attributes["label"]),
					e.Attributes["cx"],
					e.Attributes["cy"],
					e.Attributes["r"],
				)
			}
		}
	}
}

func snakeToPascal(name string) string {
	new := ""
	cap := false
	for i, r := range name {
		switch {
		case r == '_':
			cap = true
		case cap:
			new += strings.ToUpper(string(r))
			cap = false
		case i == 0:
			new += strings.ToUpper(string(r))
		default:
			new += string(r)
		}
	}

	return new
}
