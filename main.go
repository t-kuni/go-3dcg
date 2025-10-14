package main

import (
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/t-kuni/go-3dcg/domain"
)

const (
	winWidth, winHeight int32 = 800, 600
)

func main() {
	world := domain.World{
		Camera: domain.Camera{
			Location:  domain.Point3D{X: 0, Y: 0, Z: -1.0},
			Direction: domain.Point3D{X: math.Pi / 16, Y: 0, Z: 0},
		},
		LocatedObjects: []domain.LocatedObject{
			{
				X: 0, Y: 0, Z: 1,
				Object: domain.Object{
					Vertices: []domain.Vertex{
						{Point3D: domain.Point3D{X: -1.0, Y: 0.0, Z: -0.5}},
						{Point3D: domain.Point3D{X: 1.0, Y: 0.0, Z: -0.5}},
						{Point3D: domain.Point3D{X: 0.0, Y: 0.0, Z: 0.5}},
						{Point3D: domain.Point3D{X: 0.0, Y: 1.0, Z: 0.0}},
					},
				},
			},
		},
		Viewport: domain.Viewport{
			Width:      winWidth,
			Height:     winHeight,
			ScaleRatio: 0.25,
		},
	}

	// 3D世界を2D座標に変換
	discreteWorld := world.Transform()

	// 画像コンテキストを作成
	dc := gg.NewContext(int(winWidth), int(winHeight))

	// 画像をレンダリング
	render(dc, discreteWorld)

	// 画像ファイルとして保存
	err := dc.SavePNG("render.png")
	if err != nil {
		log.Fatalf("画像の保存に失敗しました: %s", err)
	}

	log.Println("render.png を出力しました")
}

func render(dc *gg.Context, discreteWorld domain.DiscreteWorld) {
	// 背景色を白に設定
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// 線の色を黒に設定
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)

	// 各DiscreteObjectについて、全ての頂点を結ぶ直線を描画
	for _, discreteObject := range discreteWorld.DiscreteObjects {
		vertices := discreteObject.Vertices
		vertexCount := len(vertices)

		// 全ての頂点の組み合わせで直線を描画
		for i := 0; i < vertexCount; i++ {
			for j := i + 1; j < vertexCount; j++ {
				start := vertices[i]
				end := vertices[j]
				dc.DrawLine(float64(start.X), float64(start.Y), float64(end.X), float64(end.Y))
			}
		}
	}

	// 描画を実行
	dc.Stroke()
}
