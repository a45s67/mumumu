package main

import (
	"bufio"
	"fmt"
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	imgManip "github.com/TheZoraiz/ascii-image-converter/image_manipulation"
	"image"
	"image/gif"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type GifFrame struct {
	asciiCharSet [][]imgManip.AsciiChar
	delay        int
}

func renderMessage(imageWidth int, startTime time.Time) {
	elapsed := time.Since(startTime)
	msg := fmt.Sprintf("You have mumumued for %d seconds", int(elapsed.Seconds()))
	msg_len := len(msg)

	msg_left_pos := (imageWidth - msg_len) / 2
	fmt.Printf("\033[2K")                // Clear line
	fmt.Printf("\033[%dG", msg_left_pos) // Move cursor to column
	fmt.Print(msg)
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

// https://opensourcedoc.com/golang-programming/class-object/
type EventCatcher struct {
	windowChange bool
	stop         bool
}

func (e *EventCatcher) listenEnter() {
	// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
	ch := make(chan string)
	go func(ch chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, _ := reader.ReadString('\n')
			ch <- s
		}
	}(ch)
	<-ch
	e.stop = true
}

func (e *EventCatcher) listenSignal() {
	// https://stackoverflow.com/questions/18106749/golang-catch-signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGWINCH,
	)
	go func() {
		sig := <-sigc
		switch sig {
		case syscall.SIGINT:
			e.stop = true
            os.Exit(1)
		case syscall.SIGWINCH:
			e.windowChange = true
		}
	}()
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

func main() {
	ec := EventCatcher{stop: false, windowChange: false}
	//ec.listenEnter()
	ec.listenSignal()

	// If file is in current directory. This can also be a URL to an image or gif.
	filePath := "./gif/bocchi-the-rock-bocchi-the-rock-gif.gif"

	flags := aic_package.DefaultFlags()

	flags.Braille = true
	flags.Colored = true
	// flags.CustomMap = " .-=+#@"

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	bochhiGif := loadGif(filePath)
	gifFramesSlice := gif2Ascii(bochhiGif, flags)
	asciiArtSet := flattenAsciiImages(gifFramesSlice, flags.Colored || flags.Grayscale)

	imageWidth := len(gifFramesSlice[0].asciiCharSet[0])

	startTime := time.Now()

	// Hide cursor: https://stackoverflow.com/questions/30126490/how-to-hide-console-cursor-in-c
	// Clear screen: https://stackoverflow.com/a/22892171/12764484
	fmt.Print("\033[?25l")
	fmt.Print("\033[H\033[2J")
	// Display the gif
	for {
		for i, asciiFrame := range asciiArtSet[0 : len(asciiArtSet)-1] {
			fmt.Print("\033[1;1H") // Move cursor to pos (1,1): https://en.wikipedia.org/wiki/ANSI_escape_code
			os.Stdout.Write([]byte(asciiFrame))

			renderMessage(imageWidth, startTime)
			time.Sleep(time.Duration((time.Second * time.Duration(gifFramesSlice[i].delay)) / 100))
		}
	}
}
