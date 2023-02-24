package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	intents = discordgo.IntentGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildVoiceStates
)

type Bot struct {
	session *discordgo.Session
}

//func init() {
//	if err := loadSound(); err != nil {
//		log.Fatalln(fmt.Errorf("load sound: %w", err))
//	}
//}

func New() (*Bot, error) {
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("discordgo new: %w", err)
	}

	b := &Bot{dg}
	b.setIntents()
	b.addHandlers()

	return b, nil
}

func (b *Bot) RegisterCommands() error {
	for _, c := range commands {
		if _, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", c); err != nil {
			return fmt.Errorf("session application command create: %w", err)
		}
	}

	return nil
}

func (b *Bot) addHandlers() {
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Println("handling interaction: ", i.ApplicationCommandData().Name)
		if h, ok := commandHandlers[command(i.ApplicationCommandData().Name)]; ok {
			if err := h(s, i); err != nil {
				log.Println("ERROR: ", err.Error())
			}
		}
	})
}

func (b *Bot) setIntents() {
	b.session.Identify.Intents = intents
}

func (b *Bot) Open() error {
	return b.session.Open()
}

func (b *Bot) Close() error {
	return b.session.Close()
}
