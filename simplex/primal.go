package simplex

import (
	"fmt"
	"math"
	"os"
	"time"

	"example.com/solver/lp"
	utils "example.com/solver/utils"
	mat "gonum.org/v1/gonum/mat"
)

type Result int

const (
	Optimal Result = iota
	Unbounded
	Infeasible
)

func PrimalSimplex(l *lp.LP, DEBUG bool) (Result, float64, []float64) {
	iteration := 0

	// Compute X
	l.X_vec = utils.Set_V(l.Make_X_B(), l.X_vec, l.B)
	l.X_vec = utils.Set_V(mat.NewVecDense(len(l.N), nil), l.X_vec, l.N)

	if mat.Min(l.X_B()) < 0 {
		return Infeasible, 0, nil
	}

	for {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "\niteration %v-----------------\n\n", iteration)
			l.Print()
		}

		// Compute Z
		l.Z_vec = utils.Set_V(mat.NewVecDense(len(l.B), nil), l.Z_vec, l.B)
		l.Z_vec = utils.Set_V(l.Make_Z_N(), l.Z_vec, l.N)

		// Check for optimality
		if mat.Min(l.Z_N()) >= 0 {
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
		// zn_i := utils.Min_Index(lpi.Z_N())
		// j := lpi.N[zn_i]

		// Bland's rule
		var j int
		for _, ind := range l.N {
			if l.Z_vec.AtVec(ind) < 0 {
				j = ind
				break
			}
		}

		// Choose a leaving variable

		// Construct theta x vector
		dX := mat.NewVecDense(l.X_vec.Len(), nil)
		dX = utils.Set_V(l.Make_DX_B(j), dX, l.B)

		// Find min for i and t
		t := math.MaxFloat64
		i := 0
		for _, bVal := range l.B {
			xi := l.X_vec.AtVec(bVal)
			dxi := dX.AtVec(bVal)

			if dxi > 0 {
				val := xi / dxi
				if val < t {
					t = val
					i = bVal
				}
			}
		}

		// Check for unboundedness
		if mat.Max(dX) <= 0 {
			return Unbounded, 0, nil
		}

		dXB := utils.Get_V(dX, l.B)
		XB := utils.Get_V(l.X_vec, l.B)

		// Update xb
		v2 := mat.NewVecDense(XB.Len(), nil)
		dXB.ScaleVec(t, dXB)
		v2.SubVec(XB, dXB)
		l.X_vec = utils.Set_V(v2, l.X_vec, l.B)

		// Set t at xj
		l.X_vec.SetVec(j, t)

		// Update B and N indices
		l.B = utils.Swap(j, i, l.B)
		l.N = utils.Swap(i, j, l.N)

		iteration++

		if DEBUG {
			fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, l.Z_vec.AtVec(j))
			fmt.Fprintf(os.Stderr, "i = %v, xi = %v\n", i, l.X_vec.AtVec(i))
			fmt.Fprintf(os.Stderr, "t = %v\n", t)

			time.Sleep(1 * time.Second)
		}
	}
}
