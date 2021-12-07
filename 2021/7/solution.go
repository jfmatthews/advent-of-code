package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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
	fmt.Printf("Part 2 solution: %d\n", part2(scanner))
}

var (
	inputFormat = regexp.MustCompile("(?P<sData>.*)")
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

func part1(input *scanner) int64 {
	// Setup
	var crabs []int64
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		crabs = parseIntArrayOrDie(line)
		break
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	minCrab := int64(999999999)
	maxCrab := int64(-99999999)
	for _, c := range crabs {
		if minCrab > c {
			minCrab = c
		}
		if maxCrab < c {
			maxCrab = c
		}
	}

	//	bestStart := int64(0)
	bestScore := int64(9999999999999)
	for center := minCrab; center <= maxCrab; center++ {
		thisScore := int64(0)
		for _, c := range crabs {
			thisScore += int64(math.Abs(float64(c - center)))
		}
		if thisScore < bestScore {
			bestScore = thisScore
		}
	}

	return bestScore
}

func part2(input *scanner) int64 {
	// Setup
	var crabs []int64
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		crabs = parseIntArrayOrDie(line)
		break
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	minCrab := int64(999999999)
	maxCrab := int64(-99999999)
	for _, c := range crabs {
		if minCrab > c {
			minCrab = c
		}
		if maxCrab < c {
			maxCrab = c
		}
	}

	//	bestStart := int64(0)
	bestScore := int64(9999999999999)
	for center := minCrab; center <= maxCrab; center++ {
		thisScore := int64(0)
		for _, c := range crabs {
			distance := int64(math.Abs(float64(c - center)))
			cost := distance * (distance + 1) / 2
			thisScore += cost
		}
		if thisScore < bestScore {
			bestScore = thisScore
		}
	}

	return bestScore
}
