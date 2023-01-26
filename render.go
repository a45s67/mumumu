package main

import (
	"fmt"
	"os"
	"time"
)

func hideCursor() {
	// Hide cursor: https://stackoverflow.com/questions/30126490/how-to-hide-console-cursor-in-c
	fmt.Print("\033[?25l")
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

func renderGif(asciiArtSet []string, gifFramesSlice []GifFrame, startTime time.Time) {
	imageWidth := len(gifFramesSlice[0].asciiCharSet[0])
	hideCursor()
	clearScreen()
	// Display the gif
	for {
		for i, asciiFrame := range asciiArtSet[0 : len(asciiArtSet)-1] {
			renderImage(asciiFrame)
			renderMessage(imageWidth, startTime)
			time.Sleep(time.Duration((time.Second * time.Duration(gifFramesSlice[i].delay)) / 100))
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
