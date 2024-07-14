package font

import (
	"fmt"
	"os"
	"strings"

	"github.com/JoshVarga/svgparser"
)

// stroke represents an svg stroke
type stroke struct {
	// TODO: would this just be better to include a transformation matrix
	// TODO: maybe a v2 would be nice here as well
	ox float64
	oy float64

	dx float64
	dy float64

	sx float64
	sy float64

	path string
}

// strokeFromElem creates a new path from an *svgparser.Element. The element must be a
// path element or an error will be returned
func strokeFromElem(elem *svgparser.Element) (*stroke, error) {
	if elem.Name != "path" {
		return nil, fmt.Errorf("can not create stroke from %s element", elem.Name)
	}

	return &stroke{
		sx:   1,
		sy:   1,
		path: elem.Attributes["d"],
	}, nil
}

// clone creates a clone of the current stroke
func (s *stroke) clone() *stroke {
	return &stroke{
		ox: s.ox,
		oy: s.oy,

		sx: s.sx,
		sy: s.sy,

		dx: s.dx,
		dy: s.dy,

		path: s.path,
	}
}

// svg renders the stroke as a valid svg string
func (s *stroke) svg() string {
	transform := ""
	if s.dx != 0 || s.dy != 0 || s.sx != 1 || s.sy != 1 {
		transform = fmt.Sprintf(
			"transform=\"scale(%.2f %.2f) translate(%.2f %.2f)\" ",
			s.sx, s.sy, s.dx, s.dy,
		)
	}

	origin := ""
	if s.ox != 0 || s.oy != 0 {
		origin = fmt.Sprintf(
			"transform-origin=\"%.2f %.2f\" ",
			s.ox, s.oy,
		)
	}

	return fmt.Sprintf(
		"<path %s%sd=\"%s\"/>",
		origin, transform, s.path,
	)
}

// cluster is a slice of paths that make up a silbex stroke cluster
type cluster []stroke

// clone creates a clone of the current cluster
func (c *cluster) clone() *cluster {
	ret := cluster{}
	for _, path := range *c {
		ret = append(ret, *path.clone())
	}

	return &ret
}

// svg returns the cluster as a valid svg string
func (c *cluster) svg() string {
	if c == nil {
		return ""
	}

	ret := ""
	for _, stroke := range *c {
		ret += stroke.svg()
	}
	return ret
}

// translate sets the translation for all the strokes that make up the cluster to dx and dy
func (c *cluster) translate(dx, dy float64) *cluster {
	for i := range *c {
		(*c)[i].dx = dx
		(*c)[i].dy = dy
	}

	return c
}

// scale sets the scale for all the strokes that make up the cluster to sx and  sy
func (c *cluster) scale(sx, sy float64) *cluster {
	for i := range *c {
		(*c)[i].sx = sx
		(*c)[i].sy = sy
	}

	return c
}

// mirriorX is the same as calling scale(-1, 1)
func (c *cluster) mirriorX() *cluster {
	return c.scale(-1, 1)
}

// mirrorY is the same as calling scale(1, -1)
func (c *cluster) mirriorY() *cluster {
	return c.scale(1, -1)
}

// mirrorXY is the same as calling scale(-1, -1)
func (c *cluster) mirrorXY() *cluster {
	return c.scale(-1, -1)
}

// origin sets the transform origin for all the strokes that make up the cluster to ox and oy
func (c *cluster) origin(ox, oy float64) *cluster {
	for i := range *c {
		(*c)[i].ox = ox
		(*c)[i].oy = oy
	}

	return c
}

// add adds the given cluster to the current cluster
func (c *cluster) add(cluster *cluster) *cluster {
	*c = append(*c, *cluster...)
	return c
}

// Char is the representation of single rune in silbex
// TODO: char should be [4]byte so it works similarly to utf-8
type Char struct {
	initial *cluster
	vowel   *cluster
	final   *cluster
}

