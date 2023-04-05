package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
	"fmt"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/tevino/abool"
)

type FlagsEx struct {
	aic_package.Flags

	halfBlock bool
	mode      string
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

func initArgv(config *string, target *string) {
	flag.StringVar(config, "c", "config.json", "Config file path.")
	flag.StringVar(target, "g", "mumumu", "Load the gif setting in config file.")
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
		config string
		target string
	)
	initArgv(&config, &target)

	gifSetting := loadConfig(config)[target]
	flagsEx := readFlags(gifSetting)

	// Note: For environments where a terminal isn't available (such as web servers), you MUST
	// specify atleast one of flags.Width, flags.Height or flags.Dimensions

	// Conversion for an image

	gr := GifRenderer{
		filePath:      gifSetting.Path,
		flags: flagsEx,
		startTime:     time.Now(),
		message:       gifSetting.Message,
	}

	gr.loadGifToAscii()
	gr.renderGif(&ec)
}
