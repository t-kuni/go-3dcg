package domain

import (
	"image/color"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestWorld_Transform_オブジェクトが原点に配置されている(t *testing.T) {
// 	world := World{
// 		Camera: Camera{
// 			Location:  Vector3D{0, 0, -1.0},
// 			Direction: Vector3D{0, 0, 0},
// 		},
// 		LocatedObjects: []LocatedObject{
// 			{
// 				X: 0.0,
// 				Y: 0.0,
// 				Z: 0.0,
// 				Object: Object{
// 					Vertices: NewVertices([]Vector3D{
// 						{0, 0, 0},
// 					}),
// 				},
// 			},
// 		},
// 		Viewport: Viewport{
// 			Width:      100,
// 			Height:     100,
// 			ScaleRatio: 0.25,
// 		},
// 	}

// 	result := world.Transform()

// 	assert.Len(t, result.DiscreteObjects, 1)
// 	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

// 	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].X)
// 	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].Y)
// }

// func TestWorld_Transform_オブジェクトの移動(t *testing.T) {
// 	world := World{
// 		Camera: Camera{
// 			Location:  Vector3D{0, 0, -1.0},
// 			Direction: Vector3D{0, 0, 0},
// 		},
// 		LocatedObjects: []LocatedObject{
// 			{
// 				X: 1.0,  // オブジェクトを移動
// 				Y: -1.0, // オブジェクトを移動
// 				Z: 0.0,
// 				Object: Object{
// 					Vertices: NewVertices([]Vector3D{
// 						{0, 0, 0},
// 					}),
// 				},
// 			},
// 		},
// 		Viewport: Viewport{
// 			Width:      100,
// 			Height:     100,
// 			ScaleRatio: 0.25,
// 		},
// 	}

// 	result := world.Transform()

// 	assert.Len(t, result.DiscreteObjects, 1)
// 	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

// 	assert.Equal(t, int32(75), result.DiscreteObjects[0].Vertices[0].X)
// 	assert.Equal(t, int32(75), result.DiscreteObjects[0].Vertices[0].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
// }

// func TestWorld_Transform_カメラの移動(t *testing.T) {
// 	world := World{
// 		Camera: Camera{
// 			Location:  Vector3D{1.0, -1.0, -1.0}, // カメラを移動
// 			Direction: Vector3D{0, 0, 0},
// 		},
// 		LocatedObjects: []LocatedObject{
// 			{
// 				X: 0.0,
// 				Y: 0.0,
// 				Z: 0.0,
// 				Object: Object{
// 					Vertices: NewVertices([]Vector3D{
// 						{0, 0, 0},
// 					}),
// 				},
// 			},
// 		},
// 		Viewport: Viewport{
// 			Width:      100,
// 			Height:     100,
// 			ScaleRatio: 0.25,
// 		},
// 	}

// 	result := world.Transform()

// 	assert.Len(t, result.DiscreteObjects, 1)
// 	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

// 	assert.Equal(t, int32(25), result.DiscreteObjects[0].Vertices[0].X)
// 	assert.Equal(t, int32(25), result.DiscreteObjects[0].Vertices[0].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
// }

// func TestWorld_Transform_カメラの向き(t *testing.T) {
// 	world := World{
// 		Camera: Camera{
// 			Location:  Vector3D{0, 0, -1.0},
// 			Direction: Vector3D{math.Pi / 16, 0, 0}, // カメラの向きを変更（少し前傾にする）
// 		},
// 		LocatedObjects: []LocatedObject{
// 			{
// 				X: 0.0,
// 				Y: 0.0,
// 				Z: 0.0,
// 				Object: Object{
// 					Vertices: NewVertices([]Vector3D{
// 						{0, 0, 0},
// 					}),
// 				},
// 			},
// 		},
// 		Viewport: Viewport{
// 			Width:      100,
// 			Height:     100,
// 			ScaleRatio: 0.25,
// 		},
// 	}

// 	result := world.Transform()

// 	assert.Len(t, result.DiscreteObjects, 1)
// 	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

// 	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].X)
// 	assert.Equal(t, int32(45), result.DiscreteObjects[0].Vertices[0].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
// }

// func TestWorld_Transform_三角形のオブジェクト(t *testing.T) {
// 	world := World{
// 		Camera: Camera{
// 			Location:  Vector3D{0, 0, -1.0},
// 			Direction: Vector3D{0, 0, 0},
// 		},
// 		LocatedObjects: []LocatedObject{
// 			{
// 				X: 0.0,
// 				Y: 0.0,
// 				Z: 0.0,
// 				Object: Object{
// 					Vertices: NewVertices([]Vector3D{
// 						{-0.5, 0.0, 1.0},
// 						{0.5, 0.0, 1.0},
// 						{0.0, 1.0, 1.0},
// 					}),
// 				},
// 			},
// 		},
// 		Viewport: Viewport{
// 			Width:      100,
// 			Height:     100,
// 			ScaleRatio: 0.25,
// 		},
// 	}

// 	result := world.Transform()

// 	assert.Len(t, result.DiscreteObjects, 1)
// 	assert.Len(t, result.DiscreteObjects[0].Vertices, 3)

// 	assert.Equal(t, int32(38), result.DiscreteObjects[0].Vertices[0].X)
// 	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].Y)
// 	assert.Equal(t, int32(63), result.DiscreteObjects[0].Vertices[1].X)
// 	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[1].Y)
// 	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[2].X)
// 	assert.Equal(t, int32(25), result.DiscreteObjects[0].Vertices[2].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
// }

func TestWorld_ViewVolume_基本的な計算(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  200,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	result := world.ViewVolume()

	assert.InDelta(t, 0.82, result.NearClippingHeight, 0.05)
	assert.InDelta(t, 1.64, result.NearClippingWidth, 0.05)
	assert.InDelta(t, 1.64, result.FarClippingHeight, 0.05)
	assert.InDelta(t, 3.28, result.FarClippingWidth, 0.05)

	// 手前・右上
	assert.InDelta(t, 0.82, result.NearTopRight.X(), 0.05)
	assert.InDelta(t, 0.41, result.NearTopRight.Y(), 0.05)
	assert.InDelta(t, 1.0, result.NearTopRight.Z(), 0.05)

	// 手前・左上
	assert.InDelta(t, -0.82, result.NearTopLeft.X(), 0.05)
	assert.InDelta(t, 0.41, result.NearTopLeft.Y(), 0.05)
	assert.InDelta(t, 1.0, result.NearTopLeft.Z(), 0.05)

	// 手前・右下
	assert.InDelta(t, 0.82, result.NearBottomRight.X(), 0.05)
	assert.InDelta(t, -0.41, result.NearBottomRight.Y(), 0.05)
	assert.InDelta(t, 1.0, result.NearBottomRight.Z(), 0.05)

	// 手前・左下
	assert.InDelta(t, -0.82, result.NearBottomLeft.X(), 0.05)
	assert.InDelta(t, -0.41, result.NearBottomLeft.Y(), 0.05)
	assert.InDelta(t, 1.0, result.NearBottomLeft.Z(), 0.05)

	// 奥・右上
	assert.InDelta(t, 1.64, result.FarTopRight.X(), 0.05)
	assert.InDelta(t, 0.82, result.FarTopRight.Y(), 0.05)
	assert.InDelta(t, 2.0, result.FarTopRight.Z(), 0.05)

	// 奥・左上
	assert.InDelta(t, -1.64, result.FarTopLeft.X(), 0.05)
	assert.InDelta(t, 0.82, result.FarTopLeft.Y(), 0.05)
	assert.InDelta(t, 2.0, result.FarTopLeft.Z(), 0.05)

	// 奥・右下
	assert.InDelta(t, 1.64, result.FarBottomRight.X(), 0.05)
	assert.InDelta(t, -0.82, result.FarBottomRight.Y(), 0.05)
	assert.InDelta(t, 2.0, result.FarBottomRight.Z(), 0.05)

	// 奥・左下
	assert.InDelta(t, -1.64, result.FarBottomLeft.X(), 0.05)
	assert.InDelta(t, -0.82, result.FarBottomLeft.Y(), 0.05)
	assert.InDelta(t, 2.0, result.FarBottomLeft.Z(), 0.05)

	// 法線・手前のクリップ面
	assert.InDelta(t, 0, result.NearPlaneNormal.X(), 0.05)
	assert.InDelta(t, 0, result.NearPlaneNormal.Y(), 0.05)
	assert.InDelta(t, -1, result.NearPlaneNormal.Z(), 0.05)

	// 法線・奥のクリップ面
	assert.InDelta(t, 0, result.FarPlaneNormal.X(), 0.05)
	assert.InDelta(t, 0, result.FarPlaneNormal.Y(), 0.05)
	assert.InDelta(t, 1, result.FarPlaneNormal.Z(), 0.05)

	// 法線・左のクリップ面
	assert.InDelta(t, -0.77, result.LeftPlaneNormal.X(), 0.05)
	assert.InDelta(t, 0, result.LeftPlaneNormal.Y(), 0.05)
	assert.InDelta(t, -0.63, result.LeftPlaneNormal.Z(), 0.05)

	// 法線・右のクリップ面
	assert.InDelta(t, 0.77, result.RightPlaneNormal.X(), 0.05)
	assert.InDelta(t, 0, result.RightPlaneNormal.Y(), 0.05)
	assert.InDelta(t, -0.63, result.RightPlaneNormal.Z(), 0.05)

	// 法線・下のクリップ面
	assert.InDelta(t, 0, result.BottomPlaneNormal.X(), 0.05)
	assert.InDelta(t, -0.93, result.BottomPlaneNormal.Y(), 0.05)
	assert.InDelta(t, -0.38, result.BottomPlaneNormal.Z(), 0.05)

	// 法線・上のクリップ面
	assert.InDelta(t, 0, result.TopPlaneNormal.X(), 0.05)
	assert.InDelta(t, 0.93, result.TopPlaneNormal.Y(), 0.05)
	assert.InDelta(t, -0.38, result.TopPlaneNormal.Z(), 0.05)
}

func TestViewVolume_SutherlandHodgman_ビューボリュームを突き抜ける三角形(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  100,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	viewVolume := world.ViewVolume()

	// 上と右に突き出した三角形
	triangle := [3]Vector3D{
		{0.0, 1.0, 1.5}, // 上
		{1.0, 0.1, 1.5}, // 右下
		{0.0, 0.0, 1.5}, // 左下
	}

	result := viewVolume.SutherlandHodgman(triangle)

	assert.Len(t, result, 5)

	// 左下（変動なし）
	assert.InDelta(t, 0.0, result[0].X(), 1e-2)
	assert.InDelta(t, 0.0, result[0].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[0].Z(), 1e-2)

	// 上（クリップされている）
	assert.InDelta(t, 0.0, result[1].X(), 1e-2)
	assert.InDelta(t, 0.62, result[1].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[1].Z(), 1e-2)

	// 右上１（クリップされている）
	assert.InDelta(t, 0.42, result[2].X(), 1e-2)
	assert.InDelta(t, 0.62, result[2].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[2].Z(), 1e-2)

	// 右上２（クリップされている）
	assert.InDelta(t, 0.62, result[3].X(), 1e-2)
	assert.InDelta(t, 0.44, result[3].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[3].Z(), 1e-2)

	// 右下（クリップされている）
	assert.InDelta(t, 0.62, result[4].X(), 1e-2)
	assert.InDelta(t, 0.06, result[4].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[4].Z(), 1e-2)
}

func TestViewVolume_SutherlandHodgman_ビューボリューム内の三角形(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  100,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	viewVolume := world.ViewVolume()

	triangle := [3]Vector3D{
		{0.0, 0.5, 1.5}, // 上
		{0.5, 0.1, 1.5}, // 右下
		{0.0, 0.0, 1.5}, // 左下
	}

	result := viewVolume.SutherlandHodgman(triangle)

	assert.Len(t, result, 3)

	// すべて変動なし
	assert.InDelta(t, 0.0, result[0].X(), 1e-2)
	assert.InDelta(t, 0.5, result[0].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[0].Z(), 1e-2)

	assert.InDelta(t, 0.5, result[1].X(), 1e-2)
	assert.InDelta(t, 0.1, result[1].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[1].Z(), 1e-2)

	assert.InDelta(t, 0.0, result[2].X(), 1e-2)
	assert.InDelta(t, 0.0, result[2].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[2].Z(), 1e-2)
}

func TestViewVolume_SutherlandHodgman_ビューボリュームを覆う三角形(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  100,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	viewVolume := world.ViewVolume()

	triangle := [3]Vector3D{
		{0.0, 10, 1.5},  // 上
		{10, -10, 1.5},  // 右下
		{-10, -10, 1.5}, // 左下
	}

	result := viewVolume.SutherlandHodgman(triangle)

	assert.Len(t, result, 4)

	// クリップ面に沿った四角形になる
	assert.InDelta(t, 0.62, result[0].X(), 1e-2)
	assert.InDelta(t, 0.62, result[0].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[0].Z(), 1e-2)

	assert.InDelta(t, 0.62, result[1].X(), 1e-2)
	assert.InDelta(t, -0.62, result[1].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[1].Z(), 1e-2)

	assert.InDelta(t, -0.62, result[2].X(), 1e-2)
	assert.InDelta(t, -0.62, result[2].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[2].Z(), 1e-2)

	assert.InDelta(t, -0.62, result[3].X(), 1e-2)
	assert.InDelta(t, 0.62, result[3].Y(), 1e-2)
	assert.InDelta(t, 1.5, result[3].Z(), 1e-2)
}

func TestDynamicObject_AddTriangle_正常系(t *testing.T) {
	obj := NewDynamicObject()

	triangle := [3]Vector3D{
		{0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
	}

	obj.AddTriangle(triangle)

	assert.Len(t, obj.Vertices, 3)
	assert.Equal(t, Vector3D{0.0, 0.0, 0.0}, obj.Vertices[0])
	assert.Equal(t, Vector3D{1.0, 0.0, 0.0}, obj.Vertices[1])
	assert.Equal(t, Vector3D{0.0, 1.0, 0.0}, obj.Vertices[2])

	assert.Len(t, obj.Edges, 3)
	assert.Equal(t, [2]int{0, 1}, obj.Edges[0])
	assert.Equal(t, [2]int{1, 2}, obj.Edges[1])
	assert.Equal(t, [2]int{2, 0}, obj.Edges[2])

	assert.Len(t, obj.Triangles, 1)
	assert.Equal(t, [3]int{0, 1, 2}, obj.Triangles[0])
}

func TestVertexGrid_AddVertex_正常系(t *testing.T) {
	epsilon := 0.1
	vg := NewVertexGrid(epsilon)

	v1 := Vector3D{1.0, 2.0, 3.0}
	v2 := Vector3D{1.2, 2.2, 3.2}    // epsilon以上離れているため追加される頂点
	v3 := Vector3D{1.05, 2.05, 3.05} // epsilon以内の座標は追加されない頂点

	idx := vg.AddVertex(v1)
	assert.Equal(t, 0, idx)

	idx = vg.AddVertex(v2)
	assert.Equal(t, 1, idx)

	idx = vg.AddVertex(v3)
	assert.Equal(t, 0, idx)

	assert.Len(t, vg.vertices, 2)
	assert.Equal(t, Vector3D{1.0, 2.0, 3.0}, vg.vertices[0])
	assert.Equal(t, Vector3D{1.2, 2.2, 3.2}, vg.vertices[1])
}

func TestVertexGrid_SearchVertex(t *testing.T) {
	vg := NewVertexGrid(0.1)

	v1 := Vector3D{1.0, 2.0, 3.0}
	vg.AddVertex(v1)

	// 同じ位置の頂点を検索
	isExist, index := vg.SearchVertex(Vector3D{1.0, 2.0, 3.0})
	assert.True(t, isExist)
	assert.Equal(t, 0, index)

	// 存在しない頂点を検索
	isExist, _ = vg.SearchVertex(Vector3D{10.0, 20.0, 30.0})
	assert.False(t, isExist)
}

func TestViewVolume_MargeVertices_頂点がマージされること１(t *testing.T) {
	viewVolume := ViewVolume{}

	// 同じ三角形を２つ用意する
	obj := Object{
		VertexMatrix: NewVertexMatrix([]Vector3D{
			{-0.5, 0.0, 1.0},
			{0.5, 0.0, 1.0},
			{0.0, 1.0, 1.0},

			{-0.5, 0.0, 1.0},
			{0.5, 0.0, 1.0},
			{0.0, 1.0, 1.0},
		}),
		Edges: [][2]int{
			{0, 1},
			{1, 2},
			{2, 3},
			{3, 4},
			{4, 5},
			{5, 3},
		},
		Triangles: [][3]int{
			{0, 1, 2},
			{3, 4, 5},
		},
	}

	result := viewVolume.MargeVertices(obj)

	assert.Equal(t, 3, result.VertexMatrix.Len())
	assert.Equal(t, Vector3D{-0.5, 0.0, 1.0}, result.VertexMatrix.GetVertex(0))
	assert.Equal(t, Vector3D{0.5, 0.0, 1.0}, result.VertexMatrix.GetVertex(1))
	assert.Equal(t, Vector3D{0.0, 1.0, 1.0}, result.VertexMatrix.GetVertex(2))

	assert.Len(t, result.Edges, 3)
	assert.Equal(t, [2]int{0, 1}, result.Edges[0])
	assert.Equal(t, [2]int{1, 2}, result.Edges[1])
	assert.Equal(t, [2]int{2, 0}, result.Edges[2])

	assert.Len(t, result.Triangles, 1)
	assert.Equal(t, [3]int{0, 1, 2}, result.Triangles[0])
}

func TestViewVolume_MargeVertices_頂点がマージされること２(t *testing.T) {
	viewVolume := ViewVolume{}

	// １辺が同じ座標の三角形を２つ用意する
	obj := Object{
		VertexMatrix: NewVertexMatrix([]Vector3D{
			{-0.5, 0.0, 1.0},
			{0.5, 0.0, 1.0},
			{0.0, 1.0, 1.0}, // 上に凸の三角形

			{-0.5, 0.0, 1.0},
			{0.5, 0.0, 1.0},
			{0.0, -1.0, 1.0}, // 下に凸の三角形
		}),
		Edges: [][2]int{
			{0, 1},
			{1, 2},
			{2, 3},
			{3, 4},
			{4, 5},
			{5, 3},
		},
		Triangles: [][3]int{
			{0, 1, 2},
			{3, 4, 5},
		},
	}

	result := viewVolume.MargeVertices(obj)

	assert.Equal(t, 4, result.VertexMatrix.Len())
	assert.Equal(t, Vector3D{-0.5, 0.0, 1.0}, result.VertexMatrix.GetVertex(0))
	assert.Equal(t, Vector3D{0.5, 0.0, 1.0}, result.VertexMatrix.GetVertex(1))
	assert.Equal(t, Vector3D{0.0, 1.0, 1.0}, result.VertexMatrix.GetVertex(2))
	assert.Equal(t, Vector3D{0.0, -1.0, 1.0}, result.VertexMatrix.GetVertex(3))

	assert.Len(t, result.Edges, 5)
	assert.Equal(t, [2]int{0, 1}, result.Edges[0])
	assert.Equal(t, [2]int{1, 2}, result.Edges[1])
	assert.Equal(t, [2]int{2, 0}, result.Edges[2])
	assert.Equal(t, [2]int{1, 3}, result.Edges[3])
	assert.Equal(t, [2]int{3, 0}, result.Edges[4])

	assert.Len(t, result.Triangles, 2)
	assert.Equal(t, [3]int{0, 1, 2}, result.Triangles[0])
	assert.Equal(t, [3]int{0, 1, 3}, result.Triangles[1])
}

func TestViewVolume_ClipObject_ビューボリュームを覆う三角形(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  100,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	viewVolume := world.ViewVolume()

	// 上・右・奥に突き抜ける三角形を作成
	// 法線は外向き
	obj := Object{
		VertexMatrix: NewVertexMatrix([]Vector3D{
			{0.0, 0.0, 1.5},  // 下
			{0.0, 10.0, 1.5}, // 上
			{0.0, 0.0, 10.0}, // 奥
			{10.0, 0, 1.5},   // 右下
		}),
		Edges: [][2]int{
			{0, 1},
			{0, 2},
			{0, 3},
			{1, 2},
			{1, 3},
			{2, 3},
		},
		Triangles: [][3]int{
			// 法線の向きに注意（右ねじの法則）
			{0, 1, 2},
			{0, 3, 1},
			{0, 2, 3},
			{3, 2, 1}, // 消失する面
		},
	}

	result := viewVolume.ClipObject(obj)

	// クリッピングされずに三角形が保持されることを確認
	assert.Equal(t, 7, result.VertexMatrix.Len())
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(0).X(), 1e-2)
	assert.InDelta(t, 0.82, result.VertexMatrix.GetVertex(0).Y(), 1e-2)
	assert.InDelta(t, 2.0, result.VertexMatrix.GetVertex(0).Z(), 1e-2)

	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(1).X(), 1e-2)
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(1).Y(), 1e-2)
	assert.InDelta(t, 2.0, result.VertexMatrix.GetVertex(1).Z(), 1e-2)

	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(2).X(), 1e-2)
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(2).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(2).Z(), 1e-2)

	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(3).X(), 1e-2)
	assert.InDelta(t, 0.62, result.VertexMatrix.GetVertex(3).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(3).Z(), 1e-2)

	assert.InDelta(t, 0.62, result.VertexMatrix.GetVertex(4).X(), 1e-2)
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(4).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(4).Z(), 1e-2)

	assert.InDelta(t, 0.62, result.VertexMatrix.GetVertex(5).X(), 1e-2)
	assert.InDelta(t, 0.62, result.VertexMatrix.GetVertex(5).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(5).Z(), 1e-2)

	assert.InDelta(t, 0.82, result.VertexMatrix.GetVertex(6).X(), 1e-2)
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(6).Y(), 1e-2)
	assert.InDelta(t, 2.0, result.VertexMatrix.GetVertex(6).Z(), 1e-2)

	assert.Len(t, result.Edges, 12)
	assert.Equal(t, [2]int{0, 1}, result.Edges[0])
	assert.Equal(t, [2]int{1, 2}, result.Edges[1])
	assert.Equal(t, [2]int{2, 0}, result.Edges[2])
	assert.Equal(t, [2]int{2, 3}, result.Edges[3])
	assert.Equal(t, [2]int{3, 0}, result.Edges[4])
	assert.Equal(t, [2]int{2, 4}, result.Edges[5])
	assert.Equal(t, [2]int{4, 3}, result.Edges[6])
	assert.Equal(t, [2]int{4, 5}, result.Edges[7])
	assert.Equal(t, [2]int{5, 3}, result.Edges[8])
	assert.Equal(t, [2]int{1, 6}, result.Edges[9])
	assert.Equal(t, [2]int{6, 4}, result.Edges[10])
	assert.Equal(t, [2]int{4, 1}, result.Edges[11])

	// 法線の向きが外向きであること（右ねじの法則）
	assert.Len(t, result.Triangles, 6)
	assert.Equal(t, [3]int{0, 1, 2}, result.Triangles[0])
	assert.Equal(t, [3]int{0, 2, 3}, result.Triangles[1])
	assert.Equal(t, [3]int{3, 2, 4}, result.Triangles[2])
	assert.Equal(t, [3]int{3, 4, 5}, result.Triangles[3])
	assert.Equal(t, [3]int{1, 6, 4}, result.Triangles[4])
	assert.Equal(t, [3]int{1, 4, 2}, result.Triangles[5])
}

