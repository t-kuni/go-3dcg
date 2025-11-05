package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t-kuni/go-3dcg/domain"
)

const (
	width, height int32 = 800, 600
)

type Game struct {
	world domain.World
}

func (g *Game) Update() error {
	// キー入力によるカメラ移動
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.world.Camera.Location[2] += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.world.Camera.Location[2] -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.world.Camera.Location[0] += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.world.Camera.Location[0] -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.world.Camera.Location[1] += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.world.Camera.Location[1] -= 0.1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	frameBuffer := g.world.Transform()

	// 画面をクリア（背景色を白に設定）
	screen.Fill(color.RGBA{255, 255, 255, 255})

	// 各DiscreteObjectについて、Edgesに従って直線を描画
	for key, value := range frameBuffer {
		screen.Set(int(key.X), int(key.Y), value.Color)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(width), int(height)
}

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
					VertexMatrix: domain.NewVertexMatrix([]domain.Vector3D{
						{-0.2, -0.1, 1.1}, // 左下
						{0.2, -0.1, 1.1},  // 右下
						{0.0, 0.1, 1.9},   // 奥
						{0.0, 0.1, 1.1},   // 上
					}),
					Edges: [][2]int{
						{0, 1},
						{0, 2},
						{0, 3},
						{1, 2},
						{1, 3},
						{2, 3},
					},
					Triangles: [][3]int{
						{0, 1, 3},
						{0, 2, 1},
						{0, 3, 2},
						{1, 2, 3},
					},
					TriangleColors: []color.RGBA{
						{255, 0, 0, 255},
						{0, 255, 0, 255},
						{0, 0, 255, 255},
						{255, 255, 0, 255},
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
			NearDistance: 0.1,
			FarDistance:  10.0,
			FieldOfView:  math.Pi / 4,
		},
	}

	game := &Game{
		world: world,
	}

	ebiten.SetWindowSize(int(width), int(height))
	ebiten.SetWindowTitle("3D CG with Ebiten")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
