# Mumumu
A terminal GIF renderer. Inspired by [nyancat-cli](https://github.com/klange/nyancat).

# Demo

| Action                                                                    | Demo                                                                       |
|---------------------------------------------------------------------------|----------------------------------------------------------------------------|
| Default (`./mumumu`)                                                      | ![demo!](https://drive.google.com/uc?id=1G9C6tMqoVM2oTVxnoa2VYvLyFzhak0g9) |
| Load customized config (e.g., `./mumumu -c config.json -g kita-kirakira`) | ![demo](https://drive.google.com/uc?id=1XT8orFf_f5IPHvw9VleEiPUoEvPP6MEY)  |
| Terminal resize                                                           | ![demo](https://drive.google.com/uc?id=1sR8pC2mD9stwcvSA1LnaYAK7Ztx04v05)  |

# Usage
``` bash
# build 
go run mumumu

# show help
./mumumu -h

Usage of ./mumumu:
  -c string
        Config file path. (default "config.json")
  -g string
        Load the gif setting in config file. (default "mumumu")

# Start mumumu (Equal to ./mumumu -c config.json -g mumumu)
./mumumu

# Render other gif (e.g., "kite-kirakira" in config.json)
./mumumu -c config.json -g kita-kirakira
```

# Configure
I predefined many gif settings in `config.json`. Ones can check it and give it a try. 

config.json
``` json
[
    {
        "name": "mumumu", // This gif setting name
        "path": "https://media.tenor.com/nIfKxqBUqQQAAAAC/shake-head-anime.gif", // Support url or file path
        "flags": {
            "color": true, // Color or gray scale
            "halfblock": false, // Render with ▀
            "braille": true, // Render with ⣿
            "threshold": 50, // Gray scale threashold 0-255
            "maxwidth": 100  // Max width (in char length) when rendering in terminal
        },
        "message": "You have mumumued for %d seconds..." // Message below the rendered git in terminal
    },
    ...
]
```

# Feature
- Supporting three rendering modes: ascii, braille, halfblock
- Center align when rendering gif
- Customizable message for rendering gif

# TODO
- config file 
    - gif infos array
        - repeat or not

# Troubleshooting
If you find the color rendered on terminal is not correct, check whether the true color mode is enabled.
```
# in .zshrc
export COLORTERM=truecolor

# in .tmux.conf
set -g default-terminal "xterm-256color"
```
You can check whether the true color is enabled successfully with [rich](https://github.com/Textualize/rich).

| True color enabled                                                         | True coloe not enabled                                                     |
|----------------------------------------------------------------------------|----------------------------------------------------------------------------|
| ![demo!](https://drive.google.com/uc?id=16OgQg7c0OBRrPFseKCK68x_fY-HT-TjV) | ![demo!](https://drive.google.com/uc?id=16Oa85bzUp5qCXPFzZrSWNsgkkncr8iwD) |

Some reference:
- [alacritty-tmux-vim_truecolor.md](https://gist.github.com/andersevenrud/015e61af2fd264371032763d4ed965b6)
- [rich](https://github.com/Textualize/rich)

# Reference
- [klange/nyancat](https://github.com/klange/nyancat)
- [TheZoraiz/ascii-image-converter](https://github.com/TheZoraiz/ascii-image-converter)
- [Textualize/rich](https://github.com/Textualize/rich)
- [viu](https://github.com/atanunq/viu)
