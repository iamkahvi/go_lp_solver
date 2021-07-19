package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	s "strings"
	"time"

	"example.com/m/lp"
	mat "gonum.org/v1/gonum/mat"
)

const DEBUG bool = false

type Result int

const (
	Optimal Result = iota
	Unbounded
	Infeasible
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func print_arr(arr []float64) string {
	var str string
	for i, xi := range arr {
		str += fmt.Sprintf("%.7g", xi)
		if i < len(arr)-1 {
			str += " "
		}
	}
	return str
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

	l := lp.New(m, r, c)

	var res Result
	var opt float64
	var x []float64

	res = Infeasible

	if l.Is_Primal_Feasible() {
		// is primal feasible
		fmt.Fprintf(os.Stderr, "Primal Feasible\n")
		res, opt, x = PrimalSimplex(l)
	} else if l.Is_Dual_Feasible() {
		// is dual feasible
		fmt.Fprintf(os.Stderr, "Dual Feasible\n")
		res, opt, x = DualSimplex(l)
	} else {
		fmt.Fprintf(os.Stderr, "SOLVE AUX\n")
		l_aux := l.CloneAux()
		_, _, _ = DualSimplex(l_aux)
		fmt.Fprintf(os.Stderr, "SOLVED AUX\n")
		l.B = l_aux.B
		l.N = l_aux.N
		res, opt, x = PrimalSimplex(l)
	}

	switch res {
	case Optimal:
		fmt.Fprintf(os.Stdout, "optimal\n%.7g\n%v\n", opt, print_arr(x))
	case Unbounded:
		fmt.Fprintf(os.Stdout, "unbounded\n")
	case Infeasible:
		fmt.Fprintf(os.Stdout, "infeasible\n")
	}
}

func PrimalSimplex(l *lp.LP) (Result, float64, []float64) {
	iteration := 0

	// Setting xB
	l.X_vec = lp.Set_V(l.Make_X_B(), l.X_vec, l.B)
	// Setting xN
	l.X_vec = lp.Set_V(mat.NewVecDense(len(l.N), nil), l.X_vec, l.N)

	if mat.Min(l.X_B()) < 0 {
		return Infeasible, 0, nil
	}

	for {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "\niteration %v-----------------\n\n", iteration)
			l.Print()
		}

		// zb <- 0
		l.Z_vec = lp.Set_V(mat.NewVecDense(len(l.B), nil), l.Z_vec, l.B)
		// zn <- complicated shit
		l.Z_vec = lp.Set_V(l.Make_Z_N(), l.Z_vec, l.N)

		if mat.Min(l.Z_N()) > 0 {
			matr := mat.NewDense(1, 1, nil)
			matr.Mul(l.C_B().T(), l.Make_X_B())

			v := l.X_vec.SliceVec(0, len(l.N))
			row := make([]float64, len(l.N))

			for i := 0; i < len(l.N); i++ {
				row[i] = v.AtVec(i)
			}

			return Optimal, matr.At(0, 0), row
		}

		// Choose entering variable

		// Largest Increase
		// zn_i := lp.Min_Index(lpi.Z_N())
		// j := lpi.N[zn_i]

		// Bland's rule
		var j int
		for _, ind := range l.N {
			if l.Z_vec.AtVec(ind) < 0 {
				j = ind
				break
			}
		}

		// Choosing a leaving variable

		// Construct theta x vector
		l.DX_vec = lp.Set_V(l.Make_TX_B(j), l.DX_vec, l.B)

		dXB := lp.Get_V(l.DX_vec, l.B)
		XB := lp.Get_V(l.X_vec, l.B)

		// Find min index for t
		t := math.MaxFloat64
		i := 0
		for _, bVal := range l.B {
			x := l.X_vec.AtVec(bVal)
			dx := l.DX_vec.AtVec(bVal)

			if dx > 0 {
				val := x / dx
				if val < t {
					t = val
					i = bVal
				}
			}
		}

		if l.Is_Unbounded() {
			return Unbounded, 0, nil
		}

		if DEBUG {
			fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, l.Z_vec.AtVec(j))
			fmt.Fprintf(os.Stderr, "i = %v, xi =  %v\n", i, l.X_vec.AtVec(i))
			fmt.Fprintf(os.Stderr, "t = %v\n", t)
		}

		// j = 0
		// i = 3

		// Updating xb
		v2 := mat.NewVecDense(XB.Len(), nil)
		dXB.ScaleVec(t, dXB)
		v2.SubVec(XB, dXB)
		l.X_vec = lp.Set_V(v2, l.X_vec, l.B)

		l.X_vec.SetVec(j, t)

		l.B = lp.Swap(j, i, l.B)
		l.N = lp.Swap(i, j, l.N)

		iteration++

		if DEBUG {
			time.Sleep(1 * time.Second)
		}
	}
}

