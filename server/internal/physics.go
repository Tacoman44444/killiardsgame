package internal

import (
	"math"
)

type Vector2 struct {
	x float64
	y float64
}

func (v Vector2) Add(u Vector2) Vector2      { return Vector2{v.x + u.x, v.y + u.y} }
func (v Vector2) Subtract(u Vector2) Vector2 { return Vector2{v.x - u.x, v.y - u.y} }
func (v Vector2) Multiply(s float64) Vector2 { return Vector2{v.x * s, v.y * s} }
func (v Vector2) Dot(u Vector2) float64      { return v.x*u.x + v.y*u.y }
func (v Vector2) Length_Squared() float64    { return v.Dot(v) }
func (v Vector2) Length() float64            { return math.Sqrt(v.Length_Squared()) }
func (v Vector2) Norm() Vector2 {
	if l := v.Length(); l != 0 {
		return v.Multiply(1 / l)
	}
	return Vector2{}
}

type Rect struct {
	topleft Vector2
	width   int
	height  int
}
