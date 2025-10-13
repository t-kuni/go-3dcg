package util

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// TransformTranslate
// 第一引数mは４行である必要がある
func TransformTranslate(m mat.Dense, x, y, z float64) mat.Dense {
	translateMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	})

	var translated mat.Dense
	translated.Mul(translateMatrix, &m)
	return translated
}

// TransformRotate は3つの軸（X、Y、Z）での回転を適用します
// x, y, z はそれぞれの軸での回転角度（ラジアン）です
// 第一引数mは４行である必要がある
func TransformRotate(m mat.Dense, x, y, z float64) mat.Dense {
	// X軸回転行列
	rotateXMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, math.Cos(x), -math.Sin(x), 0,
		0, math.Sin(x), math.Cos(x), 0,
		0, 0, 0, 1,
	})

	// Y軸回転行列
	rotateYMatrix := mat.NewDense(4, 4, []float64{
		math.Cos(y), 0, math.Sin(y), 0,
		0, 1, 0, 0,
		-math.Sin(y), 0, math.Cos(y), 0,
		0, 0, 0, 1,
	})

	// Z軸回転行列
	rotateZMatrix := mat.NewDense(4, 4, []float64{
		math.Cos(z), -math.Sin(z), 0, 0,
		math.Sin(z), math.Cos(z), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	// 回転を順番に適用: Z軸 -> Y軸 -> X軸の順で回転
	var temp1, temp2, rotated mat.Dense
	temp1.Mul(rotateZMatrix, &m)
	temp2.Mul(rotateYMatrix, &temp1)
	rotated.Mul(rotateXMatrix, &temp2)

	return rotated
}

// TransformParallelProjection 平行投影
// 第一引数mは４行である必要がある
func TransformParallelProjection(m mat.Dense) mat.Dense {
	projectionMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0, // X軸
		0, 1, 0, 0, // Y軸
		0, 0, 0, 0, // Z軸（無視）
		0, 0, 0, 1, // 同次座標
	})

	var projected mat.Dense
	projected.Mul(projectionMatrix, &m)

	return projected
}

// transformScale は3次元オブジェクトの拡大・縮小変換を行います
// scaleX, scaleY, scaleZ はそれぞれのX、Y、Z軸方向の拡大率です
// 第一引数mは４行である必要がある
func TransformScale(m mat.Dense, scaleX, scaleY, scaleZ float64) mat.Dense {
	scaleMatrix := mat.NewDense(4, 4, []float64{
		scaleX, 0, 0, 0,
		0, scaleY, 0, 0,
		0, 0, scaleZ, 0,
		0, 0, 0, 1,
	})

	var scaled mat.Dense
	scaled.Mul(scaleMatrix, &m)
	return scaled
}

// transformScaleUniform は均等な拡大・縮小変換を行います
// scale は全軸方向の拡大率です
// 第一引数mは４行である必要がある
func TransformScaleUniform(m mat.Dense, scale float64) mat.Dense {
	return TransformScale(m, scale, scale, scale)
}

// TransformViewport はビューポート変換を行います
// 第一引数mは４行である必要がある
func TransformViewport(m mat.Dense, width, height int32) mat.Dense {
	// 短辺を基準にスケールを決める
	scale := math.Min(float64(width), float64(height)) / 2

	m = TransformScale(m, scale, scale, 1)

	m = TransformTranslate(m, float64(width)/2, float64(height)/2, 0)

	return m
}

func T(m mat.Dense) mat.Dense {
	return *mat.DenseCopyOf((&m).T())
}
