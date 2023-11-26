package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth, winHeight int32 = 800, 600
)

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

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Circle in Center", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.SetDrawColor(255, 255, 255, 255) // 白色で背景を塗りつぶす
	renderer.Clear()

	renderer.SetDrawColor(255, 0, 0, 255)              // 赤色で円を描画
	drawCircle(renderer, winWidth/2, winHeight/2, 100) // ウィンドウの中央に半径100の円を描画

	renderer.Present()

	sdl.Delay(5000) // 5秒間ウィンドウを表示
}
