/*
The MIT License (MIT)

Copyright (c) 2022 Brian Pursley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, CopyAnalyzer, modify, merge, publish, distribute, sublicense, and/or sell
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

package comparisonnetwork

import (
	"github.com/brianpursley/sorting-network-go/pkg/comparator"
	"math/bits"
	"strings"
)

type ComparisonNetwork struct {
	Comparators []*comparator.Comparator
}

func Parse(data string) (*ComparisonNetwork, error) {
	cn := &ComparisonNetwork{}
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		for _, c := range strings.Split(line, ",") {
			c = strings.TrimSpace(c)
			if c == "" {
				continue
			}
			parsed, err := comparator.Parse(c)
			if err != nil {
				return nil, err
			}
			cn.Comparators = append(cn.Comparators, parsed)
		}
	}
	return cn, nil
}

func (cn *ComparisonNetwork) String() string {
	var s []string
	for _, c := range cn.Comparators {
		s = append(s, c.String())
	}
	return strings.Join(s, ",")
}

func (cn *ComparisonNetwork) MaxInput() int {
	maxInput := 0
	for _, c := range cn.Comparators {
		if c.Input2() > maxInput {
			maxInput = c.Input2()
		}
	}
	return maxInput
}

// IsSortingNetwork uses the zero-one principle to determine if this comparison network is a sorting network
func (cn *ComparisonNetwork) IsSortingNetwork() bool {
	numberOfInputs := cn.MaxInput() + 1
	maxSequenceToCheck := uint((1 << numberOfInputs) - 1)

	results := make(chan bool)
	for sequence := uint(0); sequence <= maxSequenceToCheck; sequence++ {
		go func(sequence uint) {
			onesCount := bits.OnesCount(sequence)
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
		if !<-results {
			return false
		}
	}

	return true
}

// sortBinarySequence applies all Comparators to the binary sequence
func (cn *ComparisonNetwork) sortBinarySequence(sequence uint) uint {
	result := sequence
	for _, c := range cn.Comparators {
		// Compare the two bits at the comparator's input positions and swap if needed
		pos0 := (result >> c.Input1()) & 1
		pos1 := (result >> c.Input2()) & 1
		if pos0 > pos1 {
			result ^= (1 << c.Input1()) | (1 << c.Input2())
		}
	}
	return result
}
