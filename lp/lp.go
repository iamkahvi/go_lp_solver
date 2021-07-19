package lp

import (
	"fmt"
	"os"

	"example.com/solver/utils"
	mat "gonum.org/v1/gonum/mat"
)

type LP struct {
	A     *mat.Dense
	r     int
	c     int
	B_vec *mat.VecDense
	C_vec *mat.VecDense
	X_vec *mat.VecDense
	Z_vec *mat.VecDense
	B     []int
	N     []int
}

func New(matr [][]float64, r int, c int) *LP {
	m := r - 1
	n := c - 1

	A := mat.NewDense(m, n+m, nil)
	B_vec := mat.NewVecDense(m, nil)
	C_vec := mat.NewVecDense(n+m, nil)
	X_vec := mat.NewVecDense(n+m, nil)
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

	return &LP{
		A:     A,
		r:     m,
		c:     n + m,
		B_vec: B_vec,
		C_vec: C_vec,
		X_vec: X_vec,
		Z_vec: Z_vec,
		B:     B,
		N:     N,
	}
}

func (lp LP) CloneAux() *LP {
	var B2_vec mat.VecDense
	B2_vec.CloneFromVec(lp.B_vec)

	C2_vec := mat.NewVecDense(lp.C_vec.Len(), nil)

	var X2_vec mat.VecDense
	X2_vec.CloneFromVec(lp.X_vec)

	var Z2_vec mat.VecDense
	Z2_vec.CloneFromVec(lp.Z_vec)

	B2 := make([]int, len(lp.B))
	copy(B2, lp.B)

	N2 := make([]int, len(lp.N))
	copy(N2, lp.N)

	return &LP{
		A:     lp.A,
		r:     lp.r,
		c:     lp.c,
		B_vec: &B2_vec,
		C_vec: C2_vec,
		X_vec: &X2_vec,
		Z_vec: &Z2_vec,
		B:     B2,
		N:     N2,
	}
}

func (lp LP) Print() {
	fmt.Fprintf(os.Stderr, " B = %v\n\n", lp.B)
	fmt.Fprintf(os.Stderr, " N = %v\n\n", lp.N)

	Debug("A", lp.A)

	Debug("A_B", lp.A_B())
	Debug("A_N", lp.A_N())

	Debug("x_B", lp.X_B())

	Debug("b", lp.B_vec)
	Debug("c", lp.C_vec)

	Debug("X", lp.X_vec)
	Debug("Z", lp.Z_vec)
}

func Debug(s string, m mat.Matrix) {
	fm := mat.Formatted(m, mat.Prefix("       "), mat.Squeeze())
	fmt.Fprintf(os.Stderr, "%4s = %v\n\n", s, fm)
}

func (lp LP) A_B() *mat.Dense {
	return utils.Get_M(lp.A, lp.B)
}

func (lp LP) A_N() *mat.Dense {
	return utils.Get_M(lp.A, lp.N)
}

func (lp LP) X_B() *mat.VecDense {
	return utils.Get_V(lp.X_vec, lp.B)
}

func (lp LP) X_N() *mat.VecDense {
	return utils.Get_V(lp.X_vec, lp.N)
}

func (lp LP) C_B() *mat.VecDense {
	return utils.Get_V(lp.C_vec, lp.B)
}

func (lp LP) C_N() *mat.VecDense {
	return utils.Get_V(lp.C_vec, lp.N)
}

func (lp LP) Z_B() *mat.VecDense {
	return utils.Get_V(lp.Z_vec, lp.B)
}

func (lp LP) Z_N() *mat.VecDense {
	return utils.Get_V(lp.Z_vec, lp.N)
}

func (lp LP) Is_Primal_Feasible() bool {
	return mat.Min(lp.B_vec) >= 0
}

func (lp LP) Is_Dual_Feasible() bool {
	return mat.Max(lp.C_vec) <= 0
}

func (lp LP) Make_Z_N() *mat.VecDense {
	var n mat.Dense
	n.Inverse(lp.A_B())

	var n2 mat.Dense
	n2.Mul(&n, lp.A_N())

	var nv mat.VecDense
	nv.MulVec(n2.T(), lp.C_B())

	var nv2 mat.VecDense
	nv2.SubVec(&nv, lp.C_N())

	return &nv2
}

func (lp LP) Make_TX_B(j int) *mat.VecDense {
	n := mat.NewDense(lp.r, lp.r, nil)
	n.Inverse(lp.A_B())

	aj := lp.A.ColView(j)
	//  Get_M(lp.A, []int{j}).ColView(0)

	v := mat.NewVecDense(aj.Len(), nil)
	v.MulVec(n, aj)

	return v
}

func (lp LP) Make_X_B() *mat.VecDense {
	// Setting xb
	n1 := mat.NewDense(lp.r, lp.r, nil)
	n1.Inverse(lp.A_B())

	xb := mat.NewVecDense(lp.r, nil)
	xb.MulVec(n1, lp.B_vec)

	return xb
}

func (lp LP) Make_Theta_X_B(j int) *mat.VecDense {
	col := make([]float64, lp.r)
	mat.Col(col, j, lp.A)

	nv := mat.NewVecDense(lp.r, col)
	n := mat.NewDense(lp.r, lp.r, nil)

	n.Inverse(lp.A_B())
	nv.MulVec(n, nv)

	return nv
}
