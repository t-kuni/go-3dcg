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
	Width  int32
	Height int32
	// 実数値の単位座標(1.0f)を「画面の短辺」の{scaleRatio}%分、拡大する
	// 短辺が100でScaleRatioが0.25なら、100*0.25=25となり、1.0fが25pxとなる
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
	for _, obj := range w.LocatedObjects {
		m := *obj.Object.VertexMatrix.Dense

		// ワールド座標変換
		m = TransformTranslate(m, obj.X, obj.Y, obj.Z)

		// カメラ座標変換
		m = TransformTranslate(m, -w.Camera.Location.X(), -w.Camera.Location.Y(), -w.Camera.Location.Z())
		m = TransformRotate(m, -w.Camera.Direction.X(), -w.Camera.Direction.Y(), -w.Camera.Direction.Z())

		// 透視投影
		obj.Object.VertexMatrix.Dense = &m
		o := w.TransformPerspectiveProjection(obj.Object)
		m = *o.VertexMatrix.Dense

		// ビューポート変換
		m = TransformViewport(m, w.Viewport.Width, w.Viewport.Height, w.Viewport.ScaleRatio)

		m = T(m) // 転置を戻す

		rowCnt, _ := m.Dims()
		obj := NewDiscreteObject()
		obj.Edges = o.Edges
		for r := 0; r < rowCnt; r++ {
			rx := m.At(r, 0)
			ry := m.At(r, 1)
			obj.AddVertex(DiscretePoint2D{X: int32(math.Round(rx)), Y: int32(math.Round(ry))})
		}
		discreateWorld.AddObject(obj)
	}

	return discreateWorld
}