func TestWorld_TransformPerspectiveProjection_正常系(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  100,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	// ビューボリューム内の三角形オブジェクトを作成
	locatedObject := LocatedObject{
		Object: Object{
			VertexMatrix: NewVertexMatrix([]Vector3D{
				{-0.3, 0.0, 1.1}, // 左下
				{0.3, 0.0, 1.1},  // 右下
				{0.0, 0.0, 1.9},  // 奥
				{0.0, 0.3, 1.1},  // 上
			}),
			Edges: [][2]int{
				{0, 1},
				{0, 2},
				{0, 3},
				{1, 2},
				{1, 3},
				{2, 3},
			},
			Triangles: [][3]int{
				{0, 1, 2},
				{0, 1, 3},
				{0, 2, 3},
				{1, 2, 3},
			},
		},
	}

	// 透視変換を実行
	resultObject := world.TransformPerspectiveProjection(locatedObject.Object)

	// 結果の検証
	assert.Equal(t, 4, resultObject.VertexMatrix.Len())

	assert.InDelta(t, -0.65, resultObject.VertexMatrix.GetVertex(0).X(), 1e-2)
	assert.InDelta(t, 0.0, resultObject.VertexMatrix.GetVertex(0).Y(), 1e-2)
	assert.InDelta(t, 0.18, resultObject.VertexMatrix.GetVertex(0).Z(), 1e-2)

	assert.InDelta(t, 0.65, resultObject.VertexMatrix.GetVertex(1).X(), 1e-2)
	assert.InDelta(t, 0.0, resultObject.VertexMatrix.GetVertex(1).Y(), 1e-2)
	assert.InDelta(t, 0.18, resultObject.VertexMatrix.GetVertex(1).Z(), 1e-2)

	assert.InDelta(t, 0.0, resultObject.VertexMatrix.GetVertex(2).X(), 1e-2)
	assert.InDelta(t, 0.0, resultObject.VertexMatrix.GetVertex(2).Y(), 1e-2)
	assert.InDelta(t, 0.94, resultObject.VertexMatrix.GetVertex(2).Z(), 1e-2)

	assert.InDelta(t, 0.0, resultObject.VertexMatrix.GetVertex(3).X(), 1e-2)
	assert.InDelta(t, 0.65, resultObject.VertexMatrix.GetVertex(3).Y(), 1e-2)
	assert.InDelta(t, 0.18, resultObject.VertexMatrix.GetVertex(3).Z(), 1e-2)

	assert.Len(t, resultObject.Edges, 6)
	assert.Len(t, resultObject.Triangles, 4)
}

