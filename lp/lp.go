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

	fmt.Println(lp.B)
	fmt.Println(lp.N)
}

func (lp lp) Get_Ab() *mat.Dense {
	n := mat.NewDense(lp.r, len(lp.B), nil)
	for _, ind := range lp.B {
		fmt.Fprintln(os.Stderr, lp.A.ColView(ind).(*mat.VecDense).RawVector().Data)
		// n.SetCol(i, lp.A.ColView(ind).(*mat.VecDense).RawVector().Data)
	}
	return n
}

// func (lp lp) Get_An() *mat.Dense {

// }

// func (lp lp) Get_xb() *mat.VecDense {

// }

// func (lp lp) Get_xn() *mat.VecDense {

// }

// func (lp lp) Get_cb() *mat.VecDense {

// }

// func (lp lp) Get_cn() *mat.VecDense {

// }

func New(m [][]float64, r int, c int) *lp {
	width := r + c - 2
	height := r - 1

	A := mat.NewDense(height, width, nil)
	b_vec := mat.NewVecDense(height, nil)
	c_vec := mat.NewVecDense(width, nil)
	B := make([]int, r)
	N := make([]int, width-r)

	for i := 0; i < width; i++ {
		if i <= len(m[0])-1 {
			c_vec.SetVec(i, m[0][i])
		} else {
			c_vec.SetVec(i, 0)
		}

		if i < r {
			B[i] = i
		} else {
			N[i%r] = i
		}
	}

	for i := 1; i <= height; i++ {
		for j := 0; j < width; j++ {
			if j < c-1 {
				A.Set(i-1, j, m[i][j])
			} else {
				if j%(c-1) == i-1 {
					A.Set(i-1, j, 1)
				}
			}
		}
		b_vec.SetVec(i-1, m[i][c-1])
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
