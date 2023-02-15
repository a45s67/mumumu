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
	Name    string
	Path    string
	Flags   map[string]interface{}
	Message string
}

func loadConfig(configPath string) map[string]Option {
	configFile, err := os.Open(configPath)
	defer configFile.Close()
	if err != nil {
		panic(err.Error())
	}

	var gifSettings []Option
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&gifSettings)

	gifSettingMap := map[string]Option{}
	for _, config := range gifSettings {
		gifSettingMap[config.Name] = config
	}
	return gifSettingMap
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
		stopEvent:         new(abool.AtomicBool),
		windowChangeEvent: new(abool.AtomicBool),
	}
	ec.listenKeystroke()
	ec.listenSignal()

	gifSetting := loadConfig("config.json")["mumumu"]
	flagsEx := readFlags(gifSetting)

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	gr := GifRenderer{
		filePath:      gifSetting.Path,
		renderFlagsEx: flagsEx,
		startTime:     time.Now(),
		message:       gifSetting.Message,
	}

	gr.loadGifToAscii()
	gr.renderGif(&ec)
}
