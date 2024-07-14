package linalg

import (
	"fmt"
	"math"
)

// Mat3x is a 3x3 floating point matrix
type Mat3x struct {
	Data [3][3]float64
}

// String returns the string representation of a Mat3x
func (m Mat3x) String() string {
	return fmt.Sprintf(
		"[%.2f, %.2f, %.2f]\n[%.2f, %.2f, %.2f]\n[%.2f, %.2f, %.2f]",
		m.Data[0][0], m.Data[0][1], m.Data[0][2],
		m.Data[1][0], m.Data[1][1], m.Data[2][2],
		m.Data[2][0], m.Data[2][1], m.Data[1][2],
	)
}

// Identity returns the identity Mat3x
func Identity() Mat3x {
	return Mat3x{
		Data: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
	}
}

// Scale returns a matrix that scales the x and y axes by sx and sy
func Scale(sx, sy float64) Mat3x {
	return Mat3x{
		Data: [3][3]float64{
			{sx, 0, 0},
			{0, sy, 0},
			{0, 0, 1},
		},
	}
}

// Rotate returns a matrix that rotates the x and y axes by angle
// angle is in radians
func Rotate(angle float64) Mat3x {
	return Mat3x{
		Data: [3][3]float64{
			{math.Cos(angle), -math.Sin(angle), 0},
			{math.Sin(angle), math.Cos(angle), 0},
			{0, 0, 1},
		},
	}
}

// Translate returns a matrix that translates the x and y axes by tx and ty
func Translate(tx, ty float64) Mat3x {
	return Mat3x{
		Data: [3][3]float64{
			{1, 0, tx},
			{0, 1, ty},
			{0, 0, 1},
		},
	}
}

// MatMul returns the result of multiplying two Mat3x
func MatMul(a, b Mat3x) Mat3x {
	ret := Mat3x{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				ret.Data[i][j] += a.Data[i][k] * b.Data[k][j]
			}
		}
	}

	return ret
}

// MatAdd returns the result of adding two Mat3x
func MatAdd(a, b Mat3x) Mat3x {
	ret := Mat3x{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			ret.Data[i][j] = a.Data[i][j] + b.Data[i][j]
		}
	}
	return ret
}

// Transform combines a series of transformations into a single matrix
func Transform(transforms ...Mat3x) Mat3x {
	ret := Identity()
	for _, tansform := range transforms {
		ret = MatMul(ret, tansform)
	}

	return ret
}
