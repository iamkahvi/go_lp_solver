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

func DualSimplex(l *lp.LP, DEBUG bool) (Result, float64, []float64) {
	iteration := 0

	// Setting zn and zb
	l.Z_vec = utils.Set_V(mat.NewVecDense(len(l.B), nil), l.Z_vec, l.B)
	l.Z_vec = utils.Set_V(l.Make_Z_N(), l.Z_vec, l.N)

	if mat.Min(l.Z_N()) < 0 {
		return Infeasible, 0, nil
	}

	for {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "\niteration %v-----------------\n\n", iteration)
		}

		// Setting xb and xn
		l.X_vec = utils.Set_V(l.Make_X_B(), l.X_vec, l.B)
		l.X_vec = utils.Set_V(mat.NewVecDense(len(l.N), nil), l.X_vec, l.N)

		if mat.Min(l.X_B()) >= 0 {
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
		l.Print()

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

		dZ = utils.Set_V(dZB, dZ, l.N)
		dZ = utils.Set_V(mat.NewVecDense(len(l.B), nil), dZ, l.B)

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
		v.SubVec(l.Z_N(), utils.Get_V(dZ, l.N))
		l.Z_vec = utils.Set_V(&v, l.Z_vec, l.N)
		l.Z_vec.SetVec(i, s)

		l.B = utils.Swap(j, i, l.B)
		l.N = utils.Swap(i, j, l.N)

		if DEBUG {
			fmt.Fprintf(os.Stderr, "i = %v, xi = %v\n", i, l.X_vec.AtVec(i))
			fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, l.Z_vec.AtVec(j))
			fmt.Fprintf(os.Stderr, "s = %v\n", s)
		}

		iteration++

		if DEBUG {
			time.Sleep(1 * time.Second)
		}
	}
}
