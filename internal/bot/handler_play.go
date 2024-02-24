package bot

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

func (b *Bot) handlePlay(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	youtubeURL, err := url.Parse(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		return fmt.Errorf("url parse: %w", err)
	}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		return fmt.Errorf("state guild: %w", err)
	}

	var channelID string
	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Interaction.Member.User.ID {
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

	if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Playing song " + v.Title,
		},
	}); err != nil {
		return fmt.Errorf("interaction respond: %w", err)
	}

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
