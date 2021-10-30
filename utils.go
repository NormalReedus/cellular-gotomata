package main

// func loadImage(path string) *ebiten.Image {
// 	image, _, err := ebitenutil.NewImageFromFile(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return image
// }

// Includes min and max values
func between(num, min, max int) bool {
	return num >= min && num <= max
}
