package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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

var matches = map[rune]rune{
	'{': '}',
	'[': ']',
	'<': '>',
	'(': ')',
}

var scores = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

func part1(input *scanner) int {
	// Setup
	lineNo := 0

	totalScore := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		stack := []rune{}

		thisLineScore := 0
		for i, char := range line {
			if _, isOpener := matches[char]; isOpener {
				//				log.Printf("adding %s to stack", string(char))
				stack = append(stack, char)
			} else if char == matches[stack[len(stack)-1]] {
				//				log.Printf("popping %s from stack", string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			} else {
				// this one is invalid
				thisLineScore = scores[char]
				log.Printf("Invalid char %s at pos %d of line %d; scoring %d", string(char), i, lineNo, thisLineScore)
				break
			}
		}

		totalScore += thisLineScore
		if thisLineScore == 0 {
			log.Printf("found no invalid in %s", line)
		}
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	return totalScore
}

var completionScores = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

func part2(input *scanner) int {
	// Setup
	lineNo := 0

	lineScores := []int{}
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		stack := []rune{}

		isInvalid := false
		for i, char := range line {
			if _, isOpener := matches[char]; isOpener {
				//				log.Printf("adding %s to stack", string(char))
				stack = append(stack, char)
			} else if char == matches[stack[len(stack)-1]] {
				//				log.Printf("popping %s from stack", string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			} else {
				// this one is invalid
				isInvalid = true
				log.Printf("Invalid char %s at pos %d of line %d", string(char), i, lineNo)
				break
			}
		}

		if isInvalid {
			continue
		}

		log.Printf("remaining: %s", string(stack))
		thisLineScore := 0
		for len(stack) > 0 {
			thisLineScore *= 5
			thisLineScore += completionScores[matches[stack[len(stack)-1]]]
			stack = stack[:len(stack)-1]
			log.Printf("score: %d", thisLineScore)
		}
		log.Printf("completion score: %d", thisLineScore)
		lineScores = append(lineScores, thisLineScore)
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(sort.IntSlice(lineScores))
	log.Printf("got %d incomplete lines", len(lineScores))

	middle := len(lineScores) / 2
	log.Printf("middle index is %d", middle)

	return lineScores[middle]
}
