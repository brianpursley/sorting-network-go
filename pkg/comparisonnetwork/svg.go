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
	"fmt"
	"github.com/brianpursley/sorting-network-go/pkg/comparator"
)

func (cn *ComparisonNetwork) Svg() string {
	scale := 1
	xScale := scale * 35
	yScale := scale * 20

	// Generate comparator SVG elements
	comparatorsSvg := ""
	w := xScale
	group := map[comparator.Comparator]int{}
	for _, c := range cn.Comparators {

		// If the comparator inputs are the same position as any other comparator in the group, then start a new group
		for other := range group {
			if comparatorsHaveAnySameInput(c, &other) {
				for _, pos := range group {
					if pos > w {
						w = pos
					}
				}
				w += xScale
				group = map[comparator.Comparator]int{}
				break
			}
		}

		// Adjust the comparator x position to avoid overlapping any existing Comparators in the group
		cx := w
		for other, otherPos := range group {
			if otherPos >= cx && comparatorsOverlap(c, &other) {
				cx = otherPos + xScale/3
			}
		}

		// Generate two circles and a line representing the comparator
		y0 := yScale + int(c.Input1())*yScale
		y1 := yScale + int(c.Input2())*yScale
		comparatorsSvg +=
			fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" style="stroke:black;stroke-width:1;fill=black" />`, cx, y0, 3) +
				fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" style="stroke:black;stroke-width:%d" />`, cx, y0, cx, y1, 1) +
				fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" style="stroke:black;stroke-width:1;fill=black" />`, cx, y1, 3)

		// Add this comparator to the current group
		group[*c] = cx
	}

	// Generate line SVG elements
	linesSvg := ""
	w += xScale
	n := int(cn.MaxInput() + 1)
	for i := 0; i < n; i++ {
		y := yScale + i*yScale
		linesSvg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" style="stroke:black;stroke-width:%d" />`, 0, y, w, y, 1)
	}

	h := (n + 1) * yScale
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<!DOCTYPE svg>
		<svg width="%dpx" height="%dpx" xmlns="http://www.w3.org/2000/svg">
			<rect width="100%%" height="100%%" fill="white" />
			%s
			%s
		</svg>`,
		w, h, comparatorsSvg, linesSvg)
}

func comparatorsOverlap(c, other *comparator.Comparator) bool {
	return (c.Input1() < other.Input1() && other.Input1() < c.Input2()) ||
		(c.Input1() < other.Input2() && other.Input2() < c.Input2()) ||
		(other.Input1() < c.Input1() && c.Input1() < other.Input2()) ||
		(other.Input1() < c.Input2() && c.Input2() < other.Input2())
}

func comparatorsHaveAnySameInput(c, other *comparator.Comparator) bool {
	return c.Input1() == other.Input1() ||
		c.Input1() == other.Input2() ||
		c.Input2() == other.Input1() ||
		c.Input2() == other.Input2()
}
