package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorld_Transform(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Point3D{X: 0, Y: 0, Z: 0},
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
	}

	viewWidth := int32(100)
	viewHeight := int32(100)

	result := world.Transform(viewWidth, viewHeight)

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].Y)
}

func TestWorld_Transform_2(t *testing.T) {
	world := World{
		Camera: Camera{
			Location:  Point3D{X: 0, Y: 0, Z: 0},
			Direction: Point3D{X: 0, Y: 0, Z: 0},
		},
		LocatedObjects: []LocatedObject{
			{
				X: 1.0,
				Y: 0.0,
				Z: 0.0,
				Object: Object{
					Vertices: []Vertex{
						{Point3D{X: 0, Y: 0, Z: 0}},
					},
				},
			},
		},
	}

	viewWidth := int32(100)
	viewHeight := int32(100)

	result := world.Transform(viewWidth, viewHeight)

	assert.Len(t, result.DiscreteObjects, 1)
	assert.Len(t, result.DiscreteObjects[0].Vertices, 1)

	assert.Equal(t, int32(75), result.DiscreteObjects[0].Vertices[0].X)
	assert.Equal(t, int32(50), result.DiscreteObjects[0].Vertices[0].Y)
}