// SVG renders the char as an svg image
func (c *Char) SVG() string {
	g := "<svg version=\"1.1\" width=\"1000\" height=\"1000\" viewBox=\"0 0 1000 1000\" xmlns=\"http://www.w3.org/2000/svg\">"
	g += c.initial.svg()
	g += c.vowel.svg()
	g += c.final.svg()
	g += "</svg>"

	return g
}

// Font can be used to generate any silbex character based on the provided
// svg font template
type Font struct {
	solo    map[string]*cluster
	initial map[string]*cluster
	vowels  map[string]*cluster
	final   map[string]*cluster
}

// New creates a new Font
func New(file string) (*Font, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	svg, err := svgparser.Parse(f, false)
	if err != nil {
		return nil, err
	}

	// TODO: this should get removed eventually
	v3 := findElem(svg, "v3")

	vowels, err := getVowels(v3)
	if err != nil {
		return nil, err
	}

	initials, err := getInitials(v3)
	if err != nil {
		return nil, err
	}

	return &Font{
		vowels:  vowels,
		initial: initials,
	}, nil
}

func (f *Font) NewChar(initial, vowel, final string) Char {
	initial = strings.ToLower(initial)
	vowel = strings.ToLower(vowel)
	final = strings.ToLower(final)

	if final == "" {
		return Char{
			initial: f.solo[initial],
			vowel:   f.vowels[vowel],
		}
	} else {
		return Char{
			initial: f.initial[initial],
			vowel:   f.vowels[vowel],
			final:   f.final[final],
		}
	}
}

// getVowels reads all the vowels from the "vowel" group and then adds all the derived vowels
// so that all possible vowel combinations are accounted for
func getVowels(base *svgparser.Element) (map[string]*cluster, error) {
	elem := findElem(base, "vowel")
	vowels, err := getClusters(elem, []string{"a", "ao", "e", "eu", "aoeu"})
	if err != nil {
		return nil, err
	}

	// add derived vowels
	vowels["o"] = vowels["e"].clone().translate(0, 880)
	vowels["u"] = vowels["a"].clone().translate(880, 0)
	vowels["ae"] = vowels["a"].clone().add(vowels["e"].clone())
	vowels["au"] = vowels["a"].clone().add(vowels["u"].clone())
	vowels["eo"] = vowels["e"].clone().add(vowels["o"].clone())
	vowels["ou"] = vowels["o"].clone().add(vowels["u"].clone())
	vowels["aoe"] = vowels["ao"].clone().add(vowels["e"].clone())
	vowels["aou"] = vowels["ao"].clone().add(vowels["u"].clone())
	vowels["aeu"] = vowels["ae"].clone().add(vowels["u"].clone())
	vowels["eou"] = vowels["eo"].clone().add(vowels["u"].clone())

	return vowels, nil
}

