package main

// C code courtesy of https://gist.github.com/ericek111/abe5829f6e52e4b25b3b97a0efd0b22b

/*
#include <unistd.h>
#include <stdlib.h>
#include <math.h>
#include <stdio.h>
#include <stdint.h>
#include <sys/time.h>

#include <celt/celt.h>

#define BUF_SIZE 1024*1024

int celt2wav() {
	unsigned char buf[BUF_SIZE];

	const unsigned int FRAME_SIZE = 512;
	const unsigned int SAMPLE_RATE = 22050;

	FILE *f = fopen("voicedata.dat", "r");
	if (f == NULL) {
		  return 1;
	}

	size_t read = fread(buf, 1, BUF_SIZE, f);
	fclose(f);

	CELTMode *dm = celt_mode_create(SAMPLE_RATE, FRAME_SIZE, NULL);
	CELTDecoder *d = celt_decoder_create_custom(dm, 1, NULL);

	size_t outsize = (read / 64) * FRAME_SIZE * sizeof(celt_int16);
	celt_int16* pcmout = malloc(outsize);

	size_t done = 0;
	int frames = 0;

	for (; done < read; done += 64, frames++) {
		int ret = 0;
		if ((ret = celt_decode(d, buf + done, 64, pcmout + frames * FRAME_SIZE, FRAME_SIZE)) < 0) {
			fprintf(stderr, "unable to decode... > %d (at %d/%d)\n", ret, done, read);
			return 1;
		}
	}

	FILE* file_p = fopen("out.celt", "w");
	size_t written = fwrite(pcmout, outsize, 1, file_p);
	fclose(file_p);

	free(pcmout);
	return 0;
}
*/
import "C"

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/proto"

	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	demoinfocs "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
)

// Run like this: go run capture_voice.go -demo /path/to/demo.dem
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	checkError(err)

	defer f.Close()

	p := demoinfocs.NewParserWithConfig(f, demoinfocs.ParserConfig{
		MsgQueueBufferSize: -1,
		AdditionalNetMessageCreators: map[int]demoinfocs.NetMessageCreator{
			int(msg.SVC_Messages_svc_VoiceData): func() proto.Message { return &msg.CSVCMsg_VoiceData{} },
		},
		IgnoreErrBombsiteIndexNotFound: false,
		NetMessageDecryptionKey:        nil,
	})
	defer p.Close()

	outF, err := os.Create("voicedata.dat")
	checkError(err)

	defer outF.Close()

	var lastXUID uint64 = 0

	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_VoiceData) {
		voicePlayerXUID := msg.GetXuid()

		players := p.GameState().Participants().All()

		if voicePlayerXUID != lastXUID {
			for _, player := range players {
				if player.SteamID64 == voicePlayerXUID {
					lastXUID = voicePlayerXUID
					fmt.Printf("%s: %d:%d - player %q is speaking\n", p.CurrentTime().Round(time.Second), p.GameState().TeamCounterTerrorists().Score(), p.GameState().TeamTerrorists().Score(), player.Name)
				}
			}
		}

		_, err := outF.Write(msg.GetVoiceData())
		checkError(err)
	})

	err = p.ParseToEnd()
	checkError(err)

	res := C.celt2wav()
	if res != 0 {
		panic("celt2wav failed")
	}

	fmt.Println()

	fmt.Println("saved voice chat audio to out.celt")
	fmt.Println("play via: play -t raw -r 22050 -e signed -b 16 -c 1 out.celt")
	fmt.Println("or convert to .wav via: sox -t raw -r 22050 -e signed -b 16 -c 1 out.celt out.wav")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
