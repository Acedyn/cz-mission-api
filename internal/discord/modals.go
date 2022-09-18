package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/missions"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

func getMissionOptions(defaultOption string) []discordgo.SelectMenuOption {
	missionClasses := missions.GetMissionsClasses()
	missionClassKeys := missions.GetMissionClassKeys()
	missionOptions := make([]discordgo.SelectMenuOption, 0)
	for _, classKey := range missionClassKeys {
		missionOptions = append(missionOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.Replace(classKey, "-", " ", -1)),
			Value:       classKey,
			Default:     classKey == defaultOption,
			Description: missionClasses[classKey].Description,
		})
	}
	return missionOptions
}

func getMissionModal(mission *models.Mission) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:       "Name",
					Style:       discordgo.TextInputShort,
					Placeholder: "Insert new name",
					Value:       mission.Name,
					CustomID:    "name",
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
					CustomID:    "short-description",
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
					CustomID:    "long-description",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					Label:       "Class",
					Style:       discordgo.TextInputShort,
					Placeholder: "Insert new class",
					Value:       mission.Class,
					CustomID:    "class",
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
					CustomID:    "reward",
				},
			},
		},
	}
}
