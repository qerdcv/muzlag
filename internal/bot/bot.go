package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

type Bot struct {
	session *discordgo.Session
	client  *youtube.Client
}

func New(token string) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("discordgo new: %w", err)
	}

	b := &Bot{
		session: dg,
		client:  new(youtube.Client),
	}

	dg.AddHandler(b.messageHandler)

	return b, nil
}

func (b *Bot) Run() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("session open: %w", err)
	}

	return nil
}

func (b *Bot) Close() error {
	if err := b.session.Close(); err != nil {
		return fmt.Errorf("session close: %w", err)
	}

	return nil
}
