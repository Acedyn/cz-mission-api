package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func initializeSession() (*discordgo.Session, error) {
	session, createSessionErr := discordgo.New("Bot " + *botToken)
	if createSessionErr != nil {
		log.Fatalf("Could not open discord session: %v", createSessionErr)
		return nil, createSessionErr
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	openSessionErr := session.Open()
	if openSessionErr != nil {
		log.Fatalf("Could not open discord session: %v", openSessionErr)
		return nil, openSessionErr
	}

	return session, nil
}

func RegisterCommands(session discordgo.Session) error {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for index, command := range Commands {
		appCommand, createCommandErr := session.ApplicationCommandCreate(session.State.User.ID, *guildID, command.Data)
		if createCommandErr != nil {
			log.Fatalf("Could not create command '%v': %v", command.Data.Name, createCommandErr)
			return createCommandErr
		}
		registeredCommands[index] = appCommand
	}

    return nil
}
