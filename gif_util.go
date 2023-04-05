package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/TheZoraiz/ascii-image-converter/aic_package/winsize"
	imgManip "github.com/a45s67/ascii-image-converter/image_manipulation"
)

type GifFrame struct {
	asciiCharSet [][]imgManip.AsciiChar
	delay        int
}

func isURL(str string) bool {
	if len(str) < 8 {
		return false
	} else if str[:7] == "http://" || str[:8] == "https://" {
		return true
	}
	return false
}

func requestGIF(gifUrl string) []byte {
	fmt.Printf("Fetching file from url...\r")

	// Request the gif data
	retrievedImage, err := http.Get(gifUrl)
	if err != nil {
		panic(fmt.Errorf("can't fetch content: %v", err))
	}

	urlImgBytes, err := ioutil.ReadAll(retrievedImage.Body)
	if err != nil {
		panic(fmt.Errorf("failed to read fetched content: %v", err))
	}
	defer retrievedImage.Body.Close()
	return urlImgBytes

}

func loadGIFFromURL(gifUrl string) *gif.GIF {
	urlImgBytes := requestGIF(gifUrl)

	decodedGif, err := gif.DecodeAll(bytes.NewReader(urlImgBytes))
	if err != nil {
		fmt.Printf("Decode gif file stream error: %v", err)
		os.Exit(1)
	}
	return decodedGif
}

func downloadGIF(gifUrl string, cache string) {
	urlImgBytes := requestGIF(gifUrl)

	// Write the gif raw data to cache file
	if err := os.MkdirAll(filepath.Dir(cache), 0770); err != nil {
		panic(err)
	}
	err := os.WriteFile(cache, urlImgBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func loadGif(filePath string) *gif.GIF {
	var (
		fileStream *os.File
		decodedGif *gif.GIF
	)

	fileStream, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open gif file error: %v", err)
		os.Exit(1)
	}
	defer fileStream.Close()

	decodedGif, err = gif.DecodeAll(fileStream)
	if err != nil {
		fmt.Printf("Decode gif file stream error: %v", err)
		os.Exit(1)
	}

	return decodedGif
}

func getIdealRenderSize(image_size image.Rectangle, widthLimit int) []int {
	imgWidth := float64(image_size.Dx())
	imgHeight := float64(image_size.Dy())
	aspectRatio := imgWidth / imgHeight

	tWidth, tHeight, _ := winsize.GetTerminalSize()

	if widthLimit == 0 {
		widthLimit = tWidth
	}
	idealWidth := math.Min(float64(tWidth), float64(widthLimit))
	idealHeight := 2*float64(tHeight) - 1

	if float64(idealWidth)/aspectRatio > float64(idealHeight) {
		idealWidth = idealHeight * aspectRatio
	} else {
		idealHeight = idealWidth / aspectRatio
	}
	return []int{int(idealWidth), int(idealHeight)}
}

func flattenAsciiImages(gifFramesSlice []GifFrame, colored bool) []string {
	var asciiArtSet []string
	for _, gifFrame := range gifFramesSlice {
		ascii := flattenAscii(gifFrame.asciiCharSet, colored)
		asciiArtSet = append(asciiArtSet, strings.Join(ascii, "\n"))
	}
	return asciiArtSet
}

func gif2Ascii(gifData *gif.GIF, flags FlagsEx) []GifFrame {
	var (
		err                 error
		gifFramesSlice      = make([]GifFrame, len(gifData.Image))
		counter             = 0
		concurrentProcesses = 0
		wg                  sync.WaitGroup
		hostCpuCount        = runtime.NumCPU()
	)

	fmt.Printf("Generating ascii art... 0%%\r")

	var (
		dimensions = flags.Dimensions
		width      = flags.Width
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
	for i, frame := range gifData.Image {

		wg.Add(1)
		concurrentProcesses++

		func(i int, frame *image.Paletted) {
			var imgSet [][]imgManip.AsciiPixel

			frameImage := frame.SubImage(frame.Rect)
			dimensions = getIdealRenderSize(frameImage.Bounds(), width)
			if !flags.halfBlock {
				dimensions[1] /= 2
			}

			imgSet, err = imgManip.ConvertToAsciiPixels(frameImage, dimensions, 0, 0, flipX, flipY, full, braille, dither)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(0)
			}

			var asciiCharSet [][]imgManip.AsciiChar
			if flags.halfBlock {
				asciiCharSet, err = imgManip.ConvertToHalfBlockChars(imgSet, negative, colored, grayscale)
			} else if flags.Braille {
				asciiCharSet, err = imgManip.ConvertToBrailleChars(imgSet, negative, colored, grayscale, colorBg, fontColor, threshold)
			} else {
				asciiCharSet, err = imgManip.ConvertToAsciiChars(imgSet, negative, colored, grayscale, complex, colorBg, customMap, fontColor)
			}
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(0)
			}

			gifFramesSlice[i].asciiCharSet = asciiCharSet
			gifFramesSlice[i].delay = gifData.Delay[i]

			counter++
			percentage := int((float64(counter) / float64(len(gifData.Image))) * 100)
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

	var AsciiImageString []string

	for _, line := range asciiSet {
		var coloredLine string

		for _, char := range line {
			if colored {
				coloredLine += char.OriginalColor
			} else {
				coloredLine += char.Simple
			}
		}

		AsciiImageString = append(AsciiImageString, coloredLine)
	}

	return AsciiImageString
}
