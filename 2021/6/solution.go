package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Part 1 solution: %d\n", part1())

	fmt.Printf("Part 2 solution: %d\n", part2())
}

var (
	input = []int{
		1, 1, 3, 5, 3, 1, 1, 4, 1, 1, 5, 2, 4, 3, 1, 1, 3, 1, 1, 5, 5, 1, 3, 2, 5, 4, 1, 1, 5, 1, 4, 2, 1, 4, 2, 1, 4, 4, 1, 5, 1, 4, 4, 1, 1, 5, 1, 5, 1, 5, 1, 1, 1, 5, 1, 2, 5, 1, 1, 3, 2, 2, 2, 1, 4, 1, 1, 2, 4, 1, 3, 1, 2, 1, 3, 5, 2, 3, 5, 1, 1, 4, 3, 3, 5, 1, 5, 3, 1, 2, 3, 4, 1, 1, 5, 4, 1, 3, 4, 4, 1, 2, 4, 4, 1, 1, 3, 5, 3, 1, 2, 2, 5, 1, 4, 1, 3, 3, 3, 3, 1, 1, 2, 1, 5, 3, 4, 5, 1, 5, 2, 5, 3, 2, 1, 4, 2, 1, 1, 1, 4, 1, 2, 1, 2, 2, 4, 5, 5, 5, 4, 1, 4, 1, 4, 2, 3, 2, 3, 1, 1, 2, 3, 1, 1, 1, 5, 2, 2, 5, 3, 1, 4, 1, 2, 1, 1, 5, 3, 1, 4, 5, 1, 4, 2, 1, 1, 5, 1, 5, 4, 1, 5, 5, 2, 3, 1, 3, 5, 1, 1, 1, 1, 3, 1, 1, 4, 1, 5, 2, 1, 1, 3, 5, 1, 1, 4, 2, 1, 2, 5, 2, 5, 1, 1, 1, 2, 3, 5, 5, 1, 4, 3, 2, 2, 3, 2, 1, 1, 4, 1, 3, 5, 2, 3, 1, 1, 5, 1, 3, 5, 1, 1, 5, 5, 3, 1, 3, 3, 1, 2, 3, 1, 5, 1, 3, 2, 1, 3, 1, 1, 2, 3, 5, 3, 5, 5, 4, 3, 1, 5, 1, 1, 2, 3, 2, 2, 1, 1, 2, 1, 4, 1, 2, 3, 3, 3, 1, 3, 5,
	}
)

func part1() int {
	state := make([]int, len(input))
	copy(state, input)

	for day := 0; day < 80; day++ {
		startingFish := len(state)
		// log.Printf("starting day %d with %d fish", day, startingFish)
		for fish := 0; fish < startingFish; fish++ {
			if state[fish] == 0 {
				state[fish] = 6
				state = append(state, 8)
			} else {
				state[fish]--
			}
		}
		//		log.Printf("ending day %d with %d fish", day, len(state))
	}
	return len(state)
}

func part2() int {
	states := make(map[int]int)
	for _, fishState := range input {
		states[fishState]++
	}
	for day := 0; day < 256; day++ {
		newState := make(map[int]int)
		splits := 0
		for state, count := range states {
			if state == 0 {
				splits = count
			} else {
				newState[state-1] = count
			}
		}
		newState[6] += splits
		newState[8] += splits
		states = newState

		// 		fmt.Printf("After day %d: %v\n", day, states)
	}

	totalCount := 0
	for _, count := range states {
		totalCount += count
	}
	return totalCount
}
