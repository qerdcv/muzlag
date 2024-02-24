package bot

import "github.com/bwmarrin/discordgo"

const (
	commandPing = "ping"
	commandPlay = "play"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        commandPing,
		Description: "Ping-pong healthcheck",
	},
	{
		Name:        commandPlay,
		Description: "Play song by youtube video URL",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ytb-url",
				Description: "YouTube Video URL",
				Required:    true,
			},
		},
	},
}
