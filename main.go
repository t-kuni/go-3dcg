package main

import (
	"log"
	"math"

	"github.com/t-kuni/go-3dcg/domain"
	"github.com/veandco/go-sdl2/sdl"
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
		Camera: domain.Camera{
			Location:  domain.Point3D{X: 0, Y: 0, Z: -1.0},
			Direction: domain.Point3D{X: math.Pi / 16, Y: 0, Z: 0},
		},
		LocatedObjects: []domain.LocatedObject{
			{
				X: 0.1, Y: 0.1, Z: 1,
				Object: domain.Object{
					Vertices: []domain.Vertex{
						{domain.Point3D{X: -1.0, Y: 0.0, Z: -0.5}},
						{domain.Point3D{X: 1.0, Y: 0.0, Z: -0.5}},
						{domain.Point3D{X: 0.5, Y: 0.0, Z: 0.5}},
						{domain.Point3D{X: 0.0, Y: 1.0, Z: 0.0}},
					},
				},
			},
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
