package lp

import (
	mat "gonum.org/v1/gonum/mat"
)

type lp struct {
	Ab *mat.Dense
	An *mat.Dense
	x  *mat.VecDense
	B  []int
	N  []int
}

func (lp lp) Get_xb() *mat.VecDense {
	return lp.x
}

func New(rows int, cols int) *lp {
	return &lp{
		Ab: mat.NewDense(rows, cols, nil),
		An: mat.NewDense(rows, cols, nil),
		x:  mat.NewVecDense(cols, nil),
		B:  make([]int, cols),
		N:  make([]int, cols),
	}
}
