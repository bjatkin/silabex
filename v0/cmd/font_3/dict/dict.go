package dict

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Rune struct {
	Initial string
	Vowel   string
	Final   string
}

type Dict map[string][]Rune

func NewDict(file string) (Dict, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	jsonDict := map[string]string{}
	err = json.Unmarshal(raw, &jsonDict)
	if err != nil {
		return nil, err
	}

	dict := Dict{}
	for k, v := range jsonDict {
		runes := []Rune{}
		for _, sylable := range strings.Split(v, "/") {
			r, err := Runes(sylable)
			if err != nil {
				return nil, err
			}
			runes = append(runes, r)
		}

		dict[k] = runes
	}

	return dict, nil
}

func Runes(sylable string) (Rune, error) {
	runes := []rune(sylable)

	const initial = "SKTWPRH*"
	initialCount := 0
	for i := 0; i < len(sylable); i++ {
		if !strings.Contains(initial, string(runes[i])) {
			break
		}
		initialCount++
	}

	const vowels = "AEOU"
	vowelCount := initialCount
	for i := initialCount; i < len(sylable); i++ {
		if !strings.Contains(vowels, string(runes[i])) {
			break
		}
		vowelCount++
	}

	const final = "RFBPGLSTZD"
	finalCount := vowelCount
	for i := vowelCount; i < len(sylable); i++ {
		if !strings.Contains(final, string(runes[i])) {
			break
		}
		finalCount++
	}

	if len(sylable) != finalCount {
		return Rune{}, errors.New("unmatched rune " + string(sylable[finalCount]) + " in sylabel:" + sylable)
	}

	return Rune{
		Initial: string(runes[:initialCount]),
		Vowel:   string(runes[initialCount:vowelCount]),
		Final:   string(runes[vowelCount:]),
	}, nil
}
