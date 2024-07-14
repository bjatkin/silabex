package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/JoshVarga/svgparser"
)

func main() {
	font, err := NewFont("examples/basic_font.svg")
	if err != nil {
		fmt.Println("failed to create new font", err)
		return
	}

	for k, v := range font.vowels {
		fmt.Println("key:", k, v)
	}
}

type Font struct {
	initialAccents map[string]string
	initialCore    map[string]string
	vowels         map[string]string
	finalAccents   map[string]string
	finalCore      map[string]string
}

var groups = []string{
	"init_accent",
	"init_core",
	"init_core_top",
	"init_core_bottom",
	"init_core_mid",
	"vowels",
	"fin_accent",
	"fin_core",
	"fin_core_top",
	"fin_core_bottom",
	"fin_core_mid",
	"full_accent",
	"full_core",
	"full_core_top",
	"full_core_bottom",
	"full_core_mid",
}

func NewFont(file string) (*Font, error) {
	f, err := os.Open("examples/basic_font.svg")
	if err != nil {
		return nil, errors.New("failed to open font file " + err.Error())
	}

	svg, err := svgparser.Parse(f, false)
	if err != nil {
		return nil, errors.New("failed to parse svg file" + err.Error())
	}

	font := &Font{}
	for _, group := range groups {
		found := findGroup(svg, group)
		if found == nil {
			return nil, errors.New("missing required group " + group)
		}

		fragments := map[string]string{}
		for _, frag := range found.Children {
			rawSvg := ""
			switch frag.Name {
			case "path":
				rawSvg = fmt.Sprintf(
					"<path style=\"%s\" d=\"%s\" />",
					frag.Attributes["style"],
					frag.Attributes["d"],
				)
			case "circle":
				rawSvg = fmt.Sprintf(
					"<circle cx=\"%s\" cy=\"%s\" r=\"%s\"/>",
					frag.Attributes["cx"],
					frag.Attributes["cy"],
					frag.Attributes["r"],
				)
			}

			fragments[sortString(frag.Attributes["label"])] = rawSvg
		}

		switch group {
		case "init_accent":
			font.initialAccents = fragments
		case "init_core":
			font.initialCore = fragments
		case "vowels":
			font.vowels = fragments
		case "fin_accent":
			font.finalAccents = fragments
		case "fin_core":
			font.finalCore = fragments
		}
	}

	return font, nil
}

func findGroup(svg *svgparser.Element, groupLabel string) *svgparser.Element {
	for _, child := range svg.Children {
		if child.Name == "g" && child.Attributes["label"] == groupLabel {
			return child
		}
		if found := findGroup(child, groupLabel); found != nil {
			return found
		}
	}

	return nil
}

func sortString(s string) string {
	r := []rune(s)
	slices.Sort(r)
	return strings.ToUpper(string(r))
}
