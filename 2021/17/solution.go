package main

import "fmt"

func main() {
	fmt.Printf("Part 1 solution: %d\n", part1())

	fmt.Printf("Part 2 solution: %d\n", part2())
}

const (
	minX = 257
	maxX = 286
	minY = -101
	maxY = -57
)

type Velocity struct {
	X int
	Y int
}

type Point struct {
	X int
	Y int
}

func slow(xVelo int) int {
	if xVelo < 0 {
		return xVelo + 1
	} else if xVelo > 0 {
		return xVelo - 1
	} else {
		return xVelo
	}
}
func NextStep(p Point, v Velocity) (Point, Velocity) {
	return Point{
		X: p.X + v.X,
		Y: p.Y + v.Y,
	}, Velocity {
		X: slow(v.X),
		Y: v.Y - 1,
	}
}

func part1() int {
	bestYMax := -200

	// maximize y subject to x constraint
	for dx := 1; dx < maxX; dx++ {
		for dy := -101; dy minY; dx < maxX && dy <  ; xVeloGuess < maxX && yVeloGuess > 
	return 0
}
func part2() int {
	return 0
}
