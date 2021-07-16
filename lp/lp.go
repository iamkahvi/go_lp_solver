package lp

import (
	"fmt"
	"os"

	mat "gonum.org/v1/gonum/mat"
)

type lp struct {
	A     *mat.Dense
	r     int
	c     int
	b_vec *mat.VecDense
	c_vec *mat.VecDense
	B     []int
	N     []int
}

func (lp lp) Print() {
	fm := mat.Formatted(lp.A, mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "A = %v\n\n", fm)

	fm = mat.Formatted(lp.b_vec, mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "b = %v\n\n", fm)

	fm = mat.Formatted(lp.c_vec.T(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "c = %v\n\n", fm)

	fmt.Fprintf(os.Stderr, "B = %v\n\n", lp.B)
	fmt.Fprintf(os.Stderr, "N = %v\n\n", lp.N)
}

func (lp lp) Get_Ab() *mat.Dense {
	n := mat.NewDense(lp.r, len(lp.B), nil)
	for i, ind := range lp.B {
		col := make([]float64, lp.r)

		mat.Col(col, ind, lp.A)
		n.SetCol(i, col)
	}
	return n
}

func (lp lp) Get_An() *mat.Dense {
	n := mat.NewDense(lp.r, len(lp.N), nil)
	for i, ind := range lp.N {
		col := make([]float64, lp.r)

		mat.Col(col, ind, lp.A)
		n.SetCol(i, col)
	}
	return n
}

// func (lp lp) Get_xb() *mat.VecDense {
// }

// func (lp lp) Get_xn() *mat.VecDense {

// }

// func (lp lp) Get_cb() *mat.VecDense {

// }

// func (lp lp) Get_cn() *mat.VecDense {

// }

func New(matr [][]float64, r int, c int) *lp {
	m := r - 1
	n := c - 1

	A := mat.NewDense(m, n+m, nil)
	b_vec := mat.NewVecDense(m, nil)
	c_vec := mat.NewVecDense(n+m, nil)
	B := make([]int, m)
	N := make([]int, n)

	for i := 0; i < n+m; i++ {
		if i < n {
			c_vec.SetVec(i, matr[0][i])
			N[i] = i
		} else {
			c_vec.SetVec(i, 0)
			B[i-n] = i
		}
	}

	for i := 1; i <= m; i++ {
		for j := 0; j < n+m; j++ {
			if j < n {
				A.Set(i-1, j, matr[i][j])
			} else {
				if j%n == i-1 {
					A.Set(i-1, j, 1)
				}
			}
		}
		b_vec.SetVec(i-1, matr[i][n])
	}

	rA, cA := A.Dims()

	return &lp{
		A:     A,
		r:     rA,
		c:     cA,
		b_vec: b_vec,
		c_vec: c_vec,
		B:     B,
		N:     N,
	}
}
