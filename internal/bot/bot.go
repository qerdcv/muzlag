package bot

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"

	"github.com/qerdcv/muzlag.go/internal/bot/handler"
	"github.com/qerdcv/muzlag.go/internal/logger"
)

type Bot struct {
	session *discordgo.Session
	client  *youtube.Client

	logger logger.Logger[*slog.Logger]

	commandHandlers map[string]func(session *discordgo.Session, i *discordgo.InteractionCreate) error
}

func New(
	l logger.Logger[*slog.Logger],
	token string,
	healthHandler *handler.HealthHandler,
	playerHandler *handler.PlayerHandler,
) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("discordgo new: %w", err)
	}

	b := &Bot{
		session: dg,
		client:  new(youtube.Client),
		logger:  l,
	}

	b.commandHandlers = map[string]func(session *discordgo.Session, i *discordgo.InteractionCreate) error{
		// health
		commandPing: healthHandler.Ping,
		// ======

		// player
		commandPlay:   playerHandler.Play,
		commandStop:   playerHandler.Stop,
		commandPause:  playerHandler.Pause,
		commandResume: playerHandler.Resume,
		commandSkip:   playerHandler.Skip,
		// ======
	}

	dg.AddHandler(b.commandHandler)

	return b, nil
}

func (b *Bot) Run() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("session open: %w", err)
	}

	if err := b.registerCommands(); err != nil {
		return fmt.Errorf("register commands: %w", err)
	}

	return nil
}

func (b *Bot) registerCommands() error {
	for _, v := range commands {
		if _, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", v); err != nil {
			return fmt.Errorf("session application command create: %w", err)
		}
	}

	return nil
}

func (b *Bot) Close() error {
	if err := b.session.Close(); err != nil {
		return fmt.Errorf("session close: %w", err)
	}

	return nil
}
