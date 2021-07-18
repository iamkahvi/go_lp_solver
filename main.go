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

	lpi := lp.New(m, r, c)

	if !lpi.Is_Feasible() {
		panic("Initial basis is not feasible")
	}

	iteration := 0

	// Setting X_B
	lpi.X_vec = lp.Set_V(lpi.Make_X_B(), lpi.X_vec, lpi.B)

	for {
		fmt.Fprintf(os.Stderr, "\niteration %v-----------------\n\n", iteration)

		lpi.Print()
		// zb <- 0
		lpi.Z_vec = lp.Set_V(mat.NewVecDense(len(lpi.B), nil), lpi.Z_vec, lpi.B)
		// zn <- complicated shit
		lpi.Z_vec = lp.Set_V(lpi.Make_Z_N(), lpi.Z_vec, lpi.N)

		if mat.Min(lpi.Z_N()) > 0 {
			matr := mat.NewDense(1, 1, nil)
			matr.Mul(lpi.C_B().T(), lpi.Make_X_B())

			v := lpi.X_vec.SliceVec(0, len(lpi.N))
			row := make([]int, len(lpi.N))

			for i := 0; i < len(lpi.N); i++ {
				row[i] = int(v.AtVec(i))
			}

			fmt.Fprintf(os.Stderr, "Optimal value: %v\n%v\n", matr.At(0, 0), row)

			break
		}

		// Choose entering variable

		// Largest Increase
		// zn_i := lp.Min_Index(lpi.Z_N())
		// j := lpi.N[zn_i]

		// Bland's rule
		var j int
		for _, ind := range lpi.N {
			if lpi.Z_vec.AtVec(ind) < 0 {
				j = ind
				break
			}
		}

		// Change T (theta) to D (delta) everywhere

		lp.Debug("Z", lpi.Z_vec)
		lp.Debug("Zn", lpi.Z_N())

		// Choosing a leaving variable

		// Construct theta x vector
		lpi.DX_vec = lp.Set_V(lpi.Make_TX_B(j), lpi.DX_vec, lpi.B)

		dXB := lp.Get_V(lpi.DX_vec, lpi.B)
		XB := lp.Get_V(lpi.X_vec, lpi.B)

		// Find min index for t
		t := math.MaxFloat64
		i := 0
		for _, bVal := range lpi.B {
			x := lpi.X_vec.AtVec(bVal)
			dx := lpi.DX_vec.AtVec(bVal)

			if dx > 0 {
				val := x / dx
				if val < t {
					t = val
					i = bVal
				}
			}
		}

		lp.Debug("X", lpi.X_vec)
		lp.Debug("XB", XB)
		lp.Debug("dXB", dXB)

		fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, lpi.Z_vec.AtVec(j))
		fmt.Fprintf(os.Stderr, "i = %v, xi =  %v\n", i, lpi.X_vec.AtVec(i))
		fmt.Fprintf(os.Stderr, "t = %v\n", t)

		// j = 0
		// i = 3

		// Updating xb
		v2 := mat.NewVecDense(XB.Len(), nil)
		dXB.ScaleVec(t, dXB)
		v2.SubVec(XB, dXB)
		lpi.X_vec = lp.Set_V(v2, lpi.X_vec, lpi.B)

		lpi.X_vec.SetVec(j, t)

		lpi.B = lp.Swap(j, i, lpi.B)
		lpi.N = lp.Swap(i, j, lpi.N)

		iteration++
		time.Sleep(1 * time.Second)
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
