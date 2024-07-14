package main

import (
	"reflect"
	"testing"
)

func Test_sortBy(t *testing.T) {
	type args struct {
		value     []rune
		reference []rune
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{
			"sort all vowels",
			args{
				value:     []rune("UOAE"),
				reference: []rune("AOEU"),
			},
			[]rune("AOEU"),
		},
		{
			"sort initial clusters",
			args{
				value:     []rune("*KPTS"),
				reference: []rune("SKTWPRH*"),
			},
			[]rune("SKTP*"),
		},
		{
			"sort final clusters",
			args{
				value:     []rune("BSFDLZ"),
				reference: []rune("RFBPGLSTZD*"),
			},
			[]rune("FBLSZD"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortBy(tt.args.value, tt.args.reference); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortBy() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
