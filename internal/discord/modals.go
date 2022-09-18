package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

func MissionModal(mission *models.Mission) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:       "Name",
					Style:       discordgo.TextInputShort,
					Placeholder: "Insert new name",
					Value:       mission.Name,
					CustomID:    fmt.Sprintf("update-mission:%d:name", mission.ID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:       "Short Description",
					Style:       discordgo.TextInputShort,
					Placeholder: "Insert new short description",
					Value:       mission.ShortDescription,
					CustomID:    fmt.Sprintf("update-mission:%d:short_description", mission.ID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:       "Long Description",
					Style:       discordgo.TextInputParagraph,
					Placeholder: "Insert new long description",
					Value:       mission.LongDescription,
					CustomID:    fmt.Sprintf("update-mission:%d:long_description", mission.ID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					Placeholder: "Select new class",
					MaxValues:   1,
					Options:     getMissionOptions(mission.Class),
					CustomID:    fmt.Sprintf("update-mission:%d:reward", mission.ID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:       "Reward",
					Style:       discordgo.TextInputShort,
					Placeholder: "Insert new reward",
					Value:       fmt.Sprintf("%f", mission.Reward),
					CustomID:    fmt.Sprintf("update-mission:%d:reward", mission.ID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "🔧",
					},
					Label:    "Update",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("update-mission:%d", mission.ID),
				},
				discordgo.Button{
					Emoji: discordgo.ComponentEmoji{
						Name: "🗑️",
					},
					Label:    "Cancel",
					Style:    discordgo.DangerButton,
					CustomID: fmt.Sprintf("cancel-mission:%d", mission.ID),
				},
			},
		},
	}
}
