package derive

import (
	"reflect"
	"testing"
)

func Test_combinations(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"[0, 1, 2, 3]",
			args{
				names: []string{"0", "1", "2", "3"},
			},
			[]string{
				"0", "1", "01", "2", "02", "12", "012", "3", "03", "13", "013", "23", "023", "123", "0123",
			},
		},
		{
			"[AZ, CD, BB]",
			args{
				names: []string{"AZ", "CD", "BB"},
			},
			[]string{
				"AZ", "CD", "AZCD", "BB", "AZBB", "CDBB", "AZCDBB",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combinations(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combinations() = %v, want %v", got, tt.want)
			}
		})
	}
}
