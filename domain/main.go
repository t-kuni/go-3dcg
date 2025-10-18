package domain

import (
	"math"

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

		m = T(m) // 転置(行列計算の次元を揃えるため)

		// ワールド座標変換
		m = TransformTranslate(m, locatedObject.X, locatedObject.Y, locatedObject.Z)

		// カメラ座標変換
		m = TransformTranslate(m, -w.Camera.Location.X(), -w.Camera.Location.Y(), -w.Camera.Location.Z())
		m = TransformRotate(m, -w.Camera.Direction.X(), -w.Camera.Direction.Y(), -w.Camera.Direction.Z())

		// 投影変換
		m = TransformParallelProjection(m)

		// ビューポート変換
		m = TransformViewport(m, w.Viewport.Width, w.Viewport.Height, w.Viewport.ScaleRatio)

		m = T(m) // 転置を戻す

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
	// クリッピング面の幅と高さ
	NearClippingWidth  float64
	NearClippingHeight float64
	FarClippingWidth   float64
	FarClippingHeight  float64

	// 頂点
	NearTopRight    Vector3D
	NearTopLeft     Vector3D
	NearBottomRight Vector3D
	NearBottomLeft  Vector3D
	FarTopRight     Vector3D
	FarTopLeft      Vector3D
	FarBottomRight  Vector3D
	FarBottomLeft   Vector3D

	// 法線
	NearPlaneNormal   Vector3D
	FarPlaneNormal    Vector3D
	LeftPlaneNormal   Vector3D
	RightPlaneNormal  Vector3D
	BottomPlaneNormal Vector3D
	TopPlaneNormal    Vector3D
}

func (w World) ViewVolume() ViewVolume {
	aspectRatio := float64(w.Viewport.Width) / float64(w.Viewport.Height)
	nearClippingHeight := 2.0 * math.Tan(w.Clipping.FieldOfView/2.0) * w.Clipping.NearDistance
	nearClippingWidth := nearClippingHeight * aspectRatio
	farClippingHeight := 2.0 * math.Tan(w.Clipping.FieldOfView/2.0) * w.Clipping.FarDistance
	farClippingWidth := farClippingHeight * aspectRatio

	nearClippingHeightHalf := nearClippingHeight / 2
	nearClippingWidthHalf := nearClippingWidth / 2
	farClippingHeightHalf := farClippingHeight / 2
	farClippingWidthHalf := farClippingWidth / 2

	nearTopRight := Vector3D{nearClippingWidthHalf, nearClippingHeightHalf, w.Clipping.NearDistance}
	nearTopLeft := Vector3D{-nearClippingWidthHalf, nearClippingHeightHalf, w.Clipping.NearDistance}
	nearBottomRight := Vector3D{nearClippingWidthHalf, -nearClippingHeightHalf, w.Clipping.NearDistance}
	nearBottomLeft := Vector3D{-nearClippingWidthHalf, -nearClippingHeightHalf, w.Clipping.NearDistance}
	farTopRight := Vector3D{farClippingWidthHalf, farClippingHeightHalf, w.Clipping.FarDistance}
	farTopLeft := Vector3D{-farClippingWidthHalf, farClippingHeightHalf, w.Clipping.FarDistance}
	farBottomRight := Vector3D{farClippingWidthHalf, -farClippingHeightHalf, w.Clipping.FarDistance}
	farBottomLeft := Vector3D{-farClippingWidthHalf, -farClippingHeightHalf, w.Clipping.FarDistance}

	nearPlaneNormal := CalcNormalFromPoints(nearTopRight, nearTopLeft, nearBottomLeft)
	farPlaneNormal := CalcNormalFromPoints(farTopLeft, farTopRight, farBottomRight)
	leftPlaneNormal := CalcNormalFromPoints(nearTopLeft, farTopLeft, farBottomLeft)
	rightPlaneNormal := CalcNormalFromPoints(farTopRight, nearTopRight, nearBottomRight)
	bottomPlaneNormal := CalcNormalFromPoints(nearBottomLeft, farBottomLeft, farBottomRight)
	topPlaneNormal := CalcNormalFromPoints(farTopRight, farTopLeft, nearTopLeft)

	return ViewVolume{
		NearClippingHeight: nearClippingHeight,
		NearClippingWidth:  nearClippingWidth,
		FarClippingHeight:  farClippingHeight,
		FarClippingWidth:   farClippingWidth,
		NearTopRight:       nearTopRight,
		NearTopLeft:        nearTopLeft,
		NearBottomRight:    nearBottomRight,
		NearBottomLeft:     nearBottomLeft,
		FarTopRight:        farTopRight,
		FarTopLeft:         farTopLeft,
		FarBottomRight:     farBottomRight,
		FarBottomLeft:      farBottomLeft,
		NearPlaneNormal:    nearPlaneNormal,
		FarPlaneNormal:     farPlaneNormal,
		LeftPlaneNormal:    leftPlaneNormal,
		RightPlaneNormal:   rightPlaneNormal,
		BottomPlaneNormal:  bottomPlaneNormal,
		TopPlaneNormal:     topPlaneNormal,
	}
}

type Camera struct {
	Location  Vector3D
	Direction Vector3D
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
	// Triangles 三角形を表す。3つの頂点の添字番号を保持する
	Triangles [][3]int
}

func (o Object) Matrix() mat.Dense {
	vertices := []float64{}
	for _, vertex := range o.Vertices {
		vertices = append(vertices, vertex.X(), vertex.Y(), vertex.Z(), 1)
	}
	return *mat.NewDense(len(o.Vertices), 4, vertices)
}

type Vertex struct {
	Vector3D
}

type Vector3D [3]float64

func (p Vector3D) X() float64 {
	return p[0]
}

func (p Vector3D) Y() float64 {
	return p[1]
}

func (p Vector3D) Z() float64 {
	return p[2]
}

func (p Vector3D) Vec() mat.VecDense {
	return *mat.NewVecDense(3, p[:])
}

func (p Vector3D) Matrix() mat.Dense {
	return *mat.NewDense(1, 4, append(p[:], 1))
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
	// Triangles 三角形を表す。3つの頂点の添字番号を保持する
	Triangles [][3]int
}

func NewDiscreteObject() DiscreteObject {
	return DiscreteObject{
		Vertices: []DiscretePoint2D{},
	}
}

func (o *DiscreteObject) AddVertex(vertex DiscretePoint2D) {
	o.Vertices = append(o.Vertices, vertex)
}
