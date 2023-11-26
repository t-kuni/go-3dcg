package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"gonum.org/v1/gonum/mat"
	"log"
)

const (
	winWidth, winHeight int32 = 800, 600
)

func main() {
	m := make3dModel()
	inverted := invertMatrix(m)
	projected := transformParallelProjection(inverted)
	windowCoords := invertMatrix(transformWindowCoords(projected))

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

	// レンダリングを表示
	renderer.Present()

	// ウィンドウを閉じるまで待機
	sdl.Delay(5000)
}

func drawCircle(renderer *sdl.Renderer, x, y, r int32) {
	for w := 0; w < int(r*2); w++ {
		for h := 0; h < int(r*2); h++ {
			dx := r - int32(w) // horizontal offset
			dy := r - int32(h) // vertical offset
			if dx*dx+dy*dy <= r*r {
				renderer.DrawPoint(x+int32(w)-r, y+int32(h)-r)
			}
		}
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
