package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/utils"
)

func InitializeSession(botToken string) (session *discordgo.Session, err error) {
	session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		utils.Log.Error("Could not open discord session :", err)
		return nil, err
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		utils.Log.Info("Discord session opened as :", s.State.User.Username, "#", s.State.User.Discriminator)
	})

	err = session.Open()
	if err != nil {
		utils.Log.Error("Could not open discord session :", err)
		return nil, err
	}

	return session, nil
}

func RegisterCommands(session *discordgo.Session, guildID string) (err error) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for index, command := range Commands {
		appCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, command.Data)
		if err != nil {
			utils.Log.Error("Could not create command", command.Data.Name, ":", err)
			return err
		}
		registeredCommands[index] = appCommand
	}

    utils.Log.Info("Discord commands registered :", registeredCommands)
	return nil
}
