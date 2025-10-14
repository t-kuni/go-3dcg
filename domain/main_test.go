package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorld_Transform_オブジェクトが原点に配置されている(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Point3D{X: 0, Y: 0, Z: -1.0},
			Direction: Point3D{X: 0, Y: 0, Z: 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Point3D{X: 0, Y: 0, Z: 0}},
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
			Location:  Point3D{X: 0, Y: 0, Z: -1.0},
			Direction: Point3D{X: 0, Y: 0, Z: 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 1.0,  // オブジェクトを移動
				Y: -1.0, // オブジェクトを移動
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Point3D{X: 0, Y: 0, Z: 0}},
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
			Location:  Point3D{X: 1.0, Y: -1.0, Z: -1.0}, // カメラを移動
			Direction: Point3D{X: 0, Y: 0, Z: 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Point3D{X: 0, Y: 0, Z: 0}},
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
			Location:  Point3D{X: 0, Y: 0, Z: -1.0},
			Direction: Point3D{X: math.Pi / 16, Y: 0, Z: 0}, // カメラの向きを変更（少し前傾にする）
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Point3D{X: 0, Y: 0, Z: 0}},
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
			Location:  Point3D{X: 0, Y: 0, Z: -1.0},
			Direction: Point3D{X: 0, Y: 0, Z: 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 0.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Point3D{X: -0.5, Y: 0.0, Z: 1.0}},
						{Point3D{X: 0.5, Y: 0.0, Z: 1.0}},
						{Point3D{X: 0.0, Y: 1.0, Z: 1.0}},
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
