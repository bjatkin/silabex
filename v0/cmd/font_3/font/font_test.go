package font

import (
	"os"
	"strings"
	"testing"

	"github.com/JoshVarga/svgparser"
)

func Test_getVowels(t *testing.T) {
	// TODO: move this into a test dir
	f, err := os.Open("/home/brandon/src/silabex/examples/basic_font_2.svg")
	if err != nil {
		t.Fatal("failed to open svg file", err)
	}

	svg, err := svgparser.Parse(f, false)
	if err != nil {
		t.Fatal("failed to parse svg file", err)
	}

	got, err := getVowels(findElem(svg, "v3"))
	if err != nil {
		t.Fatal("failed to get vowels", err)
	}

	// TODO: read this from a test file too...
	want := []string{
		"AOEU",
		"AOE",
		"AOU",
		"AO",
		"AEU",
		"AE",
		"AU",
		"A",
		"EOU",
		"EO",
		"OU",
		"O",
		"EU",
		"E",
		"U",
	}
	for _, name := range want {
		if got[strings.ToLower(name)] == nil {
			t.Error("missing vowel ", name)
		}
	}
}

func Test_getInitials(t *testing.T) {
	// TODO: move this into a test dir
	f, err := os.Open("/home/brandon/src/silabex/examples/basic_font_2.svg")
	if err != nil {
		t.Fatal("failed to open svg file", err)
	}

	svg, err := svgparser.Parse(f, false)
	if err != nil {
		t.Fatal("failed to parse svg file", err)
	}

	got, err := getInitials(findElem(svg, "v3"))
	if err != nil {
		t.Fatal("failed to get vowels", err)
	}

	rawClusterList, err := os.ReadFile("/home/brandon/src/silabex/left.txt")
	if err != nil {
		t.Fatal("failed to get left", err)
	}

	want := strings.Split(string(rawClusterList), "\n")
	for _, name := range want {
		if name == "" {
			continue
		}

		if got[strings.ToLower(name)] == nil {
			t.Error("missing initial cluster", name)
		}
	}
}
