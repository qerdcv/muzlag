package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/qerdcv/muzlag/internal/downloader"
	"github.com/qerdcv/muzlag/internal/output"
)

func playCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Playing",
		},
	}); err != nil {
		return fmt.Errorf("interaction respond: %w", err)
	}

	var chID string
	// Find the channel that the message came from.
	c, err := s.State.Channel(i.ChannelID)
	if err != nil {
		// Could not find channel.
		// TODO: return error
		return errors.New("no channel found")
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		// TODO: return error
		return errors.New("no guild found")
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			chID = vs.ChannelID
		}
	}

	if chID == "" {
		return errors.New("no channel id found")
	}

	var url string
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "url" {
			url = opt.Value.(string)
		}
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(g.ID, chID, false, true)
	if err != nil {
		return fmt.Errorf("join voice channel: %w", err)
	}

	o := output.NewChOut(vc.OpusSend)
	d := downloader.NewYoutubeDownloader(url, o)

	// Start speaking.
	if err = vc.Speaking(true); err != nil {
		return fmt.Errorf("voice channel enable speak: %w", err)
	}

	if err = d.Download(); err != nil {
		log.Println("ERROR: download ", err.Error())
	}

	// Stop speaking
	if err = vc.Speaking(false); err != nil {
		return fmt.Errorf("voice channel disable speak: %w", err)
	}

	// Disconnect from the provided voice channel.
	if err = vc.Disconnect(); err != nil {
		return fmt.Errorf("voice channel disconnect: %w", err)
	}

	return nil
}
