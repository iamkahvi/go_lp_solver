package lp

import (
	"fmt"
	"os"

	"example.com/solver/utils"
	mat "gonum.org/v1/gonum/mat"
)

const EPSILON = 1e-5

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
	Debug("z_N", lp.Z_N())

	Debug("b", lp.B_vec)
	Debug("c", lp.C_vec)

	Debug("c_B", lp.C_B())

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
	return mat.Min(lp.B_vec) >= -EPSILON
}

func (lp LP) Is_Dual_Feasible() bool {
	return mat.Max(lp.C_vec) <= -EPSILON
}

func (lp LP) Make_Z_N() *mat.VecDense {
	var v mat.VecDense
	v.SolveVec(lp.A_B().T(), lp.C_B())

	var temp mat.VecDense
	temp.MulVec(lp.A_N().T(), &v)

	var zn mat.VecDense
	zn.SubVec(&temp, lp.C_N())

	return &zn
}

func (lp LP) Make_DX_B(j int) *mat.VecDense {
	var dxB mat.VecDense

	dxB.SolveVec(lp.A_B(), lp.A.ColView(j))

	return &dxB
}

func (lp LP) Make_DZ_N(u *mat.VecDense) *mat.VecDense {
	var temp mat.Dense
	temp.Inverse(lp.A_B().T())

	var rh mat.VecDense
	rh.MulVec(&temp, u)

	var negA_NT mat.Dense
	negA_NT.Scale(-1, lp.A_N().T())

	var res mat.VecDense
	res.MulVec(&negA_NT, &rh)

	return &res

	// var negA_NT mat.Dense
	// negA_NT.Scale(-1, lp.A_N().T())

	// var temp mat.VecDense
	// temp.MulVec(&negA_NT, u)

	// var dzN mat.VecDense
	// dzN.SolveVec(lp.A_B().T(), &temp)

	// return &dzN
}

func (lp LP) Make_X_B() *mat.VecDense {
	xb := mat.NewVecDense(lp.r, nil)
	xb.SolveVec(lp.A_B(), lp.B_vec)

	return xb
}
