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

var (
	inputFormat = regexp.MustCompile("(?P<sData>.*)")
)

func part1(input *scanner) int {
	// Setup
	lineNo := 0
	zeroCounts := make([]int, 12)
	oneCounts := make([]int, 12)
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		// Parse line into components
		for i, char := range line {
			if char == '1' {
				oneCounts[i]++
			} else {
				zeroCounts[i]++
			}
		}
		lineNo++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	gamma := 0
	epsilon := 0
	for i := len(zeroCounts) - 1; i >= 0; i-- {
		bitMask := 1 << (len(zeroCounts) - i - 1)
		if oneCounts[i] > zeroCounts[i] {
			gamma += bitMask
		} else {
			epsilon += bitMask
		}
	}

	return gamma * epsilon
}

func part2(input *scanner) int64 {
	// Setup
	nums := make([]int64, 0)
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		n, err := strconv.ParseInt(line, 2, 64)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, n)
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}

	// Get filters:
	oxygenOptions := make([]int64, len(nums))
	copy(oxygenOptions, nums)
	for bit := 11; bit > 0 && len(oxygenOptions) > 1; bit-- {
		mostCommonBits := getMostCommonBits(oxygenOptions, 12)
		newOptions := make([]int64, 0)
		for _, option := range oxygenOptions {
			if option&(1<<bit) == mostCommonBits&(1<<bit) {
				newOptions = append(newOptions, option)
			}
		}

		oxygenOptions = newOptions
	}

	carbonOptions := make([]int64, len(nums))
	copy(carbonOptions, nums)
	for bit := 11; bit > 0 && len(carbonOptions) > 1; bit-- {
		mostCommonBits := getMostCommonBits(carbonOptions, 12)
		newOptions := make([]int64, 0)
		for _, option := range carbonOptions {
			// options here are /mismatches/
			if option&(1<<bit) != mostCommonBits&(1<<bit) {
				newOptions = append(newOptions, option)
			}
		}

		carbonOptions = newOptions
	}

	log.Printf("%v\n", oxygenOptions)
	log.Printf("%v\n", carbonOptions)

	return oxygenOptions[0] * carbonOptions[0]
}

func getMostCommonBits(nums []int64, numBits int) int64 {
	ones := make([]int, numBits)
	for _, n := range nums {
		for bit := 0; bit < numBits; bit++ {
			if n&(1<<bit) > 0 {
				ones[bit]++
			}
		}
	}

	var res int64 = 0
	for i, count := range ones {
		if count >= len(nums)/2 {
			res += (1 << i)
		}
	}
	return res
}
