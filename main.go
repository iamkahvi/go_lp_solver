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

	i := lp.New(m, r, c)

	i.Print()

	if !i.Is_Feasible() {
		panic("Initial basis is not feasible")
	}

	for {
		// zb <- 0
		i.Z_vec = lp.Set_V(mat.NewVecDense(len(i.B), nil), i.Z_vec, i.B)
		// zn <- complicated shit
		i.Z_vec = lp.Set_V(i.Make_Z_N(), i.Z_vec, i.N)

		if mat.Min(i.Z_N()) > 0 {
			fmt.Fprintf(os.Stderr, "Found optimal")
			break
		}

		// Choose entering variable
		_ = lp.Max_Index(i.Z_N())

		break

		// tX_B := i.Make_Theta_X_B(j)

		// theta_xn := mat.NewVecDense(len(lp.N), nil)

	}

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
