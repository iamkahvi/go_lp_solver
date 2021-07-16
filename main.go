package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"

	"example.com/m/lp"
	mat "gonum.org/v1/gonum/mat"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) <= 1 {
		fmt.Printf("Usage 'go run main.go < [file]'\n")
		os.Exit(0)
	}

	r, c := getDims(lines)
	m := parseElements(lines, r, c)

	// fmt.Fprintf(os.Stderr, "%v\n", m)
	// for _, row := range m {
	// 	fmt.Fprintf(os.Stderr, "%v\n", row)
	// }

	lp := lp.New(m, r, c)

	fm := mat.Formatted(lp.Get_Ab(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "Ab= %v\n\n", fm)

	lp.Print()
}

func makeNegMatrix(rows int, cols int) *mat.Dense {
	negativeMatrix := mat.NewDense(rows, cols, make([]float64, rows*cols))
	negativeMatrix.Apply(func(i, j int, v float64) float64 {
		return -1
	}, negativeMatrix)
	return negativeMatrix
}

func parseElements(lines []string, rows int, cols int) [][]float64 {
	numbers := make([][]float64, rows)

	for i, line := range lines {
		els := s.Fields(line)

		numbers[i] = make([]float64, cols)

		for j, str := range els {
			val, err := strconv.ParseFloat(str, 64)
			check(err)
			numbers[i][j] = val
		}
	}

	return numbers
}

func getDims(lines []string) (int, int) {
	rows := 0
	cols := 0
	for _, line := range lines {
		l := len(s.Fields(line))
		if l > 1 {
			rows += 1
		}
		if l > cols {
			cols = l
		}
	}

	return rows, cols
}