func getInitials(base *svgparser.Element) (map[string]*cluster, error) {
	clusters := []string{
		"t", "w", "p",
		"tp", "th", "tr", "tw", "tk", "hr", "hw",
		"twp", "tph", "tpr", "twh", "tkw", "tkh",
		"twph",
		"ktwph", "twprh", "ktwprh",
	}

	elem := findElem(base, "initial", "full")
	initial, err := getClusters(elem, append(clusters, "null", "s", "*"))
	if err != nil {
		return nil, err
	}

	elem = findElem(base, "initial", "2/3")
	initial2_3, err := getClusters(elem, clusters)
	if err != nil {
		return nil, err
	}

	/* not yet ready to use this

	elem = findElem(base, "initial", "1/2")
	initial1_2, err := getClusters(elem, []string{
		"t", "tp", "th", "tph", "twh", "twph", "ktwph", "ktwprh", "twprh", "tr", "tw", "twp",
	})
	if err != nil {
		return nil, err
	}

	*/

	// TODO: get rid of all the 300, 500 stuff...
	// TODO: check if the derived consonants exist before replacing them, that way
	//       if you want to customize the clusters you can, but you don't need to
	// TODO: I wonder if a little dsl for describing how to make derived clusters would be useful/ easy to build

	// add derived initial consonant clusters
	initial["h"] = initial["t"].clone().origin(300, 500).mirriorY()
	initial["r"] = initial["t"].clone().origin(300, 500).mirrorXY()
	initial["k"] = initial["t"].clone().origin(300, 500).mirriorX()

	initial["hk"] = initial["tr"].clone().origin(300, 500).mirriorY()
	initial["kr"] = initial["th"].clone().origin(300, 500).mirriorX()
	initial["kw"] = initial["tp"].clone().origin(300, 500).mirriorX()
	initial["*k"] = initial2_3["t"].clone().origin(300, 500).mirriorX().translate(0, -140).add(initial["*"].clone())
	initial["*h"] = initial2_3["t"].clone().origin(300, 500).mirriorY().add(initial["*"].clone())
	initial["*r"] = initial2_3["t"].clone().origin(300, 500).mirrorXY().add(initial["*"].clone())
	initial["*w"] = initial2_3["w"].clone().translate(0, -140).add(initial["*"].clone())
	initial["rw"] = initial["tp"].clone().origin(300, 500).mirrorXY()

	initial["hkw"] = initial["tpr"].clone().origin(300, 500).mirriorX()
	initial["hrw"] = initial["tkw"].clone().origin(300, 500).mirriorY()
	initial["hkr"] = initial["tkh"].clone().origin(300, 500).mirrorXY()
	initial["krw"] = initial["tph"].clone().origin(300, 500).mirriorX()
	initial["*kr"] = initial2_3["th"].clone().origin(300, 500).mirriorX().translate(0, -140).add(initial["*"].clone())
	initial["*kw"] = initial2_3["tp"].clone().origin(300, 500).mirriorX().translate(0, -140).add(initial["*"].clone())
	initial["*rw"] = initial2_3["tp"].clone().origin(300, 500).mirrorXY().add(initial["*"].clone())
	initial["*hr"] = initial2_3["hr"].clone().translate(0, -140).add(initial["*"].clone())
	initial["*hw"] = initial2_3["hw"].clone().translate(0, -140).add(initial["*"].clone())
	initial["*hk"] = initial2_3["tr"].clone().origin(300, 500).mirriorY().add(initial["*"].clone())

	initial["*krw"] = initial2_3["tph"].clone().origin(300, 500).mirriorX().translate(0, -140).add(initial["*"].clone())
	initial["*hrw"] = initial2_3["tkw"].clone().origin(300, 500).mirriorY().add(initial["*"].clone())
	initial["*hkr"] = initial2_3["tkh"].clone().origin(300, 500).mirrorXY().add(initial["*"].clone())
	initial["*hkw"] = initial2_3["tpr"].clone().origin(300, 500).mirriorX().translate(0, -140).add(initial["*"].clone())

	return initial, nil
}

// getClusters reads through the given svg base element looking for path elements labeled with the names
// provided in the clusters list. These paths are then converted to clusters and mapped to their labeled names
func getClusters(base *svgparser.Element, clusters []string) (map[string]*cluster, error) {
	clusterMap := map[string]*cluster{}

	for _, k := range clusters {
		elem := findElem(base, k)
		if elem == nil {
			return nil, fmt.Errorf("failed to find cluster %s", k)
		}

		p, err := strokeFromElem(elem)
		if err != nil {
			fmt.Println(elem.Attributes)
			return nil, fmt.Errorf("failed to creates stroke for %s: %w", k, err)
		}

		c := cluster([]stroke{*p})
		clusterMap[k] = &c
	}

	return clusterMap, nil
}

// findElem searches the children of elem until it finds the element with the first name in the names list
// it then searches that elements child searching for the next name in the list, this continues until
// the full list of names has been searched and the final *svgparser.Element is returned
func findElem(elem *svgparser.Element, names ...string) *svgparser.Element {
	base := elem
	for _, name := range names {
		for _, e := range base.Children {
			if e.Attributes["label"] == name {
				base = e
				break
			}
		}

		if base == nil {
			return nil
		}
	}

	return base
}
