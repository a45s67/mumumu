package main

import (
	"fmt"
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	imgManip "github.com/TheZoraiz/ascii-image-converter/image_manipulation"
	"os"
	"time"
)

func renderMessage(imageWidth int, startTime time.Time) {
	elapsed := time.Since(startTime)
	msg := fmt.Sprintf("You have mumumued for %d seconds", int(elapsed.Seconds()))
	msg_len := len(msg)

	msg_left_pos := (imageWidth - msg_len) / 2
	fmt.Printf("\033[2K")                // Clear line
	fmt.Printf("\033[%dG", msg_left_pos) // Move cursor to column
	fmt.Print(msg)
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
