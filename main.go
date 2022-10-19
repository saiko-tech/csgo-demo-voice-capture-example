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

int raw2celt(char* in, char* out) {
	unsigned char buf[BUF_SIZE];

	const unsigned int FRAME_SIZE = 512;
	const unsigned int SAMPLE_RATE = 22050;

	FILE *f = fopen(in, "r");
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

	FILE* file_p = fopen(out, "w");
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
	"path/filepath"
	"strconv"
	"strings"

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
			int(msg.SVC_Messages_svc_VoiceInit): func() proto.Message { return &msg.CSVCMsg_VoiceInit{} },
			int(msg.SVC_Messages_svc_VoiceData): func() proto.Message { return &msg.CSVCMsg_VoiceData{} },
		},
		IgnoreErrBombsiteIndexNotFound: false,
		NetMessageDecryptionKey:        nil,
	})
	defer p.Close()

	outFilesBySectionNr := make(map[uint32]*os.File)

	defer func() {
		for _, f := range outFilesBySectionNr {
			f.Close()
		}
	}()

	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_VoiceInit) {
		fmt.Println("init:", msg.String())
	})

	var sections []string

	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_VoiceData) {
		sectionNr := msg.GetSectionNumber()

		outFile, ok := outFilesBySectionNr[sectionNr]
		if !ok {
			sections = append(sections, strconv.Itoa(int(sectionNr)))
			outFile, err = os.Create(fmt.Sprintf("out/%d.raw", sectionNr))
			checkError(err)

			outFilesBySectionNr[sectionNr] = outFile
		}

		_, err := outFile.Write(msg.GetVoiceData())
		checkError(err)

		msg.VoiceData = nil
	})

	err = p.ParseToEnd()
	checkError(err)

	for _, rawFile := range outFilesBySectionNr {
		rawFileName := rawFile.Name()
		celtFileName := strings.TrimSuffix(rawFileName, filepath.Ext(rawFileName)) + ".celt"

		err = rawFile.Close()
		checkError(err)

		res := C.raw2celt(C.CString(rawFileName), C.CString(celtFileName))
		if res != 0 {
			panic("raw2celt failed")
		}
	}

	fmt.Println()

	fmt.Printf("order:\n\t%s\n", strings.Join(sections, "\n\t"))

	fmt.Println("saved voice chat audio to out/*.celt")
	fmt.Println("play via: play -t raw -r 22050 -e signed -b 16 -c 1 out/1.celt")
	fmt.Println("or convert to .wav via: sox -t raw -r 22050 -e signed -b 16 -c 1 out/1.celt out/1.wav")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
