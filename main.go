package main

import (
	"fmt"
	"io/ioutil"
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
	filename := os.Args[1]

	dat, err := ioutil.ReadFile(filename)
	check(err)

	lines := s.Split(string(dat), "\n")
	r, c := getDims(lines)

	m := makeMatrix(lines, r, c)

	fa := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("a = %v\n\n", fa)
	rows, cols := m.Dims()
	fmt.Printf("(%v, %v)\n", rows, cols)

	fmt.Printf("matrix: %T\n", m)
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
