package bot

import "github.com/bwmarrin/discordgo"

type (
	command        string
	commandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error
)

const (
	play command = "play"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        play.String(),
			Description: "play the music!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "url of clip to play",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[command]commandHandler{
		play: playCommand,
	}
)

func (c command) String() string {
	return string(c)
}
