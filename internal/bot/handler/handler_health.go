package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return new(HealthHandler)
}

func (h *HealthHandler) Ping(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	}); err != nil {
		return fmt.Errorf("session channel message send")
	}

	return nil
}
