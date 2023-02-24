package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"mccoy.space/g/ogg"
)

const (
	samplingRate    = 48000
	channels        = 2
	frameLength     = 20
	sampleSize      = 2 * channels
	samplesPerFrame = int(samplingRate / 960 * frameLength)
	frameSize       = samplesPerFrame * sampleSize
)

type YoutubeDownloader struct {
	url    string
	output io.Writer
}

func NewYoutubeDownloader(url string, output io.Writer) *YoutubeDownloader {
	return &YoutubeDownloader{
		url:    url,
		output: output,
	}
}

// Download youtube-dl -o - BaW_jenozKc
// bestaudio | opus?
//
//	ytdl_format_options = {
//	   'format': 'bestaudio/best',
//	   'outtmpl': '%(extractor)s-%(id)s-%(title)s.%(ext)s',
//	   'restrictfilenames': True,
//	   'noplaylist': True,
//	   'nocheckcertificate': True,
//	   'ignoreerrors': False,
//	   'logtostderr': False,
//	   'quiet': True,
//	   'no_warnings': True,
//	   'default_search': 'auto',
//	   'source_address': '0.0.0.0'  # bind to ipv4 since ipv6 addresses cause issues sometimes
//	}
func (d *YoutubeDownloader) Download() error {
	cmd := exec.Command("youtube-dl", "-g", d.url)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("cmd run: %w", err)
	}

	links := bytes.Split(out, []byte{'\n'})

	cmd = exec.Command("ffmpeg",
		"-i", string(links[1]),
		"-map_metadata",
		"-1",
		"-f",
		"opus",
		"-c:a",
		"libopus",
		"-ar",
		"48000",
		"-ac",
		"2",
		"-b:a",
		"128k",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg start: %w", err)
	}

	defer cmd.Wait()

	dec := ogg.NewDecoder(stdout)
	for {
		p, err := dec.Decode()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("ogg decoder: %w", err)
		}

		for _, packet := range p.Packets {
			d.output.Write(packet)
		}
	}
}
