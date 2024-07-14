package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Cluster [10]rune

func NewCluster(sylable []rune, start, end int) Cluster {
	cluster := Cluster{}
	for i := start; i < end; i++ {
		cluster[i-start] = sylable[i]
	}

	return cluster
}

type Symbol struct {
	Initial Cluster
	Vowel   Cluster
	Final   Cluster
}

func NewSymbol(sylable string) (Symbol, error) {
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
		return Symbol{}, errors.New("unmatched rune " + string(sylable[finalCount]) + " in sylabel:" + sylable)
	}

	return Symbol{
		Initial: NewCluster(runes, 0, initialCount),
		Vowel:   NewCluster(runes, initialCount, vowelCount),
		Final:   NewCluster(runes, vowelCount, finalCount),
	}, nil
}

func (s Symbol) String() string {
	initial := ""
	for _, r := range s.Initial {
		if r == 0 {
			continue
		}

		initial += string(r)
	}

	vowel := ""
	for _, r := range s.Vowel {
		if r == 0 {
			continue
		}

		vowel += string(r)
	}

	final := ""
	for _, r := range s.Final {
		if r == 0 {
			continue
		}
		final += string(r)
	}

	segments := []string{}
	for _, s := range []string{initial, vowel, final} {
		if len(s) > 0 {
			segments = append(segments, s)
		}
	}

	return strings.Join(segments, "|")
}

type Dict struct {
	rawWords map[string]string
	// TODO: add an LRU word cache here
}

func NewDict(file string) (*Dict, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	rawWords := map[string]string{}
	err = json.Unmarshal(raw, &rawWords)
	if err != nil {
		return nil, err
	}

	return &Dict{
		rawWords: rawWords,
	}, nil
}

func (d *Dict) Lookup(word string) ([]Symbol, error) {
	spelling, ok := d.rawWords[word]
	if !ok {
		return nil, errors.New("Unknown word: " + word)
	}

	sylables := strings.Split(spelling, "/")
	symboles := []Symbol{}
	for _, sylable := range sylables {
		got, err := NewSymbol(sylable)
		if err != nil {
			return nil, err
		}
		symboles = append(symboles, got)
	}

	return symboles, nil
}
