package util

import "gonum.org/v1/gonum/mat"

// Transpose : 転置行列の作成
// TODO : Matrixをdenseに変換して、Tを実行すれば正常に転置されるかも(よってこのAPIは不要)
func Transpose(a mat.Matrix) mat.Matrix {
	r, c := a.Dims()
	dense := mat.NewDense(c, r, nil)
	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
			dense.Set(i, j, a.At(j, i))
		}
	}
	return dense
}

// Reshape : 行列を指定した形式の行列に変換する
// ex) 6 * 6行列 => 2 * 19行列
func Reshape(base mat.Matrix, r, c int) mat.Matrix {
	bR, bC := base.Dims()
	if bR*bC != r*c {
		panic("Reshape error : 変換前と変換後の要素数に差分があります")
	}

	rawValues := RawValues(base)
	return mat.NewDense(r, c, rawValues)
}

// RawValues : matrixの元データをfloat64の配列に格納
// 配列サイズはrow * colとなる
func RawValues(a mat.Matrix) []float64 {
	dense := mat.DenseCopyOf(a)
	r, c := dense.Dims()
	rawValues := make([]float64, 0, r*c)

	for i := 0; i < r; i++ {
		rawValues = append(rawValues, dense.RawRowView(i)...)
	}
	return rawValues
}

// AddVecToMatrixCol : vectorを行方向にmatrixに加算する
// ex) 3次元のベクトルを4*3の行列に加算するケース
func AddVecToMatrixCol(d *mat.Dense, vec mat.Vector) {
	_, c := d.Dims()
	if c != vec.Len() {
		panic("AddVecToMatrixCol Error : ベクトルの次元数と行列の列数が一致していません")
	}

	d.Apply(func(i, j int, v float64) float64 {
		return vec.AtVec(j) + v
	}, d)
}

// SumEachCol : 行列の各列について、合計値を計算し、ベクトルに変換する
// ex) 3*3の行列の場合、1*3のベクトルとなる。
// [1,1,1]
// [2,2,2]
// [3,3,3]
// = [6,6,6]
func SumEachCol(x mat.Matrix) mat.Vector {
	r, c := x.Dims()
	vec := mat.NewVecDense(c, nil)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			tmpVal := vec.AtVec(j)
			vec.SetVec(j, tmpVal+x.At(i, j))
		}
	}
	return vec
}

// SumEachRow : 行列の各行について、合計値を計算し、ベクトルに変換する
// ex) 3*3の行列の場合、3*1のベクトルとなる
func SumEachRow(x mat.Matrix) mat.Vector {
	r, c := x.Dims()
	vec := mat.NewVecDense(r, nil)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			tmpVal := vec.AtVec(i)
			vec.SetVec(i, tmpVal+x.At(i, j))
		}
	}
	return vec
}

// MaxEachRow : 行列の各行について、最大値を選択し、ベクトルに変換する
// また各最大値の列数を格納する
// ex) 3*3の行列の場合、3*1のベクトルとなる
// 一つ目の戻り値 : 各行の最大値を格納したベクトル
// 二つ目の戻り値 : 各行の最大値の列数を格納した配列
func MaxEachRow(x mat.Matrix) (mat.Vector, []int) {
	r, _ := x.Dims()
	vec := mat.NewVecDense(r, nil)
	args := make([]int, 0, r)
	dense := mat.DenseCopyOf(x)
	for i := 0; i < r; i++ {
		// 各行の最大値と該当する値の列数を取得する
		rawRows := dense.RawRowView(i)
		k, v := MaxValue(rawRows)

		// 値を設定する
		vec.SetVec(i, v)
		args = append(args, k)
	}
	return vec, args
}
