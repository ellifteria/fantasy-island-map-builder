package main

import (
	"image"
	"log"

	"github.com/ellifteria/opensimplex2d-go"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	width  = 1000
	height = 1000
)

type Game struct {
	gameImage  *image.RGBA
	pixelArray []int
}

func (g *Game) Update() error {
	length := width * height
	for i := 0; i < length; i++ {
		g.gameImage.Pix[4*i] = uint8(g.pixelArray[4*i+0])
		g.gameImage.Pix[4*i+1] = uint8(g.pixelArray[4*i+1])
		g.gameImage.Pix[4*i+2] = uint8(g.pixelArray[4*i+2])
		g.gameImage.Pix[4*i+3] = uint8(g.pixelArray[4*i+3])
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.gameImage.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func getBiomeColor(height float64) [4]int {
	switch {
	case height < -0.5:
		return [4]int{1, 1, 122, 255}
	case height < 0:
		return [4]int{3, 138, 255, 255}
	case height < 0.25:
		return [4]int{243, 225, 107, 255}
	case height < 0.5:
		return [4]int{22, 160, 133, 255}
	case height < 0.75:
		return [4]int{108, 122, 137, 255}
	default:
		return [4]int{255, 255, 255, 255}
	}
}

func main() {
	var heightArray [width][height]float64

	openSimplexNoise := opensimplex2d.NewNoise(13)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			nx := (0.01*width)*float64(x)/width - 0.5
			ny := (0.01*height)*float64(y)/height - 0.5
			heightArray[x][y] = openSimplexNoise.Noise2D(
				float64(nx),
				float64(ny),
			)
		}
	}

	var colorArray [width * height * 4]int

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// color := getBiomeColor(heightArray[x][y])
			color := [4]int{
				int((heightArray[x][y] + 1.0) / 2.0 * 255),
				int((heightArray[x][y] + 1.0) / 2.0 * 255),
				int((heightArray[x][y] + 1.0) / 2.0 * 255),
				int((heightArray[x][y] + 1.0) / 2.0 * 255),
			}
			colorArray[(x+y*width)*4+0] = color[0]
			colorArray[(x+y*width)*4+1] = color[1]
			colorArray[(x+y*width)*4+2] = color[2]
			colorArray[(x+y*width)*4+3] = color[3]
		}
	}

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Go Fantasy Map Builder")

	g := &Game{
		gameImage:  image.NewRGBA(image.Rect(0, 0, width, height)),
		pixelArray: colorArray[:],
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
