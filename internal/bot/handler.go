package bot

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ErrUnknownCommand = errors.New("unknown command")

func (b *Bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error
	command := i.ApplicationCommandData().Name

	if h, ok := b.commandHandlers[command]; ok {
		err = h(s, i)
	} else {
		err = ErrUnknownCommand
	}

	if err != nil {
		if sendErr := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error while handling command: %s", err.Error()),
				},
			},
		); sendErr != nil {
			panic(sendErr.Error())
		}
	}

}
