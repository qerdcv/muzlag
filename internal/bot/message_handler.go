package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	msgContent := m.Content
	switch {
	case strings.HasPrefix(msgContent, ".ping"):
		s.ChannelMessageSend(m.ChannelID, "pong")
	case strings.HasPrefix(msgContent, ".play"):
		print("handling play command")
		if err := b.playCommand(s, m); err != nil {
			log.Println("play command: ", err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	}
}
