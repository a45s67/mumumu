package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/tevino/abool"
)

type FlagsEx struct {
	aic_package.Flags

	halfBlock bool
	mode      string
}

type Option struct {
	Name    string                 `json:"name"`
	URL     string                 `json:"url"`
	Cache   string                 `json:"cache"`
	Flags   map[string]interface{} `json:"flags"`
	Message string                 `json:"message"`
}

func fileExist(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func loadConfig(configPath string) map[string]Option {
	fo, err := os.Open(configPath)
	defer fo.Close()
	if err != nil {
		panic(err.Error())
	}

	var gifSettings []Option
	jsonParser := json.NewDecoder(fo)
	err = jsonParser.Decode(&gifSettings)
	if err != nil {
		panic(err.Error())
	}

	gifSettingMap := map[string]Option{}
	for _, config := range gifSettings {
		gifSettingMap[config.Name] = config
	}
	return gifSettingMap
}

func readFlags(gifOption Option) FlagsEx {
	flags := FlagsEx{
		aic_package.DefaultFlags(),
		false,
		"ascii",
	}
	if val, ok := gifOption.Flags["color"]; ok {
		flags.Colored = val.(bool)
	}
	if val, ok := gifOption.Flags["mode"]; ok {
		flags.halfBlock = false
		flags.Braille = false

		mode := val.(string)
		switch mode {
		case "braille":
			flags.Braille = true
		case "halfblock":
			flags.halfBlock = true
		case "ascii":
		default:
			panic(fmt.Errorf("Error: Unknown mode %s", mode))
		}
	}
	if val, ok := gifOption.Flags["threshold"]; ok {
		flags.Threshold = int(val.(float64))
	}
	if val, ok := gifOption.Flags["maxwidth"]; ok {
		flags.Width = int(val.(float64))
	}
	return flags
}

func initArgv(config *string, target *string, nocache *bool) {
	flag.StringVar(config, "c", "config.json", "Config file path.")
	flag.StringVar(target, "g", "mumumu", "Load the gif setting in config file.")
	flag.BoolVar(nocache, "n", false, "Ignore the gif cache stored in local and download it again.")
	flag.Parse()
}

func main() {
	ec := EventCatcher{
		stopEvent:         new(abool.AtomicBool),
		windowChangeEvent: new(abool.AtomicBool),
	}
	ec.listenKeystroke()
	ec.listenSignal()

	var (
		config  string
		target  string
		nocache bool
	)
	initArgv(&config, &target, &nocache)

	gifOption := loadConfig(config)[target]
	flagsEx := readFlags(gifOption)

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	gr := GifRenderer{
		flags:     flagsEx,
		startTime: time.Now(),
		message:   gifOption.Message,
	}

	if gifOption.Cache == "" {
		gr.decodedGifData = loadGIFFromURL(gifOption.URL)
	} else if nocache || !fileExist(gifOption.Cache) {
		downloadGIF(gifOption.URL, gifOption.Cache)
		gr.decodedGifData = loadGif(gifOption.Cache)
	} else if fileExist(gifOption.Cache) {
		gr.decodedGifData = loadGif(gifOption.Cache)
	} else {
		panic(fmt.Errorf("Invalid setting for url '%s', cache '%s' in option '%s'",
			gifOption.URL, gifOption.Cache, gifOption.Name))
	}

	gr.loadGifToAscii()
	gr.renderGif(&ec)
}
