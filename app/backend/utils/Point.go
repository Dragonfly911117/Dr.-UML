package utils

import "math"

type Point struct {
	X int
	Y int
}

func (p Point) Magnitude() float64 {
	return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

func (p Point) MagnitudeInt() int {
	return int(p.Magnitude())
}

func Equal(p1, p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

func Add(p1, p2 Point) Point {
	return Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

func Sub(p1, p2 Point) Point {
	return Point{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}
