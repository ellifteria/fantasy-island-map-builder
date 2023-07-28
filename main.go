package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"

	"github.com/ellifteria/opensimplex2d-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Width             = 1200
	Height            = 1000
	ElevationExponent = 1.5
	MoistureExponent  = 1.5
	NoiseAmplitude    = 4.5
	IslandPercent     = 0.1
)

var (
	ElevationSeed int64 = 13
	MoistureSeed  int64 = 259
	colorArray    []int = make([]int, Width*Height*4)
)

type Game struct {
	gameImage *image.RGBA
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		GenerateRandomSeeds()
		GenerateMap()
	}

	length := Width * Height
	for i := 0; i < length; i++ {
		g.gameImage.Pix[4*i] = uint8(colorArray[4*i+0])
		g.gameImage.Pix[4*i+1] = uint8(colorArray[4*i+1])
		g.gameImage.Pix[4*i+2] = uint8(colorArray[4*i+2])
		g.gameImage.Pix[4*i+3] = uint8(colorArray[4*i+3])
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.gameImage.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func GetBiomeColor(elevation, moisture float64) [4]int {
	switch {
	case elevation < 0.2:
		return [4]int{1, 1, 122, 255}
	case elevation < 0.3:
		return [4]int{3, 138, 255, 255}
	case elevation < 0.45:
		return [4]int{243, 225, 107, 255}
	case elevation < 0.6:
		return [4]int{22, 160, 133, 255}
	case elevation < 0.75:
		return [4]int{108, 122, 137, 255}
	default:
		return [4]int{255, 255, 255, 255}
	}
}

func GenerateRandomSeeds() {
	ElevationSeed = rand.Int63()
	MoistureSeed = rand.Int63()

	for ElevationSeed == MoistureSeed {
		MoistureSeed = rand.Int63()
	}

	fmt.Printf("Elevation seed: %d\n", ElevationSeed)
	fmt.Printf("Moisture seed: %d\n", MoistureSeed)
}

func GenerateMap() {
	var elevationArray [Width][Height]float64
	var moistureArray [Width][Height]float64

	elevationNoise := opensimplex2d.NewNoise(ElevationSeed)
	moistureNoise := opensimplex2d.NewNoise(MoistureSeed)

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			nx := NoiseAmplitude * (2 * (float64(x) - Width/2) / Width)
			ny := NoiseAmplitude * (2 * (float64(y) - Height/2) / Height)

			dx := nx / NoiseAmplitude
			dy := ny / NoiseAmplitude
			d := math.Min(1, (dx*dx+dy*dy)/math.Sqrt(2))

			var e float64 = 0.0
			var eSum float64 = 0.0
			eGains := [3]float64{1, 2, 4}
			for i := range eGains {
				e += (1 / eGains[i]) * elevationNoise.NormalizedNoise2D(eGains[i]*nx, eGains[i]*ny)
				eSum += (1 / eGains[i])
			}
			e = e / (eSum)
			e = (1-IslandPercent)*e + IslandPercent*(1-d)
			elevationArray[x][y] = math.Pow(e, ElevationExponent)

			var m float64 = 0.0
			var mSum float64 = 0.0
			mGains := [3]float64{1, 2, 4}
			for i := range mGains {
				m += (1 / mGains[i]) * moistureNoise.NormalizedNoise2D(mGains[i]*nx, mGains[i]*ny)
				mSum += (1 / mGains[i])
			}
			m = m / (mSum)
			moistureArray[x][y] = math.Pow(m, MoistureExponent)
		}
	}

	colorArray = make([]int, Width*Height*4)

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			color := GetBiomeColor(elevationArray[x][y], moistureArray[x][y])
			// color := [4]int{
			// 	int(elevationArray[x][y] * 255),
			// 	int(elevationArray[x][y] * 255),
			// 	int(elevationArray[x][y] * 255),
			// 	int(elevationArray[x][y] * 255),
			// }
			colorArray[(x+y*Width)*4+0] = color[0]
			colorArray[(x+y*Width)*4+1] = color[1]
			colorArray[(x+y*Width)*4+2] = color[2]
			colorArray[(x+y*Width)*4+3] = color[3]
		}
	}
}

func main() {
	GenerateMap()

	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("Go Fantasy Map Builder")

	g := &Game{
		gameImage: image.NewRGBA(image.Rect(0, 0, Width, Height)),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
