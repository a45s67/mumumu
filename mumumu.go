package main

import (
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/tevino/abool"
	"time"
)

func main() {
	ec := EventCatcher{stop: new(abool.AtomicBool), windowChange: new(abool.AtomicBool)}
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

	startTime := time.Now()
	renderGif(asciiArtSet, gifFramesSlice, startTime, &ec)
}
