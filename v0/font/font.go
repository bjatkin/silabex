package font

import (
	"fmt"
	"os"
	"strings"

	"github.com/JoshVarga/svgparser"
	"github.com/bjatkins/silabex/linalg"
)

type Font struct {
	metadata *metadata
}

func New(svgPath, derivedPath string) (*Font, error) {
	metadata, err := newMetadata(svgPath)
	if err != nil {
		return nil, err
	}

	clusterBuilder, err := newClusterBuilder(derivedPath)
	if err != nil {
		return nil, err
	}

	for _, initial := range clusterBuilder.initial {
		name := initial.name
		if _, ok := metadata.initial.full[name]; !ok {
			fmt.Println("assign full initial cluster", name)
			metadata.initial.full[name] = initial.build(metadata.initial.full)
		}
		if _, ok := metadata.initial.half[name]; !ok {
			fmt.Println("assign half initial cluster", name)
			metadata.initial.half[name] = initial.build(metadata.initial.half)
		}
		if _, ok := metadata.initial.twoThrids[name]; !ok {
			fmt.Println("assign twoThirds initial cluster", name)
			metadata.initial.twoThrids[name] = initial.build(metadata.initial.twoThrids)
		}

		if _, ok := metadata.solo.full[name]; !ok {
			fmt.Println("assign full solo cluster", name)
			metadata.solo.full[name] = initial.build(metadata.solo.full)
		}
		if _, ok := metadata.solo.half[name]; !ok {
			fmt.Println("assign half solo cluster", name)
			metadata.solo.half[name] = initial.build(metadata.solo.half)
		}
		if _, ok := metadata.solo.twoThrids[name]; !ok {
			fmt.Println("assign twoThirds solo cluster", name)
			metadata.solo.twoThrids[name] = initial.build(metadata.solo.twoThrids)
		}
	}

	for _, vowel := range clusterBuilder.vowel {
		if _, ok := metadata.vowel[vowel.name]; !ok {
			metadata.vowel[vowel.name] = vowel.build(metadata.vowel)
			fmt.Println("assign vowel cluster", vowel.name)
		}
	}

	for _, final := range clusterBuilder.final {
		if _, ok := metadata.final.full[final.name]; !ok {
			fmt.Println("assign full final cluster", final.name)
			metadata.final.full[final.name] = final.build(metadata.final.full)
		}
		if _, ok := metadata.final.half[final.name]; !ok {
			fmt.Println("assign half final cluster", final.name)
			metadata.final.half[final.name] = final.build(metadata.final.half)
		}
		if _, ok := metadata.final.twoThrids[final.name]; !ok {
			fmt.Println("assign twoThirds final cluster", final.name)
			metadata.final.twoThrids[final.name] = final.build(metadata.final.twoThrids)
		}
	}

	return &Font{
		metadata: metadata,
	}, nil
}

func (f *Font) NewChars(word string) ([]Char, error) {
	validateCluster := func(s, check string) error {
		for _, r := range s {
			if !strings.ContainsRune(check, r) {
				return fmt.Errorf("invalid character %s in cluster", string(r))
			}
		}

		return nil
	}

	chars := []Char{}
	sylables := strings.Split(word, "/")
	for _, sylable := range sylables {
		initial, vowel, final := splitSylable(sylable)
		err := validateCluster(initial, "SKTWPRH*")
		if err != nil {
			return nil, fmt.Errorf("invalid initial cluster %s in sylable %s", initial, sylable)
		}

		err = validateCluster(vowel, "AOEU")
		if err != nil {
			return nil, fmt.Errorf("invalid vowel cluster %s in sylable %s", vowel, sylable)
		}

		err = validateCluster(final, "RFBPGLSTZD")
		if err != nil {
			return nil, fmt.Errorf("invalid final cluster %s in sylable %s", final, sylable)
		}

		if final == "" {
			initialCluster := f.metadata.solo.full[initial]
			vowelCluster := f.metadata.vowel[vowel]
			chars = append(chars, Char{
				name:    sylable,
				initial: initialCluster,
				vowel:   vowelCluster,
			})
		} else {
			initialCluster := f.metadata.initial.full[initial]
			vowelCluster := f.metadata.vowel[vowel]
			finalCluster := f.metadata.final.full[final]
			chars = append(chars, Char{
				name:    sylable,
				initial: initialCluster,
				vowel:   vowelCluster,
				final:   finalCluster,
			})
		}
	}

	return chars, nil
}

