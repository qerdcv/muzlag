package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"

	"github.com/qerdcv/muzlag.go/internal/bot/player"
	"github.com/qerdcv/muzlag.go/internal/logger"
)

type eventSource struct {
	voiceChannelID string
	source         chan string
}

type PlayerHandler struct {
	client *youtube.Client
	logger logger.Logger[*slog.Logger]

	player            *player.Player
	playerEventSource map[string]eventSource
}

func NewPlayerHandler(
	logger logger.Logger[*slog.Logger],
	player *player.Player,
) *PlayerHandler {
	return &PlayerHandler{
		client: new(youtube.Client),
		logger: logger,
		player: player,

		playerEventSource: map[string]eventSource{},
	}
}

func (h *PlayerHandler) Play(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	youtubeURL, err := url.Parse(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		return fmt.Errorf("url parse: %w", err)
	}

	gID, vcID, err := h.resolveIDs(s, i)
	if err != nil {
		return fmt.Errorf("resolve ids: %w", err)
	}

	es, ok := h.playerEventSource[gID]
	if !ok {
		eventStream := make(chan string)

		go h.processAudioQueue(s, gID, i.ChannelID, vcID, eventStream)

		es = eventSource{
			voiceChannelID: vcID,
			source:         eventStream,
		}

		h.playerEventSource[gID] = es
	}

	if es.voiceChannelID != vcID {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must be in the same voice channel as the bot!",
			},
		})
	}

	h.player.AddToQueue(gID, youtubeURL.String())

	if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Song added to queue",
		},
	}); err != nil {
		return fmt.Errorf("session interaction respond")
	}

	return nil
}

func (h *PlayerHandler) processAudioQueue(
	s *discordgo.Session,
	gID, chID, vcID string,
	eventSource chan string,
) {
	defer func() {
		delete(h.playerEventSource, gID)
		close(eventSource)
	}()

	vc, chVoiceJoinErr := s.ChannelVoiceJoin(gID, vcID, false, true)
	if chVoiceJoinErr != nil {
		h.logger.Error("channel voice join", "err", chVoiceJoinErr)
		return
	}

	defer vc.Disconnect()

	if err := vc.Speaking(true); err != nil {
		h.logger.Error("initiate speaking", "err", err)
		return
	}

	defer vc.Speaking(false)

outerLoop:
	for h.player.Next(gID) {
		youtubeURL, err := h.player.PopQueue(gID)
		if err != nil && errors.Is(err, player.ErrEmptyQueue) {
			return
		}

		audioStream, title, err := h.fetchAudioStream(youtubeURL)
		if _, err = s.ChannelMessageSend(chID, fmt.Sprintf("Playing song [%s](%s)", title, youtubeURL)); err != nil {
			h.logger.Error("channel message", "err", err)
			return
		}

		enc, err := dca.EncodeMem(audioStream, dca.StdEncodeOptions)
		if err != nil {
			h.logger.Error("dca encode mem", "err", err)
			return
		}

		done := make(chan error)
		stream := dca.NewStream(enc, vc, done)

		cleanedUP := false
		cleanup := func() {
			if !cleanedUP {
				enc.Cleanup()
				audioStream.Close()
				cleanedUP = true
			}
		}

		for {
			select {
			case streamErr := <-done:
				cleanup()
				if streamErr != nil && !errors.Is(streamErr, io.EOF) {
					h.logger.Error("stream ended with an error", "err", streamErr)
					return
				}

				continue outerLoop
			case event := <-eventSource:
				switch event {
				case "stop":
					h.player.CleanQueue(gID)
					cleanup()
					return
				case "pause":
					stream.SetPaused(true)
				case "resume":
					stream.SetPaused(false)
				case "skip":
					cleanup()
					continue
				default:
					h.logger.Warn("received unknown event", "event", event)
				}
			}
		}
	}
}

func (h *PlayerHandler) Stop(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return h.handleEvent(s, i, "stop")
}

func (h *PlayerHandler) Pause(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return h.handleEvent(s, i, "pause")
}

func (h *PlayerHandler) Resume(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return h.handleEvent(s, i, "resume")
}

func (h *PlayerHandler) Skip(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return h.handleEvent(s, i, "skip")
}

func (h *PlayerHandler) handleEvent(s *discordgo.Session, i *discordgo.InteractionCreate, event string) error {
	gID, vcID, err := h.resolveIDs(s, i)
	if err != nil {
		return fmt.Errorf("resolve id: %w", err)
	}

	if es, ok := h.playerEventSource[gID]; ok {
		if es.voiceChannelID != vcID {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You must be in the same voice channel as the bot!",
				},
			})
		}

		if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "OK!",
			},
		}); err != nil {
			return fmt.Errorf("session interaction respond: %w", err)
		}

		es.source <- event
		return nil
	}

	return fmt.Errorf("not playing any songs now")
}

func (h *PlayerHandler) resolveIDs(s *discordgo.Session, i *discordgo.InteractionCreate) (string, string, error) {
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		return "", "", fmt.Errorf("state guild: %w", err)
	}

	gID := g.ID
	var vcID string
	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Interaction.Member.User.ID {
			vcID = vs.ChannelID
			break
		}
	}

	if vcID == "" {
		return "", "", fmt.Errorf("not connected to the voice channel")
	}

	return gID, vcID, nil
}

func (h *PlayerHandler) fetchAudioStream(url string) (io.ReadCloser, string, error) {
	v, err := h.client.GetVideo(url)
	if err != nil {
		return nil, "", fmt.Errorf("client get video: %w", err)
	}

	audioFormats := v.Formats.WithAudioChannels().Select(func(f youtube.Format) bool {
		return f.AudioQuality == "AUDIO_QUALITY_MEDIUM" && f.AudioSampleRate == "48000"
	})

	audioStream, _, err := h.client.GetStream(v, &audioFormats[0])
	if err != nil {
		return nil, "", fmt.Errorf("client get stream: %w", err)
	}

	return audioStream, v.Title, nil
}
