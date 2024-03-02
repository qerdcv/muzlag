package bot

import "github.com/bwmarrin/discordgo"

const (
	commandPing = "ping"

	// queue
	commandPlay   = "play"
	commandStop   = "stop"
	commandPause  = "pause"
	commandResume = "resume"
	commandSkip   = "skip"
	commandQueue  = "queue"
	// =====
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
	{
		Name:        commandStop,
		Description: "Stop current music queue",
	},
	{
		Name:        commandPause,
		Description: "Pauses current music queue",
	},
	{
		Name:        commandResume,
		Description: "Resumes current music queue",
	},
	{
		Name:        commandSkip,
		Description: "Skips current song in the music queue",
	},
	{
		Name:        commandQueue,
		Description: "Display current state of the queue",
	},
}
