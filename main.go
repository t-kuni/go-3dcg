package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
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

	rotateXTheta := float64(0)
	rotateYTheta := float64(0)

	camera := makeCamera()

	running := true
	for running {
		m := invertMatrix(make3dModel())

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						rotateXTheta += math.Pi / 16
					case sdl.K_DOWN:
						rotateXTheta -= math.Pi / 16
					case sdl.K_LEFT:
						rotateYTheta += math.Pi / 16
					case sdl.K_RIGHT:
						rotateYTheta -= math.Pi / 16
					}
				}
			}
		}

		// レンダリング
		m2 := transform(camera, m, rotateXTheta, rotateYTheta)
		render(renderer, m2)
		renderer.Present()
		sdl.Delay(16) // 少し遅延を入れてCPU使用率を下げる
	}
}

func render(renderer *sdl.Renderer, m *mat.Dense) {
	projected := transformParallelProjection(m)
	windowCoords := invertMatrix(transformWindowCoords(projected))

	// ウィンドウの背景色を設定
	renderer.SetDrawColor(255, 255, 255, 255) // 白色
	renderer.Clear()

	// 頂点を直線で結ぶ
	renderer.SetDrawColor(0, 0, 0, 255) // 黒色

	edges := [][2]int{
		{0, 1}, {0, 2}, {0, 3}, // 頂点0から各頂点への辺
		{1, 2}, {1, 3}, // 頂点1から頂点2と3への辺
		{2, 3}, // 頂点2から頂点3への辺
	}

	for _, edge := range edges {
		start := sdl.Point{
			X: int32(windowCoords.At(edge[0], 0)),
			Y: int32(windowCoords.At(edge[0], 1)),
		}
		end := sdl.Point{
			X: int32(windowCoords.At(edge[1], 0)),
			Y: int32(windowCoords.At(edge[1], 1)),
		}
		renderer.DrawLine(start.X, start.Y, end.X, end.Y)
	}
}

type Camera struct {
	Loc *mat.Dense
	Dir *mat.Dense
	Up  *mat.Dense
}

func makeCamera() *Camera {
	loc := mat.NewDense(1, 4, []float64{0.3, -0.5, 0.5, 1})
	dirTo := mat.NewDense(1, 4, []float64{0, 0, 0, 1})
	dir := calcDirection(loc, dirTo)
	return &Camera{
		Loc: loc,
		Dir: dir,
		Up:  mat.NewDense(1, 4, []float64{0, 1, 0, 1}),
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

func makeViewMatrix(camera *Camera) *mat.Dense {
	// 前方ベクトル (Forward Vector)
	f := camera.Dir

	// 右方ベクトル (Right Vector) = 上方向ベクトル (Up) x 前方ベクトル (Forward)
	r := crossProduct(camera.Up, f)

	// 新しい上方ベクトル (New Up Vector) = 前方ベクトル (Forward) x 右方ベクトル (Right)
	u := crossProduct(f, r)

	// カメラの位置ベクトル
	loc := camera.Loc

	// ビュー変換行列を作成
	viewMatrix := mat.NewDense(4, 4, []float64{
		r.At(0, 0), r.At(0, 1), r.At(0, 2), -dotProduct(r, loc),
		u.At(0, 0), u.At(0, 1), u.At(0, 2), -dotProduct(u, loc),
		f.At(0, 0), f.At(0, 1), f.At(0, 2), -dotProduct(f, loc),
		0, 0, 0, 1,
	})

	return viewMatrix
}

// crossProduct は2つのベクトルの外積を計算します。
func crossProduct(a, b *mat.Dense) *mat.Dense {
	return mat.NewDense(1, 4, []float64{
		a.At(0, 1)*b.At(0, 2) - a.At(0, 2)*b.At(0, 1),
		a.At(0, 2)*b.At(0, 0) - a.At(0, 0)*b.At(0, 2),
		a.At(0, 0)*b.At(0, 1) - a.At(0, 1)*b.At(0, 0),
		1, // 同次座標成分
	})
}

// dotProduct は2つのベクトルの内積を計算します。
func dotProduct(a, b *mat.Dense) float64 {
	return a.At(0, 0)*b.At(0, 0) +
		a.At(0, 1)*b.At(0, 1) +
		a.At(0, 2)*b.At(0, 2)
}

func make3dModel() *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		0, -0.35355339059327373, -0.2886751345948129, 1,
		-0.5, -0.35355339059327373, 0.2886751345948129, 1,
		0.5, -0.35355339059327373, 0.2886751345948129, 1,
		0, 0.7071067811865476, 0, 1,
	})
}

func invertMatrix(m *mat.Dense) *mat.Dense {
	var inverted mat.Dense
	inverted.CloneFrom(m.T())

	return &inverted
}

func transform(camera *Camera, m *mat.Dense, rotateXTheta float64, rotateYTheta float64) *mat.Dense {
	m1 := transformRotateX(m, rotateXTheta)
	m2 := transformRotateY(m1, rotateYTheta)
	m3 := transformViewCoords(camera, m2)
	return m3
}

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

func transformParallelProjection(m *mat.Dense) *mat.Dense {
	projectionMatrix := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0, // X軸
		0, 1, 0, 0, // Y軸
		0, 0, 0, 0, // Z軸（無視）
		0, 0, 0, 1, // 同次座標
	})

	var projected mat.Dense
	projected.Mul(projectionMatrix, m)

	return &projected
}

func transformWindowCoords(m *mat.Dense) *mat.Dense {
	// スケーリングと平行移動を行う変換行列
	// スケーリング： 画面の幅、高さ（ピクセル）の値に変換
	// 平行移動： ウィンドウ座標系の原点を画面の中心に移動
	transformMatrix := mat.NewDense(4, 4, []float64{
		float64(winWidth) / 2, 0, 0, float64(winWidth) / 2,
		0, -float64(winHeight) / 2, 0, float64(winHeight) / 2, // Y軸は反転
		0, 0, 1, 0, // Z軸はそのまま
		0, 0, 0, 1, // 同次座標
	})

	var transformed mat.Dense
	transformed.Mul(transformMatrix, m)

	return &transformed
}

func transformViewCoords(camera *Camera, m *mat.Dense) *mat.Dense {
	viewMatrix := makeViewMatrix(camera)

	var transformed mat.Dense
	transformed.Mul(viewMatrix, m)

	return &transformed
}
