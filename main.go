package main

import (
	"bufio"
	"fmt"
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

	for {
		lpi.Print()
		// zb <- 0
		lpi.Z_vec = lp.Set_V(mat.NewVecDense(len(lpi.B), nil), lpi.Z_vec, lpi.B)
		// zn <- complicated shit
		lpi.Z_vec = lp.Set_V(lpi.Make_Z_N(), lpi.Z_vec, lpi.N)

		if mat.Min(lpi.Z_N()) > 0 {
			fmt.Fprintf(os.Stderr, "Found optimal")
			break
		}

		// Choose entering variable
		zn_i := lp.Min_Index(lpi.Z_N())
		j := lpi.N[zn_i]

		lp.Debug("Z", lpi.Z_vec)
		lp.Debug("Zn", lpi.Z_N())

		// Choosing a leaving variable

		// Construct theta x vector
		lpi.TX_vec = lp.Set_V(lpi.Make_TX_B(j), lpi.TX_vec, lpi.B)

		tXB := lp.Get_V(lpi.TX_vec, lpi.B)
		XB := lp.Get_V(lpi.X_vec, lpi.B)

		// Creating xb/txb
		v := mat.NewVecDense(tXB.Len(), nil)
		v.DivElemVec(XB, tXB)

		lp.Debug("X", lpi.X_vec)
		lp.Debug("XB", XB)
		lp.Debug("tXB", tXB)

		// Find min index for t
		xb_i := lp.Min_Index(v)
		t := v.AtVec(xb_i)
		i := lpi.B[xb_i]

		fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, lpi.Z_vec.AtVec(j))
		fmt.Fprintf(os.Stderr, "i = %v, xi =  %v\n", i, lpi.X_vec.AtVec(i))
		fmt.Fprintf(os.Stderr, "t = %v\n", t)

		// j = 0
		// i = 3

		// Updating xb
		v2 := mat.NewVecDense(XB.Len(), nil)
		tXB.ScaleVec(t, tXB)
		v2.SubVec(XB, tXB)
		lpi.X_vec = lp.Set_V(v2, lpi.X_vec, lpi.B)

		lpi.X_vec.SetVec(j, t)

		lpi.B = lp.Swap(j, i, lpi.B)
		lpi.N = lp.Swap(i, j, lpi.N)

		time.Sleep(1 * time.Second)
		fmt.Fprintf(os.Stderr, "-----------------\n\n")
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
