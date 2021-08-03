package utils

import (
	"math"
	"strconv"

	s "strings"

	"gonum.org/v1/gonum/blas/blas64"
	mat "gonum.org/v1/gonum/mat"
)

func ParseLines(lines []string) ([][]float64, int, int) {
	rows := 0
	cols := 0
	for _, line := range lines {
		l := len(s.Fields(line))
		if l > 1 {
			rows += 1
		}
		if l > cols {
			cols = l
		}
	}

	numbers := make([][]float64, rows)

	for i, line := range lines {
		els := s.Fields(line)

		numbers[i] = make([]float64, cols)

		for j, str := range els {
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				panic(err)
			}
			numbers[i][j] = val
		}
	}

	return numbers, rows, cols
}

func Swap(in int, out int, arr []int) []int {
	for i, el := range arr {
		if el == out {
			arr[i] = in
		}
	}

	return arr
}

func Max_Index(v *mat.VecDense) int {
	return blas64.Iamax(v.RawVector())
}

func Min_Non_Neg_Index(v *mat.VecDense) int {
	ind := 0
	min := math.MaxFloat64
	for i := 1; i < v.Len(); i++ {
		x := v.AtVec(i)
		if x < min && x > 0 {
			min = x
			ind = i
		}
	}

	return ind
}

func Get_M(m *mat.Dense, ind []int) *mat.Dense {
	r, _ := m.Dims()
	n := mat.NewDense(r, len(ind), nil)

	for i, ind := range ind {
		col := make([]float64, r)

		mat.Col(col, ind, m)
		n.SetCol(i, col)
	}

	return n
}

func Get_V(v *mat.VecDense, ind []int) *mat.VecDense {
	n := mat.NewVecDense(len(ind), nil)

	for i, ind := range ind {
		n.SetVec(i, v.AtVec(ind))
	}

	return n
}

func Set_M(new *mat.Dense, m *mat.Dense, ind []int) *mat.Dense {
	r, _ := m.Dims()

	for i, ind := range ind {
		col := make([]float64, r)

		mat.Col(col, i, new)
		m.SetCol(ind, col)
	}

	return m
}

func Set_V(new *mat.VecDense, m *mat.VecDense, ind []int) *mat.VecDense {
	for i, ind := range ind {
		m.SetVec(ind, new.AtVec(i))
	}

	return m
}
