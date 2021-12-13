package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
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
	fmt.Printf("Part 2 solution: \n%s\n", part2(scanner))
}

var (
	inputFormat = regexp.MustCompile("(?P<iX>\\d+),(?P<iY>\\d+)")
	foldFormat  = regexp.MustCompile("fold along (?P<sDim>\\w)=(?P<iVal>\\d+)")
)

func parseIntArrayOrDie(text string) []int64 {
	res := make([]int64, 0)
	text = strings.TrimSpace(text)
	elements := strings.Split(text, ",")

	for _, e := range elements {
		i, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			panic(err)
		}
		res = append(res, i)
	}

	return res
}

type parsedParams struct {
	FullMatch string
	Strings   map[string]string
	Numbers   map[string]int
}

func extractRegexp(pattern *regexp.Regexp, str string) (*parsedParams, error) {
	params := &parsedParams{
		Strings: make(map[string]string),
		Numbers: make(map[string]int),
	}

	subMatches := pattern.FindStringSubmatch(str)
	if len(subMatches) != len(pattern.SubexpNames()) {
		return nil, fmt.Errorf("%s does not match pattern %s", str, pattern)
	}
	for i, name := range pattern.SubexpNames() {
		if name == "" {
			params.FullMatch = subMatches[i]
		} else if name[0] == 's' {
			params.Strings[name[1:]] = subMatches[i]
		} else if name[0] == 'i' {
			num, err := strconv.Atoi(subMatches[i])
			if err != nil {
				return nil, err
			}
			params.Numbers[name[1:]] = num
		} else {
			return nil, fmt.Errorf("unknown parse type for capture group %s", name)
		}
	}

	return params, nil
}

type Point struct {
	X int
	Y int
}

type Dir int

const (
	HORIZ Dir = iota
	VERT
)

func (d Dir) String() string {
	switch d {
	case HORIZ:
		return "y="
	case VERT:
		return "x="
	default:
		return "err"
	}
}

type Fold struct {
	Dir Dir
	Val int
}

func part1(input *scanner) int {
	// Setup
	lineNo := 0
	dots := make(map[Point]struct{})
	folds := []Fold{}
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		if parsedLine, err := extractRegexp(inputFormat, line); err == nil {
			dots[Point{X: parsedLine.Numbers["X"], Y: parsedLine.Numbers["Y"]}] = struct{}{}
		} else if instr, err := extractRegexp(foldFormat, line); err == nil {
			if instr.Strings["Dim"] == "x" {
				folds = append(folds, Fold{Dir: VERT, Val: instr.Numbers["Val"]})
			} else {
				folds = append(folds, Fold{Dir: HORIZ, Val: instr.Numbers["Val"]})
			}
		} else if len(line) == 0 {
			//meh
		} else {
			log.Fatalf("couldn't parse line %d: %s", lineNo+1, line)
		}

		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	log.Printf("starting with %d dots", len(dots))

	newDots := make(map[Point]struct{})
	for p := range dots {
		if folds[0].Dir == VERT && p.X < folds[0].Val {
			newDots[p] = struct{}{}
			log.Printf("keeping %v for instruction %v", p, folds[0])
		} else if folds[0].Dir == HORIZ && p.Y < folds[0].Val {
			newDots[p] = struct{}{}
			log.Printf("keeping %v for instruction %v", p, folds[0])
		} else if folds[0].Dir == VERT && p.X > folds[0].Val {
			distance := p.X - folds[0].Val
			newPoint := Point{X: folds[0].Val - distance, Y: p.Y}
			newDots[newPoint] = struct{}{}
			log.Printf("flipping %v to %v for instruction %v", p, newPoint, folds[0])
		} else if folds[0].Dir == HORIZ && p.Y > folds[0].Val {
			distance := p.Y - folds[0].Val
			newPoint := Point{X: p.X, Y: folds[0].Val - distance}
			newDots[newPoint] = struct{}{}
			log.Printf("flipping %v to %v for instruction %v", p, newPoint, folds[0])
		} else {
			log.Printf("point %v is on the line for instruction %v", p, folds[0])
		}
	}

	return len(newDots)
}

func part2(input *scanner) string {
	// Setup
	lineNo := 0
	dots := make(map[Point]struct{})
	folds := []Fold{}
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		if parsedLine, err := extractRegexp(inputFormat, line); err == nil {
			dots[Point{X: parsedLine.Numbers["X"], Y: parsedLine.Numbers["Y"]}] = struct{}{}
		} else if instr, err := extractRegexp(foldFormat, line); err == nil {
			if instr.Strings["Dim"] == "x" {
				folds = append(folds, Fold{Dir: VERT, Val: instr.Numbers["Val"]})
			} else {
				folds = append(folds, Fold{Dir: HORIZ, Val: instr.Numbers["Val"]})
			}
		} else if len(line) == 0 {
			//meh
		} else {
			log.Fatalf("couldn't parse line %d: %s", lineNo+1, line)
		}

		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	log.Printf("starting with %d dots", len(dots))

	for _, fold := range folds {
		newDots := make(map[Point]struct{})
		for p := range dots {
			if fold.Dir == VERT && p.X < fold.Val {
				newDots[p] = struct{}{}
				//				log.Printf("keeping %v for instruction %v", p, fold)
			} else if fold.Dir == HORIZ && p.Y < fold.Val {
				newDots[p] = struct{}{}
				//				log.Printf("keeping %v for instruction %v", p, fold)
			} else if fold.Dir == VERT && p.X > fold.Val {
				distance := p.X - fold.Val
				newPoint := Point{X: fold.Val - distance, Y: p.Y}
				newDots[newPoint] = struct{}{}
				//				log.Printf("flipping %v to %v for instruction %v", p, newPoint, fold)
			} else if fold.Dir == HORIZ && p.Y > fold.Val {
				distance := p.Y - fold.Val
				newPoint := Point{X: p.X, Y: fold.Val - distance}
				newDots[newPoint] = struct{}{}
				//				log.Printf("flipping %v to %v for instruction %v", p, newPoint, fold)
			} else {
				//				log.Printf("point %v is on the line for instruction %v", p, fold)
			}
		}
		dots = newDots
		log.Printf("%d dots after instruction %v", len(dots), fold)
	}

	log.Printf("%v", dots)

	grid := [][]string{}
	maxX := 0
	maxY := 0
	for p := range dots {
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	for row := 0; row <= maxY; row++ {
		grid = append(grid, make([]string, maxX+1))
		grid = append(grid, make([]string, maxX+1))
		for col := 0; col <= maxX; col++ {
			if _, hit := dots[Point{Y: row, X: col}]; hit {
				grid[row*2][col] = "XX"
				grid[(row*2)+1][col] = "XX"
			} else {
				grid[row*2][col] = "  "
				grid[(row*2)+1][col] = "  "
			}
		}
	}

	result := []string{}
	for _, r := range grid {
		result = append(result, strings.Join(r, ""))
	}

	return strings.Join(result, "\n")
}
