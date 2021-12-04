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
	fmt.Printf("Part 2 solution: %d\n", part2(scanner))
}

type Board struct {
	Id      int
	Size    int
	Numbers [][]int64
	Marked  [][]bool
}

func NewBoard(id int, size int) *Board {
	nums := make([][]int64, size)
	marks := make([][]bool, size)
	for i := 0; i < size; i++ {
		nums[i] = make([]int64, size)
		marks[i] = make([]bool, size)
		for j := 0; j < size; j++ {
			marks[i][j] = false
		}
	}
	return &Board{
		Id:      id,
		Size:    size,
		Numbers: nums,
		Marked:  marks,
	}
}

func (b *Board) HasWon() bool {
	// horiz
	for row := 0; row < b.Size; row++ {
		rowWin := true
		for col := 0; col < b.Size; col++ {
			if !b.Marked[row][col] {
				rowWin = false
				break
			}
		}
		if rowWin {
			log.Printf("board %d won horiz on line %d", b.Id, row)
			return true
		}
	}

	// vert
	for col := 0; col < b.Size; col++ {
		colWin := true
		for row := 0; row < b.Size; row++ {
			if !b.Marked[row][col] {
				colWin = false
				break
			}
		}
		if colWin {
			log.Printf("board %d won vert on col %d", b.Id, col)
			return true
		}
	}
	return false
}

func (b *Board) Mark(num int64) {
	for row := 0; row < b.Size; row++ {
		for col := 0; col < b.Size; col++ {
			if b.Numbers[row][col] == num {
				log.Printf("Board %d marks row %d, col %d\n", b.Id, row, col)
				b.Marked[row][col] = true
			}
		}
	}
}

func (b *Board) SumOfUnmarked() int64 {
	var sum int64 = 0
	for row := 0; row < b.Size; row++ {
		for col := 0; col < b.Size; col++ {
			if !b.Marked[row][col] {
				sum += b.Numbers[row][col]
			}
		}
	}
	return sum
}

var (
	whitespace = regexp.MustCompile("\\s+")
)

func part1(input *scanner) int64 {
	firstLine, ok := input.NextLine()
	if !ok {
		log.Fatal("adsf")
	}
	nums := strings.Split(firstLine, ",")
	numbersCalled := make([]int64, len(nums))
	for i, n := range nums {
		parsed, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		numbersCalled[i] = parsed
	}

	boards := make([]*Board, 0)
	lineInBoard := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		line = strings.TrimSpace(line)
		// log.Printf("handling line %s\n", line)
		if len(line) == 0 {
			// 	log.Printf("between boards...\n")
			lineInBoard = 0
			continue
		}

		boardLine := whitespace.Split(line, -1)
		if lineInBoard == 0 {
			boards = append(boards, NewBoard(len(boards), len(boardLine)))
		}
		for i, n := range boardLine {
			//			log.Printf("parsing num %s\n", n)
			parsed, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			boards[len(boards)-1].Numbers[lineInBoard][i] = parsed
		}
		lineInBoard++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d numbers to call and %d boards to play...", len(numbersCalled), len(boards))

	for _, n := range numbersCalled {
		log.Printf("calling %d", n)
		for i, b := range boards {
			b.Mark(n)
			if b.HasWon() {
				log.Printf("board %d has won!", i)
				return b.SumOfUnmarked() * n
			}
		}
	}

	return 0
}

func part2(input *scanner) int64 {
	firstLine, ok := input.NextLine()
	if !ok {
		log.Fatal("adsf")
	}
	nums := strings.Split(firstLine, ",")
	numbersCalled := make([]int64, len(nums))
	for i, n := range nums {
		//		log.Printf("parsing num %s", n)
		parsed, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		numbersCalled[i] = parsed
	}

	boards := make([]*Board, 0)
	lineInBoard := 0
	for line, ok := input.NextLine(); ok; line, ok = input.NextLine() {
		line = strings.TrimSpace(line)
		// log.Printf("handling line %s\n", line)
		if len(line) == 0 {
			// 	log.Printf("between boards...\n")
			lineInBoard = 0
			continue
		}

		boardLine := whitespace.Split(line, -1)
		if lineInBoard == 0 {
			boards = append(boards, NewBoard(len(boards), len(boardLine)))
		}
		for i, n := range boardLine {
			//			log.Printf("parsing num %s\n", n)
			parsed, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			boards[len(boards)-1].Numbers[lineInBoard][i] = parsed
		}
		lineInBoard++
	}
	if err := input.Finish(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d numbers to call and %d boards to play...", len(numbersCalled), len(boards))

	boardsRemaining := make(map[*Board]struct{})
	for _, b := range boards {
		boardsRemaining[b] = struct{}{}
	}
	for _, n := range numbersCalled {
		log.Printf("calling %d", n)

		winners := make([]*Board, 0)
		for b := range boardsRemaining {
			b.Mark(n)
			if b.HasWon() {
				winners = append(winners, b)

				if len(boardsRemaining) == 1 {
					log.Printf("board %d has finally won!", b.Id)
					return b.SumOfUnmarked() * n
				}
			}
		}
		log.Printf("%d boards won and have been eliminated", len(winners))
		for _, b := range winners {
			delete(boardsRemaining, b)
		}
		log.Printf("%d boards remain", len(boardsRemaining))
	}

	return 0
}
