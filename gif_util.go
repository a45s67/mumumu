package main

import (
	"fmt"
	"image"
	"image/gif"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	imgManip "github.com/TheZoraiz/ascii-image-converter/image_manipulation"
)

type GifFrame struct {
	asciiCharSet [][]imgManip.AsciiChar
	delay        int
}

func loadGif(filePath string) *gif.GIF {
	var (
		fileStream *os.File
		bochhiGif  *gif.GIF
	)

	fileStream, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open gif file error: %v", err)
		os.Exit(1)
	}
	defer fileStream.Close()

	bochhiGif, err = gif.DecodeAll(fileStream)
	if err != nil {
		fmt.Printf("Decode gif file stream error: %v", err)
		os.Exit(1)
	}

	return bochhiGif
}

func flattenAsciiImages(gifFramesSlice []GifFrame, colored bool) []string {
	var asciiArtSet []string
	for _, gifFrame := range gifFramesSlice {
		ascii := flattenAscii(gifFrame.asciiCharSet, colored)
		asciiArtSet = append(asciiArtSet, strings.Join(ascii, "\n"))
	}
	return asciiArtSet
}

func gif2Ascii(bochhiGif *gif.GIF, flags aic_package.Flags) []GifFrame {
	var (
		err            error
		gifFramesSlice = make([]GifFrame, len(bochhiGif.Image))

		counter             = 0
		concurrentProcesses = 0
		wg                  sync.WaitGroup
		hostCpuCount        = runtime.NumCPU()
	)

	fmt.Printf("Generating ascii art... 0%%\r")

	// Get first frame of gif and its dimensions
	firstGifFrame := bochhiGif.Image[0].SubImage(bochhiGif.Image[0].Rect)
	firstGifFrameWidth := firstGifFrame.Bounds().Dx()
	firstGifFrameHeight := firstGifFrame.Bounds().Dy()

	var (
		dimensions = flags.Dimensions
		width      = flags.Width
		height     = flags.Height
		complex    = flags.Complex
		negative   = flags.Negative
		colored    = flags.Colored
		colorBg    = flags.CharBackgroundColor
		grayscale  = flags.Grayscale
		customMap  = flags.CustomMap
		flipX      = flags.FlipX
		flipY      = flags.FlipY
		full       = flags.Full
		fontColor  = flags.FontColor
		braille    = flags.Braille
		threshold  = flags.Threshold
		dither     = flags.Dither
	)

	// Multi-threaded loop to decrease execution time
	for i, frame := range bochhiGif.Image {

		wg.Add(1)
		concurrentProcesses++

		go func(i int, frame *image.Paletted) {

			frameImage := frame.SubImage(frame.Rect)

			// If a frame is found that is smaller than the first frame, then this gif contains smaller subimages that are
			// positioned inside the original gif. This behavior isn't supported by this app
			if firstGifFrameWidth != frameImage.Bounds().Dx() || firstGifFrameHeight != frameImage.Bounds().Dy() {
				fmt.Printf("Error: Gif contains subimages smaller than default width and height\n\nProcess aborted because ascii-image-converter doesn't support subimage placement and transparency in GIFs\n\n")
				os.Exit(0)
			}

			var imgSet [][]imgManip.AsciiPixel

			imgSet, err = imgManip.ConvertToAsciiPixels(frameImage, dimensions, width, height, flipX, flipY, full, braille, dither)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(0)
			}

			var asciiCharSet [][]imgManip.AsciiChar
			if flags.Braille {
				asciiCharSet, err = imgManip.ConvertToBrailleChars(imgSet, negative, colored, grayscale, colorBg, fontColor, threshold)
			} else {
				asciiCharSet, err = imgManip.ConvertToAsciiChars(imgSet, negative, colored, grayscale, complex, colorBg, customMap, fontColor)
			}
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(0)
			}

			gifFramesSlice[i].asciiCharSet = asciiCharSet
			gifFramesSlice[i].delay = bochhiGif.Delay[i]

			counter++
			percentage := int((float64(counter) / float64(len(bochhiGif.Image))) * 100)
			fmt.Printf("Generating ascii art... " + strconv.Itoa(percentage) + "%%\r")

			wg.Done()

		}(i, frame)

		// Limit concurrent processes according to host's CPU count to avoid overwhelming memory
		if concurrentProcesses == hostCpuCount {
			wg.Wait()
			concurrentProcesses = 0
		}
	}

	wg.Wait()

	return gifFramesSlice
}

func flattenAscii(asciiSet [][]imgManip.AsciiChar, colored bool) []string {

	var ascii []string

	for _, line := range asciiSet {
		var tempAscii string

		for _, char := range line {
			if colored {
				tempAscii += char.OriginalColor
			} else {
				tempAscii += char.Simple
			}
		}

		ascii = append(ascii, tempAscii)
	}

	return ascii
}
