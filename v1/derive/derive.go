package derive

import (
	"slices"
	"strings"

	"github.com/bjatkin/silabex/stroke"
)

type builder func(stroke.StrokeSlice) stroke.StrokeSlice

func buildSimpleVowels(s stroke.StrokeSlice) stroke.StrokeSlice {
	ret := stroke.StrokeSlice{}

	ret = append(ret, find(s, "1", filterInCluster(stroke.Vowel)).Up().SetName("2")...)
	ret = append(ret, find(s, "0", filterInCluster(stroke.Vowel)).Right().SetName("3")...)

	return ret
}

func buildComplexVowels(s stroke.StrokeSlice) stroke.StrokeSlice {
	ret := stroke.StrokeSlice{}

	for _, name := range combinations([]string{"0", "1", "2", "3"}) {
		switch len(name) {
		case 1:
			// nothing to do here, if the stoke is of len 1 and can be found it's already derived
		case 2:
			ret = append(ret, join(
				name,
				find(s, string(name[0])),
				find(s, string(name[1])),
			)...)
		case 3:
			ret = append(ret, join(
				name,
				find(s, string(name[0])),
				find(s, string(name[1])),
				find(s, string(name[2])),
			)...)
		case 4:
			ret = append(ret, join(
				name,
				find(s, string(name[0])),
				find(s, string(name[1])),
				find(s, string(name[2])),
				find(s, string(name[3])),
			)...)
		}
	}

	return ret
}

func buildComplexIntitials(s stroke.StrokeSlice) stroke.StrokeSlice {
	ret := stroke.StrokeSlice{}

	ret = append(ret, find(s, "3", filterInCluster(stroke.Initial, stroke.Solo)).FlipX().SetName("2")...)
	ret = append(ret, find(s, "3", filterInCluster(stroke.Initial, stroke.Solo)).FlipY().SetName("7")...)
	ret = append(ret, find(s, "3", filterInCluster(stroke.Initial, stroke.Solo)).FlipX().FlipY().SetName("6")...)

	return ret
}

func buildConsonantHeaderAndFooter(s stroke.StrokeSlice) stroke.StrokeSlice {
	ret := stroke.StrokeSlice{}

	for _, name := range combinations([]string{"2", "3", "4", "5", "6", "7"}) {
		for _, head := range combinations([]string{"0", "1"}) {
			for _, foot := range combinations([]string{"8", "9"}) {
				ret = append(ret, join(
					head+name,
					find(s, head, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Head)),
					find(s, name, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Stand)),
				)...)

				ret = append(ret, join(
					name+foot,
					find(s, name, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Stand)).Up(),
					find(s, foot, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Foot)),
				)...)

				ret = append(ret, join(
					head+name+foot,
					find(s, head, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Head)),
					find(s, name, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Core)),
					find(s, foot, filterInCluster(stroke.Initial, stroke.Solo), filterOnlyInSegment(stroke.Foot)),
				)...)
			}
		}
	}

	return ret
}

type CharBuilder struct {
	builders []builder
}

func NewCharBuilder() *CharBuilder {
	return &CharBuilder{
		builders: []builder{
			buildSimpleVowels,
			buildComplexVowels,
			buildConsonantHeaderAndFooter,
		},
	}
}

type filter func(*stroke.Stroke) bool

func filterInCluster(cluster ...stroke.Cluster) filter {
	return func(s *stroke.Stroke) bool {
		for _, c := range cluster {
			if s.Cluster() == c {
				return true
			}
		}
		return false
	}
}

func filterOnlyInSegment(segment stroke.Segment) filter {
	return func(s *stroke.Stroke) bool {
		segments := s.Segments()
		if len(segments) > 1 {
			return false
		}

		return segment == segments[0]
	}
}

func find(strokes stroke.StrokeSlice, name string, filters ...filter) stroke.StrokeSlice {
	ret := stroke.StrokeSlice{}

	for _, stroke := range strokes {
		ok := true
		for _, filter := range filters {
			if !filter(stroke) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		if stroke.Name() != name {
			continue
		}

		ret = append(ret, stroke)
	}

	return ret
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

func join(name string, strokeSlices ...stroke.StrokeSlice) stroke.StrokeSlice {
	if len(strokeSlices) == 0 {
		return stroke.StrokeSlice{}
	}

	ret := strokeSlices[0]
	for _, strokeSlice := range strokeSlices[1:] {
		ret = joinTwo(name, ret, strokeSlice)
	}

	return ret
}

func joinTwo(name string, a, b stroke.StrokeSlice) stroke.StrokeSlice {
	ret := stroke.StrokeSlice{}
	for _, outer := range a {
		for _, inner := range b {
			if !stroke.Match(outer, inner) {
				continue
			}

			ret = append(ret, stroke.Join(name, outer, inner))
			break
		}
	}

	return ret
}
