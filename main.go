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
						rotateXTheta += math.Pi / 8
					case sdl.K_DOWN:
						rotateXTheta -= math.Pi / 8
					case sdl.K_LEFT:
						rotateYTheta += math.Pi / 8
					case sdl.K_RIGHT:
						rotateYTheta -= math.Pi / 8
					}
				}
			}
		}

		// レンダリング
		m2 := transform(m, rotateXTheta, rotateYTheta)
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

func transform(m *mat.Dense, rotateXTheta float64, rotateYTheta float64) *mat.Dense {
	m1 := transformRotateX(m, rotateXTheta)
	m2 := transformRotateY(m1, rotateYTheta)
	return m2
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
