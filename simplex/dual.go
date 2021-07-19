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

	// Compute Z
	l.Z_vec = utils.Set_V(mat.NewVecDense(len(l.B), nil), l.Z_vec, l.B)
	l.Z_vec = utils.Set_V(l.Make_Z_N(), l.Z_vec, l.N)

	if mat.Min(l.Z_N()) < -EPSILON {
		return Infeasible, 0, nil
	}

	for {
		if DEBUG {
			fmt.Fprintf(os.Stderr, "\niteration %v-----------------\n\n", iteration)
			l.Print()
		}

		// Compute X
		l.X_vec = utils.Set_V(l.Make_X_B(), l.X_vec, l.B)
		l.X_vec = utils.Set_V(mat.NewVecDense(len(l.N), nil), l.X_vec, l.N)

		// Check for optimality
		if mat.Min(l.X_B()) >= -EPSILON {
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
			if l.X_vec.AtVec(ind) < -EPSILON {
				i = ind
				break
			}
		}

		// Choose entering variable

		// Compute vector u
		u := mat.NewVecDense(len(l.B), nil)

		for idx, val := range l.B {
			if val == i {
				u.SetVec(idx, 1)
			}
		}

		// Compute delta Z vector
		dZ := mat.NewVecDense(l.Z_vec.Len(), nil)
		dZ = utils.Set_V(l.Make_DZ_N(u), dZ, l.N)

		// Find min for s and j
		s := math.MaxFloat64
		j := 0
		for _, nVal := range l.N {
			x := l.Z_vec.AtVec(nVal)
			dx := dZ.AtVec(nVal)

			if dx > EPSILON {
				val := x / dx
				if val < s {
					s = val
					j = nVal
				}
			}
		}

		// Check for unboundedness/infeasibility
		if mat.Max(dZ) <= EPSILON {
			return Infeasible, 0, nil
		}

		// Update zn
		var v mat.VecDense
		v.SubVec(l.Z_N(), utils.Get_V(dZ, l.N))
		l.Z_vec = utils.Set_V(&v, l.Z_vec, l.N)

		// Set s at zi
		l.Z_vec.SetVec(i, s)

		l.B = utils.Swap(j, i, l.B)
		l.N = utils.Swap(i, j, l.N)

		// Update B and N indices
		iteration++

		if DEBUG {
			fmt.Fprintf(os.Stderr, "i = %v, xi = %v\n", i, l.X_vec.AtVec(i))
			fmt.Fprintf(os.Stderr, "j = %v, zj = %v\n", j, l.Z_vec.AtVec(j))
			fmt.Fprintf(os.Stderr, "s = %v\n", s)
			time.Sleep(1 * time.Second)
		}
	}
}
