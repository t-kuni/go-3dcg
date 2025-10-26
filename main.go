package main

import (
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/t-kuni/go-3dcg/domain"
)

const (
	width, height int32 = 800, 600
)

func main() {
	world := domain.World{
		Camera: domain.Camera{
			Location:  domain.Vector3D{0, 0, 0},
			Direction: domain.Vector3D{0, 0, 0},
		},
		LocatedObjects: []domain.LocatedObject{
			{
				X: 0, Y: 0, Z: 0,
				Object: domain.Object{
					Vertices: []domain.Vertex{
						{Vector3D: domain.Vector3D{-0.3, 0.0, 1.1}}, // 左下
						{Vector3D: domain.Vector3D{0.3, 0.0, 1.1}},  // 右下
						{Vector3D: domain.Vector3D{0.0, 0.0, 1.9}},  // 奥
						{Vector3D: domain.Vector3D{0.0, 0.3, 1.1}},  // 上
					},
					Edges: [][2]int{
						{0, 1},
						{0, 2},
						{0, 3},
						{1, 2},
						{1, 3},
						{2, 3},
					},
					Triangles: [][3]int{
						{0, 1, 2},
						{0, 1, 3},
						{0, 2, 3},
						{1, 2, 3},
					},
				},
			},
		},
		Viewport: domain.Viewport{
			Width:      width,
			Height:     height,
			ScaleRatio: 0.5,
		},
		Clipping: domain.Clipping{
			NearDistance: 1.0,
			FarDistance:  2.0,
			FieldOfView:  math.Pi / 4,
		},
	}

	// 3D世界を2D座標に変換
	discreteWorld := world.Transform()

	// 画像コンテキストを作成
	dc := gg.NewContext(int(width), int(height))

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

	// 各DiscreteObjectについて、Edgesに従って直線を描画
	for _, discreteObject := range discreteWorld.DiscreteObjects {
		vertices := discreteObject.Vertices
		edges := discreteObject.Edges

		// Edgesに従って直線を描画
		for _, edge := range edges {
			if edge[0] < len(vertices) && edge[1] < len(vertices) {
				start := vertices[edge[0]]
				end := vertices[edge[1]]
				dc.DrawLine(float64(start.X), float64(start.Y), float64(end.X), float64(end.Y))
			}
		}
	}

	// 描画を実行
	dc.Stroke()
}
