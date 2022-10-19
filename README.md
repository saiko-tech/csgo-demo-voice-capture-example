# Capturing CS:GO Demo Voice Data

This is an example of how to use the [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang) library to capture voice chat in CELT format.


## Join us on Discord

[![Discord Chat](https://img.shields.io/discord/901824796302643281?color=%235865F2&label=discord&style=for-the-badge)](https://discord.gg/eTVBgKeHnh)


## Prerequisites

You need to have the following installed:

- Linux (macOS or WSL may work, but are not tested)
- CS:GO Linux Binaries
- CELT - Audio Codec Library
- Sox - Sound Processing Tools (for playback and conversion to `.wav`)

## Running the example

Adjust the paths in the below example before running.

```terminal
STEAM_LIBRARY="..." # <--- insert path to steam library here
CSGO_BIN="$STEAM_LIBRARY/steamapps/common/Counter-Strike Global Offensive/bin/linux64"
export CGO_LDFLAGS="-L \"$CSGO_BIN\" -l:vaudio_celt_client.so"
export LD_LIBRARY_PATH="$CSGO_BIN:$LD_LIBRARY_PATH"

go run capture_voice.go -demo /path/to/demo.dem # <--- replace with your demo
```

This will create a different files in an `out/` directory. Each file is a separate sequence of voice audio.

With Sox installed, you can play these files back via the following command (replace 1 with the sequence to play back):

    play -t raw -r 22050 -e signed -b 16 -c 1 out/1.celt

Or convert it to `.wav` via:

    sox -t raw -r 22050 -e signed -b 16 -c 1 -L out/1.celt out/1.wav

## Acknowledgements

Thanks to [@ericek111](https://github.com/ericek111) for [this gist](https://gist.github.com/ericek111/abe5829f6e52e4b25b3b97a0efd0b22b)
