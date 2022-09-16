package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/missions"
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

type DiscordCommand struct {
	Data    *discordgo.ApplicationCommand
	Handler func(*discordgo.Session, *discordgo.InteractionCreate)
}

var Commands = []*DiscordCommand{
	{
		Data: &discordgo.ApplicationCommand{
			Name:        "test-command",
			Description: "Command for testing purpose",
		},
		Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Command registered successfully",
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
					Name:        "description",
					Description: "Short description of the mission (255 characters)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "class",
					Description: "Mission class that will define the rules of completion",
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
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Mission created successfully",
				},
			})
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
					Name:        "description",
					Description: "Description of the mission",
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "type",
					Description: "Type of mission",
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
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Mission canceled successfully",
				},
			})
		},
	},
}
