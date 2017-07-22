/**
 * FAST intends for "Features from Accelerated Segment Test". This method
 * performs a point segment test corner detection. The segment test
 * criterion operates by considering a circle of sixteen pixels around the
 * corner candidate p. The detector classifies p as a corner if there exists
 * a set of n contiguous pixelsin the circle which are all brighter than the
 * intensity of the candidate pixel Ip plus a threshold t, or all darker
 * than Ip − t.
 *
 *       15 00 01
 *    14          02
 * 13                03
 * 12       []       04
 * 11                05
 *    10          06
 *       09 08 07
 *
 * For more reference:
 * http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.60.3991&rep=rep1&type=pdf
 */

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png" // register the PNG format with the image package
	"os"
)

/**
 * Finds corners coordinates on the graysacaled image.
 */
func findCorners(pixels []int, width, height, threshold int) []int {
	var circleOffsets = getCircleOffsets(width)
	var circlePixels [16]int
	var corners []int

	// When looping through the image pixels, skips the first three lines from
	// the image boundaries to constrain the surrounding circle inside the image
	// area.
	for i := 3; i < height-3; i++ {
		for j := 3; j < width-3; j++ {
			var w = i*width + j
			var p = pixels[w]

			// Loops the circle offsets to read the pixel value for the sixteen
			// surrounding pixels.
			for k := 0; k < 16; k++ {
				circlePixels[k] = pixels[w+circleOffsets[k]]
			}

			if isCorner(p, circlePixels, threshold) {
				// The pixel p is classified as a corner, as optimization increment j
				// by the circle radius 3 to skip the neighbor pixels inside the
				// surrounding circle. This can be removed without compromising the
				// result.
				corners = append(corners, j, i)
				j += 3
			}
		}
	}

	return corners
}

/**
 * Checks if the circle pixel is within the corner of the candidate pixel p
 * by a threshold.
 */
func isCorner(p int, circlePixels [16]int, threshold int) bool {
	if isTriviallyExcluded(circlePixels, p, threshold) {
		return false
	}

	for x := 0; x < 16; x++ {
		var darker = true
		var brighter = true

		for y := 0; y < 9; y++ {
			var circlePixel = circlePixels[(x+y)&15]

			if !isBrighter(p, circlePixel, threshold) {
				brighter = false
				if darker == false {
					break
				}
			}

			if !isDarker(p, circlePixel, threshold) {
				darker = false
				if brighter == false {
					break
				}
			}
		}

		if brighter || darker {
			return true
		}
	}

	return false
}

/**
 * Fast check to test if the candidate pixel is a trivially excluded value.
 * In order to be a corner, the candidate pixel value should be darker or
 * brighter than 9-12 surrounding pixels, when at least three of the top,
 * bottom, left and right pixels are brighter or darker it can be
 * automatically excluded improving the performance.
 */
func isTriviallyExcluded(circlePixels [16]int, p int, threshold int) bool {
	var count = 0
	var circleBottom = circlePixels[8]
	var circleLeft = circlePixels[12]
	var circleRight = circlePixels[4]
	var circleTop = circlePixels[0]

	if isBrighter(circleTop, p, threshold) {
		count++
	}
	if isBrighter(circleRight, p, threshold) {
		count++
	}
	if isBrighter(circleBottom, p, threshold) {
		count++
	}
	if isBrighter(circleLeft, p, threshold) {
		count++
	}

	if count < 3 {
		count = 0
		if isDarker(circleTop, p, threshold) {
			count++
		}
		if isDarker(circleRight, p, threshold) {
			count++
		}
		if isDarker(circleBottom, p, threshold) {
			count++
		}
		if isDarker(circleLeft, p, threshold) {
			count++
		}
		if count < 3 {
			return true
		}
	}

	return false
}

/**
 * Checks if the circle pixel is brighter than the candidate pixel p by
 * a threshold.
 */
func isBrighter(circlePixel int, p int, threshold int) bool {
	return circlePixel-p > threshold
}

/**
 * Checks if the circle pixel is darker than the candidate pixel p by
 * a threshold.
 */
func isDarker(circlePixel int, p int, threshold int) bool {
	return p-circlePixel > threshold
}

/**
 * Gets the sixteen offset values of the circle surrounding pixel.
 */
func getCircleOffsets(width int) [16]int {
	var circle [16]int
	circle[0] = -width - width - width
	circle[1] = circle[0] + 1
	circle[2] = circle[1] + width + 1
	circle[3] = circle[2] + width + 1
	circle[4] = circle[3] + width
	circle[5] = circle[4] + width
	circle[6] = circle[5] + width - 1
	circle[7] = circle[6] + width - 1
	circle[8] = circle[7] - 1
	circle[9] = circle[8] - 1
	circle[10] = circle[9] - width - 1
	circle[11] = circle[10] - width - 1
	circle[12] = circle[11] - width
	circle[13] = circle[12] - width
	circle[14] = circle[13] - width + 1
	circle[15] = circle[14] - width + 1

	return circle
}

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
	gray := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	pixList := make([]int, w*h)
	for index := 0; index < w*h; index++ {
		pixList[index] = int(gray.Pix[index])
	}

	cornerList := findCorners(pixList, w, h, 20)
	newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			newImage.Set(x, y, grayColor)
		}
	}
	for i := 0; i < len(cornerList); i += 2 {
		newImage.Set(cornerList[i], cornerList[i+1], color.RGBA{255, 0, 0, 255})
	}

	// Encode the grayscale image to the output file
	outfile, err := os.Create("image_2.png")
	if err != nil {
		// replace this with real error handling
		panic(err)
	}

	defer outfile.Close()
	png.Encode(outfile, newImage)
	fmt.Println("done")
}