func TestWorld_ClipWithViewVolume_正常系(t *testing.T) {
	world := World{
		Viewport: Viewport{
			Width:  100,
			Height: 100,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	// 左に突き抜けるオブジェクト
	obj := Object{
		VertexMatrix: NewVertexMatrix([]Vector3D{
			{0.0, 0.0, 1.5},   // 中央
			{-10.0, 0.0, 1.5}, // 左
			{0.0, 0.3, 1.5},   // 上
		}),
		Edges: [][2]int{
			{0, 1},
			{1, 2},
			{2, 0},
		},
		Triangles: [][3]int{
			{0, 1, 2},
		},
	}

	// ClipWithViewVolumeを実行
	result := world.ClipWithViewVolume(obj)

	// ビューボリューム内なので変更されずに保持されることを確認
	assert.Equal(t, 4, result.VertexMatrix.Len())
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(0).X(), 1e-2)
	assert.InDelta(t, 0.3, result.VertexMatrix.GetVertex(0).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(0).Z(), 1e-2)

	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(1).X(), 1e-2)
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(1).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(1).Z(), 1e-2)

	assert.InDelta(t, -0.62, result.VertexMatrix.GetVertex(2).X(), 1e-2)
	assert.InDelta(t, 0.0, result.VertexMatrix.GetVertex(2).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(2).Z(), 1e-2)

	assert.InDelta(t, -0.62, result.VertexMatrix.GetVertex(3).X(), 1e-2)
	assert.InDelta(t, 0.28, result.VertexMatrix.GetVertex(3).Y(), 1e-2)
	assert.InDelta(t, 1.5, result.VertexMatrix.GetVertex(3).Z(), 1e-2)

	assert.Len(t, result.Edges, 5)
	assert.Equal(t, [2]int{0, 1}, result.Edges[0])
	assert.Equal(t, [2]int{1, 2}, result.Edges[1])
	assert.Equal(t, [2]int{2, 0}, result.Edges[2])
	assert.Equal(t, [2]int{2, 3}, result.Edges[3])
	assert.Equal(t, [2]int{3, 0}, result.Edges[4])

	assert.Len(t, result.Triangles, 2)
	assert.Equal(t, [3]int{0, 1, 2}, result.Triangles[0])
	assert.Equal(t, [3]int{0, 2, 3}, result.Triangles[1])
}

