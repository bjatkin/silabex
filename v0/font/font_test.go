package font

import (
	"os"
	"strings"
	"testing"
)

func Test_splitSylable(t *testing.T) {
	type args struct {
		sylable string
	}
	tests := []struct {
		name        string
		args        args
		wantInitial string
		wantVowel   string
		wantFinal   string
	}{
		{
			"golden",
			args{"TAOEUPB"},
			"T", "AOEU", "PB",
		},
		{
			"missing initial consonant",
			args{"OULT"},
			"", "OU", "LT",
		},
		{
			"missing vowel",
			args{"KT"},
			"KT", "", "",
		},
		{
			"missing final consonant",
			args{"TPHROE"},
			"TPHR", "OE", "",
		},
		{
			"ends with a single vowel",
			args{"TO"},
			"T", "O", "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInitial, gotVowel, gotFinal := splitSylable(tt.args.sylable)
			if gotInitial != tt.wantInitial {
				t.Errorf("splitSylable() gotInitial = %v, want %v", gotInitial, tt.wantInitial)
			}
			if gotVowel != tt.wantVowel {
				t.Errorf("splitSylable() gotVowel = %v, want %v", gotVowel, tt.wantVowel)
			}
			if gotFinal != tt.wantFinal {
				t.Errorf("splitSylable() gotFinal = %v, want %v", gotFinal, tt.wantFinal)
			}
		})
	}
}

func Test_New(t *testing.T) {
	font, err := New("testdata/font.svg", "derived")
	if err != nil {
		t.Fatal("failed to build font", err)
	}

	// check for all vowel clusters
	vowelsRaw, err := os.ReadFile("testdata/vowel.txt")
	if err != nil {
		t.Fatal("failed to read vowel clusters data", err)
	}

	for _, vowel := range strings.Split(string(vowelsRaw), "\n") {
		if vowel == "" {
			continue
		}

		if _, ok := font.metadata.vowel[vowel]; !ok {
			t.Errorf("missing vowel cluster %s", vowel)
		}
	}

	// check for all initial clusters
	initialRaw, err := os.ReadFile("testdata/initial.txt")
	if err != nil {
		t.Fatal("failed to read initial clusters data", err)
	}

	for _, initial := range strings.Split(string(initialRaw), "\n") {
		if initial == "" {
			continue
		}

		if _, ok := font.metadata.initial.full[initial]; !ok {
			t.Errorf("missing initial cluster %s", initial)
		}
	}
}
