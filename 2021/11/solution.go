package main

import (
	"fmt"
	"log"
)

func main() {
	doBoth()
}

type Point struct {
	row int
	col int
}

func blinkFrom(blinker Point, octoState [][]int) {
	neighbors := []Point{
		{row: blinker.row - 1, col: blinker.col - 1},
		{row: blinker.row - 1, col: blinker.col},
		{row: blinker.row - 1, col: blinker.col + 1},
		{row: blinker.row, col: blinker.col - 1},
		{row: blinker.row, col: blinker.col + 1},
		{row: blinker.row + 1, col: blinker.col - 1},
		{row: blinker.row + 1, col: blinker.col},
		{row: blinker.row + 1, col: blinker.col + 1},
	}
	for _, n := range neighbors {
		if n.row >= 0 && n.row < 10 && n.col >= 0 && n.col < 10 {
			// real point
			octoState[n.row][n.col]++
			if octoState[n.row][n.col] == 10 {
				// if >10, it's already been triggered
				blinkFrom(n, octoState)
			}
		}
	}
}

func doBoth() {
	octoState := [][]int{
		{7, 7, 7, 7, 8, 3, 8, 3, 5, 3},
		{2, 2, 1, 7, 2, 7, 2, 4, 7, 8},
		{3, 3, 5, 5, 3, 1, 8, 6, 4, 5},
		{2, 2, 4, 2, 6, 1, 8, 1, 1, 3},
		{7, 1, 8, 2, 4, 6, 8, 6, 6, 6},
		{5, 4, 4, 1, 6, 4, 1, 1, 1, 1},
		{4, 7, 7, 3, 8, 6, 2, 3, 6, 4},
		{5, 7, 1, 7, 1, 2, 5, 5, 2, 1},
		{7, 5, 4, 2, 1, 2, 7, 7, 2, 1},
		{4, 5, 7, 6, 6, 7, 8, 3, 4, 1},
	}

	blinks := 0
	foundSynchronizedBlink := false
	for step := 1; step <= 100 || !foundSynchronizedBlink; step++ {
		blinksThisStep := 0

		for r := 0; r < 10; r++ {
			for c := 0; c < 10; c++ {
				octoState[r][c]++
				if octoState[r][c] == 10 {
					// if >10, it's already been triggered
					blinkFrom(Point{row: r, col: c}, octoState)
				}
			}
		}

		for r := 0; r < 10; r++ {
			for c := 0; c < 10; c++ {
				if octoState[r][c] > 9 {
					octoState[r][c] = 0
					blinksThisStep++
				}
			}
		}

		blinks += blinksThisStep
		for _, row := range octoState {
			log.Printf("%v\n", row)
		}
		log.Printf("%d blinks on step %d\n", blinksThisStep, step)

		if step == 100 {
			fmt.Printf("Part 1 solution: %d\n", blinks)
		}
		if blinksThisStep == 100 {
			fmt.Printf("Part 2 solution: %d\n", step)
			foundSynchronizedBlink = true
		}
	}
}