// func TestWorld_Transform_正常系(t *testing.T) {
// 	world := World{
// 		Camera: Camera{
// 			Location:  Vector3D{0, 0, 0},
// 			Direction: Vector3D{0, 0, 0},
// 		},
// 		LocatedObjects: []LocatedObject{
// 			{
// 				X: 0, Y: 0, Z: 0,
// 				Object: Object{
// 					Vertices: NewVertices([]Vector3D{
// 						{-0.3, 0.0, 1.1}, // 左下
// 						{0.3, 0.0, 1.1},  // 右下
// 						{0.0, 0.0, 1.9},  // 奥
// 						{0.0, 0.3, 1.1},  // 上
// 					}),
// 					Edges: [][2]int{
// 						{0, 1},
// 						{0, 2},
// 						{0, 3},
// 						{1, 2},
// 						{1, 3},
// 						{2, 3},
// 					},
// 					Triangles: [][3]int{
// 						{0, 1, 2},
// 						{0, 1, 3},
// 						{0, 2, 3},
// 						{1, 2, 3},
// 					},
// 				},
// 			},
// 		},
// 		Viewport: Viewport{
// 			Width:      100,
// 			Height:     100,
// 			ScaleRatio: 0.5,
// 		},
// 		Clipping: Clipping{
// 			NearDistance: 1.0,
// 			FarDistance:  2.0,
// 			FieldOfView:  math.Pi / 4,
// 		},
// 	}

