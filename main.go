package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/t-kuni/go-3dcg/domain"
)

const (
	width, height int32 = 800, 600
)

type Game struct {
	world      domain.World
	frameCount int
}

func (g *Game) Update() error {
	// キー入力によるカメラ移動
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.world.Camera.Location[2] += 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.world.Camera.Location[2] -= 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.world.Camera.Location[0] += 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.world.Camera.Location[0] -= 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.world.Camera.Location[1] += 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.world.Camera.Location[1] -= 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.world.Camera.Direction[1] += -0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.world.Camera.Direction[1] += 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		g.world.Camera.Direction[0] += -0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		g.world.Camera.Direction[0] += 0.025
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.frameCount++

	g.world.LocatedObjects[0].Rotation[1] += 0.1
	g.world.LocatedObjects[0].Rotation[0] += 0.1

	frameBuffer := g.world.Transform()

	// 画面をクリア（背景色を白に設定）
	screen.Fill(color.RGBA{255, 255, 255, 255})

	for key, value := range frameBuffer {
		screen.Set(int(key.X), int(key.Y), value.Color)
	}

	// FPSを表示
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
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
				Location: domain.Vector3D{0, 0, 2},
				Scale:    domain.Vector3D{1.0, 1.0, 1.0},
				Rotation: domain.Vector3D{0.0, 0.0, 0.0},
				Object:   domain.NewTetrahedronObject(0.15),
			},
			{
				Location: domain.Vector3D{0.0, 0.0, 2.5},
				Scale:    domain.Vector3D{1.0, 1.0, 1.0},
				Rotation: domain.Vector3D{0.0, 0.0, 0.0},
				Object:   domain.NewPlaneObject(0.3, 0.3, color.RGBA{50, 50, 50, 255}),
			},
			{
				Location: domain.Vector3D{-0.2, 0.0, 2},
				Scale:    domain.Vector3D{1.0, 1.0, 1.0},
				Rotation: domain.Vector3D{0.0, -math.Pi / 2.0, 0.0},
				Object:   domain.NewPlaneObject(0.3, 0.3, color.RGBA{50, 50, 50, 255}),
			},
			{
				Location: domain.Vector3D{0.2, 0.0, 2},
				Scale:    domain.Vector3D{1.0, 1.0, 1.0},
				Rotation: domain.Vector3D{0.0, math.Pi / 2.0, 0.0},
				Object:   domain.NewPlaneObject(0.3, 0.3, color.RGBA{50, 50, 50, 255}),
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
