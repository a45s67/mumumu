package main

import (
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/tevino/abool"
	"time"
)

func main() {
	ec := EventCatcher{
		stop:         new(abool.AtomicBool),
		windowChange: new(abool.AtomicBool),
	}
	ec.listenKeystroke()
	ec.listenSignal()

	// If file is in current directory. This can also be a URL to an image or gif.
    // filePath := "https://media.tenor.com/nIfKxqBUqQQAAAAC/shake-head-anime.gif"
    // filePath := "https://media.tenor.com/nQbnkbBw3EUAAAAC/bocchi-bocchi-the.gif"
    // filePath := "https://media.tenor.com/JFLRMx4dpScAAAAC/bocchi-the-rock-bocchi-the-rock-gif.gif"
    // filePath := "https://media.tenor.com/nB9YSFNoDi8AAAAC/bocchi-the-rock-chainsaw-man.gif"
    // filePath := "https://media.tenor.com/Kif5V8DkNgYAAAAC/bocchi-the-rock-bocchi.gif"
    // filePath := "https://media1.tenor.com/images/77a28a3d1343fcc2482fa835091418b3/tenor.gif?itemid=27477401"
    // filePath := "https://media.tenor.com/rSryJTbRc4YAAAAC/bocchi-the-rock-kita-ikuyo.gif"
    // filePath := "https://media1.tenor.com/images/2ddebc5c19dda8a602591782ab19ffcb/tenor.gif?itemid=27477556"
    // filePath := "https://media.tenor.com/526-Kgq17kMAAAAC/bocchi-the-rock-bocchi.gif"
    // filePath := "https://i0.hdslb.com/bfs/article/b8c7c7f6a34280e5f6cf2c75b422d949a52bd199.gif"
    filePath := "https://media1.tenor.com/images/ef715b505f47d21418b497f4c1d7fbb0/tenor.gif?itemid=27478087"
    // filePath := "https://i0.hdslb.com/bfs/article/43ab8d5f2a35d29d5a7ef9246c7dc0c7b7495a57.gif"
    // filePath := "https://tenor.com/zh-TW/view/bocchi-the-rock-chainsaw-man-hitori-goto-gif-27061336"
    // filePath := "https://media.tenor.com/VyugKLEBolsAAAAC/bocchi-bocchi-the-rock.gif"
    // filePath := "https://media.tenor.com/nQbnkbBw3EUAAAAC/bocchi-bocchi-the.gif"
	// filePath := "./gif/bochhi.gif"

	flags := aic_package.DefaultFlags()

	flags.Braille = false
	flags.Threshold = 30
	flags.Colored = true
	flags.CustomMap = ""

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	gr := GifRenderer{
		filePath:    filePath,
		renderFlags: flags,
		startTime:   time.Now(),
	}

	gr.loadGifToAscii()
	gr.renderGif(&ec)
}