func splitSylable(sylable string) (initial, vowel, final string) {
	isVowel := func(r rune) bool {
		return strings.ContainsRune("AOEU", r)
	}

	isSet := func(i int) bool {
		return i > -1
	}

	vowelStart, vowelEnd := -1, -1
	for i, r := range sylable {
		if isVowel(r) && isSet(vowelStart) {
			vowelEnd = i + 1
		}
		if isVowel(r) && !isSet(vowelStart) {
			vowelStart = i
			vowelEnd = i + 1
		}
	}

	if vowelStart == -1 {
		return sylable, "", ""
	}

	return sylable[:vowelStart],
		sylable[vowelStart:vowelEnd],
		sylable[vowelEnd:]
}

// stroke represents an svg stroke
type stroke struct {
	slot clusterSlot

	dx, dy float64
	sx, sy float64

	path string
}

// strokeFromElem creates a new path from an *svgparser.Element. The element must be a
// path element or an error will be returned
func strokeFromElem(slot clusterSlot, elem *svgparser.Element) (*stroke, error) {
	if elem.Name != "path" {
		return nil, fmt.Errorf("can not create stroke from %s element", elem.Name)
	}

	return &stroke{
		slot: slot,
		path: elem.Attributes["d"],
		sx:   1,
		sy:   1,
	}, nil
}

// svg renders the stroke as a valid svg string
func (s *stroke) svg() string {
	var origin linalg.Vec3
	switch s.slot {
	case initialSlot:
		origin = linalg.NewPoint2(300, 500)
	case finalSlot:
		origin = linalg.NewPoint2(700, 500)
	default:
		origin = linalg.NewPoint2(500, 500)
	}

	preScale := linalg.Translate(origin.X, origin.Y)
	scale := linalg.Scale(s.sx, s.sy)
	posScale := linalg.Translate(-origin.X, -origin.Y)

	translate := linalg.Translate(s.dx, s.dy)

	transform := linalg.Transform(preScale, scale, posScale, translate)

	transformString := fmt.Sprintf(
		"transform=\"matrix(%.2f %.2f %.2f %.2f %.2f %.2f)\"",
		transform.Data[0][0], transform.Data[1][0],
		transform.Data[0][1], transform.Data[1][1],
		transform.Data[0][2], transform.Data[1][2],
	)

	// no need to add the transform if it is identical to the identity matrix
	if transformString == "transform=\"matrix(1.00 0.00 0.00 1.00 0.00 0.00)\"" {
		return fmt.Sprintf("<path d=\"%s\"/>", s.path)
	}

	return fmt.Sprintf("<path %s d=\"%s\"/>", transformString, s.path)
}

type clusterSlot int

const (
	soloSlot clusterSlot = iota
	initialSlot
	vowelSlot
	finalSlot
)

// cluster is a slice of paths that makes up a silbex stroke cluster
type cluster struct {
	slot    clusterSlot
	strokes []stroke
}

func mergeClusters(clusters ...*cluster) *cluster {
	if len(clusters) == 0 {
		return nil
	}

	cluster := &cluster{
		slot: clusters[0].slot,
	}
	for _, c := range clusters {
		cluster.strokes = append(cluster.strokes, c.strokes...)
	}

	return cluster
}

func (c *cluster) clone() *cluster {
	ret := &cluster{
		slot: c.slot,
	}

	ret.strokes = append(ret.strokes, c.strokes...)
	return ret
}

// svg returns the cluster as a valid svg string
func (c *cluster) svg() string {
	if c == nil {
		return ""
	}

	ret := ""
	for _, stroke := range c.strokes {
		ret += stroke.svg()
	}
	return ret
}

// translate translates all the strokes that make up the cluster by dx, dy
func (c *cluster) translate(dx, dy float64) *cluster {
	for i := range c.strokes {
		c.strokes[i].dx += dx
		c.strokes[i].dy += dy
	}

	return c
}

// mirriorX flips all the strokes that make up the cluster along the x axis
func (c *cluster) mirriorX() *cluster {
	for i := range c.strokes {
		c.strokes[i].sx = -1
	}

	return c
}

// mirriorY flips all the strokes that make up the cluster along the y axis
func (c *cluster) mirriorY() *cluster {
	for i := range c.strokes {
		c.strokes[i].sy = -1
	}

	return c
}

// Char is the representation of single rune in silbex
// TODO: char should support [4]byte encoding so it works similarly to utf-8
type Char struct {
	name    string
	initial *cluster
	vowel   *cluster
	final   *cluster
}