// TransformPerspectiveProjection 透視変換を行う
// 引数oのVerticesは無視される
// 引数mをVerticesとして扱う。mは転置されて4行N列になっている。
func (w World) TransformPerspectiveProjection(o Object) Object {
	viewVolume := w.ViewVolume()
	clippedObject := viewVolume.ClipObject(o)

	aspect := float64(w.Viewport.Width) / float64(w.Viewport.Height)
	tan := math.Tan(w.Clipping.FieldOfView / 2.0)

	zn := w.Clipping.NearDistance
	zf := w.Clipping.FarDistance

	// 左手座標系なので一般的な式
	// https://www.notion.so/t-kuni/28fb12fb627480edb34ff2935f5392d5?v=28fb12fb6274802ead99000c2486b3bf&source=copy_link#297b12fb627480dca89bcb0709acc844
	projectionMatrix := mat.NewDense(4, 4, []float64{
		1 / (aspect * tan), 0, 0, 0, // X軸
		0, 1 / tan, 0, 0, // Y軸
		0, 0, zf / (zf - zn), -(zf * zn) / (zf - zn), // Z軸
		0, 0, 1, 0, // 同次座標
	})

	var projected mat.Dense
	projected.Mul(projectionMatrix, &clippedObject.VertexMatrix)

	// mは転置されて4行N列になっている
	_, colCnt := projected.Dims()
	for colIdx := 0; colIdx < colCnt; colIdx++ {
		w := projected.At(3, colIdx)
		projected.Set(0, colIdx, projected.At(0, colIdx)/w)
		projected.Set(1, colIdx, projected.At(1, colIdx)/w)
		projected.Set(2, colIdx, projected.At(2, colIdx)/w)
		projected.Set(3, colIdx, 1)
	}

	clippedObject.VertexMatrix = VartexMatrix{Dense: &projected}

	return clippedObject
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

func (v ViewVolume) PlaneNormal(clippingPlaneType ClippingPlaneType) Vector3D {
	switch clippingPlaneType {
	case Near:
		return v.NearPlaneNormal
	case Far:
		return v.FarPlaneNormal
	case Left:
		return v.LeftPlaneNormal
	case Right:
		return v.RightPlaneNormal
	case Bottom:
		return v.BottomPlaneNormal
	case Top:
		return v.TopPlaneNormal
	}
	return Vector3D{}
}

func (v ViewVolume) PlanePoint(clippingPlaneType ClippingPlaneType) Vector3D {
	switch clippingPlaneType {
	case Near:
		return v.NearTopLeft
	case Far:
		return v.FarTopLeft
	case Left:
		return v.NearTopLeft
	case Right:
		return v.NearTopRight
	case Bottom:
		return v.NearBottomLeft
	case Top:
		return v.NearTopLeft
	}
	return Vector3D{}
}

func (v ViewVolume) ClassifyEdgeByPlane(vertex Vector3D, clippingPlaneType ClippingPlaneType) bool {
	return ClassifyEdgeByPlane(vertex, v.PlaneNormal(clippingPlaneType), v.PlanePoint(clippingPlaneType))
}

type ClippingPlaneType int

const (
	Near ClippingPlaneType = iota
	Far
	Left
	Right
	Bottom
	Top
)

func ClippingPlaneTypes() []ClippingPlaneType {
	return []ClippingPlaneType{Near, Far, Left, Right, Bottom, Top}
}

func (v ViewVolume) SutherlandHodgman(triangle [3]Vector3D) []Vector3D {
	work1Vertices := make([]Vector3D, 0, 10)
	work1Vertices = append(work1Vertices, triangle[0])
	work1Vertices = append(work1Vertices, triangle[1])
	work1Vertices = append(work1Vertices, triangle[2])
	work2Vertices := make([]Vector3D, 0, 10)

	for _, clippingPlaneType := range ClippingPlaneTypes() {
		for i := 0; i < len(work1Vertices); i++ {
			fromIndex := i
			toIndex := (i + 1) % len(work1Vertices)

			fromVertex := work1Vertices[fromIndex]
			toVertex := work1Vertices[toIndex]

			fromInside := v.ClassifyEdgeByPlane(fromVertex, clippingPlaneType)
			toInside := v.ClassifyEdgeByPlane(toVertex, clippingPlaneType)

			if fromInside && toInside {
				// 内から内
				work2Vertices = append(work2Vertices, toVertex)
			} else if fromInside && !toInside {
				// 内から外
				intersectionPoint := v.IntersectPlaneIntersectionPoint(fromVertex, toVertex, clippingPlaneType)
				work2Vertices = append(work2Vertices, intersectionPoint)
			} else if !fromInside && toInside {
				// 外から内
				// 先に交点を追加する。その後、内側の頂点を追加する。（反時計周りの頂点の順番を維持するため）
				intersectionPoint := v.IntersectPlaneIntersectionPoint(fromVertex, toVertex, clippingPlaneType)
				work2Vertices = append(work2Vertices, intersectionPoint)
				work2Vertices = append(work2Vertices, toVertex)
			} else {
				// 外から外
				// 何もしない
			}
		}
		work1Vertices = work2Vertices
		work2Vertices = make([]Vector3D, 0, 10)
	}

	return work1Vertices
}

func (v ViewVolume) ClipObject(o Object) Object {
	newObject := NewDynamicObject()

	for _, triangle := range o.Triangles {
		triangleVertices := [3]Vector3D{
			o.VertexMatrix.GetVertex(triangle[0]),
			o.VertexMatrix.GetVertex(triangle[1]),
			o.VertexMatrix.GetVertex(triangle[2]),
		}
		vertices := v.SutherlandHodgman(triangleVertices)
		triangles := Triangulate(vertices)
		for _, triangle := range triangles {
			newObject.AddTriangle(triangle)
		}
	}

	return v.MargeVertices(newObject.ToObject())
}

type VertexGrid struct {
	grid     map[[3]int][]int
	vertices []Vertex
	epsilon  float64
}

func NewVertexGrid(epsilon float64) VertexGrid {
	return VertexGrid{
		grid:     make(map[[3]int][]int),
		vertices: make([]Vertex, 0, 50),
		epsilon:  epsilon,
	}
}

func (vg VertexGrid) Vertices() []Vertex {
	return vg.vertices
}

func (vg VertexGrid) makeKey(v Vector3D) [3]int {
	return [3]int{int(math.Floor(v[0] / vg.epsilon)), int(math.Floor(v[1] / vg.epsilon)), int(math.Floor(v[2] / vg.epsilon))}
}

func (vg VertexGrid) SearchVertex(v Vector3D) (bool, int) {
	baseGridKey := vg.makeKey(v)
	for _, dx := range []int{0, 1, -1} {
		for _, dy := range []int{0, 1, -1} {
			for _, dz := range []int{0, 1, -1} {
				gridKey := [3]int{baseGridKey[0] + dx, baseGridKey[1] + dy, baseGridKey[2] + dz}
				if candidateVertexIndexes, ok := vg.grid[gridKey]; ok {
					for _, candidateVertexIndex := range candidateVertexIndexes {
						candidateVertex := vg.vertices[candidateVertexIndex]
						if v.Distance(candidateVertex) < vg.epsilon {
							return true, candidateVertexIndex
						}
					}
				}
			}
		}
	}
	return false, 0
}

// AddVertex は頂点を追加します。
// 追加した頂点の新しい添字番号を返します。
func (vg *VertexGrid) AddVertex(v Vector3D) int {
	existSameLocation, sameLocationVertexIndex := vg.SearchVertex(v)
	if existSameLocation {
		return sameLocationVertexIndex
	} else {
		nextIndex := len(vg.vertices)
		gridKey := vg.makeKey(v)
		vg.grid[gridKey] = append(vg.grid[gridKey], nextIndex)
		vg.vertices = append(vg.vertices, v)
		return nextIndex
	}
}

func (v ViewVolume) MargeVertices(o Object) Object {
	grid := NewVertexGrid(1e-2)

	vertexMap := make(map[int]int, 50)
	o.VertexMatrix.EachVertex(func(i int, vertex Vertex) bool {
		vertexMap[i] = grid.AddVertex(vertex)
		return true
	})

	newEdges := make([][2]int, 0, len(o.Edges))
	for _, edge := range o.Edges {
		newEdges = append(newEdges, [2]int{vertexMap[edge[0]], vertexMap[edge[1]]})
	}

	newTriangles := make([][3]int, 0, len(o.Triangles))
	for _, triangle := range o.Triangles {
		newTriangles = append(newTriangles, [3]int{vertexMap[triangle[0]], vertexMap[triangle[1]], vertexMap[triangle[2]]})
	}

	dObj := DynamicObject{
		Vertices:  grid.Vertices(),
		Edges:     CleanEdges(newEdges),
		Triangles: CleanTriangles(newTriangles),
	}
	return dObj.ToObject()
}

func (v ViewVolume) IntersectPlaneIntersectionPoint(fromVertex, toVertex Vector3D, clippingPlaneType ClippingPlaneType) Vector3D {
	planeNormal := v.PlaneNormal(clippingPlaneType)
	planePoint := v.PlanePoint(clippingPlaneType)
	return IntersectPlaneIntersectionPoint(planeNormal, planePoint, fromVertex, toVertex)
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
	// VertexMatrix 頂点の行列
	// 4行N列（頂点数分、横に伸びていきます。同次座標を含みます）
	VertexMatrix VartexMatrix
	// Edges 辺を表す。[0]は始点の頂点の添字番号、[1]は終点の頂点の添字番号。
	Edges [][2]int
	// Triangles 三角形を表す。3つの頂点の添字番号を保持する
	// 右ねじの法則に従って法線の方向が決まります。
	Triangles [][3]int
}

func NewObject(vertices []Vertex, edges [][2]int, triangles [][3]int) Object {
	return Object{
		VertexMatrix: NewVertexMatrix(vertices),
		Edges:        edges,
		Triangles:    triangles,
	}
}

// DynamicObject は動的に頂点を追加できるオブジェクトを表します。
type DynamicObject struct {
	Vertices  []Vertex
	Edges     [][2]int
	Triangles [][3]int
}

func NewDynamicObject() DynamicObject {
	return DynamicObject{
		Vertices:  make([]Vertex, 0, 50),
		Edges:     make([][2]int, 0, 50),
		Triangles: make([][3]int, 0, 50),
	}
}

func (o *DynamicObject) AddTriangle(triangle [3]Vector3D) {
	nextIndex := len(o.Vertices)
	o.Vertices = append(o.Vertices, triangle[0], triangle[1], triangle[2])
	o.Edges = append(o.Edges, [2]int{nextIndex, nextIndex + 1}, [2]int{nextIndex + 1, nextIndex + 2}, [2]int{nextIndex + 2, nextIndex})
	o.Triangles = append(o.Triangles, [3]int{nextIndex, nextIndex + 1, nextIndex + 2})
}

func (o *DynamicObject) ToObject() Object {
	return Object{
		VertexMatrix: NewVertexMatrix(o.Vertices),
		Edges:        o.Edges,
		Triangles:    o.Triangles,
	}
}

type VartexMatrix struct {
	*mat.Dense
}

// NewVertexMatrix は頂点のスライスを行列に変換します
// 4行N列の行列を返します（頂点数分、横に伸びていきます。同次座標を含みます）
func NewVertexMatrix(vertices []Vector3D) VartexMatrix {
	m := mat.NewDense(4, len(vertices), nil)
	for i, v := range vertices {
		m.Set(0, i, v[0])
		m.Set(1, i, v[1])
		m.Set(2, i, v[2])
		m.Set(3, i, 1)
	}
	return VartexMatrix{Dense: m}
}

func (v VartexMatrix) GetVertex(i int) Vertex {
	return Vector3D{v.At(0, i), v.At(1, i), v.At(2, i)}
}

func (v VartexMatrix) EachVertex(f func(int, Vertex) bool) {
	_, colCnt := v.Dims()
	for i := 0; i < colCnt; i++ {
		if !f(i, v.GetVertex(i)) {
			break
		}
	}
}

func (v VartexMatrix) Len() int {
	_, colCnt := v.Dims()
	return colCnt
}

type Vertex = Vector3D

type Vector3D [3]float64

func (v1 Vector3D) X() float64 {
	return v1[0]
}

func (v1 Vector3D) Y() float64 {
	return v1[1]
}

func (v1 Vector3D) Z() float64 {
	return v1[2]
}

func (v1 Vector3D) Vec() mat.VecDense {
	return *mat.NewVecDense(3, v1[:])
}

func (v1 Vector3D) Matrix() mat.Dense {
	return *mat.NewDense(1, 4, append(v1[:], 1))
}

func (v1 Vector3D) Add(v2 Vector3D) Vector3D {
	return Vector3D{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

func (v1 Vector3D) Sub(v2 Vector3D) Vector3D {
	return Vector3D{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

func (v1 Vector3D) Mul(v float64) Vector3D {
	return Vector3D{v1[0] * v, v1[1] * v, v1[2] * v}
}

func (v1 Vector3D) Distance(v2 Vector3D) float64 {
	return math.Sqrt(math.Pow(v1[0]-v2[0], 2) + math.Pow(v1[1]-v2[1], 2) + math.Pow(v1[2]-v2[2], 2))
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
