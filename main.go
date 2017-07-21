package main

import (
	"image"
	"image/color"
	"image/png" // register the PNG format with the image package
	"os"
)

func main() {
	infile, err := os.Open("image_1.png")
	if err != nil {
		// replace this with real error handling
		panic(err)
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		// replace this with real error handling
		panic(err)
	}

	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	// Encode the grayscale image to the output file
	outfile, err := os.Create("image_2.png")
	if err != nil {
		// replace this with real error handling
		panic(err)
	}

	defer outfile.Close()
	png.Encode(outfile, gray)
}
