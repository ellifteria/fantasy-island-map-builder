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
	Width             = 1200
	Height            = 1000
	ElevationExponent = 1.7
	MoistureExponent  = 1.7
	NoiseAmplitude    = 5
	IslandPercent     = 0.1
)

var (
	ElevationSeed int64   = 13
	MoistureSeed  int64   = 259
	colorArray    []Color = make([]Color, Width*Height)
)

// Biome type
type Biome int

const (
	Ocean Biome = iota
	DeepOcean
	ShallowOcean
	Lake
	Beach
	Tundra
	TemperateBroadleafForest
	TemperateSteppe
	SubtropicalRainforest
	AridDesert
	ShrubLand
	DrySteppe
	SemiArdDesert
	GrassSavanna
	TreeSavanna
	DryForest
	TropicalRainforest
	AlpineTundra
	MontaneForest
)

// Color type
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Ebiten game declarations
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

// Map creation functions
// func GetBiome(elevation, moisture, latitude float64) Biome {
// 	// switch {
// 	// case elevation < 0.2:
// 	// 	return DeepOcean
// 	// case elevation < 0.3:
// 	// 	return ShallowOcean
// 	// case elevation < 0.45:
// 	// 	switch {
// 	// 	case moisture < 0.3:
// 	// 		return AridDesert
// 	// 	case moisture < 0.375:
// 	// 		return SemiArdDesert
// 	// 	case moisture < 0.45:
// 	// 		return ShrubLand
// 	// 	case moisture < 0.65:
// 	// 		return TemperateSteppe
// 	// 	case moisture < 0.8:
// 	// 		return TemperateBroadleafForest
// 	// 	default:
// 	// 		return SubtropicalRainforest
// 	// 	}
// 	// case elevation < 0.55:
// 	// 	switch {
// 	// 	case moisture < 0.3:
// 	// 		return DrySteppe
// 	// 	}
// 	// 	return DrySteppe
// 	// default:
// 	// 	return AlpineTundra
// 	// }

// 	switch {
// 	case elevation < 0.2:
// 		return Ocean
// 	}
// }

// func GetBiomeColor(elevation, moisture, latitude float64) Color {
// 	biome := GetBiome(elevation, moisture, latitude)

// 	switch biome {
// 	case DeepOcean:
// 		return Color{R: 1, G: 1, B: 122, A: 255}
// 	case ShallowOcean:
// 		return Color{R: 3, G: 138, B: 255, A: 255}
// 	case Ocean:
// 		return Color{R: 68, G: 68, B: 68, A: 255}
// 	case Beach:
// 		return Color{R: 243, G: 225, B: 107, A: 255}
// 	case Tundra:
// 		return Color{R: 154, G: 202, B: 189, A: 255}
// 	case TemperateSteppe:
// 		return Color{R: 243, G: 231, B: 113, A: 255}
// 	case TemperateBroadleafForest:
// 		return Color{R: 161, G: 214, B: 94, A: 255}
// 	case SubtropicalRainforest:
// 		return Color{R: 45, G: 102, B: 28, A: 255}
// 	case AridDesert:
// 		return Color{R: 121, G: 69, B: 46, A: 255}
// 	case ShrubLand:
// 		return Color{R: 160, G: 99, B: 68, A: 255}
// 	case DrySteppe:
// 		return Color{R: 132, G: 112, B: 60, A: 255}
// 	case SemiArdDesert:
// 		return Color{R: 207, G: 171, B: 122, A: 255}
// 	case GrassSavanna:
// 		return Color{R: 192, G: 189, B: 84, A: 255}
// 	case TreeSavanna:
// 		return Color{R: 154, G: 149, B: 50, A: 255}
// 	case DryForest:
// 		return Color{R: 101, G: 121, B: 49, A: 255}
// 	case TropicalRainforest:
// 		return Color{R: 27, G: 69, B: 14, A: 255}
// 	case AlpineTundra:
// 		return Color{R: 154, G: 173, B: 207, A: 255}
// 	case MontaneForest:
// 		return Color{R: 68, G: 129, B: 131, A: 255}
// 	default:
// 		return Color{R: 0, G: 0, B: 0, A: 255}
// 	}
// }

func GetBiomeColor(elevation, moisture, latitude float64) Color {
	switch {
	case elevation < 0.25:
		return Color{68, 68, 118, 255}
	case elevation < 0.275:
		return Color{157, 145, 122, 255}
	case elevation < 0.45:
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
	case elevation < 0.65:
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
	case elevation < 0.85:
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

	colorArray = make([]Color, Width*Height)

	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			color := GetBiomeColor(elevationArray[x][y],
				moistureArray[x][y], 2.0*float64(y)/float64(Height)-1.0)
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
