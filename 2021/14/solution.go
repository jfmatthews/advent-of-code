package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Printf("Part 1 solution: %d\n", part1())

	fmt.Printf("Part 2 solution: %d\n", part2())
}

const (
	start = "OOFNFCBHCKBBVNHBNVCP"
)

var mapping = map[string]byte{
	"PH": 'V',
	"OK": 'S',
	"KK": 'O',
	"BV": 'K',
	"CV": 'S',
	"SV": 'C',
	"CK": 'O',
	"PC": 'F',
	"SC": 'O',
	"KC": 'S',
	"KF": 'N',
	"SN": 'C',
	"SF": 'P',
	"OS": 'O',
	"OP": 'N',
	"FS": 'P',
	"FV": 'N',
	"CP": 'S',
	"VS": 'P',
	"PB": 'P',
	"HP": 'P',
	"PK": 'S',
	"FC": 'F',
	"SB": 'K',
	"NC": 'V',
	"PP": 'B',
	"PN": 'N',
	"VN": 'C',
	"NV": 'O',
	"OV": 'O',
	"BS": 'K',
	"FP": 'V',
	"NK": 'K',
	"PO": 'B',
	"HF": 'H',
	"VK": 'S',
	"ON": 'C',
	"KH": 'F',
	"HO": 'P',
	"OO": 'H',
	"BC": 'V',
	"CS": 'O',
	"OC": 'B',
	"VB": 'N',
	"OF": 'P',
	"FK": 'H',
	"OH": 'H',
	"CF": 'K',
	"CC": 'V',
	"BK": 'O',
	"BH": 'F',
	"VV": 'N',
	"KS": 'V',
	"FO": 'F',
	"SH": 'F',
	"OB": 'O',
	"VH": 'F',
	"HH": 'P',
	"PF": 'C',
	"NF": 'V',
	"VP": 'S',
	"CN": 'V',
	"SK": 'O',
	"FB": 'S',
	"FN": 'S',
	"BF": 'H',
	"FF": 'V',
	"CB": 'P',
	"NN": 'O',
	"VC": 'F',
	"HK": 'F',
	"BO": 'H',
	"KO": 'C',
	"CH": 'N',
	"KP": 'C',
	"HS": 'P',
	"NP": 'O',
	"NS": 'V',
	"NB": 'H',
	"HN": 'O',
	"BP": 'C',
	"VF": 'S',
	"KN": 'P',
	"HC": 'C',
	"PS": 'K',
	"BB": 'O',
	"NO": 'N',
	"NH": 'F',
	"BN": 'F',
	"KV": 'V',
	"SS": 'K',
	"CO": 'H',
	"KB": 'P',
	"FH": 'C',
	"SP": 'C',
	"SO": 'V',
	"PV": 'S',
	"VO": 'O',
	"HV": 'N',
	"HB": 'V',
}

func part1() int {
	theString := start
	for step := 0; step < 10; step++ {
		newString := []byte{}
		for i := 0; i < len(theString)-1; i++ {
			if insertion, exists := mapping[theString[i:i+2]]; exists {
				newString = append(newString, theString[i], insertion)
			} else {
				newString = append(newString, theString[i], theString[i+1])
			}
		}
		newString = append(newString, theString[len(theString)-1])
		theString = string(newString)
		log.Printf("after %d substitutions: %s", step+1, theString)
	}

	counts := map[rune]int{}
	for _, char := range theString {
		counts[char]++
	}

	minCount := len(theString)
	maxCount := 0
	for _, c := range counts {
		if c < minCount {
			minCount = c
		}
		if c > maxCount {
			maxCount = c
		}
	}

	return maxCount - minCount
}

func part2() int64 {
	// Similar to day 6 (the puzzle with the reproducing lanternfish), we don't
	// actually care /where/ each character is, just how many of each
	// subpattern there are. Therefore, we can handle them in bulk; this
	// version of the update operation scales with, at worst, the square of the
	// number of distinct characters in the string.

	// count of each character currently in the string. used to compute final
	// output.
	counts := map[byte]int64{}
	for i := range start {
		counts[start[i]]++
	}

	// number of each /pair/ of adjacent characters currently in the string.
	// used to do updates in each generation.
	bigrams := map[string]int64{}
	for i := 0; i < len(start)-1; i++ {
		bigrams[string(start[i:i+2])]++
	}
	log.Printf("starting bigrams: %v", bigrams)

	for step := 0; step < 40; step++ {
		newBigrams := map[string]int64{}
		for bigram, count := range bigrams {
			if insertion, exists := mapping[bigram]; exists {
				// count the char we're adding
				counts[insertion] += count

				// if AB->C, then in the next generation, we'll have bigrams AC
				// and CB instead of each instance of AB. (Note that we could
				// also end up constructing any of AC, CB, or AB from other
				// patterns; this is why we need to update into a new
				// bigram-count map for each generation rather than updating in
				// place.)
				newBigrams[string(bigram[0])+string(insertion)] += count
				newBigrams[string(insertion)+string(bigram[1])] += count
			} else {
				// if no match, these substrings won't change.
				newBigrams[bigram] += count
			}
		}
		bigrams = newBigrams
	}

	log.Printf("ending bigrams: %v", bigrams)

	minCount := int64(-1)
	maxCount := int64(0)
	for _, c := range counts {
		if c < minCount || minCount == -1 {
			minCount = c
		}
		if c > maxCount {
			maxCount = c
		}
	}

	return maxCount - minCount
}
