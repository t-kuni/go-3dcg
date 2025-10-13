package domain

import (
	"math"

	"github.com/t-kuni/go-3dcg/util"
	"gonum.org/v1/gonum/mat"
)

type World struct {
	Camera         Camera
	LocatedObjects []LocatedObject
}

func (w World) Transform(viewWidth int32, viewHeight int32) DiscreteWorld {
	discreateWorld := NewDiscreteWorld()
	for _, locatedObject := range w.LocatedObjects {
		m := locatedObject.Object.Matrix()

		// ワールド座標変換
		m = util.TransformTranslate(m, locatedObject.X, locatedObject.Y, locatedObject.Z)

		// カメラ座標変換
		m = util.TransformTranslate(m, -w.Camera.Location.X, -w.Camera.Location.Y, -w.Camera.Location.Z)
		m = util.TransformRotate(m, -w.Camera.Direction.X, -w.Camera.Direction.Y, -w.Camera.Direction.Z)

		// 投影変換
		m = util.TransformParallelProjection(m)

		// ビューポート変換
		m = util.TransformViewport(m, viewWidth, viewHeight)

		rowCnt, _ := m.Dims()
		for r := 0; r < rowCnt; r++ {
			discreateWorld.AddObject(DiscreteObject{
				Vertices: []DiscretePoint2D{
					{X: int32(math.Round(m.At(r, 0))), Y: int32(math.Round(m.At(r, 1)))},
				},
			})
		}
	}

	return discreateWorld
}

type Camera struct {
	Location  Point3D
	Direction Point3D
	// Up        *mat.Dense
}

type LocatedObject struct {
	X, Y, Z float64
	Object  Object
}

type Object struct {
	Vertices []Vertex
	// EdgeIndexes []int
}

func (o Object) Matrix() mat.Dense {
	vertices := []float64{}
	for _, vertex := range o.Vertices {
		vertices = append(vertices, vertex.X, vertex.Y, vertex.Z, 1)
	}
	return *mat.NewDense(4, 4, vertices)
}

type Vertex struct {
	Point3D
}

// func (w World) AddObject(object Object) {
// 	w.Objects = append(w.Objects, object)
// }

type Point3D struct {
	X, Y, Z float64
}

func (p Point3D) Matrix() mat.Dense {
	return *mat.NewDense(1, 4, []float64{p.X, p.Y, p.Z, 1})
}

type DiscreteWorld struct {
	DiscreteObjects []DiscreteObject
}

func NewDiscreteWorld() DiscreteWorld {
	return DiscreteWorld{
		DiscreteObjects: []DiscreteObject{},
	}
}

func (w *DiscreteWorld) AddObject(object DiscreteObject) {
	w.DiscreteObjects = append(w.DiscreteObjects, object)
}

// DiscretePoint2D 整数型の二次元座標
type DiscretePoint2D struct {
	X, Y int32
}

type DiscreteObject struct {
	Vertices []DiscretePoint2D
}
