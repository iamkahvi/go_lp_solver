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
	fmt.Fprintf(os.Stderr, "B = %v\n\n", lp.B)
	fmt.Fprintf(os.Stderr, "N = %v\n\n", lp.N)

	fm := mat.Formatted(lp.Get_Ab(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "Ab= %v\n\n", fm)

	fm = mat.Formatted(lp.Get_An(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "An= %v\n\n", fm)

	fm = mat.Formatted(lp.Get_xb(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "xb= %v\n\n", fm)

	fm = mat.Formatted(lp.Get_cb(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "cb= %v\n\n", fm)

	fm = mat.Formatted(lp.Get_cn(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "cn= %v\n\n", fm)

	fm = mat.Formatted(lp.Get_zn(), mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "zn= %v\n\n", fm)
}

func debug(m mat.Matrix) {
	fm := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "m = %v\n\n", fm)
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

func (lp lp) Get_xb() *mat.VecDense {
	n := mat.NewDense(lp.r, lp.r, nil)
	n.Inverse(lp.Get_Ab())

	nv := mat.NewVecDense(lp.r, nil)
	nv.MulVec(n, lp.b_vec)

	return nv
}

func (lp lp) Get_cb() *mat.VecDense {
	n := mat.NewVecDense(len(lp.B), nil)

	for i, ind := range lp.B {
		n.SetVec(i, lp.c_vec.AtVec(ind))
	}

	return n
}

func (lp lp) Get_zn() *mat.VecDense {
	n := mat.NewDense(lp.r, lp.r, nil)
	n.Inverse(lp.Get_Ab())

	n2 := mat.NewDense(lp.r, len(lp.N), nil)
	n2.Mul(n, lp.Get_An())

	nv := mat.NewVecDense(lp.c-lp.r, nil)
	nv.MulVec(n2.T(), lp.Get_cb())

	nv2 := mat.NewVecDense(len(lp.N), nil)
	nv2.SubVec(nv, lp.Get_cn())

	return nv2
}

func (lp lp) Get_cn() *mat.VecDense {
	n := mat.NewVecDense(len(lp.N), nil)

	for i, ind := range lp.N {
		n.SetVec(i, lp.c_vec.AtVec(ind))
	}

	return n
}

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

	return &lp{
		A:     A,
		r:     m,
		c:     n + m,
		b_vec: b_vec,
		c_vec: c_vec,
		B:     B,
		N:     N,
	}
}
