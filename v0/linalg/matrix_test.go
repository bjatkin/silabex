package linalg

import (
	"math"
	"testing"
)

func TestMatMul(t *testing.T) {
	type args struct {
		a Mat3x
		b Mat3x
	}
	tests := []struct {
		name string
		args args
		want Mat3x
	}{
		{
			"scale then translate",
			args{
				a: Scale(1, -1),
				b: Translate(5, 10),
			},
			Mat3x{
				Data: [3][3]float64{
					{1, 0, 5},
					{0, -1, -10},
					{0, 0, 1},
				},
			},
		},
		{
			"rotate then translate",
			args{
				a: Rotate(math.Pi / 4),
				b: Translate(5, 10),
			},
			Mat3x{
				Data: [3][3]float64{
					{math.Sqrt2 / 2, -math.Sqrt2 / 2, -3.535},
					{math.Sqrt2 / 2, math.Sqrt2 / 2, 10.605},
					{0, 0, 1},
				},
			},
		},
		{
			"itendity",
			args{
				a: Translate(0, 0),
				b: Scale(1, 1),
			},
			Mat3x{
				Data: [3][3]float64{
					{1, 0, 0},
					{0, 1, 0},
					{0, 0, 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatMul(tt.args.a, tt.args.b); !matEqual(got, tt.want, 0.01) {
				t.Errorf("MatMul() = \n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

// matEqual checks if two matrices are equal within a given error margin e
func matEqual(a, b Mat3x, e float64) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(a.Data[i][j]-b.Data[i][j]) > e {
				return false
			}
		}
	}

	return true
}
