package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type scanner struct {
	file *os.File
	sc   *bufio.Scanner
	err  error
}

func (s *scanner) NextLine() (string, bool) {
	if s.sc.Scan() {
		return s.sc.Text(), true
	} else {
		return "", false
	}
}

func (s *scanner) Finish() error {
	if s.file != nil {
		s.err = s.sc.Err()

		s.file.Close()
		s.file = nil
		s.sc = nil
	}
	return s.err
}

func getInputScanner() *scanner {
	_, thisFilePath, _, _ := runtime.Caller(0)
	f, err := os.OpenFile(filepath.Join(filepath.Dir(thisFilePath), "input.txt"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	return &scanner{
		file: f,
		sc:   bufio.NewScanner(f),
	}
}

func main() {
	scanner := getInputScanner()
	fmt.Printf("Part 1 solution: %d\n", part1(scanner))

	scanner = getInputScanner()
	fmt.Printf("Part 2 solution: %d\n", part2(scanner))
}

var (
	inputFormat = regexp.MustCompile("(?P<sData>.*)")
)

func part1(input *scanner) int {
	// Setup
	lineNo := 0

	space := [][]int64{}
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		nums := make([]int64, len(line))
		numStrs := strings.Split(line, "")
		for i, s := range numStrs {
			var err error
			nums[i], err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
		}
		space = append(space, nums)

		// Process line
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	totalRisk := 0
	for r, row := range space {
		for c, cell := range row {
			if (r == 0 || space[r-1][c] > cell) &&
				(c == 0 || row[c-1] > cell) &&
				(c == len(row)-1 || row[c+1] > cell) &&
				(r == len(space)-1 || space[r+1][c] > cell) {
				totalRisk += int(cell + 1)
			}
		}
	}

	return totalRisk
}

func part2(input *scanner) int {
	// Setup
	lineNo := 0

	space := [][]int64{}
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		nums := make([]int64, len(line))
		numStrs := strings.Split(line, "")
		for i, s := range numStrs {
			var err error
			nums[i], err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
		}
		space = append(space, nums)
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	width := len(space[0])
	height := len(space)
	basinAssignments := make([][]int, height)
	for r := 0; r < height; r++ {
		basinAssignments[r] = make([]int, width)
		for c := 0; c < width; c++ {
			basinAssignments[r][c] = -1
		}
	}

	nextBasin := 0
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if basinAssignments[r][c] == -1 && space[r][c] != 9 {
				floodFrom(space, basinAssignments, r, c, nextBasin)
				nextBasin++
			}
		}
	}

	basinSizes := make(map[int]int)
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if basinAssignments[r][c] != -1 {
				basinSizes[basinAssignments[r][c]]++
			}
		}
	}
	log.Printf("%v", basinSizes)

	sortedCounts := []int{}
	for _, count := range basinSizes {
		sortedCounts = append(sortedCounts, count)
	}
	sort.Sort(sort.IntSlice(sortedCounts))

	return sortedCounts[len(sortedCounts)-1] * sortedCounts[len(sortedCounts)-2] * sortedCounts[len(sortedCounts)-3]
}

type Point struct {
	row int
	col int
}

func floodFrom(space [][]int64, assignments [][]int, row int, col int, label int) {
	queue := []Point{
		{row: row, col: col},
	}
	for len(queue) > 0 {
		thisPoint := queue[0]
		queue = queue[1:]
		//		log.Println(thisPoint)

		assignments[thisPoint.row][thisPoint.col] = label
		neighbors := []Point{
			{row: thisPoint.row + 1, col: thisPoint.col},
			{row: thisPoint.row - 1, col: thisPoint.col},
			{row: thisPoint.row, col: thisPoint.col + 1},
			{row: thisPoint.row, col: thisPoint.col - 1},
		}
		for _, n := range neighbors {
			if n.row >= 0 && n.row < len(space) && n.col >= 0 && n.col < len(space[n.row]) {
				// real point
				if assignments[n.row][n.col] == -1 && space[n.row][n.col] != 9 {
					queue = append(queue, n)
				}
			}
		}
	}
}
