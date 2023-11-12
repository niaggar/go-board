package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"go-board/utils/gmath"
	"image/color"
	"math"
	"time"
)

const (
	screenWidth  = 500
	screenHeight = 500
)

type GaltonBoardWindow struct {
	camPos        gmath.Vector
	camCenterMove gmath.Vector
	camSpeed      float64
	camZoom       float64
	camZoomSpeed  float64
	camDragSpeed  float64
	camDrag       bool

	last time.Time
}

func (g *GaltonBoardWindow) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *GaltonBoardWindow) Update() error {
	dt := time.Since(g.last).Seconds()
	g.last = time.Now()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		x, y := ebiten.CursorPosition()

		g.camCenterMove = gmath.Vector{X: float32(x), Y: float32(y)}
		g.camDrag = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		g.camDrag = false
		g.camCenterMove = gmath.Vector{X: 0, Y: 0}
	}
	if g.camDrag {
		x, y := ebiten.CursorPosition()

		disX := g.camCenterMove.X - float32(x)
		disY := g.camCenterMove.Y - float32(y)

		g.camPos.X -= disX
		g.camPos.Y -= disY

		g.camCenterMove = gmath.Vector{X: float32(x), Y: float32(y)}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camPos.X += float32(g.camSpeed * dt * 1 / g.camZoom)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camPos.X -= float32(g.camSpeed * dt * 1 / g.camZoom)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camPos.Y -= float32(g.camSpeed * dt * 1 / g.camZoom)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camPos.Y += float32(g.camSpeed * dt * 1 / g.camZoom)
	}

	_, wy := ebiten.Wheel()
	g.camZoom *= math.Pow(g.camZoomSpeed, wy)

	return nil
}

func (g *GaltonBoardWindow) Camera() ebiten.GeoM {
	cam := ebiten.GeoM{}
	cam.Scale(g.camZoom, g.camZoom)
	cam.Translate(float64(g.camPos.X), float64(g.camPos.Y))
	return cam
}

func (g *GaltonBoardWindow) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(g.Camera())

	image := ebiten.NewImage(300, 300)

	ballColor := color.RGBA{R: 255, B: 255, A: 255}
	image.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	vector.DrawFilledCircle(image, 0, 0, 100, ballColor, false)

	screen.DrawImage(image, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func RunWindow() {
	ebiten.SetWindowTitle("Galton Board")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	galtonBoard := &GaltonBoardWindow{
		camPos: gmath.Vector{
			X: 0,
			Y: 0,
		},
		camCenterMove: gmath.Vector{
			X: 0,
			Y: 0,
		},
		camSpeed:     500,
		camZoom:      1,
		camZoomSpeed: 1.2,
		camDragSpeed: 1,
	}

	if err := ebiten.RunGame(galtonBoard); err != nil {
		panic(err)
	}
}
