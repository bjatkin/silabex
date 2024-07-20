package font

import (
	"os"
	"slices"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/bjatkin/silabex/svg"
)

type Cluster int

const (
	Vowel Cluster = iota
	Solo
	Initial
	Final
)

type StrokeGroup struct {
	cluster Cluster
	group   svg.Group
}

func (s StrokeGroup) SVG() string {
	return s.group.SVG()
}

type Character struct {
	initialStrokes StrokeGroup
	vowelStrokes   StrokeGroup
	finalStrokes   StrokeGroup
}

func (c Character) SVG() string {
	ret := []string{}
	ret = append(ret, "<svg width=\"1000\" height=\"1000\" viewBox=\"0 0 1000 1000\" xmlns=\"http://www.w3.org/2000/svg\">")
	ret = append(ret, c.initialStrokes.SVG())
	ret = append(ret, c.vowelStrokes.SVG())
	ret = append(ret, c.finalStrokes.SVG())
	ret = append(ret, "</svg>")

	return strings.Join(ret, "\n")
}

type Font struct {
	soloStrokes    map[string]StrokeGroup
	initialStrokes map[string]StrokeGroup
	vowelStrokes   map[string]StrokeGroup
}

func NewFont(svgPath string) (*Font, error) {
	root, err := parseSVG(svgPath)
	if err != nil {
		return nil, err
	}

	vowels := map[string]StrokeGroup{}
	for _, name := range combinations([]string{"0", "1", "2", "3"}) {
		elem := findElem(root, "vowels", name)
		group := svg.NewGroup(elem, 0, 0)
		vowels[name] = StrokeGroup{
			cluster: Vowel,
			group:   *group,
		}
	}

	initial := map[string]StrokeGroup{}
	for _, name := range combinations([]string{"2", "3", "4", "5", "6", "7"}) {
		elem := findElem(root, "initial", "tall", name)
		group := svg.NewGroup(elem, 0, 0)
		initial[name] = StrokeGroup{
			cluster: Initial,
			group:   *group,
		}

		for _, prefix := range combinations([]string{"0", "1"}) {
			elem = findElem(root, "initial", "stand", name)
			stand := svg.NewGroup(elem, 0, 0)

			elem = findElem(root, "initial", "head", prefix)
			head := svg.NewGroup(elem, 0, 0)

			group = svg.Merge(head, stand)
			initial[prefix+name] = StrokeGroup{
				cluster: Initial,
				group:   *group,
			}
		}

		for _, suffix := range combinations([]string{"8", "9"}) {
			elem = findElem(root, "initial", "stand", name)
			stand := svg.NewGroup(elem, 0, -140)

			elem = findElem(root, "initial", "foot", suffix)
			foot := svg.NewGroup(elem, 0, 0)

			group = svg.Merge(stand, foot)
			initial[name+suffix] = StrokeGroup{
				cluster: Initial,
				group:   *group,
			}
		}

		for _, prefix := range combinations([]string{"0", "1"}) {
			for _, suffix := range combinations([]string{"8", "9"}) {
				elem = findElem(root, "initial", "core", name)
				stand := svg.NewGroup(elem, 0, 0)

				elem = findElem(root, "initial", "head", prefix)
				head := svg.NewGroup(elem, 0, 0)

				elem = findElem(root, "initial", "foot", suffix)
				foot := svg.NewGroup(elem, 0, 0)

				group = svg.Merge(head, stand, foot)
				initial[prefix+name+suffix] = StrokeGroup{
					cluster: Initial,
					group:   *group,
				}
			}
		}
	}

	solo := map[string]StrokeGroup{}
	for _, name := range combinations([]string{"2", "3", "4", "5", "6", "7"}) {
		elem := findElem(root, "solos", "tall", name)
		group := svg.NewGroup(elem, 0, 0)
		solo[name] = StrokeGroup{
			cluster: Solo,
			group:   *group,
		}

		for _, prefix := range combinations([]string{"0", "1"}) {
			elem = findElem(root, "solos", "stand", name)
			stand := svg.NewGroup(elem, 0, 0)

			elem = findElem(root, "solos", "head", prefix)
			head := svg.NewGroup(elem, 0, 0)

			group = svg.Merge(head, stand)
			solo[prefix+name] = StrokeGroup{
				cluster: Solo,
				group:   *group,
			}
		}

		for _, suffix := range combinations([]string{"8", "9"}) {
			elem = findElem(root, "solos", "stand", name)
			stand := svg.NewGroup(elem, 0, -140)

			elem = findElem(root, "solos", "foot", suffix)
			foot := svg.NewGroup(elem, 0, 0)

			group = svg.Merge(stand, foot)
			solo[name+suffix] = StrokeGroup{
				cluster: Solo,
				group:   *group,
			}
		}

		for _, prefix := range combinations([]string{"0", "1"}) {
			for _, suffix := range combinations([]string{"8", "9"}) {
				elem = findElem(root, "solos", "core", name)
				stand := svg.NewGroup(elem, 0, 0)

				elem = findElem(root, "solos", "head", prefix)
				head := svg.NewGroup(elem, 0, 0)

				elem = findElem(root, "solos", "foot", suffix)
				foot := svg.NewGroup(elem, 0, 0)

				group = svg.Merge(head, stand, foot)
				solo[prefix+name+suffix] = StrokeGroup{
					cluster: Solo,
					group:   *group,
				}
			}
		}
	}

	return &Font{
		initialStrokes: initial,
		vowelStrokes:   vowels,
		soloStrokes:    solo,
	}, nil
}

func (f *Font) NewCharacter(initial, vowel, final string) *Character {
	if final == "" {
		return &Character{
			initialStrokes: f.soloStrokes[initial],
			vowelStrokes:   f.vowelStrokes[vowel],
		}
	}

	finalStroke := f.initialStrokes[final]
	finalStroke.group.Transform(390)

	return &Character{
		initialStrokes: f.initialStrokes[initial],
		vowelStrokes:   f.vowelStrokes[vowel],
		finalStrokes:   finalStroke,
	}
}

func parseSVG(svgPath string) (*svgparser.Element, error) {
	f, err := os.Open(svgPath)
	if err != nil {
		return nil, err
	}

	root, err := svgparser.Parse(f, true)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// findElem searches the children of elem until it finds the element with the first name in the names list
// it then searches that elements child searching for the next name in the list, this continues until
// the full list of names has been searched and the final *svgparser.Element is returned
func findElem(elem *svgparser.Element, names ...string) *svgparser.Element {
	base := elem
	for _, name := range names {
		var found bool
		if base == nil {
			continue
		}

		for _, e := range base.Children {
			if e.Attributes["label"] == name {
				found = true
				base = e
				break
			}
		}

		if !found {
			return nil
		}
	}

	return base
}

func combinations(names []string) []string {
	sortName := func(name, reference []string) {
		slices.SortFunc(name, func(a, b string) int {
			if slices.Index(reference, a) > slices.Index(reference, b) {
				return 1
			}
			return -1
		})
	}

	var inner func([]string) [][]string
	inner = func(innerNames []string) [][]string {
		if len(innerNames) == 0 {
			return [][]string{}
		}

		car := innerNames[0]
		cdr := innerNames[1:]
		combos := inner(cdr)

		ret := [][]string{{car}}
		for _, combo := range combos {
			old := append([]string{}, combo...)
			sortName(old, names)
			ret = append(ret, old)

			new := append(old, car)
			sortName(new, names)
			ret = append(ret, new)
		}

		return ret
	}

	ret := []string{}
	for _, name := range inner(names) {
		ret = append(ret, strings.Join(name, ""))
	}

	return ret
}
