package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/bjatkins/silabex/cmd/font_3/dict"
	"github.com/bjatkins/silabex/cmd/font_3/font"
)

func main() {
	f, err := font.New("examples/basic_font_2.svg")
	if err != nil {
		fmt.Println("font err:", err)
		return
	}

	c := f.NewChar("*hkw", "eo", "fail")
	fmt.Println("c: ", c.SVG())
	err = os.WriteFile("gen/test_3.svg", []byte(c.SVG()), 0o0655)
	if err != nil {
		fmt.Println("svg err:", err)
		return
	}
}

func main2() {
	f, err := NewFont("examples/basic_font_2.svg")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	dict, err := dict.NewDict("cmd/v2/dict.json")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	word := "tiny"
	for i, r := range dict[word] {
		svg := f.Rune(r.Initial, r.Vowel, r.Final)
		err = os.WriteFile(fmt.Sprintf("gen/%s_%d.svg", word, i+1), []byte(svg.Render()), 0o0655)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
	}

	// r := f.Rune("T", "AOEU", "")
	// err = os.WriteFile("gen/tiny_1.svg", []byte(r.Render()), 0o0655)
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return
	// }

	// r = f.Rune("TPH*", "EU", "")
	// err = os.WriteFile("gen/tiny_2.svg", []byte(r.Render()), 0o0655)
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return
	// }
}

type Path struct {
	// I don't really need transform to mark dirty...
	transform   bool
	translateDx float64
	translateDy float64
	mirrior     bool

	path string
}

func NewElement(elem *svgparser.Element) Path {
	return Path{
		path: elem.Attributes["d"],
	}
}

func (e Path) Render() string {
	if e.transform {
		return fmt.Sprintf(
			"<path transform=\"translate(%.2f %.2f)\" d=\"%s\"/>",
			e.translateDx, e.translateDy, e.path,
		)
	}

	return fmt.Sprintf("<path d=\"%s\"/>", e.path)
}

func (e Path) Copy() Path {
	return Path{
		translateDx: e.translateDx,
		translateDy: e.translateDy,
		mirrior:     e.mirrior,

		path: e.path,
	}
}

func (e Path) Translate(dx, dy float64) Path {
	return Path{
		transform:   true,
		translateDx: e.translateDx + dx,
		translateDy: e.translateDy + dy,
		mirrior:     e.mirrior,

		path: e.path,
	}
}

type Cluster []Path

func NewCluster(elems ...*svgparser.Element) Cluster {
	elements := []Path{}
	for _, elem := range elems {
		elements = append(elements, NewElement(elem))
	}

	return Cluster(elements)
}

func CombineCluster(clusters ...Cluster) Cluster {
	paths := []Path{}
	for _, cluster := range clusters {
		for _, path := range cluster {
			paths = append(paths, path)
		}
	}
	return Cluster(paths)
}

func (c Cluster) Translate(dx, dy float64) Cluster {
	cluster := Cluster{}
	for _, e := range c {
		cluster = append(cluster, e.Translate(dx, dy))
	}

	return cluster
}

func (c Cluster) Copy() Cluster {
	cluster := Cluster{}
	for _, e := range c {
		cluster = append(cluster, e.Copy())
	}

	return cluster
}

func (c Cluster) Render() string {
	if len(c) == 0 {
		return ""
	}

	g := "<g>"
	for _, p := range c {
		g += p.Render()
	}
	g += "</g>"

	return g
}

type Rune [3]Cluster

func (r Rune) Render() string {
	g := "<svg version=\"1.1\" width=\"1000\" height=\"1000\" viewBox=\"0 0 1000 1000\" xmlns=\"http://www.w3.org/2000/svg\">"
	for _, c := range r {
		g += c.Render()
	}
	g += "</svg>"

	return g
}

type Font struct {
	Vowels     map[string]Cluster
	Initial    map[string]Cluster
	Initial1_3 map[string]Cluster
	Initial1_2 map[string]Cluster
	Solo       map[string]Cluster
	Solo1_3    map[string]Cluster
	Solo1_2    map[string]Cluster
	Final      map[string]Cluster
}

