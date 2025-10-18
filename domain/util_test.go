package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestCalcNormalFromPoints(t *testing.T) {
	// 正常系テストケース: 3つの点から法線ベクトルを計算
	// 三角形の頂点: (0,0,0), (1,0,0), (0,1,0)
	// この三角形の法線ベクトルは左ネジの法則に従い (0,0,-1) になるはず
	p1 := Vector3D{0, 0, 0}
	p2 := Vector3D{1, 0, 0}
	p3 := Vector3D{0, 1, 0}

	result := CalcNormalFromPoints(p1, p2, p3)

	assert.InDelta(t, 0, result.X(), 0.001)
	assert.InDelta(t, 0, result.Y(), 0.001)
	assert.InDelta(t, -1, result.Z(), 0.001)
}

func TestCrossVecDense(t *testing.T) {
	// 正常系テストケース: 2つのベクトルの外積を計算
	// ベクトル a = (1, 0, 0), ベクトル b = (0, 1, 0)
	// 外積 a × b = (0, 0, 1)
	a := mat.NewVecDense(3, []float64{1, 0, 0})
	b := mat.NewVecDense(3, []float64{0, 1, 0})

	result := CrossVecDense(*a, *b)

	assert.InDelta(t, 0, result.At(0, 0), 0.001)
	assert.InDelta(t, 0, result.At(1, 0), 0.001)
	assert.InDelta(t, 1, result.At(2, 0), 0.001)
}

func TestNormalizeVecDense(t *testing.T) {
	// 正常系テストケース: ベクトルを正規化
	// ベクトル v = (3, 4, 0)
	// 正規化されたベクトル = (0.6, 0.8, 0)
	v := mat.NewVecDense(3, []float64{3, 4, 0})

	result := NormalizeVecDense(*v)

	assert.InDelta(t, 0.6, result.At(0, 0), 0.001)
	assert.InDelta(t, 0.8, result.At(1, 0), 0.001)
	assert.InDelta(t, 0, result.At(2, 0), 0.001)
}

func TestClassifyEdgeByPlane(t *testing.T) {
	// 正常系テストケース: 点が平面のどちら側にあるかを判定
	// 平面: XY平面 (z = 0) の法線ベクトルは (0, 0, 1)
	// 平面上の点: (0, 0, 0)
	// テスト対象の点: (1, 1, -1) は平面の負の側にあるので true を返すはず
	targetP := Vector3D{1, 1, -1}
	planeNormal := Vector3D{0, 0, 1}
	pInPlane := Vector3D{0, 0, 0}

	result := ClassifyEdgeByPlane(targetP, planeNormal, pInPlane)

	assert.True(t, result)
}
