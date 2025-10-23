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
	// テスト対象の点: (0, 0, -1) は平面の負の側にあるので true を返すはず
	targetP := Vector3D{0, 0, -0.1}
	planeNormal := Vector3D{0, 0, 1}
	pInPlane := Vector3D{0, 0, 0}

	result := ClassifyEdgeByPlane(targetP, planeNormal, pInPlane)

	assert.True(t, result)
}

func TestClassifyEdgeByPlane2(t *testing.T) {
	targetP := Vector3D{0, 0, 0.1}
	planeNormal := Vector3D{0, 0, 1}
	pInPlane := Vector3D{0, 0, 0}

	result := ClassifyEdgeByPlane(targetP, planeNormal, pInPlane)

	assert.False(t, result)
}

func TestIntersectPlaneIntersectionPoint(t *testing.T) {
	// 正常系テストケース: XY平面（z = 0）と線分の交点を計算
	// 平面: z = 0（法線ベクトル: (0, 0, 1)）
	// 平面上の点: (0, 0, 0)
	// 線分: 点A(0, 0, -1)から点B(0, 0, 1)への線分
	// 期待される交点: (0, 0, 0)
	planeNormal := Vector3D{0, 0, 1}
	planePoint := Vector3D{0, 1, 0}
	fromVertex := Vector3D{1, 0, -1}
	toVertex := Vector3D{1, 0, 1}

	result := IntersectPlaneIntersectionPoint(planeNormal, planePoint, fromVertex, toVertex)

	assert.InDelta(t, 1, result.X(), 0.001)
	assert.InDelta(t, 0, result.Y(), 0.001)
	assert.InDelta(t, 0, result.Z(), 0.001)
}

func TestTriangulate_LessThanThreeVertices(t *testing.T) {
	// 頂点が3つ未満の場合のテストケース
	// 空のスライスや2つの頂点の場合、空の三角形配列が返されるべき

	// 空のスライスの場合
	emptyVertices := []Vector3D{}
	result := Triangulate(emptyVertices)
	assert.Equal(t, [][3]Vector3D{}, result)

	// 1つの頂点の場合
	oneVertex := []Vector3D{{0, 0, 0}}
	result = Triangulate(oneVertex)
	assert.Equal(t, [][3]Vector3D{}, result)

	// 2つの頂点の場合
	twoVertices := []Vector3D{{0, 0, 0}, {1, 0, 0}}
	result = Triangulate(twoVertices)
	assert.Equal(t, [][3]Vector3D{}, result)
}

func TestTriangulate_ExactlyThreeVertices(t *testing.T) {
	// 頂点がちょうど3つの場合のテストケース
	// そのまま1つの三角形として返されるべき
	vertices := []Vector3D{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
	}

	result := Triangulate(vertices)

	expected := [][3]Vector3D{{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
	}}

	assert.Equal(t, expected, result)
	assert.Len(t, result, 1)
}

func TestTriangulate_MoreThanThreeVertices(t *testing.T) {
	// 頂点が4つ以上の場合のテストケース（ファン三角化）
	// 四角形を2つの三角形に分割することを確認
	vertices := []Vector3D{
		{0, 0, 0}, // 中心点
		{1, 0, 0}, // 右
		{1, 1, 0}, // 右上
		{0, 1, 0}, // 上
	}

	result := Triangulate(vertices)

	// 4つの頂点から2つの三角形が生成されるべき
	// 三角形1: vertices[0], vertices[1], vertices[2]
	// 三角形2: vertices[0], vertices[2], vertices[3]
	expected := [][3]Vector3D{
		{{0, 0, 0}, {1, 0, 0}, {1, 1, 0}},
		{{0, 0, 0}, {1, 1, 0}, {0, 1, 0}},
	}

	assert.Equal(t, expected, result)
	assert.Len(t, result, 2)

	// 一般的なケース：5つの頂点から3つの三角形が生成される
	fiveVertices := []Vector3D{
		{0, 0, 0},
		{1, 0, 0},
		{2, 1, 0},
		{1, 2, 0},
		{0, 1, 0},
	}

	fiveResult := Triangulate(fiveVertices)
	assert.Len(t, fiveResult, 3) // n-2 = 5-2 = 3つの三角形
}

func TestCleanEdges(t *testing.T) {
	// 正常系テストケース: 重複と無効な辺をクリーンアップ
	// 入力: 重複する辺、同じ頂点への辺を含む辺のリスト
	// 期待結果: 重複と無効な辺が除去された辺のリスト
	edges := [][2]int{
		{0, 1}, // 有効な辺
		{1, 2}, // 有効な辺
		{0, 1}, // 重複する辺（除去されるべき）
		{3, 3}, // 同じ頂点への辺（除去されるべき）
		{2, 3}, // 有効な辺
	}

	result := CleanEdges(edges)

	expected := [][2]int{
		{0, 1},
		{1, 2},
		{2, 3},
	}

	assert.Equal(t, expected, result)
	assert.Len(t, result, 3)
}

func TestCleanTriangles(t *testing.T) {
	// 正常系テストケース: 重複と無効な三角形をクリーンアップ
	// 入力: 重複する三角形、同じ頂点を持つ三角形を含む三角形のリスト
	// 期待結果: 重複と無効な三角形が除去された三角形のリスト
	triangles := [][3]int{
		{0, 1, 2}, // 有効な三角形
		{1, 2, 3}, // 有効な三角形
		{0, 1, 2}, // 重複する三角形（除去されるべき）
		{4, 4, 5}, // 同じ頂点を持つ三角形（除去されるべき）
		{2, 3, 4}, // 有効な三角形
		{5, 6, 5}, // 同じ頂点を持つ三角形（除去されるべき）
	}

	result := CleanTriangles(triangles)

	expected := [][3]int{
		{0, 1, 2},
		{1, 2, 3},
		{2, 3, 4},
	}

	assert.Equal(t, expected, result)
	assert.Len(t, result, 3)
}
