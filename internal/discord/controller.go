package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/database"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

var (
	controllerInitializers = make([]func(*DiscordController) error, 0)
	controllerListeners    = make([]func(*DiscordController) error, 0)
	controllerCleanups     = make([]func(*DiscordController) error, 0)
)

type DiscordController struct {
	*discordgo.Session
	RegisteredCommands []*discordgo.ApplicationCommand
	GuildId            string

	Commands map[string]*DiscordCommand
	Buttons  map[string]*DiscordButton
	Modals   map[string]*DiscordModal

	initialized        bool
	DatabaseController *database.DatabaseController
}

func (controller *DiscordController) Initialize(botToken string,) (err error) {
	controller.Session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		return fmt.Errorf("Could not open discord session\n\t%s", err)
	}

	controller.Session.AddHandler(
		func(session *discordgo.Session, ready *discordgo.Ready) {
			utils.Log.Info(
				"Discord session opened as :",
				session.State.User.Username,
				"#",
				session.State.User.Discriminator,
			)
		},
	)

	err = controller.Session.Open()
	if err != nil {
		return fmt.Errorf("Could not open discord session\n\t%s", err)
	}

	for _, initializer := range controllerInitializers {
		initializer(controller)
	}

	controller.initialized = true
	return err
}

func (controller *DiscordController) Listen() (err error) {
	for _, listener := range controllerListeners {
		err = listener(controller)
		if err != nil {
			return fmt.Errorf("Could not start discord listener\n\t%s", err)
		}
	}
	return err
}

func (controller *DiscordController) Shutdown() (err error) {
	for _, cleanup := range controllerCleanups {
		err = cleanup(controller)
		if err != nil {
			return fmt.Errorf("Could not start stop listener\n\t%s", err)
		}
	}

	err = controller.Session.Close()
	if err != nil {
		return fmt.Errorf("Could not close discord session\n\t%s", err)
	}
	return err
}