func DualSimplex(l *lp.LP) (Result, float64, []float64) {
	iteration := 0

	// Setting zn and zb
	l.Z_vec = lp.Set_V(mat.NewVecDense(len(l.B), nil), l.Z_vec, l.B)
	l.Z_vec = lp.Set_V(l.Make_Z_N(), l.Z_vec, l.N)

	if mat.Min(l.Z_N()) < 0 {
		return Infeasible, 0, nil
	}

	for {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "\niteration %v-----------------\n\n", iteration)
			// l.Print()
		}

		// Setting xb and xn
		l.X_vec = lp.Set_V(l.Make_X_B(), l.X_vec, l.B)
		l.X_vec = lp.Set_V(mat.NewVecDense(len(l.N), nil), l.X_vec, l.N)

		if mat.Min(l.X_B()) > 0 {
			// Computing optimal
			matr := mat.NewDense(1, 1, nil)
			matr.Mul(l.C_B().T(), l.Make_X_B())

			v := l.X_vec.SliceVec(0, len(l.N))
			row := make([]float64, len(l.N))

			for i := 0; i < len(l.N); i++ {
				row[i] = v.AtVec(i)
			}

			return Optimal, matr.At(0, 0), row
		}

		// Choose leaving variable

		// Bland's rule
		var i int
		for _, ind := range l.B {
			if l.X_vec.AtVec(ind) < 0 {
				i = ind
				break
			}
		}

		// Choose entering variable

		// Creating vector u
		u := mat.NewVecDense(len(l.B), nil)

		for idx, val := range l.B {
			if val == i {
				u.SetVec(idx, 1)
			}
		}

		// Create delta Z vector
		ab := l.A_B()
		ab.Inverse(ab.T())

		var rh mat.Dense
		rh.Mul(ab, u)

		an := l.A_N()
		var an2 mat.Dense
		an2.Scale(-1, an.T())

		var res mat.Dense
		res.Mul(&an2, &rh)
		dZB := res.ColView(0).(*mat.VecDense)

		dZ := mat.NewVecDense(l.Z_vec.Len(), nil)

		dZ = lp.Set_V(dZB, dZ, l.N)
		dZ = lp.Set_V(mat.NewVecDense(len(l.B), nil), dZ, l.B)

		// Find min index for t
		s := math.MaxFloat64
		j := 0
		for _, nVal := range l.N {
			x := l.Z_vec.AtVec(nVal)
			dx := dZ.AtVec(nVal)

			if dx > 0 {
				val := x / dx
				if val < s {
					s = val
					j = nVal
				}
			}
		}

		if mat.Max(dZ) <= 0 {
			return Infeasible, 0, nil
		}

		var v mat.VecDense
		v.SubVec(l.Z_N(), lp.Get_V(dZ, l.N))
		l.Z_vec = lp.Set_V(&v, l.Z_vec, l.N)
		l.Z_vec.SetVec(i, s)

		l.B = lp.Swap(j, i, l.B)
		l.N = lp.Swap(i, j, l.N)

		if DEBUG {
			fmt.Fprintf(os.Stderr, "i = %v, xi =  %v\n", i, l.X_vec.AtVec(i))
			fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, l.Z_vec.AtVec(j))
			fmt.Fprintf(os.Stderr, "s = %v\n", s)
		}

		iteration++

		if DEBUG {
			time.Sleep(1 * time.Second)
		}
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