func (c *Char) Name() string {
	return c.name
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

// metadata contains font metadata for constructing a full font
type metadata struct {
	vowel   map[string]*cluster
	solo    *clusterGroup
	initial *clusterGroup
	final   *clusterGroup
}

func newMetadata(svgFile string) (*metadata, error) {
	f, err := os.Open(svgFile)
	if err != nil {
		return nil, err
	}

	svg, err := svgparser.Parse(f, true)
	if err != nil {
		return nil, err
	}

	// TODO: this should get removed eventually
	v3 := findElem(svg, "v3")

	vowel, err := getVowels(v3)
	if err != nil {
		return nil, err
	}

	solo, err := getSolos(v3)
	if err != nil {
		return nil, err
	}

	initial, err := getInitials(v3)
	if err != nil {
		return nil, err
	}

	final, err := getFinals(v3)
	if err != nil {
		return nil, err
	}

	return &metadata{
		vowel:   vowel,
		solo:    solo,
		initial: initial,
		final:   final,
	}, nil
}

// getVowels reads all the vowels from the "vowel" group and returns them in a map
func getVowels(base *svgparser.Element) (map[string]*cluster, error) {
	elem := findElem(base, "vowel")

	vowels, err := getClusters(vowelSlot, elem, []string{"A", "AO", "E", "EU", "AOEU"})
	if err != nil {
		return nil, err
	}

	return vowels, nil
}

// getInitials reads all the initial consonant clusters from the "initial" group and
// returns them as a clusterGroup
func getInitials(base *svgparser.Element) (*clusterGroup, error) {
	return getInitialsOrSolos(initialSlot, base)
}

// getInitials reads all the initial consonant clusters from the "solo" group and
// returns them as a clusterGroup
func getSolos(base *svgparser.Element) (*clusterGroup, error) {
	return getInitialsOrSolos(soloSlot, base)
}

// getInitialsOrSolos reads all the initial consonant clusters from the group specified by the
// provided group param and returns them as a clusterGroup
func getInitialsOrSolos(group clusterSlot, base *svgparser.Element) (*clusterGroup, error) {
	groupName := "initial"
	if group == soloSlot {
		groupName = "solo"
	}

	clusters := []string{
		"T", "W", "P",
		"TP", "TH", "TR", "TW", "KT", "RH", "WH",
		"TWP", "TPH", "TPR", "TWH", "KTW", "KTH",
		"TWPH",
		"KTWPH", "TWPRH",
		"KTWPRH",
	}

	elem := findElem(base, groupName, "full")
	full, err := getClusters(group, elem, append(clusters, "NULL", "S", "*"))
	if err != nil {
		return nil, err
	}

	elem = findElem(base, groupName, "2/3")
	twoThrids, err := getClusters(group, elem, clusters)
	if err != nil {
		return nil, err
	}

	elem = findElem(base, groupName, "1/2")
	half, err := getClusters(group, elem, clusters)
	if err != nil {
		return nil, err
	}

	return &clusterGroup{
		full:      full,
		twoThrids: twoThrids,
		half:      half,
	}, nil
}

// getFinals reads all the final consonant clusters from the "final" group
// and returns them as a clusterGroup
func getFinals(base *svgparser.Element) (*clusterGroup, error) {
	clusters := []string{"R", "F"}

	elem := findElem(base, "final", "full")
	full, err := getClusters(finalSlot, elem, clusters)
	if err != nil {
		return nil, err
	}

	return &clusterGroup{
		full: full,
	}, nil
}

// getClusters reads through the given svg base element looking for path elements labeled with the names
// provided in the clusters list. These paths are then converted to clusters and mapped to their labeled names
func getClusters(slot clusterSlot, base *svgparser.Element, clusters []string) (map[string]*cluster, error) {
	clusterMap := map[string]*cluster{}

	for _, k := range clusters {
		elem := findElem(base, k)
		if elem == nil {
			return nil, fmt.Errorf("failed to find cluster %s", k)
		}

		p, err := strokeFromElem(slot, elem)
		if err != nil {
			return nil, fmt.Errorf("failed to create stroke for %s: %w", k, err)
		}

		clusterMap[k] = &cluster{
			slot:    slot,
			strokes: []stroke{*p},
		}
	}

	return clusterMap, nil
}

// findElem searches the children of elem until it finds the element with the first name in the names list
// it then searches that elements child searching for the next name in the list, this continues until
// the full list of names has been searched and the final *svgparser.Element is returned
func findElem(elem *svgparser.Element, names ...string) *svgparser.Element {
	fmt.Printf("searching base for %v\n", names)
	base := elem

	for _, name := range names {
		var found bool
		if base == nil {
			fmt.Println("base is nil")
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

// clusterGroup contains all the size classes for a set of clusters
type clusterGroup struct {
	full      map[string]*cluster
	twoThrids map[string]*cluster
	half      map[string]*cluster
}
