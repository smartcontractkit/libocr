package minigobberish

// Taken from gobberish package, which provides helper functions for generating
// random strings for testing. https://github.com/chrismcguire/gobberish, commit
// 1d8adb5

// LICENSE (for this file, from
// https://github.com/chrismcguire/gobberish/blob/master/LICENSE):
//
// The MIT License (MIT)
//
// Copyright (c) 2015 Chris McGuire
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"errors"
	"math/rand"
	"time"
	"unicode"
)

// Generate a random utf-8 string of a given character (not byte) length.
// The range of the random characters is the entire printable unicode range.
func GenerateString(length int) string {
	var s []rune
	for i := 0; i < length; i++ {
		s = append(s, createRandomRune())
	}

	return string(s)
}

// Generates a random rune in the printable range.
func createRandomRune() rune {
	return createRandomRuneInRange(unicode.GraphicRanges)
}

// Generates a random rune in the given RangeTable.
func createRandomRuneInRange(tables []*unicode.RangeTable) rune {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(totalInRange(tables))
	x, _ := getItemInRangeTable(i, tables)

	return rune(x)
}

// Returns the nth item contained in the array of ranges.
func getItemInRangeTable(n int, tables []*unicode.RangeTable) (int, error) {
	nPointer := n
	var picked int
	found := false

	for _, table := range tables {
		if found == false {
			for _, r16 := range table.R16 {
				countInRange := int((r16.Hi-r16.Lo)/r16.Stride) + 1
				if nPointer <= countInRange-1 {
					picked = int(r16.Lo) + (int(r16.Stride) * nPointer)
					found = true
					break
				} else {
					nPointer -= countInRange
				}
			}

			if found == false {
				for _, r32 := range table.R32 {
					countInRange := int((r32.Hi-r32.Lo)/r32.Stride) + 1
					if nPointer <= countInRange-1 {
						picked = int(r32.Lo) + (int(r32.Stride) * nPointer)
						found = true
						break
					} else {
						nPointer -= countInRange
					}
				}
			}
		}
	}

	if found == true {
		return picked, nil
	} else {
		return -1, errors.New("Value not found in range")
	}
}

// Counts the total number of items contained in the array of ranges.
func totalInRange(tables []*unicode.RangeTable) int {
	total := 0
	for _, table := range tables {
		for _, r16 := range table.R16 {
			total += int((r16.Hi-r16.Lo)/r16.Stride) + 1
		}
		for _, r32 := range table.R32 {
			total += int((r32.Hi-r32.Lo)/r32.Stride) + 1
		}
	}
	return total
}
