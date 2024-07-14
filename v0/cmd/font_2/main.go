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
	font, err := NewFont("examples/basic_font_2.svg")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	r := font.Rune("*", "aoeu", "s")
	for _, fragment := range r {
		fmt.Println(fragment)
	}

	fmt.Println(renderRune(r))
	err = os.WriteFile("test2.svg", []byte(renderRune(r)), 0o0655)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}

func renderRune(fragments []*Fragment) string {
	got := "<svg version=\"1.1\" width=\"1000\" height=\"1000\" viewBox=\"0 0 1000 1000\" xmlns=\"http://www.w3.org/2000/svg\"><g>"
	for _, f := range fragments {
		if f == nil {
			continue
		}

		got += f.Render()
	}
	got += "</g></svg>"

	return got
}

type Font struct {
	clusters map[string]*Fragment
}

func (f *Font) Rune(initial, vowel, final string) []*Fragment {
	return []*Fragment{
		f.Fragment(initial),
		f.Fragment(vowel),
		f.Fragment(final),
	}
}

func (f *Font) Fragment(search string) *Fragment {
	runes := []rune(strings.ToUpper(search))
	slices.Sort(runes)
	return f.clusters[string(runes)]
}

func NewFont(file string) (*Font, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	svg, err := svgparser.Parse(f, false)
	if err != nil {
		return nil, errors.New("failed to parse svg file " + err.Error())
	}

	clusters, err := readClusters("left.txt", "right.txt", "vowel.txt")
	if err != nil {
		return nil, err
	}

	clusterMap := map[string]*Fragment{}
	for _, cluster := range clusters {
		clusterMap[cluster] = findGroup(svg, cluster)
		if clusterMap[cluster] != nil {
			fmt.Println("got ", cluster)
		}
	}

	for _, cluster := range clusters {
		if clusterMap[cluster] != nil {
			continue
		}

		build, ok := builders[cluster]
		if !ok {
			// skip errors for now
			continue
			// return errors.New("missing cluster: " + cluster)
		}

		clusterMap[cluster] = build(clusterMap)
		fmt.Println("got ", cluster)
	}

	return &Font{
		clusters: clusterMap,
	}, nil
}

func readClusters(files ...string) ([]string, error) {
	names := []string{}
	for _, file := range files {
		clusters, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		for _, line := range strings.Split(string(clusters), "\n") {
			if line == "" {
				continue
			}

			names = append(names, line)
		}
	}

	return names, nil
}

func findGroup(svg *svgparser.Element, group string) *Fragment {
	for _, elem := range svg.Children {
		if elem.Name != "g" {
			continue
		}
		if elem.Attributes["label"] != "new" {
			continue
		}

		for _, e := range elem.Children {
			if equalGroup(e.Attributes["label"], group) {
				return &Fragment{d: e.Attributes["d"]}
			}
		}
	}

	return nil
}

func equalGroup(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	aRunes := []rune(strings.ToUpper(a))
	slices.Sort(aRunes)
	bRunes := []rune(strings.ToUpper(b))
	slices.Sort(bRunes)

	for i := range aRunes {
		if aRunes[i] != bRunes[i] {
			return false
		}
	}

	return true
}

type Cluster []Fragment

type Fragment struct {
	translate [2]float64
	rotate    [3]float64

	d string
}

func (f *Fragment) Render() string {
	translate := fmt.Sprintf("translate(%.2f %.2f)", f.translate[0], f.translate[1])
	rotate := fmt.Sprintf("rotate(%.2f %.2f %.2f)", f.rotate[0], f.rotate[1], f.rotate[2])

	transform := ""
	switch {
	case translate == "translate(0.00 0.00)" && rotate == "rotate(0.00 0.00 0.00)":
		break
	case rotate == "rotate(0.00 0.00 0.00)":
		transform = fmt.Sprintf(" transform=\"%s\"", translate)
	case translate == "translate(0.00 0.00)":
		transform = fmt.Sprintf(" transform=\"%s\"", rotate)
	default:
		transform = fmt.Sprintf(" transform=\"%s %s\"", translate, rotate)
	}

	return fmt.Sprintf("<path%s d=\"%s\" />", transform, f.d)
}

type builder func(clusterMap map[string]*Fragment) *Fragment

var builders = map[string]builder{
	"*": func(clusterMap map[string]*Fragment) *Fragment {
		return &Fragment{
			translate: [2]float64{0, 625},
			d:         clusterMap["S"].d,
		}
	},
}
