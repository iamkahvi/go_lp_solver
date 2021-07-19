package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"

	"example.com/solver/lp"
	sp "example.com/solver/simplex"
)

const DEBUG bool = false

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

	l := lp.New(m, r, c)

	var res sp.Result
	var opt float64
	var x []float64

	res = sp.Infeasible

	if l.Is_Primal_Feasible() {
		fmt.Fprintf(os.Stderr, "Primal Feasible\n")
		res, opt, x = sp.PrimalSimplex(l, DEBUG)
	} else if l.Is_Dual_Feasible() {
		fmt.Fprintf(os.Stderr, "Dual Feasible\n")
		res, opt, x = sp.DualSimplex(l, DEBUG)
	} else {
		fmt.Fprintf(os.Stderr, "Solve Aux\n")
		l_aux := l.CloneAux()
		_, _, _ = sp.DualSimplex(l_aux, DEBUG)
		l.B = l_aux.B
		l.N = l_aux.N
		res, opt, x = sp.PrimalSimplex(l, DEBUG)
	}

	switch res {
	case sp.Optimal:
		fmt.Fprintf(os.Stdout, "optimal\n%.7g\n%v\n", opt, print_arr(x))
	case sp.Unbounded:
		fmt.Fprintf(os.Stdout, "unbounded\n")
	case sp.Infeasible:
		fmt.Fprintf(os.Stdout, "infeasible\n")
	}
}

func parseElements(lines []string, rows int, cols int) [][]float64 {
	numbers := make([][]float64, rows)

	for i, line := range lines {
		els := s.Fields(line)

		numbers[i] = make([]float64, cols)

		for j, str := range els {
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				panic(err)
			}
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

func print_arr(arr []float64) string {
	var str string
	for i, xi := range arr {
		if xi < 1e-10 {
			str += "0"
		} else {
			str += fmt.Sprintf("%.7g", xi)
		}
		if i < len(arr)-1 {
			str += " "
		}
	}
	return str
}
