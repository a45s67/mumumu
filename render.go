package main

import (
	"fmt"
	"image/gif"
	"os"
	"strings"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package/winsize"
	imgManip "github.com/a45s67/ascii-image-converter/image_manipulation"
)

func hideCursor() {
	// Hide&Show cursor: https://stackoverflow.com/questions/30126490/how-to-hide-console-cursor-in-c
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func moveCursorToColumn(pos int) {
	fmt.Printf("\033[%dG", pos) // Move cursor to column
}

func clearScreen() {
	// Clear screen: https://stackoverflow.com/a/22892171/12764484
	fmt.Print("\033[H\033[2J")
}

func clearLine() {
	fmt.Printf("\033[2K") // Clear line
}

func getAsciiImageWidth(asciiCharSet [][]imgManip.AsciiChar) int {
	var asciis_of_top_line []string
	for _, asciichar := range asciiCharSet[0] {
		asciis_of_top_line = append(asciis_of_top_line, asciichar.Simple)
	}
	// Convert to rune to count the utf8 string length.
	// https://stackoverflow.com/questions/12668681/how-to-get-the-number-of-characters-in-a-string
	return len([]rune(strings.Join(asciis_of_top_line, "")))
}

type GifRenderer struct {
	cache     string
	flags     FlagsEx
	startTime time.Time
	message   string

	decodedGifData *gif.GIF
	gifFramesSlice []GifFrame
	asciiArtSet    []string

	terminalSize [2]int
}

func (gr *GifRenderer) loadGifToAscii() {
	gr.gifFramesSlice = gif2Ascii(gr.decodedGifData, gr.flags)
	gr.asciiArtSet = flattenAsciiImages(gr.gifFramesSlice,
		gr.flags.Colored || gr.flags.Grayscale)
}

func (gr *GifRenderer) reload() {
	gr.gifFramesSlice = gif2Ascii(gr.decodedGifData, gr.flags)
	gr.asciiArtSet = flattenAsciiImages(gr.gifFramesSlice,
		gr.flags.Colored || gr.flags.Grayscale)
}

func (gr *GifRenderer) renderGif(e *EventCatcher) {
	imageWidth := getAsciiImageWidth(gr.gifFramesSlice[0].asciiCharSet)
	imageHeight := len(gr.gifFramesSlice[0].asciiCharSet)
	gr.terminalSize[0], gr.terminalSize[1], _ = winsize.GetTerminalSize()
	hideCursor()
	clearScreen()
	defer showCursor()
	defer clearScreen()
	// Display the gif
	for {
		for i, asciiFrame := range gr.asciiArtSet[0:len(gr.asciiArtSet)] {
			// TODO: Move action checking below into GifRenderer method
			if e.stopEvent.IsSet() {
				return
			}
			if e.windowChangeEvent.IsSet() {
				gr.reload()
				clearScreen()
				imageWidth = getAsciiImageWidth(gr.gifFramesSlice[0].asciiCharSet)
				imageHeight = len(gr.gifFramesSlice[0].asciiCharSet)
				gr.terminalSize[0], gr.terminalSize[1], _ = winsize.GetTerminalSize()
				e.windowChangeEvent.UnSet()
				break
			}

			gr.renderImage(asciiFrame, imageWidth, imageHeight)
			gr.renderMessage(imageWidth)
			time.Sleep(time.Duration(
				(time.Second * time.Duration(gr.gifFramesSlice[i].delay)) / 100))
		}
	}
}

func (gr *GifRenderer) renderImage(asciiFrame string, imageWidth int, imageHeight int) {
	left := (gr.terminalSize[0]-imageWidth)/2 + 1
	top := (gr.terminalSize[1]-imageHeight)/2 + 1

	cursorTopLeftPos := fmt.Sprintf("\033[%d;%dH", top, left)
	cursorLeftPos := fmt.Sprintf("\033[%dG", left)

	fmt.Print(cursorTopLeftPos) // Move cursor to pos (1,1): https://en.wikipedia.org/wiki/ANSI_escape_code
	asciiFrame = strings.Replace(asciiFrame, "\n", "\n"+cursorLeftPos, -1)

	os.Stdout.Write([]byte(asciiFrame))
}

func (gr *GifRenderer) renderMessage(imageWidth int) {
	if len(gr.message) == 0 {
		return
	}

	elapsed := time.Since(gr.startTime)
	msg := fmt.Sprintf(gr.message, int(elapsed.Seconds()))

	left := (gr.terminalSize[0]-imageWidth)/2 + 1

	msg_len := len(msg)
	msg_left_pos := (imageWidth-msg_len)/2 + left

	fmt.Print("\n")
	clearLine()
	moveCursorToColumn(msg_left_pos)
	fmt.Print(msg)
}
