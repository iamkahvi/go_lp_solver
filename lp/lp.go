package lp

import (
	"fmt"
	"os"

	mat "gonum.org/v1/gonum/mat"
)

type lp struct {
	A      *mat.Dense
	r      int
	c      int
	B_vec  *mat.VecDense
	C_vec  *mat.VecDense
	X_vec  *mat.VecDense
	DX_vec *mat.VecDense
	Z_vec  *mat.VecDense
	B      []int
	N      []int
}

func New(matr [][]float64, r int, c int) *lp {
	m := r - 1
	n := c - 1

	A := mat.NewDense(m, n+m, nil)
	B_vec := mat.NewVecDense(m, nil)
	C_vec := mat.NewVecDense(n+m, nil)
	X_vec := mat.NewVecDense(n+m, nil)
	DX_vec := mat.NewVecDense(n+m, nil)
	Z_vec := mat.NewVecDense(n+m, nil)
	B := make([]int, m)
	N := make([]int, n)

	for i := 0; i < n+m; i++ {
		if i < n {
			C_vec.SetVec(i, matr[0][i])
			N[i] = i
		} else {
			C_vec.SetVec(i, 0)
			B[i-n] = i
		}
	}

	for i := 1; i <= m; i++ {
		for j := 0; j < n+m; j++ {
			if j < n {
				A.Set(i-1, j, matr[i][j])
			} else {
				if j-n == i-1 {
					A.Set(i-1, j, 1)
				}
			}
		}
		B_vec.SetVec(i-1, matr[i][n])
	}

	return &lp{
		A:      A,
		r:      m,
		c:      n + m,
		B_vec:  B_vec,
		C_vec:  C_vec,
		X_vec:  X_vec,
		DX_vec: DX_vec,
		Z_vec:  Z_vec,
		B:      B,
		N:      N,
	}
}

func (lp lp) Print() {
	fmt.Fprintf(os.Stderr, " B = %v\n\n", lp.B)
	fmt.Fprintf(os.Stderr, " N = %v\n\n", lp.N)

	// Debug("A", lp.A)

	// Debug("A_B", lp.A_B())
	// Debug("A_N", lp.A_N())

	// Debug("x_B", lp.X_B())

	// Debug("c_B", lp.C_B())
	// Debug("c_N", lp.C_N())

	// Debug("X", lp.X_vec)
}

func Debug(s string, m mat.Matrix) {
	fm := mat.Formatted(m, mat.Prefix("       "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "%4s = %v\n\n", s, fm)
}

func (lp lp) A_B() *mat.Dense {
	return Get_M(lp.A, lp.B)
}

func (lp lp) A_N() *mat.Dense {
	return Get_M(lp.A, lp.N)
}

func (lp lp) X_B() *mat.VecDense {
	return Get_V(lp.X_vec, lp.B)
}

func (lp lp) X_N() *mat.VecDense {
	return Get_V(lp.X_vec, lp.N)
}

func (lp lp) C_B() *mat.VecDense {
	return Get_V(lp.C_vec, lp.B)
}

func (lp lp) C_N() *mat.VecDense {
	return Get_V(lp.C_vec, lp.N)
}

func (lp lp) Z_B() *mat.VecDense {
	return Get_V(lp.Z_vec, lp.B)
}

func (lp lp) Z_N() *mat.VecDense {
	return Get_V(lp.Z_vec, lp.N)
}

func (lp lp) Make_Z_N() *mat.VecDense {
	n := mat.NewDense(lp.r, lp.r, nil)
	n.Inverse(lp.A_B())

	n2 := mat.NewDense(lp.r, len(lp.N), nil)
	n2.Mul(n, lp.A_N())

	nv := mat.NewVecDense(lp.c-lp.r, nil)
	nv.MulVec(n2.T(), lp.C_B())

	nv2 := mat.NewVecDense(len(lp.N), nil)
	nv2.SubVec(nv, lp.C_N())

	return nv2
}

func (lp lp) Make_TX_B(j int) *mat.VecDense {
	n := mat.NewDense(lp.r, lp.r, nil)
	n.Inverse(lp.A_B())

	aj := lp.A.ColView(j)
	//  Get_M(lp.A, []int{j}).ColView(0)

	v := mat.NewVecDense(aj.Len(), nil)
	v.MulVec(n, aj)

	return v
}

func (lp lp) Make_X_B() *mat.VecDense {
	// Setting xb
	n1 := mat.NewDense(lp.r, lp.r, nil)
	n1.Inverse(lp.A_B())

	xb := mat.NewVecDense(lp.r, nil)
	xb.MulVec(n1, lp.B_vec)

	return xb
}

func (lp lp) Is_Feasible() bool {
	return mat.Min(lp.X_B()) >= 0
}

func (lp lp) Make_Theta_X_B(j int) *mat.VecDense {
	col := make([]float64, lp.r)
	mat.Col(col, j, lp.A)

	nv := mat.NewVecDense(lp.r, col)
	n := mat.NewDense(lp.r, lp.r, nil)

	n.Inverse(lp.A_B())
	nv.MulVec(n, nv)

	return nv
}
