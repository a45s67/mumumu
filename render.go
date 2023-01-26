package main

import (
	"fmt"
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"os"
	"time"
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

type GifRenderer struct {
	filePath       string
	renderFlags    aic_package.Flags
	startTime      time.Time
	gifFramesSlice []GifFrame
	asciiArtSet    []string
}

func (gr *GifRenderer) loadGifToAscii() {
	bochhiGif := loadGif(gr.filePath)
	gr.gifFramesSlice = gif2Ascii(bochhiGif, gr.renderFlags)
	gr.asciiArtSet = flattenAsciiImages(gr.gifFramesSlice,
		gr.renderFlags.Colored || gr.renderFlags.Grayscale)
}

func (gr *GifRenderer) reload() {
	gr.loadGifToAscii()
}

func (gr *GifRenderer) renderGif(e *EventCatcher) {
	imageWidth := len(gr.gifFramesSlice[0].asciiCharSet[0])
	hideCursor()
	clearScreen()
	defer showCursor()
	defer clearScreen()
	// Display the gif
	for {
		for i, asciiFrame := range gr.asciiArtSet[0 : len(gr.asciiArtSet)-1] {
			// TODO: Move action checking below into GifRenderer method
			if e.stop.IsSet() {
				return
			}
			if e.windowChange.IsSet() {
				gr.reload()
				e.windowChange.UnSet()
			}

			renderImage(asciiFrame)
			renderMessage(imageWidth, gr.startTime)
			time.Sleep(time.Duration(
				(time.Second * time.Duration(gr.gifFramesSlice[i].delay)) / 100))
		}
	}
}

func renderImage(asciiFrame string) {
	fmt.Print("\033[1;1H") // Move cursor to pos (1,1): https://en.wikipedia.org/wiki/ANSI_escape_code
	os.Stdout.Write([]byte(asciiFrame))
}

func renderMessage(imageWidth int, startTime time.Time) {
	elapsed := time.Since(startTime)
	msg := fmt.Sprintf("You have mumumued for %d seconds", int(elapsed.Seconds()))

	msg_len := len(msg)
	msg_left_pos := (imageWidth - msg_len) / 2

	clearLine()
	moveCursorToColumn(msg_left_pos)
	fmt.Print(msg)
}
