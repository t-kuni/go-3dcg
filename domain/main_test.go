package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorld_Transform_オブジェクトが原点に配置されている(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Vector3D{0, 0, -1.0},
			Direction: Vector3D{0, 0, 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Vector3D{0, 0, 0}},
					},
				},
			},
		},
		Viewport: Viewport{
			Width:      100,
			Height:     100,
			ScaleRatio: 0.25,
		},
	}

	result := world.Transform()

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].Y)
}

func TestWorld_Transform_オブジェクトの移動(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Vector3D{0, 0, -1.0},
			Direction: Vector3D{0, 0, 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 1.0,  // オブジェクトを移動
				Y: -1.0, // オブジェクトを移動
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Vector3D{0, 0, 0}},
					},
				},
			},
		},
		Viewport: Viewport{
			Width:      100,
			Height:     100,
			ScaleRatio: 0.25,
		},
	}

	result := world.Transform()

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

	assert.Equal(t, int32(75), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(75), result.DiscreteObjects[0].Vertices[0].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
}

func TestWorld_Transform_カメラの移動(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Vector3D{1.0, -1.0, -1.0}, // カメラを移動
			Direction: Vector3D{0, 0, 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Vector3D{0, 0, 0}},
					},
				},
			},
		},
		Viewport: Viewport{
			Width:      100,
			Height:     100,
			ScaleRatio: 0.25,
		},
	}

	result := world.Transform()

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

	assert.Equal(t, int32(25), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(25), result.DiscreteObjects[0].Vertices[0].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
}

func TestWorld_Transform_カメラの向き(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Vector3D{0, 0, -1.0},
			Direction: Vector3D{math.Pi / 16, 0, 0}, // カメラの向きを変更（少し前傾にする）
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Vector3D{0, 0, 0}},
					},
				},
			},
		},
		Viewport: Viewport{
			Width:      100,
			Height:     100,
			ScaleRatio: 0.25,
		},
	}

	result := world.Transform()

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(45), result.DiscreteObjects[0].Vertices[0].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
}

func TestWorld_Transform_三角形のオブジェクト(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Vector3D{0, 0, -1.0},
			Direction: Vector3D{0, 0, 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Vector3D{-0.5, 0.0, 1.0}},
						{Vector3D{0.5, 0.0, 1.0}},
						{Vector3D{0.0, 1.0, 1.0}},
					},
				},
			},
		},
		Viewport: Viewport{
			Width:      100,
			Height:     100,
			ScaleRatio: 0.25,
		},
	}

	result := world.Transform()

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 3)

	assert.Equal(t, int32(38), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].Y)
	assert.Equal(t, int32(63), result.DiscreteObjects[0].Vertices[1].X)
	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[1].Y)
	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[2].X)
	assert.Equal(t, int32(25), result.DiscreteObjects[0].Vertices[2].Y) // SDL2の仕様に準拠するため上下逆転する点に注意する
}

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
