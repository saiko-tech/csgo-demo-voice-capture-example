# Capturing CS:GO Demo Voice Data

This is an example of how to use the [demoinfocs-golang](github.com/markus-wa/demoinfocs-golang) library to capture voice chat in CELT format.

## Prerequisites

You need to have the following installed:

- Linux (macOS or WSL may work, but are not tested)
- CS:GO Linux Binaries
- CELT - Audio Codec Library
- Sox - Sound Processing Tools (for conversion to `.wav`)

## Running the example

Adjust the paths in the below example before running.

```terminal
STEAM_LIBRARY="..." # <--- insert path to steam library here
CSGO_BIN="$STEAM_LIBRARY/steamapps/common/Counter-Strike Global Offensive/bin/linux64"
export CGO_LDFLAGS="-L \"$CSGO_BIN\" -l:vaudio_celt_client.so"
export LD_LIBRARY_PATH="$CSGO_BIN:$LD_LIBRARY_PATH"

go run capture_voice.go -demo /path/to/demo.dem # <--- replace with your demo
```

This will create a file called `out.celt`.
You can play this file via:

    play -t raw -r 22050 -e signed -b 16 -c 1 out.celt

Or convert it to `.wav` via:

    sox -t raw -r 22050 -e signed -b 16 -c 1 -L out.celt out.wav

## Acknowledgements

Thanks to [@ericek111](https://github.com/ericek111) for [this gist](https://gist.github.com/ericek111/abe5829f6e52e4b25b3b97a0efd0b22b)
