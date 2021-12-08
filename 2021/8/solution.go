package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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

func part1(input *scanner) int {
	// Setup
	lineNo := 0
	digitCount := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		inputHalves := strings.Split(line, " | ")
		if len(inputHalves) != 2 {
			panic(inputHalves)
		}

		outputDigits := strings.Split(strings.TrimSpace(inputHalves[1]), " ")
		// 		startCount := digitCount
		for _, d := range outputDigits {
			if len(d) == 2 || len(d) == 4 || len(d) == 7 || len(d) == 3 {
				digitCount++
			}
		}
		//		log.Printf("found %d matches in %v", digitCount-startCount, outputDigits)

		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	return digitCount
}

var validCombinations = map[string]int{
	"abcefg":  0,
	"cf":      1,
	"acdeg":   2,
	"acdfg":   3,
	"bcdf":    4,
	"abdfg":   5,
	"abdefg":  6,
	"acf":     7,
	"abcdefg": 8,
	"abcdfg":  9,
}

func SortChars(s string) string {
	slice := strings.Split(s, "")
	sort.Sort(sort.StringSlice(slice))
	return strings.Join(slice, "")
}

// Applies the character->character mapping in 'code' to 's'.
// Example: if s is "ab" and code is {"a": "f", "b": "d"}, this will return "fd"
func ApplyShuffle(s string, code map[string]string) string {
	slice := strings.Split(s, "")
	for i, char := range slice {
		var ok bool
		slice[i], ok = code[char]
		if !ok {
			log.Fatalf("encoding %v doesn't contain a mapping for char %s, needed for translating %s", code, char, s)
		}
	}
	return strings.Join(slice, "")
}

// Generates all permutations of 'arr'. Similar to Python's
// itertools.permutations, but materializes all permutations in memory, and
// builds them recursively without memoization, so don't use this with large
// arrays.
func PermuteArray(arr []string) [][]string {
	if len(arr) == 0 {
		return [][]string{}
	}

	if len(arr) == 1 {
		return [][]string{arr}
	}

	if len(arr) == 2 {
		return [][]string{
			{arr[0], arr[1]},
			{arr[1], arr[0]},
		}
	}

	if len(arr) == 3 {
		return [][]string{
			{arr[0], arr[1], arr[2]},
			{arr[0], arr[2], arr[1]},
			{arr[1], arr[0], arr[2]},
			{arr[1], arr[2], arr[0]},
			{arr[2], arr[0], arr[1]},
			{arr[2], arr[1], arr[0]},
		}
	}

	// recursive case: try putting each element first, and shuffling what's left.
	possibilities := [][]string{}
	for i := range arr {
		rest := make([]string, len(arr)-1)
		copy(rest, arr[:i])
		copy(rest[i:], arr[i+1:])
		for _, subPossibility := range PermuteArray(rest) {
			thisPossibility := append([]string{arr[i]}, subPossibility...)
			possibilities = append(possibilities, thisPossibility)
		}
	}
	return possibilities
}

// Returns all possible mappings from one char (in "abcdefg") to another.
func getCharacterShuffles() []map[string]string {
	shuffles := []map[string]string{}

	chars := strings.Split("abcdefg", "")
	permutations := PermuteArray(chars)
	log.Printf("Generated %d permutations of %s", len(permutations), chars)

	for _, permutation := range permutations {
		thisShuffle := make(map[string]string)
		for i := range chars {
			thisShuffle[chars[i]] = permutation[i]
		}
		shuffles = append(shuffles, thisShuffle)
	}

	return shuffles
}

// Each element of this slice is a copy of validCombinations with the segment
// IDs permuted. Each string is re-sorted after permuting, so an encoded
// (input) digit can be looked up directly if its encoded segment IDs are also
// sorted.
func getAllPossibleMappings() []map[string]int {
	shuffledChars := getCharacterShuffles()

	permutations := []map[string]int{}
	for _, shuffle := range shuffledChars {
		newPermutation := make(map[string]int)
		for s, val := range validCombinations {
			newPermutation[SortChars(ApplyShuffle(s, shuffle))] = val
		}
		permutations = append(permutations, newPermutation)
	}
	return permutations
}

func isValidMapping(mapping map[string]int, inputDigits map[string]int) bool {
	for d := range inputDigits {
		if _, ok := mapping[d]; !ok {
			return false
		}
	}
	return true
}

func part2(input *scanner) int {
	// There are only 7! = 5040 possible wirings, so we can brute-force this.
	// It would be /really/ cool to write a solver that worked out the problem
	// like a human, eliminating possibilities as we go (and could handle
	// certain input cases that don't include all 10 digits!), but this is a
	// much simpler puzzle than that.
	var possibleMappings []map[string]int = getAllPossibleMappings()

	totalOutput := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		inputHalves := strings.Split(line, " | ")
		if len(inputHalves) != 2 {
			panic(inputHalves)
		}

		inputDigits := make(map[string]int)
		outputDigits := make([]string, 0)

		for i := 0; i < 2; i++ {
			for _, d := range strings.Split(strings.TrimSpace(inputHalves[i]), " ") {
				if i == 0 {
					inputDigits[SortChars(d)] = -1
				}
				if i == 1 {
					outputDigits = append(outputDigits, SortChars(d))
				}
			}
		}
		for _, d := range outputDigits {
			if _, ok := inputDigits[d]; !ok {
				log.Fatalf("didn't see matching input for %s: %s", d, line)
			}
		}

		// Work out which mapping this case is using.
		foundValid := false
		for _, p := range possibleMappings {
			if isValidMapping(p, inputDigits) {
				for d := range inputDigits {
					inputDigits[d] = p[d]
				}
				foundValid = true
				break
			}
		}
		if !foundValid {
			log.Fatalf("couldn't find correct match for %s", line)
		}

		// Decode the output
		thisOutputNum := 0
		for _, d := range outputDigits {
			if inputDigits[d] == -1 {
				log.Fatalf("didn't decode %s in line: %s", d, line)
			}
			thisOutputNum *= 10
			thisOutputNum += inputDigits[d]
		}

		// Add to our total.
		totalOutput += thisOutputNum
	}

	return totalOutput
}
