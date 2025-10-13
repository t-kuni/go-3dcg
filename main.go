package main

import (
	"log"
	"math"

	"github.com/t-kuni/go-3dcg/domain"
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

	world := domain.World{
		Camera: makeCamera(),
		LocatedObjects: []domain.LocatedObject{
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
			discreateWorld := world.Transform(winWidth, winHeight)
			render(renderer, discreateWorld)
			renderer.Present()

			once = true
		}
		sdl.Delay(16) // 少し遅延を入れてCPU使用率を下げる
	}
}

func render(renderer *sdl.Renderer, discreateWorld domain.DiscreteWorld) {
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

func makeCamera() domain.Camera {
	return domain.Camera{
		Location:  domain.Point3D{X: 0, Y: 0, Z: -1.0},
		Direction: domain.Point3D{X: math.Pi / 16, Y: 0, Z: 0},
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

func makeObject() domain.Object {
	return domain.Object{
		Vertices: []domain.Vertex{
			{domain.Point3D{X: -1.0, Y: 0.0, Z: -0.5}},
			{domain.Point3D{X: 1.0, Y: 0.0, Z: -0.5}},
			{domain.Point3D{X: 0.5, Y: 0.0, Z: 0.5}},
			{domain.Point3D{X: 0.0, Y: 1.0, Z: 0.0}},
		},
	}
}

func invertMatrix(m *mat.Dense) *mat.Dense {
	var inverted mat.Dense
	inverted.CloneFrom(m.T())

	return &inverted
}
