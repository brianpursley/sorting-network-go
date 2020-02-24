/*
The MIT License (MIT)

Copyright (c) 2020 Brian Pursley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

type Comparator struct {
	i1 byte
	i2 byte
}

func newComparator(s string) *Comparator {
	inputs := strings.Split(s, ":")
	i1, err := strconv.Atoi(inputs[0])
	if err != nil {
		log.Fatalf("invalid comparator: %s", s)
	}
	i2, err := strconv.Atoi(inputs[1])
	if err != nil {
		log.Fatalf("invalid comparator: %s", s)
	}

	c := new(Comparator)
	if i1 < i2 {
		c.i1 = byte(i1)
		c.i2 = byte(i2)
	} else {
		c.i1 = byte(i2)
		c.i2 = byte(i1)
	}
	return c
}

func (c Comparator) overlaps(other Comparator) bool {
	return (c.i1 < other.i1 && other.i1 < c.i2) ||
		(c.i1 < other.i2 && other.i2 < c.i2) ||
		(other.i1 < c.i1 && c.i1 < other.i2) ||
		(other.i1 < c.i2 && c.i2 < other.i2)
}

func (c Comparator) hasSameInput(other Comparator) bool {
	return c.i1 == other.i1 ||
		c.i1 == other.i2 ||
		c.i2 == other.i1 ||
		c.i2 == other.i2
}

type ComparisonNetwork struct {
	comparators []Comparator
}

func newComparisonNetwork(file *os.File) *ComparisonNetwork {
	cn := new(ComparisonNetwork)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, c := range strings.Split(line, ",") {
			cn.comparators = append(cn.comparators, *newComparator(c))
		}
	}
	return cn
}

// Use the zero-one principle to determine if this comparison network is a sorting network
func (cn ComparisonNetwork) isSortingNetwork() bool {
	numberOfInputs := cn.getMaxInput() + 1
	maxSequenceToCheck := uint((1 << numberOfInputs) - 1)

	results := make(chan bool)
	for sequence := uint(0); sequence <= maxSequenceToCheck; sequence++ {
		go func(sequence uint) {
			onesCount := byte(bits.OnesCount(sequence))
			var expectedSortedSequence uint
			if onesCount > 0 {
				zerosCount := numberOfInputs - onesCount
				expectedSortedSequence = uint(1<<onesCount-1) << zerosCount
			} else {
				expectedSortedSequence = 0
			}
			results <- cn.sortBinarySequence(sequence) == expectedSortedSequence
		}(sequence)
	}

	for sequence := uint(0); sequence <= maxSequenceToCheck; sequence++ {
		if ! <- results {
			return false
		}
	}

	return true
}

// Apply all comparators to the binary sequence
func (cn ComparisonNetwork) sortBinarySequence(sequence uint) uint {
	result := sequence
	for _, c := range cn.comparators {
		// Compare the two bits at the comparator's input positions and swap if needed
		pos0 := (result >> c.i1) & 1
		pos1 := (result >> c.i2) & 1
		if pos0 > pos1 {
			result ^= (1 << c.i1) | (1 << c.i2)
		}
	}
	return result
}

// Sorts a sequence of numbers
func (cn ComparisonNetwork) sortSequence(sequence []int64) []int64 {
	result := make([]int64, len(sequence))
	copy(result, sequence)
	for _, c := range cn.comparators {
		if result[c.i1] > result[c.i2] {
			result[c.i1], result[c.i2] = result[c.i2], result[c.i1]
		}
	}
	return result
}

func stringToNumberSequence(s string) []int64 {
	var result []int64
	for _, val := range strings.Split(s, ",") {
		num, _ := strconv.ParseInt(val, 10, 64)
		result = append(result, num)
	}
	return result
}

func (cn ComparisonNetwork) getMaxInput() byte {
	maxInput := byte(0)
	for _, c := range cn.comparators {
		if c.i2 > maxInput {
			maxInput = c.i2
		}
	}
	return maxInput
}

func (cn ComparisonNetwork) svg() string {
	scale := 1
	xScale := scale * 35
	yScale := scale * 20

	// Generate comparator SVG elements
	comparatorsSvg := ""
	w := xScale
	group := map[Comparator]int{}
	for _, c := range cn.comparators {

		// If the comparator inputs are the same position as any other comparator in the group, then start a new group
		for other := range group {
			if c.hasSameInput(other) {
				for _, pos := range group {
					if pos > w {
						w = pos
					}
				}
				w += xScale
				group = map[Comparator]int{}
				break
			}
		}

		// Adjust the comparator x position to avoid overlapping any existing comparators in the group
		cx := w
		for other, otherPos := range group {
			if otherPos >= cx && c.overlaps(other) {
				cx = otherPos + xScale / 3
			}
		}

		// Generate two circles and a line representing the comparator
		y0 := yScale + int(c.i1) * yScale
		y1 := yScale + int(c.i2) * yScale
		comparatorsSvg +=
				fmt.Sprintf("<circle cx='%d' cy='%d' r='%d' style='stroke:black;stroke-width:1;fill=yellow' />", cx, y0, 3) +
				fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' style='stroke:black;stroke-width:%d' />", cx, y0, cx, y1, 1) +
				fmt.Sprintf("<circle cx='%d' cy='%d' r='%d' style='stroke:black;stroke-width:1;fill=yellow' />", cx, y1, 3)

		// Add this comparator to the current group
		group[c] = cx
	}

	// Generate line SVG elements
	linesSvg := ""
	w += xScale
	n := int(cn.getMaxInput() + 1)
	for i := 0; i < n; i++ {
		y := yScale + i *yScale
		linesSvg += fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' style='stroke:black;stroke-width:%d' />", 0, y, w, y, 1)
	}

	h := (n + 1) * yScale
	return fmt.Sprintf(
		"<?xml version='1.0' encoding='utf-8'?>" +
		"<!DOCTYPE svg>" +
		"<svg width='%dpx' height='%dpx' xmlns='http://www.w3.org/2000/svg'>%s%s</svg>",
		w, h, comparatorsSvg, linesSvg)
}

func main() {
	inputFlag := flag.String("input", "", "specify a file containing comparison network definition")
	checkFlag := flag.Bool("check", false, "check whether it is a sorting network")
	sortFlag := flag.String("sort", "", "sorts the list of numbers using the input comparison network")
	svgFlag := flag.Bool("svg", false, "generate svg output")

	flag.Parse()

	var cn *ComparisonNetwork
	if *inputFlag == "" {
		cn = newComparisonNetwork(os.Stdin)
	} else {
		file, err := os.Open(*inputFlag)
		if err != nil {
			log.Fatalf("Failed to open input file: %s", *inputFlag)
		}
		defer file.Close()
		cn = newComparisonNetwork(file)
	}

	if *checkFlag {
		if cn.isSortingNetwork() {
			fmt.Println("It is a sorting network!")
		} else {
			fmt.Println("It is not a sorting network.")
		}
	}

	if *sortFlag != "" {
		for _, num := range cn.sortSequence(stringToNumberSequence(*sortFlag)) {
			println(num)
		}
	}

	if *svgFlag {
		fmt.Println(cn.svg())
	}
}
