package main

import (
	"fmt"
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	imgManip "github.com/TheZoraiz/ascii-image-converter/image_manipulation"
	"image"
	"image/gif"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GifFrame struct {
	asciiCharSet [][]imgManip.AsciiChar
	delay        int
}

func flattenAscii(asciiSet [][]imgManip.AsciiChar, fontColor [3]int, colored, toSaveTxt bool) []string {

	var ascii []string

	for _, line := range asciiSet {
		var tempAscii string

		for _, char := range line {
			if toSaveTxt {
				tempAscii += char.Simple
				continue
			}

			if colored {
				tempAscii += char.OriginalColor
			} else if fontColor != [3]int{255, 255, 255} {
				tempAscii += char.SetColor
			} else {
				tempAscii += char.Simple
			}
		}

		ascii = append(ascii, tempAscii)
	}

	return ascii
}

func main() {
	// If file is in current directory. This can also be a URL to an image or gif.
	file_path := "./gif/bocchi-the-rock-bocchi-the-rock-gif.gif"

	flags := aic_package.DefaultFlags()

	flags.Braille = true
	flags.Colored = true
	// flags.CustomMap = " .-=+#@"

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	var (
		localFile  *os.File
		bochhi_gif *gif.GIF
	)

	localFile, err := os.Open(file_path)
	if err != nil {
		fmt.Errorf("Open gif file error: %v", err)
		return
	}
	defer localFile.Close()

	bochhi_gif, err = gif.DecodeAll(localFile)

	if err != nil {
		if file_path == "-" {
			fmt.Errorf("can't decode piped input: %v", err)
			return
		} else {
			fmt.Errorf("can't decode %v: %v", file_path, err)
			return
		}
	}

	var (
		asciiArtSet    = make([]string, len(bochhi_gif.Image))
		gifFramesSlice = make([]GifFrame, len(bochhi_gif.Image))

		counter             = 0
		concurrentProcesses = 0
		wg                  sync.WaitGroup
		hostCpuCount        = runtime.NumCPU()
	)

	fmt.Printf("Generating ascii art... 0%%\r")

	// Get first frame of gif and its dimensions
	firstGifFrame := bochhi_gif.Image[0].SubImage(bochhi_gif.Image[0].Rect)
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
	for i, frame := range bochhi_gif.Image {

		wg.Add(1)
		concurrentProcesses++

		go func(i int, frame *image.Paletted) {

			frameImage := frame.SubImage(frame.Rect)

			// If a frame is found that is smaller than the first frame, then this gif contains smaller subimages that are
			// positioned inside the original gif. This behavior isn't supported by this app
			if firstGifFrameWidth != frameImage.Bounds().Dx() || firstGifFrameHeight != frameImage.Bounds().Dy() {
				fmt.Printf("Error: " + file_path + " contains subimages smaller than default width and height\n\nProcess aborted because ascii-image-converter doesn't support subimage placement and transparency in GIFs\n\n")
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
			gifFramesSlice[i].delay = bochhi_gif.Delay[i]

			ascii := flattenAscii(asciiCharSet, fontColor, colored || grayscale, false)

			asciiArtSet[i] = strings.Join(ascii, "\n")

			counter++
			percentage := int((float64(counter) / float64(len(bochhi_gif.Image))) * 100)
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
	fmt.Printf("                              \r")

	fmt.Print("\033[H\033[2J")
	// Display the gif
	for {
		for i, asciiFrame := range asciiArtSet {
			// Clear screen: https://stackoverflow.com/a/22892171/12764484
			os.Stdout.Write([]byte(asciiFrame))

			// fmt.Println("OK......")
			time.Sleep(1)
			time.Sleep(time.Duration((time.Second * time.Duration(bochhi_gif.Delay[i])) / 100))
		}
	}

	//fmt.Printf("%v\n", asciiArt)
}
