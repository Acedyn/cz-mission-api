package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/database"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/utils"
)

type DiscordController struct {
	*discordgo.Session
	GuildId            string
	Commands           []*DiscordCommand
	RegisteredCommands []*discordgo.ApplicationCommand
}

func (controller *DiscordController) Initialize(botToken string, databaseController *database.DatabaseController) (err error) {
	controller.Session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		return fmt.Errorf("Could not open discord session\n\t%s", err)
	}

	controller.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		utils.Log.Info("Discord session opened as :", s.State.User.Username, "#", s.State.User.Discriminator)
	})

	err = controller.Session.Open()
	if err != nil {
		return fmt.Errorf("Could not open discord session\n\t%s", err)
	}

	controller.Commands = GetCommands(databaseController)
	return err
}

func (controller *DiscordController) RegisterCommands() (err error) {
	controller.RegisteredCommands = make([]*discordgo.ApplicationCommand, len(controller.Commands))
	for index, command := range controller.Commands {
		appCommand, err := controller.Session.ApplicationCommandCreate(controller.Session.State.User.ID, controller.GuildId, command.Data)
		if err != nil {
			return fmt.Errorf("Could not register command %s\n\t%s", command.Data.Name, err)
		}
		utils.Log.Info("Discord command registered :", command.Data.Name)
		controller.RegisteredCommands[index] = appCommand
	}

	utils.Log.Info("Discord commands registration completed")
	return err
}

func (controller *DiscordController) DeregisterCommands() (err error) {
	for _, command := range controller.RegisteredCommands {
		err := controller.Session.ApplicationCommandDelete(controller.Session.State.User.ID, controller.GuildId, command.ID)
		if err != nil {
			return fmt.Errorf("Could not deregister command %s\n\t%s", command.Name, err)
		}
		utils.Log.Info("Discord command deregistered :", command.Name)
	}

	utils.Log.Info("Discord commands deregistration completed")
	return err
}

func (controller *DiscordController) ListenCommands() {
	controller.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		for _, command := range controller.Commands {
			if command.Data.Name == i.ApplicationCommandData().Name {
				command.Handler(s, i)
			}
		}
	})
}
