package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	words := []string{"tiny", "moon"}

	raw, err := os.ReadFile("cmd/translate/dict.json")
	if err != nil {
		fmt.Println("dict read err:", err)
		return
	}

	dict := map[string]string{}
	err = json.Unmarshal(raw, &dict)
	if err != nil {
		fmt.Println("dict unmarshal err:", err)
		return
	}

	raw, err = os.ReadFile("cmd/translate/strokes.json")
	if err != nil {
		fmt.Println("stroke read err:", err)
		return
	}

	strokes := Strokes{}
	err = json.Unmarshal(raw, &strokes)
	if err != nil {
		fmt.Println("strokes unmarshal err:", err)
		return
	}

	for _, word := range words {
		grid, ok := dict[word]
		if !ok {
			fmt.Println("unknown word:", word)
			return
		}

		fmt.Printf("%s |%s| %+v\n", word, grid, split(grid))
	}
}

/*

Vowels
 |----2----|
 |         |
 0   (4)   3
 |         |
 |----1----|

Points
 _:0---S:1     R:0---F:1

 K:2---T:3     B:2---P:3
  |     |       |     |
 W:4---P:5     G:4---L:5
  |     |       |     |
 R:6---H:7     S:6---T:7

 *:8---_:9     Z:8---D:9

Segments
 -----0-----

 |----1----|
 2----3----4
 |----5----|

 -----6-----

~~~

 ---\-7-\---

 |---------|
 |--\-8-\--|
 |---------|

 ---\-9-\---

~~~

 -----------

 |---------|
 |---(A)---|
 |---------|

 -----------
*/

var (
	initial = "_SKTWPRH*"
	vowel   = "AOEU"
	final   = "RFBPGLSTZD"
)

func split(sPattern string) []Pattern {
	patterns := []Pattern{}

	for _, p := range strings.Split(sPattern, "/") {
		pattern := Pattern{}
		i := 0

		for ; i < len(p); i++ {
			idx := strings.Index(initial, string(p[i]))
			if idx < 0 {
				break
			}
			pattern.Initial = append(pattern.Initial, idx)
		}
		for ; i < len(p); i++ {
			idx := strings.Index(vowel, string(p[i]))
			if idx < 0 {
				break
			}
			pattern.Vowel = append(pattern.Vowel, idx)
		}
		for ; i < len(p); i++ {
			idx := strings.Index(final, string(p[i]))
			if idx < 0 {
				break
			}
			pattern.Final = append(pattern.Final, idx)
		}

		patterns = append(patterns, pattern)
	}

	return patterns
}

type Pattern struct {
	Initial []int
	Vowel   []int
	Final   []int
}

type Strokes struct {
	Initial map[string]string
	Final   map[string]string
}