// 	// 透視変換を実行
// 	actual := world.Transform()

// 	// 結果の検証
// 	assert.Len(t, actual.DiscreteObjects, 1)
// 	assert.Len(t, actual.DiscreteObjects[0].Vertices, 4)
// }

func TestCalculatedWorld_RayTrace_正常系(t *testing.T) {
	// 小さなビューポートで簡単な三角形をレイトレースするテスト
	world := World{
		Viewport: Viewport{
			Width:  20,
			Height: 20,
		},
		Clipping: Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4, // 45度
		},
	}

	// ビューボリューム内に配置された三角形オブジェクト
	obj := Object{
		VertexMatrix: NewVertexMatrix([]Vector3D{
			{-0.3, -0.3, 1.5}, // 左下
			{0.3, -0.3, 1.5},  // 右下
			{0.0, 0.3, 1.5},   // 上
		}),
		Edges: [][2]int{
			{0, 1},
			{1, 2},
			{2, 0},
		},
		Triangles: [][3]int{
			{0, 1, 2},
		},
	}

	calculatedWorld := CalculatedWorld{
		Origin:  world,
		Objects: []Object{obj},
	}

	frameBuffer := calculatedWorld.RayTrace()

	assert.Len(t, frameBuffer, 50)

	black := color.RGBA{0, 0, 0, 255}
	// 上段
	assert.NotContains(t, frameBuffer, FrameBufferKey{X: 5, Y: 6})
	assert.Equal(t, black, frameBuffer[FrameBufferKey{X: 10, Y: 6}].Color)
	assert.NotContains(t, frameBuffer, FrameBufferKey{X: 14, Y: 6})
	// 中段
	assert.NotContains(t, frameBuffer, FrameBufferKey{X: 5, Y: 10})
	assert.Equal(t, black, frameBuffer[FrameBufferKey{X: 10, Y: 10}].Color)
	assert.NotContains(t, frameBuffer, FrameBufferKey{X: 14, Y: 10})
	// 下段
	assert.Equal(t, black, frameBuffer[FrameBufferKey{X: 5, Y: 14}].Color)
	assert.Equal(t, black, frameBuffer[FrameBufferKey{X: 10, Y: 14}].Color)
	assert.Equal(t, black, frameBuffer[FrameBufferKey{X: 14, Y: 14}].Color)
}