func NewFont(file string) (*Font, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	svg, err := svgparser.Parse(f, false)
	if err != nil {
		return nil, err
	}

	v3 := findElement(svg, "v3")
	vowelElem := findElement(v3, "vowel")
	vowels, err := getClusters(vowelElem, []string{"a", "ao", "e", "eu", "aoeu"})
	if err != nil {
		return nil, err
	}

	vowels["o"] = vowels["e"].Copy().Translate(0, 650)

	initial := findDeepElement(v3, "initial", "full")
	initialClusters, err := getClusters(initial, []string{"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp", "null", "s", "*"})
	if err != nil {
		return nil, err
	}

	initial1_3 := findDeepElement(v3, "initial", "1/3")
	initial1_3Clusters, err := getClusters(initial1_3, []string{"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp"})
	if err != nil {
		return nil, err
	}

	initial1_2 := findDeepElement(v3, "initial", "1/2")
	initial1_2Clusters, err := getClusters(initial1_2, []string{"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp"})
	if err != nil {
		return nil, err
	}

	solo := findDeepElement(v3, "solo", "full")
	soloClusters, err := getClusters(solo, []string{"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp", "null", "s", "*"})
	if err != nil {
		return nil, err
	}

	solo1_3 := findDeepElement(v3, "solo", "1/3")
	solo1_3Clusters, err := getClusters(solo1_3, []string{"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp"})
	if err != nil {
		return nil, err
	}

	solo1_2 := findDeepElement(v3, "solo", "1/2")
	solo1_2Clusters, err := getClusters(solo1_2, []string{"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp"})
	if err != nil {
		return nil, err
	}

	soloClusters["tph*"] = CombineCluster(solo1_3Clusters["tph"].Copy().Translate(0, -150), soloClusters["*"])

	final := findDeepElement(v3, "final", "full")
	finalClusters, err := getClusters(final, []string{"r", "rf"})
	if err != nil {
		return nil, err
	}

	return &Font{
		Vowels:     vowels,
		Initial:    initialClusters,
		Initial1_3: initial1_3Clusters,
		Initial1_2: initial1_2Clusters,
		Solo:       soloClusters,
		Solo1_3:    solo1_3Clusters,
		Solo1_2:    solo1_2Clusters,
		Final:      finalClusters,
	}, nil
}

func (f *Font) Rune(initial, vowel, final string) Rune {
	initial = strings.ToLower(initial)
	vowel = strings.ToLower(vowel)
	final = strings.ToLower(final)

	// if initial != "" && f.Initial[initial] == nil {
	// 	for k := range f.Initial {
	// 		fmt.Print(k + ", ")
	// 	}
	// 	fmt.Println("")

	// 	panic("missing initial: " + initial)
	// }
	// if vowel != "" && f.Vowels[vowel] == nil {
	// 	for k := range f.Vowels {
	// 		fmt.Print(k + ", ")
	// 	}
	// 	fmt.Println("")

	// 	panic("missing vowel: " + vowel)
	// }
	// if final != "" && f.Final[final] == nil {
	// 	for k := range f.Final {
	// 		fmt.Print(k + ", ")
	// 	}
	// 	fmt.Println("")

	// 	panic("missing final: " + final)
	// }

	if final == "" {
		return Rune{
			f.Solo[initial],
			f.Vowels[vowel],
		}
	}

	return Rune{
		f.Initial[initial],
		f.Vowels[vowel],
		f.Final[final],
	}
}

func getClusters(base *svgparser.Element, clusters []string) (map[string]Cluster, error) {
	clusterMap := map[string]Cluster{}

	for _, k := range clusters {
		elem := findElement(base, k)
		if elem == nil {
			return nil, errors.New("failed to find cluster " + k)
		}

		fmt.Println(base.Attributes["label"], "|got|", elem.Attributes["label"])
		clusterMap[k] = NewCluster(elem)
	}

	return clusterMap, nil
}

func findDeepElement(svg *svgparser.Element, names ...string) *svgparser.Element {
	base := svg
	for _, name := range names {
		base = findElement(base, name)
		if base == nil {
			return nil
		}
	}

	return base
}

func findElement(svg *svgparser.Element, name string) *svgparser.Element {
	for _, elem := range svg.Children {
		if elem.Attributes["label"] == name {
			return elem
		}
	}

	return nil
}
