package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/database"
	"github.com/cardboard-citizens/cz-mission-api/internal/missions"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

func getMissionChoices() []*discordgo.ApplicationCommandOptionChoice {
	missionClassKeys := missions.GetMissionClassKeys()
	missionChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for _, classKey := range missionClassKeys {
		missionChoices = append(missionChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  strings.Title(strings.Replace(classKey, "-", " ", -1)),
			Value: classKey,
		})
	}
	return missionChoices
}

func getSortChoices() []*discordgo.ApplicationCommandOptionChoice {
	sortKeys := database.GetMissionFieldNames()
	missionChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for _, sortKey := range sortKeys {
		missionChoices = append(missionChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  strings.Title(strings.Replace(sortKey, "_", " ", -1)),
			Value: sortKey,
		})
	}
	return missionChoices
}

func getInteractionOptions(interaction *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionList := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(optionList))
	for _, option := range optionList {
		optionMap[option.Name] = option
	}
	return optionMap
}

type DiscordCommand struct {
	Data    *discordgo.ApplicationCommand
	Handler func(*discordgo.Session, *discordgo.InteractionCreate)
}

func getCommands(controller *DiscordController) map[string]*DiscordCommand {
	return map[string]*DiscordCommand{
		"ping-command": {
			Data: &discordgo.ApplicationCommand{
				Name:        "ping-command",
				Description: "Command for testing purpose",
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				utils.Log.Info("Ping from discord slash command")

				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Pong, connection with the server successfull",
					},
				})
			},
		},
		"create-mission": {
			Data: &discordgo.ApplicationCommand{
				Name:        "create-mission",
				Description: "Create a mission",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "name",
						Description: "Name of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "short-description",
						Description: "Short description of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "long-description",
						Description: "Long description of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "class",
						Description: "The mission class defines the rules to complete the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Choices:     getMissionChoices(),
						Required:    true,
					},
					{
						Name:        "reward",
						Description: "Reward of the mission",
						Type:        discordgo.ApplicationCommandOptionNumber,
						MinValue:    new(float64),
						Required:    true,
					},
				},
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				utils.Log.Debug("Create mission slash command received")
				options := getInteractionOptions(interaction)
				mission := models.Mission{
					Name:             options["name"].StringValue(),
					ShortDescription: options["short-description"].StringValue(),
					LongDescription:  options["long-description"].StringValue(),
					Class:            options["class"].StringValue(),
					Reward:           options["reward"].FloatValue(),
				}
				err := controller.databaseController.CreateMission(&mission)
				if err != nil {
					utils.Log.Error(fmt.Errorf("An error occured while creating the mission %s\n\t%s", mission.Format(), err))
					session.InteractionRespond(interaction.Interaction, CreateMissionError(&mission, err))
					return
				}

				err = session.InteractionRespond(interaction.Interaction, controller.CreateMissionResponse(&mission))
				if err != nil {
					utils.Log.Error(fmt.Errorf("An error occured while responding to the interaction %s\n\t%s", interaction.ID, err))
					session.InteractionRespond(interaction.Interaction, CreateMissionError(&mission, err))
					return
				}
			},
		},
		"get-mission": {
			Data: &discordgo.ApplicationCommand{
				Name:        "get-missions",
				Description: "Get multiple missions data",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "limit",
						Description: "Max missions to retrieve",
						Type:        discordgo.ApplicationCommandOptionInteger,
						MaxValue:    10.0,
					},
					{
						Name:        "sort",
						Description: "Key to sort the existing missions",
						Type:        discordgo.ApplicationCommandOptionString,
						Choices:     getSortChoices(),
					},
				},
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				utils.Log.Debug("Get missions slash command received")
				options := getInteractionOptions(interaction)
				var limit int = 1
				if limitOption, ok := options["limit"]; ok {
					limit = int(limitOption.IntValue())
				}
				var sort *string = nil
				if sortOption, ok := options["sort"]; ok {
					sortValue := sortOption.StringValue()
					sort = &sortValue
				}
				missions := controller.databaseController.GetMissions(limit, sort)

				if missions == nil || len(missions) == 0 {
					session.InteractionRespond(
						interaction.Interaction,
						GetMissionError(fmt.Errorf("Database returned 0 matched with the options \nLimit: %d\nSort: %s", limit, *sort)),
					)
					return
				}

				for _, mission := range missions {
					err := session.InteractionRespond(interaction.Interaction, GetMissionResponse(controller, &mission))
					if err != nil {
						utils.Log.Error(fmt.Errorf("An error occured while responding to the interaction %s\n\t%s", interaction.Message.ID, err))
						return
					}
				}
			},
		},
	}
}

func init() {
	controllerInitializers = append(controllerInitializers, func(controller *DiscordController) (err error) {
		controller.Commands = getCommands(controller)
		controller.RegisteredCommands = make([]*discordgo.ApplicationCommand, len(controller.Commands))
		for _, command := range controller.Commands {
			appCommand, err := controller.Session.ApplicationCommandCreate(controller.Session.State.User.ID, controller.GuildId, command.Data)
			if err != nil {
				return fmt.Errorf("Could not register command %s\n\t%s", command.Data.Name, err)
			}
			utils.Log.Info("Discord command registered :", command.Data.Name)
			controller.RegisteredCommands = append(controller.RegisteredCommands, appCommand)
		}

		utils.Log.Info("Discord commands registration completed")
		return err
	})

	controllerListeners = append(controllerListeners, func(controller *DiscordController) (err error) {
		controller.Session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			if interaction.Type != discordgo.InteractionApplicationCommand {
				return
			}
			command, ok := controller.Commands[interaction.ApplicationCommandData().Name]
			if ok {
				command.Handler(session, interaction)
			}
		})
		utils.Log.Info("Discord command listeners started")
		return err
	})

	controllerCleanups = append(controllerCleanups, func(controller *DiscordController) (err error) {
		for _, command := range controller.RegisteredCommands {
			err := controller.Session.ApplicationCommandDelete(controller.Session.State.User.ID, controller.GuildId, command.ID)
			if err != nil {
				return fmt.Errorf("Could not deregister command %s\n\t%s", command.Name, err)
			}
			utils.Log.Info("Discord command deregistered :", command.Name)
		}

		utils.Log.Info("Discord commands deregistration completed")
		return err
	})
}
