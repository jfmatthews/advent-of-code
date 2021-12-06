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
	inputFormat = regexp.MustCompile("(?P<ix1>\\d+),(?P<iy1>\\d+) -> (?P<ix2>\\d+),(?P<iy2>\\d+)")
)

type parsedParams struct {
	FullMatch string
	Strings   map[string]string
	Numbers   map[string]int64
}

func extractRegexp(pattern *regexp.Regexp, str string) (*parsedParams, error) {
	params := &parsedParams{
		Strings: make(map[string]string),
		Numbers: make(map[string]int64),
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
			num, err := strconv.ParseInt(subMatches[i], 10, 64)
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
	X int64
	Y int64
}

type Line struct {
	P1      Point
	P2      Point
	IsHoriz bool
	IsVert  bool
}

func NewLine(p1 Point, p2 Point) Line {
	//	log.Printf("line from %+v to %+v\n", p1, p2)
	isHoriz := false
	isVert := false
	if p1.Y == p2.Y {
		isHoriz = true
		if p2.X < p1.X {
			p1, p2 = p2, p1
		}
	}
	if p1.X == p2.X {
		isVert = true
		if p2.Y < p1.Y {
			p1, p2 = p2, p1
		}
	}
	if !isHoriz && !isVert {
		if p2.X < p1.X {
			p1, p2 = p2, p1
		}
	}
	return Line{
		P1:      p1,
		P2:      p2,
		IsHoriz: isHoriz,
		IsVert:  isVert,
	}
}

func (l Line) Points() map[Point]struct{} {
	points := make(map[Point]struct{})
	p := l.P1
	for {
		points[p] = struct{}{}
		if p == l.P2 {
			break
		}
		if l.P2.X > p.X {
			p.X++
		}
		if l.P2.Y > p.Y {
			p.Y++
		}
		if l.P2.Y < p.Y {
			p.Y--
		}
	}
	return points
}

func (l Line) Intersect(other Line) []Point {
	l1Points := l.Points()
	l2Points := other.Points()

	intersections := make([]Point, 0)
	for p := range l1Points {
		if _, in := l2Points[p]; in {
			intersections = append(intersections, p)
		}
	}
	return intersections
}

func part1(input *scanner) int {
	// Setup
	lineNo := 0
	lines := make([]Line, 0)
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", lineNo, err)
		}

		// Process line
		lines = append(lines, NewLine(
			Point{
				X: parsedLine.Numbers["x1"],
				Y: parsedLine.Numbers["y1"],
			},
			Point{
				X: parsedLine.Numbers["x2"],
				Y: parsedLine.Numbers["y2"],
			},
		))
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	//	log.Printf("%v", lines)

	collisions := make(map[Point]int)
	for i := 0; i < len(lines); i++ {
		if !lines[i].IsHoriz && !lines[i].IsVert {
			continue
		}
		for j := i + 1; j < len(lines); j++ {
			if !lines[j].IsHoriz && !lines[j].IsVert {
				continue
			}
			for _, p := range lines[i].Intersect(lines[j]) {
				collisions[p]++
			}
		}
	}

	return len(collisions)
}

func part2(input *scanner) int {
	// Setup
	lineNo := 0
	lines := make([]Line, 0)
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", lineNo, err)
		}

		// Process line
		lines = append(lines, NewLine(
			Point{
				X: parsedLine.Numbers["x1"],
				Y: parsedLine.Numbers["y1"],
			},
			Point{
				X: parsedLine.Numbers["x2"],
				Y: parsedLine.Numbers["y2"],
			},
		))
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	//	log.Printf("%v", lines)

	collisions := make(map[Point]int)
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			for _, p := range lines[i].Intersect(lines[j]) {
				collisions[p]++
			}
		}
	}

	return len(collisions)
}
