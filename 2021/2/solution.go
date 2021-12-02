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
	inputFormat = regexp.MustCompile("(?P<sDir>\\w+) (?P<iDistance>\\d+)")
)

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

func part1(input *scanner) int {
	// Setup
	count := 0
	x := 0
	y := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", count, err)
		}

		// Process line
		switch parsedLine.Strings["Dir"] {
		case "forward":
			x += parsedLine.Numbers["Distance"]
		case "up":
			y -= parsedLine.Numbers["Distance"]
			if y < 0 {
				y = 0
			}
		case "down":
			y += parsedLine.Numbers["Distance"]
		}
		count++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}
	return x * y
}

func part2(input *scanner) int {
	// Setup
	count := 0
	x := 0
	y := 0
	aim := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", count, err)
		}

		// Process line
		switch parsedLine.Strings["Dir"] {
		case "forward":
			x += parsedLine.Numbers["Distance"]
			y += (aim * parsedLine.Numbers["Distance"])
			if y < 0 {
				y = 0
			}
		case "up":
			aim -= parsedLine.Numbers["Distance"]
		case "down":
			aim += parsedLine.Numbers["Distance"]
		}
		count++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}
	return x * y
}
