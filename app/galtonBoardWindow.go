package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"go-board/logic"
	"go-board/logic/config"
	"go-board/utils"
	"go-board/utils/gmath"
	"image/color"
	"log"
	"math"
	"os"
	"time"
)

type GaltonBoardWindow struct {
	camPos        gmath.Vector
	camCenterMove gmath.Vector
	camSpeed      float64
	camZoom       float64
	camZoomSpeed  float64
	camDragSpeed  float64
	camDrag       bool
	last          time.Time
	galtonBoard   *logic.GaltonBoard
	imageBoard    *ebiten.Image
	scale         float32
	dt            float64
}

func (g *GaltonBoardWindow) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *GaltonBoardWindow) Update() error {
	dt := time.Since(g.last).Seconds()
	g.dt = dt
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
	g.galtonBoard.RunStep(float32(dt))

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

	//g.imageBoard.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	drawObjects(g.imageBoard, g.galtonBoard, g.scale)

	screen.DrawImage(g.imageBoard, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f DT: %f", ebiten.ActualFPS(), g.dt))
}

func RunWindow() {
	configsRoute, exportRoute := utils.GetBaseAppRoute()
	configsFiles := utils.GetNewExpConfigs(configsRoute)
	configsSelected := utils.SelectConfigsFiles(&configsFiles)

	configSelected := configsSelected[0]
	gbConfigRoute := configsRoute + "/" + configsFiles[configSelected]

	board := executeGaltonBoardWindow(gbConfigRoute, exportRoute)
	width := board.CellSize * float32(board.ColumnNumber)
	height := board.CellSize * float32(board.RowNumber)
	scale := float32(10)

	image := ebiten.NewImage(int(width*scale), int(height*scale))

	ebiten.SetWindowTitle("Galton Board")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	galtonBoard := &GaltonBoardWindow{
		camPos:        gmath.Vector{X: 0, Y: 0},
		camCenterMove: gmath.Vector{X: 0, Y: 0},
		camSpeed:      500,
		camZoom:       1,
		camZoomSpeed:  1.1,
		camDragSpeed:  1,
		galtonBoard:   board,
		imageBoard:    image,
		scale:         scale,
	}

	if err := ebiten.RunGame(galtonBoard); err != nil {
		panic(err)
	}
}

func executeGaltonBoardWindow(gbConfigRoute, exportRoute string) *logic.GaltonBoard {
	currentTime := time.Now()
	timeTxt := currentTime.Format("2006-01-02-15-04-05")
	gbConfig := config.GetConfiguration(gbConfigRoute)
	baseRoute := fmt.Sprintf("%s/exp-%s-%s", exportRoute, gbConfig.Experiment.Name, timeTxt)

	fmt.Printf("Execute of --%s-- at %s\n", gbConfig.Experiment.Name, timeTxt)
	if _, err := os.Stat(baseRoute); os.IsNotExist(err) {
		err := os.MkdirAll(baseRoute, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}

	exportHistoRoute := baseRoute + fmt.Sprintf("/histo-%d.dat", 1)
	exportPathRoute := baseRoute + fmt.Sprintf("/path-%d.dat", 1)

	gbConfig.Experiment.ExportHistogram.Route = exportHistoRoute
	gbConfig.Experiment.ExportPaths.Route = exportPathRoute

	gb := logic.NewGaltonBoard(&gbConfig)

	return gb
}

func drawObjects(image *ebiten.Image, board *logic.GaltonBoard, scale float32) {
	for i := 0; i < len(board.Balls); i++ {
		ballColor := color.RGBA{R: 255, B: 255, A: 255}

		cx := board.Balls[i].Position.X * scale
		cy := board.Balls[i].Position.Y * scale
		radius := board.Balls[i].Radius * scale

		vector.DrawFilledCircle(image, cx, cy, radius, ballColor, false)
	}

	for i := 0; i < len(board.Obstacles); i++ {
		ballColor := color.RGBA{R: 255, B: 255, A: 255}

		cx := board.Obstacles[i].Position.X * scale
		cy := board.Obstacles[i].Position.Y * scale
		radius := board.Obstacles[i].Radius * scale

		vector.DrawFilledCircle(image, cx, cy, radius, ballColor, false)
	}
}
