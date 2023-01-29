package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/tevino/abool"
)

type FlagsEx struct {
	halfBlock bool
	flags     aic_package.Flags
}

type Option struct {
	Name  string
	Path  string
	Flags map[string]interface{}
}

func loadConfig(configPath string) map[string]Option {
	var configArray []Option
	configFile, err := os.Open(configPath)
	defer configFile.Close()
	if err != nil {
		panic(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configArray)

	configSet := map[string]Option{}
	for _, config := range configArray {
		configSet[config.Name] = config
	}
	return configSet
}

func readFlags(gifOption Option) FlagsEx {
	flagsEx := FlagsEx{
		halfBlock: false,
	}
	flags := aic_package.DefaultFlags()
	if val, ok := gifOption.Flags["color"]; ok {
		flags.Colored = val.(bool)
	}
	if val, ok := gifOption.Flags["halfblock"]; ok {
		flagsEx.halfBlock = val.(bool)
	}
	if val, ok := gifOption.Flags["braille"]; ok {
		flags.Braille = val.(bool)
	}
	if val, ok := gifOption.Flags["threshold"]; ok {
		flags.Threshold = int(val.(float64))
	}
	if val, ok := gifOption.Flags["maxwidth"]; ok {
		flags.Width = int(val.(float64))
	}
    flagsEx.flags = flags
	return flagsEx
}

func main() {
	ec := EventCatcher{
		stop:         new(abool.AtomicBool),
		windowChange: new(abool.AtomicBool),
	}
	ec.listenKeystroke()
	ec.listenSignal()

	targetGifOption := loadConfig("config.json")["mumumu"]
	flagsEx := readFlags(targetGifOption)

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	gr := GifRenderer{
		filePath:      targetGifOption.Path,
		renderFlagsEx: flagsEx,
		startTime:     time.Now(),
	}

	gr.loadGifToAscii()
	gr.renderGif(&ec)
}
