package main

import (
	"bufio"
	"fmt"
	"flag"
	"github.com/projectdiscovery/subfinder/pkg/log"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

type Comparator struct {
	inputs [2]byte
}

type ComparisonNetwork struct {
	comparators []Comparator
}

// Use the zero-one principle to determine if this comparison network is a sorting network
func (cn ComparisonNetwork) isSortingNetwork() bool {
	numberOfInputs := cn.getMaxInput() + 1
	maxSequenceToCheck := uint((1 << numberOfInputs) - 1)

	results := make(chan bool)
	for sequence := uint(0); sequence <= maxSequenceToCheck; sequence++ {
		go func(sequence uint) {
			onesCount := byte(bits.OnesCount(sequence))
			zerosCount := numberOfInputs - onesCount
			expectedSortedSequence := uint(1 << onesCount - 1) << zerosCount
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
		c.applyComparatorToBinarySequence(&result)
	}
	return result
}

// Compare the two bits at the comparator's input positions and swap if needed
func (c Comparator) applyComparatorToBinarySequence(result *uint) {
	pos0 := (*result >> c.inputs[0]) & 1
	pos1 := (*result >> c.inputs[1]) & 1
	if pos0 > pos1 {
		*result ^= (1 << c.inputs[0]) | (1 << c.inputs[1])
	}
}

// Sorts a sequence of numbers
func (cn ComparisonNetwork) sortSequence(sequence []int64) []int64 {
	result := make([]int64, len(sequence))
	copy(result, sequence)
	for _, c := range cn.comparators {
		if result[c.inputs[0]] > result[c.inputs[1]] {
			result[c.inputs[0]], result[c.inputs[1]] = result[c.inputs[1]], result[c.inputs[0]]
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
		if c.inputs[0] > maxInput {
			maxInput = c.inputs[0]
		}
		if c.inputs[1] > maxInput {
			maxInput = c.inputs[1]
		}
	}
	return maxInput
}

func (cn ComparisonNetwork) svg() string {
	scale := 1
	xscale := scale * 35
	yscale := scale * 20

	// Generate comparator SVG elements
	comparatorsSvg := ""
	w := xscale
	group := map[Comparator]int{}
	for _, c := range cn.comparators {

		// If the comparator inputs are the same position as any other comparator in the group, then start a new group
		for gc := range group {
			if gc.inputs[0] == c.inputs[0] || gc.inputs[0] == c.inputs[1] || gc.inputs[1] == c.inputs[0] || gc.inputs[1] == c.inputs[1] {
				max := 0
				for _, pos := range group {
					if pos > max {
						max = pos
					}
				}
				w = max + xscale
				group = map[Comparator]int{}
				break
			}
		}

		// Adjust the comparator x position to avoid overlapping any existing comparators in the group
		cx := w
		for g, gx := range group {
			if comparatorsOverlap(c, g) {
				if gx >= cx {
					cx = gx + xscale / 3
				}
			}
		}
		y0 := yscale + int(c.inputs[0]) * yscale
		y1 := yscale + int(c.inputs[1]) * yscale

		// Generate two circles and a line representing the comparator
		comparatorsSvg +=
				fmt.Sprintf("<circle cx='%d' cy='%d' r='%d' style='stroke:black;stroke-width:1;fill=yellow' />", cx, y0, 3) +
				fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' style='stroke:black;stroke-width:%d' />", cx, y0, cx, y1, 1) +
				fmt.Sprintf("<circle cx='%d' cy='%d' r='%d' style='stroke:black;stroke-width:1;fill=yellow' />", cx, y1, 3)

		// Add this comparator to the current group
		group[c] = cx
	}

	// Generate line SVG elements
	linesSvg := ""
	w += xscale
	n := int(cn.getMaxInput() + 1)
	for i := 0; i < n; i++ {
		y := yscale + i * yscale
		linesSvg += fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' style='stroke:black;stroke-width:%d' />", 0, y, w, y, 1)
	}

	h := (n + 1) * yscale
	return fmt.Sprintf("<?xml version='1.0' encoding='utf-8'?><!DOCTYPE svg><svg width='%dpx' height='%dpx' xmlns='http://www.w3.org/2000/svg'>%s%s</svg>",
		w, h, comparatorsSvg, linesSvg)
}

func comparatorsOverlap(c1 Comparator, c2 Comparator) bool {
	return (c2.inputs[0] > c1.inputs[0] && c2.inputs[0] < c1.inputs[1]) ||
			(c2.inputs[1] > c1.inputs[0] && c2.inputs[1] < c1.inputs[1]) ||
			(c1.inputs[0] > c2.inputs[0] && c1.inputs[0] < c2.inputs[1]) ||
			(c1.inputs[1] > c2.inputs[0] && c1.inputs[1] < c2.inputs[1])
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

func newComparator(s string) *Comparator {
	inputs := strings.Split(s, ":")
	input0, err := strconv.Atoi(inputs[0])
	if err != nil {
		log.Fatalf("Failed to parse input file")
	}
	input1, err := strconv.Atoi(inputs[1])
	if err != nil {
		log.Fatalf("Failed to parse input file")
	}

	c := new(Comparator)
	if input0 < input1 {
		c.inputs[0] = byte(input0)
		c.inputs[1] = byte(input1)
	} else {
		c.inputs[0] = byte(input1)
		c.inputs[1] = byte(input0)
	}
	return c
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
		//cn.svg()
		fmt.Println(cn.svg())
	}
}
