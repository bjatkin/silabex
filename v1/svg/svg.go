package svg

import (
	"fmt"
	"strings"

	"github.com/JoshVarga/svgparser"
)

// Group represents a group of strokes that make up a silbex character
type Group struct {
	dx       float64
	elements []*svgparser.Element
}

func (g *Group) Transform(dx float64) {
	g.dx += dx
}

func (g Group) SVG() string {
	if len(g.elements) == 0 {
		return ""
	}

	ret := []string{}
	if g.dx == 0 {
		ret = append(ret, "<g>")
	} else {
		ret = append(ret, fmt.Sprintf("<g transform=\"translate(%.2f 0)\">", g.dx))
	}

	for _, elem := range g.elements {
		if elem.Attributes["transform"] != "" {
			ret = append(ret, fmt.Sprintf("<path transform=\"%s\" d=\"%s\" />", elem.Attributes["transform"], elem.Attributes["d"]))
		} else {
			ret = append(ret, fmt.Sprintf("<path d=\"%s\" />", elem.Attributes["d"]))
		}
	}

	ret = append(ret, "</g>")

	return strings.Join(ret, "\n")
}

// NewGroup creates a new Group from an svgparser.Element
func NewGroup(root *svgparser.Element, dx, dy float64) *Group {
	copyElement := func(elem *svgparser.Element) *svgparser.Element {
		attrs := map[string]string{}
		for k, v := range elem.Attributes {
			attrs[k] = v
		}

		if dy == 0 {
			attrs["transform"] = fmt.Sprintf("translate(0 %.2f)", dy)
		}

		return &svgparser.Element{
			Name:       root.Name,
			Attributes: attrs,
			Content:    root.Content,
		}
	}

	elements := []*svgparser.Element{}
	for _, elem := range root.Children {
		elements = append(elements, copyElement(elem))
	}

	return &Group{
		elements: elements,
	}
}

// Merge merges multiple groups together
func Merge(groups ...*Group) *Group {
	elements := []*svgparser.Element{}
	for _, g := range groups {
		elements = append(elements, g.elements...)
	}

	return &Group{
		dx:       groups[0].dx,
		elements: elements,
	}
}
