package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
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
	inputFormat = regexp.MustCompile("(?P<Data>.*)")
)

func extractRegexp(pattern *regexp.Regexp, str string) (map[string]string, error) {
	params := make(map[string]string)

	subMatches := pattern.FindStringSubmatch(str)
	if len(subMatches) != len(pattern.SubexpNames()) {
		return nil, fmt.Errorf("%s does not match pattern %s", str, pattern)
	}
	for i, name := range pattern.SubexpNames() {
		params[name] = subMatches[i]
	}

	return params, nil
}

func part1(input *scanner) int {
	// Setup
	count := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", count, err)
		}

		// Process line
		fmt.Printf("got line: %v\n", parsedLine)
		count++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}
	return count
}

func part2(input *scanner) int {
	// Setup
	count := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", count, err)
		}

		// Process line
		fmt.Printf("got line: %v\n", parsedLine)
		count++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}
	return count
}
