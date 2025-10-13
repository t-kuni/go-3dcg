package main

import (
	"log"
	"math"

	"github.com/t-kuni/go-3dcg/util"
	"github.com/veandco/go-sdl2/sdl"
	"gonum.org/v1/gonum/mat"
)

const (
	winWidth, winHeight int32 = 800, 600
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("SDLを初期化できませんでした: %s", err)
	}
	defer sdl.Quit()

	// ウィンドウとレンダラの作成
	window, renderer, err := sdl.CreateWindowAndRenderer(winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("ウィンドウとレンダラを作成できませんでした: %s", err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	world := World{
		Camera: makeCamera(),
		LocatedObjects: []LocatedObject{
			{X: 0.1, Y: 0.1, Z: 0, Object: makeObject()},
		},
	}

	once := false

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					// case sdl.K_UP:
					// 	rotateXTheta += math.Pi / 16
					// case sdl.K_DOWN:
					// 	rotateXTheta -= math.Pi / 16
					// case sdl.K_LEFT:
					// 	rotateYTheta += math.Pi / 16
					// case sdl.K_RIGHT:
					// 	rotateYTheta -= math.Pi / 16
					}
				}
			}
		}
		if !once {
			discreateWorld := world.Transform()
			render(renderer, discreateWorld)
			renderer.Present()

			once = true
		}
		sdl.Delay(16) // 少し遅延を入れてCPU使用率を下げる
	}
}

func render(renderer *sdl.Renderer, discreateWorld DiscreteWorld) {
	// ウィンドウの背景色を設定
	renderer.SetDrawColor(255, 255, 255, 255) // 白色
	renderer.Clear()

	// 頂点を直線で結ぶ
	renderer.SetDrawColor(0, 0, 0, 255) // 黒色

	// 各DiscreteObjectについて、全ての頂点を結ぶ直線を描画
	for _, discreteObject := range discreateWorld.DiscreteObjects {
		vertices := discreteObject.Vertices
		vertexCount := len(vertices)

		// 全ての頂点の組み合わせで直線を描画
		for i := 0; i < vertexCount; i++ {
			for j := i + 1; j < vertexCount; j++ {
				start := vertices[i]
				end := vertices[j]
				renderer.DrawLine(start.X, start.Y, end.X, end.Y)
			}
		}
	}
}

type Camera struct {
	Location  Point3D
	Direction Point3D
	// Up        *mat.Dense
}

func makeCamera() Camera {
	return Camera{
		Location:  Point3D{X: 0, Y: 0, Z: -1.0},
		Direction: Point3D{X: math.Pi / 16, Y: 0, Z: 0},
	}
}

func calcDirection(a, b *mat.Dense) *mat.Dense {
	// aからbへのベクトルを計算
	dx := b.At(0, 0) - a.At(0, 0)
	dy := b.At(0, 1) - a.At(0, 1)
	dz := b.At(0, 2) - a.At(0, 2)

	// ベクトルの長さを計算
	length := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// ベクトルを正規化
	if length != 0 {
		dx /= length
		dy /= length
		dz /= length
	}

	// 正規化されたベクトルを返す
	return mat.NewDense(1, 4, []float64{dx, dy, dz, 1})
}

// func makeViewMatrix(camera *Camera) *mat.Dense {
// 	// 前方ベクトル (Forward Vector)
// 	f := camera.Dir

// 	// 右方ベクトル (Right Vector) = 上方向ベクトル (Up) x 前方ベクトル (Forward)
// 	r := crossProduct(camera.Up, f)

// 	// 新しい上方ベクトル (New Up Vector) = 前方ベクトル (Forward) x 右方ベクトル (Right)
// 	u := crossProduct(f, r)

// 	// カメラの位置ベクトル
// 	loc := camera.Loc

// 	// ビュー変換行列を作成
// 	viewMatrix := mat.NewDense(4, 4, []float64{
// 		r.At(0, 0), r.At(0, 1), r.At(0, 2), -dotProduct(r, loc),
// 		u.At(0, 0), u.At(0, 1), u.At(0, 2), -dotProduct(u, loc),
// 		f.At(0, 0), f.At(0, 1), f.At(0, 2), -dotProduct(f, loc),
// 		0, 0, 0, 1,
// 	})

// 	return viewMatrix
// }

// // crossProduct は2つのベクトルの外積を計算します。
// func crossProduct(a, b *mat.Dense) *mat.Dense {
// 	return mat.NewDense(1, 4, []float64{
// 		a.At(0, 1)*b.At(0, 2) - a.At(0, 2)*b.At(0, 1),
// 		a.At(0, 2)*b.At(0, 0) - a.At(0, 0)*b.At(0, 2),
// 		a.At(0, 0)*b.At(0, 1) - a.At(0, 1)*b.At(0, 0),
// 		1, // 同次座標成分
// 	})
// }

// // dotProduct は2つのベクトルの内積を計算します。
// func dotProduct(a, b *mat.Dense) float64 {
// 	return a.At(0, 0)*b.At(0, 0) +
// 		a.At(0, 1)*b.At(0, 1) +
// 		a.At(0, 2)*b.At(0, 2)
// }

func makeObject() Object {
	return Object{
		Vertices: []Vertex{
			{Point3D{X: -1.0, Y: 0.0, Z: -0.5}},
			{Point3D{X: 1.0, Y: 0.0, Z: -0.5}},
			{Point3D{X: 0.5, Y: 0.0, Z: 0.5}},
			{Point3D{X: 0.0, Y: 1.0, Z: 0.0}},
		},
	}
}

func invertMatrix(m *mat.Dense) *mat.Dense {
	var inverted mat.Dense
	inverted.CloneFrom(m.T())

	return &inverted
}

// func transform(camera *Camera, m *mat.Dense) *mat.Dense {
// 	m3 := transformViewCoords(camera, m)
// 	return m3
// }

func transformRotateX(m *mat.Dense, theta float64) *mat.Dense {
	rotateMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, math.Cos(theta), -math.Sin(theta), 0,
		0, math.Sin(theta), math.Cos(theta), 0,
		0, 0, 0, 1,
	})

	var rotated mat.Dense
	rotated.Mul(rotateMatrix, m)
	return &rotated
}

func transformRotateY(m *mat.Dense, theta float64) *mat.Dense {
	rotateMatrix := mat.NewDense(4, 4, []float64{
		math.Cos(theta), 0, math.Sin(theta), 0,
		0, 1, 0, 0,
		-math.Sin(theta), 0, math.Cos(theta), 0,
		0, 0, 0, 1,
	})

	var rotated mat.Dense
	rotated.Mul(rotateMatrix, m)
	return &rotated
}

// func transformViewCoords(camera *Camera, m *mat.Dense) *mat.Dense {
// 	viewMatrix := makeViewMatrix(camera)

// 	var transformed mat.Dense
// 	transformed.Mul(viewMatrix, m)

// 	return &transformed
// }

// type CameraView struct {
// 	Camera Camera
// 	Objects []Object
// }

type TransformedWorld struct {
	TransformedObjects []TransformedObject
}

type TransformedObject struct {
	Vertices mat.Dense
}

type World struct {
	Camera         Camera
	LocatedObjects []LocatedObject
}

func (w World) Transform() DiscreteWorld {
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
		m = util.TransformViewport(m, winWidth, winHeight)

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
