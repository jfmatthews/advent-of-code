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
	"unicode"
	"unicode/utf8"
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
	inputFormat = regexp.MustCompile("(?P<sCaveA>[a-zA-Z]+)-(?P<sCaveB>[a-zA-Z]+)")
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

type Cave struct {
	isBig bool
	name  string

	neighbors map[string]*Cave
}

func getAllPaths(start *Cave, end *Cave, visited map[*Cave]struct{}, allowDoubleVisit bool, depth int, path string) int {
	if start == end {
		return 1
	}

	paths := 0
	for _, n := range start.neighbors {
		if _, seen := visited[n]; n.isBig || !seen || allowDoubleVisit {
			visitedClone := make(map[*Cave]struct{})
			for visitedAlready := range visited {
				visitedClone[visitedAlready] = struct{}{}
			}
			if !start.isBig {
				visitedClone[start] = struct{}{}
			}

			_, doneDoubleVisit := visitedClone[n]
			if doneDoubleVisit {
				log.Printf("visiting %s for the second time at depth %d: %s", n.name, depth, path)
			}

			paths += getAllPaths(n, end, visitedClone, allowDoubleVisit && !doneDoubleVisit, depth+1, path+"-"+n.name)
		}
	}

	log.Printf("%d paths starting from %s at depth %d\n", paths, start.name, depth)

	return paths
}

func part1(input *scanner) int {
	// Setup
	lineNo := 0
	caves := make(map[string]*Cave)
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", lineNo, err)
		}

		// Process line
		caveA := parsedLine.Strings["CaveA"]
		caveB := parsedLine.Strings["CaveB"]
		if _, exists := caves[caveA]; !exists {
			firstRune, _ := utf8.DecodeRuneInString(caveA)
			isBig := unicode.IsUpper(firstRune)
			log.Printf("creating cave %s (big? %t)", caveA, isBig)
			caves[caveA] = &Cave{
				isBig:     isBig,
				name:      caveA,
				neighbors: make(map[string]*Cave),
			}
		}
		if _, exists := caves[caveB]; !exists {
			firstRune, _ := utf8.DecodeRuneInString(caveB)
			isBig := unicode.IsUpper(firstRune)
			log.Printf("creating cave %s (big? %t)", caveB, isBig)
			caves[caveB] = &Cave{
				isBig:     isBig,
				name:      caveB,
				neighbors: make(map[string]*Cave),
			}
		}
		caves[caveA].neighbors[caveB] = caves[caveB]
		caves[caveB].neighbors[caveA] = caves[caveA]

		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	// search
	return getAllPaths(caves["start"], caves["end"], map[*Cave]struct{}{}, false, 0, "")
}

func part2(input *scanner) int {
	// Setup
	lineNo := 0
	caves := make(map[string]*Cave)
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		parsedLine, err := extractRegexp(inputFormat, line)
		if err != nil {
			log.Fatalf("input line %d: %v", lineNo, err)
		}

		// Process line
		caveA := parsedLine.Strings["CaveA"]
		caveB := parsedLine.Strings["CaveB"]
		if _, exists := caves[caveA]; !exists {
			firstRune, _ := utf8.DecodeRuneInString(caveA)
			isBig := unicode.IsUpper(firstRune)
			log.Printf("creating cave %s (big? %t)", caveA, isBig)
			caves[caveA] = &Cave{
				isBig:     isBig,
				name:      caveA,
				neighbors: make(map[string]*Cave),
			}
		}
		if _, exists := caves[caveB]; !exists {
			firstRune, _ := utf8.DecodeRuneInString(caveB)
			isBig := unicode.IsUpper(firstRune)
			log.Printf("creating cave %s (big? %t)", caveB, isBig)
			caves[caveB] = &Cave{
				isBig:     isBig,
				name:      caveB,
				neighbors: make(map[string]*Cave),
			}
		}
		if caveA != "end" && caveB != "start" {
			log.Printf("path from %s to %s", caveA, caveB)
			caves[caveA].neighbors[caveB] = caves[caveB]
		}
		if caveA != "start" && caveB != "end" {
			log.Printf("path from %s to %s", caveB, caveA)
			caves[caveB].neighbors[caveA] = caves[caveA]
		}

		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	// search
	return getAllPaths(caves["start"], caves["end"], map[*Cave]struct{}{}, true, 0, "")
}
