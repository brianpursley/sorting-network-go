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

package comparator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Comparator struct {
	i1 int
	i2 int
}

func (c *Comparator) Input1() int {
	return c.i1
}

func (c *Comparator) Input2() int {
	return c.i2
}

func (c *Comparator) String() string {
	return fmt.Sprintf("%d:%d", c.i1, c.i2)
}

func NewComparator(input1, input2 int) (*Comparator, error) {
	if input1 == input2 {
		return nil, errors.New("input1 and input2 cannot be the same")
	}
	c := &Comparator{}
	if input1 < input2 {
		c.i1 = input1
		c.i2 = input2
	} else {
		c.i1 = input2
		c.i2 = input1
	}
	return c, nil
}

func Parse(s string) (*Comparator, error) {
	inputs := strings.Split(s, ":")
	i1, err := strconv.Atoi(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid comparator: %s", s)
	}
	i2, err := strconv.Atoi(inputs[1])
	if err != nil {
		return nil, fmt.Errorf("invalid comparator: %s", s)
	}

	return NewComparator(i1, i2)
}
