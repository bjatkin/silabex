package linalg

// Vec3 is a 3-dimensional vector
type Vec3 struct {
	X, Y, Z float64
}

// NewPoint2 returns a new Vec3 that repesents a 2d point located at x, y
func NewPoint2(x, y float64) Vec3 {
	return Vec3{
		X: x,
		Y: y,
		Z: 1,
	}
}

// VecMul returns the result of multiplying the Mat3x m by the Vec3 v
func VecMul(m Mat3x, v Vec3) Vec3 {
	return Vec3{
		X: m.Data[0][0]*v.X + m.Data[0][1]*v.Y,
		Y: m.Data[1][0]*v.X + m.Data[1][1]*v.Y,
		Z: m.Data[2][0]*v.X + m.Data[2][1]*v.Y,
	}
}
