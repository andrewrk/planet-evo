package main

import "math"

type Vec2f struct {
	X float64
	Y float64
}

func (v *Vec2f) Add(other *Vec2f) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2f) Subtract(other *Vec2f) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vec2f) Clear() {
	v.X = 0
	v.Y = 0
}

func (v *Vec2f) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v *Vec2f) Negate() {
	v.X = -v.X
	v.Y = -v.Y
}

func (v *Vec2f) Scale(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2f) FloorEql(other *Vec2f) bool {
	return math.Floor(v.X) == math.Floor(other.X) && math.Floor(v.Y) == math.Floor(other.Y)
}

func (v *Vec2f) Normalize() {
	l := v.Length()
	v.X /= l
	v.Y /= l
}

func (v *Vec2f) LengthSqrd() float64 {
	return v.X * v.X + v.Y * v.Y
}

func (v *Vec2f) Length() float64 {
	return math.Sqrt(v.LengthSqrd())
}
