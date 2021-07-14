package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"

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
	m := makeMatrix(lines, r, c)
	nMatrix := makeNegMatrix(r, c)

	fm := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("a = %v\n\n", fm)

	// inverseMatrix := mat.NewDense(r, c, make([]float64, r*c))
	// inverseMatrix.Inverse(m)
	// fim := mat.Formatted(inverseMatrix, mat.Prefix("    "), mat.Squeeze())
	// fmt.Printf("a = %v\n\n", fim)

	negMatrix := mat.NewDense(r, c, make([]float64, r*c))
	negMatrix.MulElem(m, nMatrix)
	fat := mat.Formatted(negMatrix.T(), mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("a = %v\n\n", fat)

	vct := mat.Formatted(m.ColView(0))
	fmt.Printf("%v\n\n", vct)

	rw := mat.Formatted(m.RowView(0))
	fmt.Printf("%v\n\n", rw)

	rows, cols := m.Dims()
	fmt.Printf("(%v, %v)\n", rows, cols)

	fmt.Printf("matrix: %T\n", m)
}

func makeNegMatrix(rows int, cols int) *mat.Dense {
	negativeMatrix := mat.NewDense(rows, cols, make([]float64, rows*cols))
	negativeMatrix.Apply(func(i, j int, v float64) float64 {
		return -1
	}, negativeMatrix)
	return negativeMatrix
}

func makeMatrix(lines []string, rows int, cols int) *mat.Dense {
	m := mat.NewDense(rows, cols, nil)

	for i, line := range lines {
		els := s.Fields(line)

		for j, str := range els {
			val, err := strconv.ParseFloat(str, 64)
			check(err)
			m.Set(i, j, val)
		}
	}

	return m
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
