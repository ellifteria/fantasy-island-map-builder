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

// Map parameters
const (
	Width                      = 1200
	Height                     = 1000
	ElevationExponent          = 2
	MoistureExponent           = 1
	ElevationAmplitude         = 8
	MoistureAmplitude          = 4
	IslandPercent              = 0.25
	WaterLevel                 = 0.25
	Gamma              float64 = 0.5
	ShadowSlope        float64 = 0.002
)

var (
	ElevationSeed  int64 = 13
	MoistureSeed   int64 = 259
	ElevationGains       = [3]float64{1, 2, 4}
	MoistureGains        = [3]float64{1, 2, 4}
)

// // Biome type
// type Biome int

// const (
// 	Ocean Biome = iota
// 	DeepOcean
// 	ShallowOcean
// 	Lake
// 	Beach
// 	Tundra
// 	TemperateBroadleafForest
// 	TemperateSteppe
// 	SubtropicalRainforest
// 	AridDesert
// 	ShrubLand
// 	DrySteppe
// 	SemiArdDesert
// 	GrassSavanna
// 	TreeSavanna
// 	DryForest
// 	TropicalRainforest
// 	AlpineTundra
// 	MontaneForest
// )

// Color type
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Image arrays
var (
	ElevationArray []float64 = make([]float64, Width*Height)
	MoistureArray  []float64 = make([]float64, Width*Height)
	colorArray     []Color   = make([]Color, Width*Height)
)

// Ebiten game declarations
type Game struct {
	gameImage *image.RGBA
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		GenerateRandomSeeds()
		GenerateMap()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		AddShadow()
	}

	length := Width * Height
	for i := 0; i < length; i++ {
		g.gameImage.Pix[4*i+0] = colorArray[i].R
		g.gameImage.Pix[4*i+1] = colorArray[i].G
		g.gameImage.Pix[4*i+2] = colorArray[i].B
		g.gameImage.Pix[4*i+3] = colorArray[i].A
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.gameImage.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func AddShadow() {

	var shadowArray []float64 = make([]float64, Width*Height)
	for y := 0; y < Height; y++ {
		shadowArray[y*Width] = 0
	}
	for x := 1; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if ElevationArray[((x-1)+y*Width)] > WaterLevel {
				shadowArray[(x + y*Width)] = math.Max(shadowArray[((x-1)+y*Width)],
					ElevationArray[((x-1)+y*Width)]) - ShadowSlope
			} else {
				shadowArray[(x + y*Width)] = shadowArray[((x-1)+y*Width)] -
					ShadowSlope
			}
		}
	}
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if shadowArray[(x+y*Width)] > ElevationArray[(x+y*Width)] {
				currentColor := colorArray[(x + y*Width)]
				newR := uint8(math.Pow(float64(currentColor.R)/255.0,
					1.0/Gamma) * 255.0)
				newB := uint8(math.Pow(float64(currentColor.B)/255.0,
					1.0/Gamma) * 255.0)
				newG := uint8(math.Pow(float64(currentColor.G)/255.0,
					1.0/Gamma) * 255.0)
				newColor := Color{newR, newG, newB, currentColor.A}
				colorArray[(x + y*Width)] = newColor
			}
		}
	}
}

func GetBiomeColor(elevation, moisture, latitude float64) Color {
	switch {
	case elevation <= WaterLevel:
		return Color{20, 52, 164, 255}
	case elevation < WaterLevel+0.025:
		return Color{157, 145, 122, 255}
	case elevation < 0.5:
		switch {
		case moisture < 0.16:
			return Color{203, 210, 161, 255}
		case moisture < 0.33:
			return Color{143, 169, 96, 255}
		case moisture < 0.66:
			return Color{101, 151, 79, 255}
		default:
			return Color{69, 117, 88, 255}
		}
	case elevation < 0.75:
		switch {
		case moisture < 0.16:
			return Color{203, 210, 161, 255}
		case moisture < 0.5:
			return Color{143, 169, 96, 255}
		case moisture < 0.83:
			return Color{113, 147, 95, 255}
		default:
			return Color{85, 134, 90, 255}
		}
	case elevation < 0.9:
		switch {
		case moisture < 0.33:
			return Color{203, 210, 161, 255}
		case moisture < 0.66:
			return Color{139, 152, 122, 255}
		default:
			return Color{156, 170, 124, 255}
		}
	default:
		switch {
		case moisture < 0.1:
			return Color{85, 85, 85, 255}
		case moisture < 0.2:
			return Color{136, 136, 136, 255}
		case moisture < 0.5:
			return Color{187, 187, 172, 255}
		default:
			return Color{221, 221, 227, 255}
		}
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
	ElevationArray = make([]float64, Width*Height)
	MoistureArray = make([]float64, Width*Height)

	elevationNoise := opensimplex2d.NewNoise(ElevationSeed)
	moistureNoise := opensimplex2d.NewNoise(MoistureSeed)

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			nXE := ElevationAmplitude * (2 * (float64(x) - Width/2) / Width)
			nYE := ElevationAmplitude * (2 * (float64(y) - Height/2) / Height)

			dXE := nXE / ElevationAmplitude
			dYE := nYE / ElevationAmplitude
			dE := math.Min(1, (dXE*dXE+dYE*dYE)/math.Sqrt(2))

			var e float64 = 0.0
			var eSum float64 = 0.0
			ElevationGains := [3]float64{1, 2, 4}
			for i := range ElevationGains {
				e += (1 / ElevationGains[i]) *
					elevationNoise.NormalizedNoise2D(ElevationGains[i]*nXE,
						ElevationGains[i]*nYE)
				eSum += (1 / ElevationGains[i])
			}
			e = e / (eSum)
			e = (1-IslandPercent)*e + IslandPercent*(1-dE)
			ElevationArray[y*Width+x] = math.Pow(e, ElevationExponent)

			nXM := MoistureAmplitude * (2 * (float64(x) - Width/2) / Width)
			nYM := MoistureAmplitude * (2 * (float64(y) - Height/2) / Height)

			var m float64 = 0.0
			var mSum float64 = 0.0
			for i := range MoistureGains {
				m += (1 / MoistureGains[i]) *
					moistureNoise.NormalizedNoise2D(MoistureGains[i]*nXM,
						MoistureGains[i]*nYM)
				mSum += (1 / MoistureGains[i])
			}
			m = m / (mSum)
			MoistureArray[y*Width+x] = math.Pow(m, MoistureExponent)
		}
	}

	colorArray = make([]Color, Width*Height)

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			if ElevationArray[x+y*Width] < WaterLevel {
				ElevationArray[x+y*Width] = WaterLevel
			}
			color := GetBiomeColor(ElevationArray[x+y*Width],
				MoistureArray[x+y*Width], 2.0*float64(y)/float64(Height)-1.0)
			colorArray[(x + y*Width)] = color
		}
	}
}

// Main function
func main() {
	GenerateMap()

	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("Fantasy Map Builder")

	g := &Game{
		gameImage: image.NewRGBA(image.Rect(0, 0, Width, Height)),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
