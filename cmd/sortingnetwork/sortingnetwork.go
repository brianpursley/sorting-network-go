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

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/brianpursley/sorting-network-go/pkg/comparisonnetwork"
)

func main() {
	inputFlag := flag.String("input", "", "specify a file containing comparison network definition")
	checkFlag := flag.Bool("check", false, "check whether it is a sorting network")
	sortFlag := flag.String("sort", "", "sorts the list of numbers using the input comparison network")
	svgFlag := flag.Bool("svg", false, "generate svg output")

	flag.Parse()

	var err error
	var data []byte
	if *inputFlag == "" {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v", err)
		}
	} else {
		data, err = os.ReadFile(*inputFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from file %s: %v", *inputFlag, err)
		}
	}
	cn, err := comparisonnetwork.Parse(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse comparison network: %v", err)
	}

	if *checkFlag {
		if cn.IsSortingNetwork() {
			fmt.Println("It is a sorting network!")
		} else {
			fmt.Println("It is not a sorting network.")
		}
	}

	if *sortFlag != "" {
		var numbers []int64
		for _, val := range strings.Split(*sortFlag, ",") {
			num, _ := strconv.ParseInt(val, 10, 64)
			numbers = append(numbers, num)
		}
		comparisonnetwork.Sort(cn, numbers, func(i, j int64) bool { return i < j })
		for _, num := range numbers {
			fmt.Println(num)
		}
	}

	if *svgFlag {
		fmt.Println(cn.Svg())
	}
}
