package main

import (
	"bufio"
	"fmt"
	"os"

	"example.com/solver/lp"
	sp "example.com/solver/simplex"
	utils "example.com/solver/utils"
)

const DEBUG bool = false

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	m, r, c := utils.ParseLines(lines)

	if r <= 1 {
		fmt.Printf("Usage 'go run main.go < [file]'\n")
		os.Exit(0)
	}

	l := lp.New(m, r, c)

	// Perturbation method?
	// l.B_vec.ScaleVec(lp.EPSILON, l.B_vec)

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
		fmt.Fprintf(os.Stderr, "Solved Aux\n")
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

func print_arr(arr []float64) string {
	var str string
	for i, xi := range arr {
		if xi < 1e-3 {
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
