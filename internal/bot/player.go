package bot

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

var (
	ErrNoURLProvided = errors.New("no url provided")
)

func (b *Bot) playCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	sc := strings.Split(m.Content, " ")
	if len(sc) < 2 {
		return ErrNoURLProvided
	}

	youtubeURL, err := url.Parse(sc[1])
	if err != nil {
		return fmt.Errorf("url parse: %w", err)
	}

	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		return fmt.Errorf("state guild: %w", err)
	}

	var channelID string
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			if vs.ChannelID == "" {
				return fmt.Errorf("you need to be in a voice channel")
			}

			channelID = vs.ChannelID
			break
		}
	}

	vc, err := s.ChannelVoiceJoin(g.ID, channelID, false, true)
	if err != nil {
		return fmt.Errorf("channel voice join: %w", err)
	}

	defer vc.Disconnect()

	v, err := b.client.GetVideo(youtubeURL.String())
	if err != nil {
		return fmt.Errorf("client get video: %w", err)
	}

	audioFormats := v.Formats.WithAudioChannels().Select(func(f youtube.Format) bool {
		return f.AudioQuality == "AUDIO_QUALITY_MEDIUM" && f.AudioSampleRate == "48000"
	})

	audioStream, _, err := b.client.GetStream(v, &audioFormats[0])
	if err != nil {
		return fmt.Errorf("client get stream: %w", err)
	}

	defer audioStream.Close()

	if err = vc.Speaking(true); err != nil {
		return fmt.Errorf("failed initiate speaking: %w", err)
	}

	defer vc.Speaking(false)

	enc, err := dca.EncodeMem(audioStream, dca.StdEncodeOptions)
	if err != nil {
		return fmt.Errorf("dca encode mem: %w", err)
	}
	defer enc.Cleanup()

	done := make(chan error)
	dca.NewStream(enc, vc, done)
	if err = <-done; err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("stream encoding: %w", err)
	}

	return nil
}
