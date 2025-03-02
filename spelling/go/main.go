package main

import "fmt"

var phonemes = []string{"W", "Y", "L", "R", "M", "N", "NG", "DH", "F", "V", "TH", "S", "Z", "SH", "ZH", "HH", "B", "D", "G", "K", "P", "T", "CH", "JH"}
var patterns = []string{"1 2", "1 3", "1 4", "1 5", "1 6", "1 7", "1 8", "2 3", "2 4", "2 5", "2 6", "2 7", "2 8", "3 4", "3 5", "3 6", "3 7", "3 8", "4 5", "4 6", "4 7", "4 8", "5 6", "5 7", "5 8", "6 7", "6 8", "7 8"}

func main() {
	max := 7
	phonemes = phonemes[:max]
	patterns = patterns[:max]
	fmt.Printf("%#v\n", phonemes)
	fmt.Printf("%#v\n", patterns)
	got := build_alphabet(phonemes, patterns)
	fmt.Printf("len(%d)\n", len(got))
	fmt.Printf("got[0]:%#v\n", got[0])
}

func build_alphabet(phonemes []string, patterns []string) []map[string]string {
	alphabets := []map[string]string{}

	for _, p := range phonemes {
		for _, t := range patterns {
			new_phonemes := filter(phonemes, p)
			new_patterns := filter(patterns, t)
			built := build_alphabet(new_phonemes, new_patterns)

			for _, b := range built {
				new_b := copyMap(b)
				new_b[p] = t
				alphabets = append(alphabets, new_b)
			}

			if len(built) == 0 {
				alphabets = append(alphabets, map[string]string{p: t})
			}
		}
	}

	return alphabets
}

func copyMap(m map[string]string) map[string]string {
	new := map[string]string{}
	for k, v := range m {
		new[k] = v
	}
	return new
}

func filter(list []string, remove string) []string {
	new := []string{}
	for _, l := range list {
		if l == remove {
			continue
		}

		new = append(new, l)
	}

	return new
}
