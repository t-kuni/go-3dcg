package domain

import (
	"math"

	"github.com/t-kuni/go-3dcg/util"
	"gonum.org/v1/gonum/mat"
)

type World struct {
	Camera         Camera
	LocatedObjects []LocatedObject
	Viewport       Viewport
	Clipping       Clipping
}

type Viewport struct {
	Width      int32
	Height     int32
	ScaleRatio float64
}

type Clipping struct {
	// NearDistance 前方クリップ面までの距離
	NearDistance float64
	// FarDistance 後方クリップ面までの距離
	FarDistance float64
	// FieldOfView 視野角(単位：ラジアン)
	FieldOfView float64
}

func (w World) Transform() DiscreteWorld {
	discreateWorld := NewDiscreteWorld()
	for _, locatedObject := range w.LocatedObjects {
		m := locatedObject.Object.Matrix()

		m = util.T(m) // 転置(行列計算の次元を揃えるため)

		// ワールド座標変換
		m = util.TransformTranslate(m, locatedObject.X, locatedObject.Y, locatedObject.Z)

		// カメラ座標変換
		m = util.TransformTranslate(m, -w.Camera.Location.X, -w.Camera.Location.Y, -w.Camera.Location.Z)
		m = util.TransformRotate(m, -w.Camera.Direction.X, -w.Camera.Direction.Y, -w.Camera.Direction.Z)

		// 投影変換
		m = util.TransformParallelProjection(m)

		// ビューポート変換
		m = util.TransformViewport(m, w.Viewport.Width, w.Viewport.Height, w.Viewport.ScaleRatio)

		m = util.T(m) // 転置を戻す

		rowCnt, _ := m.Dims()
		obj := NewDiscreteObject()
		obj.Edges = locatedObject.Object.Edges
		for r := 0; r < rowCnt; r++ {
			obj.AddVertex(DiscretePoint2D{X: int32(math.Round(m.At(r, 0))), Y: int32(math.Round(m.At(r, 1)))})
		}
		discreateWorld.AddObject(obj)
	}

	return discreateWorld
}

type ViewVolume struct {
	NearClippingWidth  float64
	NearClippingHeight float64
	FarClippingWidth   float64
	FarClippingHeight  float64
}

func (w World) ViewVolume() ViewVolume {
	aspectRatio := float64(w.Viewport.Width) / float64(w.Viewport.Height)
	nearClippingHeight := 2.0 * math.Tan(w.Clipping.FieldOfView/2.0) * w.Clipping.NearDistance
	nearClippingWidth := nearClippingHeight * aspectRatio
	farClippingHeight := 2.0 * math.Tan(w.Clipping.FieldOfView/2.0) * w.Clipping.FarDistance
	farClippingWidth := farClippingHeight * aspectRatio

	return ViewVolume{
		NearClippingHeight: nearClippingHeight,
		NearClippingWidth:  nearClippingWidth,
		FarClippingHeight:  farClippingHeight,
		FarClippingWidth:   farClippingWidth,
	}
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
	// Edges 辺を表す。[0]は始点の頂点の添字番号、[1]は終点の頂点の添字番号。
	Edges [][2]int
}

func (o Object) Matrix() mat.Dense {
	vertices := []float64{}
	for _, vertex := range o.Vertices {
		vertices = append(vertices, vertex.X, vertex.Y, vertex.Z, 1)
	}
	return *mat.NewDense(len(o.Vertices), 4, vertices)
}

type Vertex struct {
	Point3D
}

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
	// Edges 辺を表す。[0]は始点の頂点の添字番号、[1]は終点の頂点の添字番号。
	Edges [][2]int
}

func NewDiscreteObject() DiscreteObject {
	return DiscreteObject{
		Vertices: []DiscretePoint2D{},
	}
}

func (o *DiscreteObject) AddVertex(vertex DiscretePoint2D) {
	o.Vertices = append(o.Vertices, vertex)
}
