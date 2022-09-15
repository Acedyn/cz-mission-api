package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/missions"
	"strings"
)

func getMissionChoices() []*discordgo.ApplicationCommandOptionChoice {
	var missionKeys = make([]*discordgo.ApplicationCommandOptionChoice, len(missions.Missions))
	for key := range missions.Missions {
		missionKeys = append(missionKeys, &discordgo.ApplicationCommandOptionChoice{
			Name:  strings.Title(strings.Replace(key, "-", " ", -1)),
			Value: key,
		})
	}
	return missionKeys
}

type DiscordCommand struct {
	Data    discordgo.ApplicationCommand
	Handler func(*discordgo.Session, *discordgo.InteractionCreate)
}

var Commands = []*DiscordCommand{
	{
		Data: discordgo.ApplicationCommand{
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
		Data: discordgo.ApplicationCommand{
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
					Description: "Description of the mission",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "type",
					Description: "Type of mission",
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
		Data: discordgo.ApplicationCommand{
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
		Data: discordgo.ApplicationCommand{
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
