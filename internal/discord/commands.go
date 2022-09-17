package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/database"
	"github.com/cardboard-citizens/cz-mission-api/internal/missions"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
	"strings"
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

func GetCommands(dbController *database.DatabaseController) (commands []*DiscordCommand) {
	return []*DiscordCommand{
		{
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
		{
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
				err := dbController.CreateMission(&mission)
				if err != nil {
					utils.Log.Error(fmt.Errorf("An error occured while creating the mission %s\n\t%s", mission.Format(), err))
					session.InteractionRespond(interaction.Interaction, CreateMissionError(&mission, err))
					return
				}

				err = session.InteractionRespond(interaction.Interaction, CreateMissionResponse(&mission))
				if err != nil {
					utils.Log.Error(fmt.Errorf("An error occured while responding to the interaction %s\n\t%s", interaction.Message.ID, err))
					session.InteractionRespond(interaction.Interaction, CreateMissionError(&mission, err))
					return
				}
			},
		},
		{
			Data: &discordgo.ApplicationCommand{
				Name:        "update-mission",
				Description: "Update a mission",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "id",
						Description: "ID of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "name",
						Description: "Name of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
					},
					{
						Name:        "short-description",
						Description: "Short description of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
					},
					{
						Name:        "long-description",
						Description: "Long description of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
					},
					{
						Name:        "class",
						Description: "The mission class defines the rules to complete the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Choices:     getMissionChoices(),
					},
					{
						Name:        "reward",
						Description: "Reward of the mission",
						Type:        discordgo.ApplicationCommandOptionNumber,
						MinValue:    new(float64),
					},
				},
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				utils.Log.Debug("Update mission slash command received")

				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Mission updated successfully",
					},
				})
			},
		},
		{
			Data: &discordgo.ApplicationCommand{
				Name:        "get-missions",
				Description: "Get all missions data",
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				utils.Log.Debug("Get missions slash command received")
				missions := dbController.GetMissions()

				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Existing missions: %+v", missions),
					},
				})
			},
		},
		{
			Data: &discordgo.ApplicationCommand{
				Name:        "cancel-mission",
				Description: "Cancel a mission",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "id",
						Description: "ID of the mission",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				utils.Log.Debug("Cancel mission slash command received")

				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Mission canceled successfully",
					},
				})
			},
		},
	}
}
