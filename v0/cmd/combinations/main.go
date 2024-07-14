package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	left := []rune("STPHKWR*")
	leftNames := []string{}
	for _, c := range combos(left) {
		leftNames = append(leftNames, string(sortBy(c, left)))
	}

	slices.Sort(leftNames)
	err := os.WriteFile("left.txt", []byte(strings.Join(leftNames, "\n")), 0o0655)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	right := []rune("RFBPGLSTZD")
	rightNames := []string{}
	for _, c := range combos(right) {
		rightNames = append(rightNames, string(sortBy(c, right)))
	}

	slices.Sort(rightNames)
	err = os.WriteFile("right.txt", []byte(strings.Join(rightNames, "\n")), 0o0655)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	vowel := []rune("AOEU")
	vowelNames := []string{}
	for _, c := range combos(vowel) {
		vowelNames = append(vowelNames, string(sortBy(c, vowel)))
	}

	slices.Sort(vowelNames)
	err = os.WriteFile("vowel.txt", []byte(strings.Join(vowelNames, "\n")), 0o0655)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}

func sortBy[T comparable](value, reference []T) []T {
	slices.SortFunc(value, func(a, b T) int {
		if slices.Index(reference, a) > slices.Index(reference, b) {
			return 1
		}
		return -1
	})

	return value
}

func bestMatch[T comparable](keys [][]T, check []T) []T {
	best := int(^uint(0) >> 1)
	bestKey := []T{}
	for _, key := range keys {
		remain, ok := match(key, check)
		if !ok {
			continue
		}

		if len(remain) == 0 {
			return key
		}

		if len(remain) < best {
			best = len(remain)
			bestKey = key
		}
	}

	return bestKey
}

func match[T comparable](key, check []T) ([]T, bool) {
	for _, c := range check {
		if slices.Contains(key, c) {
			continue
		}
		// not a match because some values in check can not be found in the key
		return nil, false
	}

	remainder := []T{}
	for _, k := range key {
		if slices.Contains(check, k) {
			continue
		}
		remainder = append(remainder, k)
	}

	return remainder, true
}

func inverseIntersection[T comparable](left, right []T) []T {
	inverse := []T{}
	for _, a := range left {
		if slices.Contains(right, a) {
			continue
		}
		inverse = append(inverse, a)
	}

	for _, b := range right {
		if slices.Contains(left, b) {
			continue
		}
		inverse = append(inverse, b)
	}

	return inverse
}

func combos[T cmp.Ordered](set []T) [][]T {
	if len(set) <= 0 {
		return [][]T{set}
	}

	car := set[0]
	cdr := set[1:]
	combos := combos(cdr)

	ret := make([][]T, len(combos)*2)
	for i := range combos {
		s := []T{}
		s = append(s, combos[i]...)
		s = append(s, car)

		ret[i] = s
		slices.Sort(ret[i])
		ret[i+len(combos)] = combos[i]
		slices.Sort(ret[i+len(combos)])
	}

	return ret
}
