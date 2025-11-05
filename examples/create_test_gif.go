package main

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
)

func main() {
	// Create a simple animated GIF for testing
	const (
		width  = 64
		height = 48
		frames = 8
	)

	var images []*image.Paletted
	var delays []int

	palette := []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // black
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // white
		color.RGBA{0xff, 0x00, 0x00, 0xff}, // red
		color.RGBA{0x00, 0xff, 0x00, 0xff}, // green
		color.RGBA{0x00, 0x00, 0xff, 0xff}, // blue
		color.RGBA{0xff, 0xff, 0x00, 0xff}, // yellow
	}

	for frame := 0; frame < frames; frame++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

		// Fill background
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.SetColorIndex(x, y, 0) // black background
			}
		}

		// Draw animated circle
		angle := float64(frame) * 2 * math.Pi / float64(frames)
		centerX := width / 2
		centerY := height / 2
		radius := 15.0

		// Circle position
		cx := centerX + int(math.Cos(angle)*10)
		cy := centerY + int(math.Sin(angle)*10)

		// Draw circle
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dx := float64(x - cx)
				dy := float64(y - cy)
				if dx*dx+dy*dy < radius*radius {
					img.SetColorIndex(x, y, uint8(1+frame%5)) // cycling colors
				}
			}
		}

		// Draw border
		for x := 0; x < width; x++ {
			img.SetColorIndex(x, 0, 1)
			img.SetColorIndex(x, height-1, 1)
		}
		for y := 0; y < height; y++ {
			img.SetColorIndex(0, y, 1)
			img.SetColorIndex(width-1, y, 1)
		}

		images = append(images, img)
		delays = append(delays, 10) // 100ms per frame
	}

	// Save GIF
	f, err := os.Create("test_animation.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})

	println("Created test_animation.gif")
}
