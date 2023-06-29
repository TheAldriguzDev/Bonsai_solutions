package datadome

import (
	"math"
	"math/rand"
	"time"
)

func GenerateNumber(start float64, end float64, decimalPlaces int) float64 {
	rand.Seed(time.Now().UnixNano())

	randomFloat := rand.Float64()

	rangeStart := start
	rangeEnd := end
	rangeSize := rangeEnd - rangeStart

	randomNumber := rangeStart + (randomFloat * rangeSize)

	pow := math.Pow(10, float64(decimalPlaces))
	randomNumber = math.Round(randomNumber*pow) / pow

	return randomNumber
}

func generateCpu() int {
	rand.Seed(time.Now().UnixNano())

	randomCpu := rand.Intn(16) + 1

	return randomCpu
}

func GenerateWindow() (int, int) {
	rand.Seed(time.Now().UnixNano())

	averageHeight := 800
	averageWidth := 1200

	heightVariation := rand.Intn(200) - 100
	widthVariation := rand.Intn(200) - 100

	outerHeight := averageHeight + heightVariation
	outerWidth := averageWidth + widthVariation

	return outerHeight, outerWidth
}

func GenerateScreen() (int, int) {
	rand.Seed(time.Now().UnixNano())

	averageHeight := 1080
	averageWidth := 1920

	heightVariation := rand.Intn(200) - 100
	widthVariation := rand.Intn(200) - 100

	screenHeight := averageHeight + heightVariation
	screenWidth := averageWidth + widthVariation

	return screenHeight, screenWidth
}
