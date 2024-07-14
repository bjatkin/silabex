package stroke

import (
	"github.com/JoshVarga/svgparser"
)

type Cluster int

const (
	Solo Cluster = iota
	Initial
	Vowel
	Final
)

type Segment int

const (
	Border Segment = iota
	Head
	Hang
	Core
	Stand
	Foot
	Tall
)

type element struct {
	segment Segment

	elements []svgparser.Element

	dx, dy float64
	sx, sy float64
}

func (e *element) up(cluster Cluster) {
	if cluster == Vowel {
		e.dx -= 880
	}

	if e.segment == Stand {
		e.dx -= 140
		e.segment = Hang
	}

	if e.segment == Foot {
		e.dx -= 620
		e.segment = Head
	}
}

func (e *element) down(cluster Cluster) {
	if cluster == Vowel {
		e.dx += 880
	}

	if e.segment == Hang {
		e.dx += 140
		e.segment = Stand
	}

	if e.segment == Head {
		e.dx += 620
		e.segment = Foot
	}
}

func (e *element) left(cluster Cluster) {
	if cluster == Vowel {
		e.dy -= 880
	}

	if cluster == Final {
		e.dy -= 390
	}
}

func (e *element) right(cluster Cluster) {
	if cluster == Vowel {
		e.dy += 880
	}

	if cluster == Initial {
		e.dy += 390
		cluster = Final
	}
}

func (e *element) flipX() {
	e.sx *= -1
}

func (e *element) FlipY() {
	e.sy *= -1
}

func (e *element) copy() element {
	return element{
		segment:  e.segment,
		elements: append([]svgparser.Element{}, e.elements...),
		dx:       e.dx,
		dy:       e.dy,
		sx:       e.sx,
		sy:       e.sy,
	}
}

type Stroke struct {
	name     string
	cluster  Cluster
	elements []element
}

func (s *Stroke) copy() *Stroke {
	elements := []element{}
	for _, e := range s.elements {
		elements = append(elements, e.copy())
	}

	return &Stroke{
		name:     s.name,
		elements: elements,
	}
}

func (s *Stroke) Name() string {
	return s.name
}

func (s *Stroke) Cluster() Cluster {
	return s.cluster
}

func (s *Stroke) Segments() []Segment {
	segments := map[Segment]struct{}{}

	for _, element := range s.elements {
		segments[element.segment] = struct{}{}
	}

	ret := []Segment{}
	for key := range segments {
		ret = append(ret, key)
	}

	return ret
}

func (s *Stroke) Up() *Stroke {
	for _, e := range s.elements {
		e.up(s.cluster)
	}

	return s
}

func (s *Stroke) Down() *Stroke {
	for _, e := range s.elements {
		e.down(s.cluster)
	}

	return s
}

func (s *Stroke) Left() *Stroke {
	for _, e := range s.elements {
		e.left(s.cluster)
	}

	if s.cluster == Final {
		s.cluster = Initial
	}

	return s
}

func (s *Stroke) Right() *Stroke {
	for _, e := range s.elements {
		e.right(s.cluster)
	}

	if s.cluster == Initial {
		s.cluster = Final
	}

	return s
}

func (s *Stroke) FlipX() *Stroke {
	for _, e := range s.elements {
		e.flipX()
	}

	return s
}

func (s *Stroke) FlipY() *Stroke {
	for _, e := range s.elements {
		e.FlipY()
	}

	return s
}

func Match(a, b *Stroke) bool {
	return a.cluster == b.cluster
}

func Join(name string, a, b *Stroke) *Stroke {
	return &Stroke{
		name:     name,
		elements: append(a.copy().elements, b.copy().elements...),
	}
}

type StrokeSlice []*Stroke

func NewSlice(strokes ...*Stroke) StrokeSlice {
	return StrokeSlice(strokes)
}

func (s StrokeSlice) SetName(name string) StrokeSlice {
	for _, stroke := range s {
		stroke.name = name
	}

	return s
}

func (s StrokeSlice) Up() StrokeSlice {
	for _, stroke := range s {
		stroke.Up()
	}

	return s
}

func (s StrokeSlice) Down() StrokeSlice {
	for _, stroke := range s {
		stroke.Down()
	}

	return s
}

func (s StrokeSlice) Left() StrokeSlice {
	for _, stroke := range s {
		stroke.Left()
	}

	return s
}

func (s StrokeSlice) Right() StrokeSlice {
	for _, stroke := range s {
		stroke.Right()
	}

	return s
}

func (s StrokeSlice) FlipX() StrokeSlice {
	for _, stroke := range s {
		stroke.FlipX()
	}

	return s
}

func (s StrokeSlice) FlipY() StrokeSlice {
	for _, stroke := range s {
		stroke.FlipY()
	}

	return s
}
