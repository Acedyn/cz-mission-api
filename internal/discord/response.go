package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/missions"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

const (
	THUMBNAIL_SUCCESS = "https://icons.iconarchive.com/icons/paomedia/small-n-flat/32/sign-check-icon.png"
	THUMBNAIL_ERROR   = "https://icons.iconarchive.com/icons/paomedia/small-n-flat/32/sign-error-icon.png"
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

func MissionEmbed(mission *models.Mission) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Mission %s", mission.Format()),
		Description: mission.ShortDescription,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Description",
				Value:  mission.LongDescription,
				Inline: false,
			},
			{
				Name:   "Category",
				Value:  mission.GetClassData().Category,
				Inline: true,
			},
			{
				Name:   "Class",
				Value:  mission.Class,
				Inline: true,
			},
			{
				Name:   "Reward",
				Value:  fmt.Sprintf("%f", mission.Reward),
				Inline: true,
			},
			{
				Name:   "Close date",
				Value:  mission.CloseAt.Format("2 jan 2006 15:04"),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: mission.GetClassData().Logo,
		},
		URL: "https://cardboardcitizen.com",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Mission last updated on %s", mission.UpdatedAt.Format("2 jan 2006 15:04")),
		},
	}
}

func MissionResponse(mission *models.Mission, message string, embeds []*discordgo.MessageEmbed, components []discordgo.MessageComponent) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf("Mission %s#%d %s", mission.Name, mission.ID, message),
			Embeds:     embeds,
			Components: components,
		},
	}
}

func ErrorResponse(message string, err error, retry string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Error report",
				Description: fmt.Sprintln(err),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: THUMBNAIL_ERROR,
				},
			}},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "🪃",
							},
							Label:    "Retry",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("retry-command:%s", retry),
						},
					},
				},
			},
		},
	}
}

func CreateMissionResponse(controller *DiscordController, mission *models.Mission) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 771906 // #0BC742
	missionEmbed.Footer.IconURL = THUMBNAIL_SUCCESS

	response := MissionResponse(mission,
		"successfully created",
		[]*discordgo.MessageEmbed{missionEmbed},
		[]discordgo.MessageComponent{getMissionActionRow(controller, mission)},
	)
	response.Data.Title = "Mission created"
	return response
}

func CreateMissionError(mission *models.Mission, err error) *discordgo.InteractionResponse {
	response := ErrorResponse(fmt.Sprintf("An error occured during the creation of the mission %s#%d", mission.Name, mission.ID), err, "create-mission")
	return response
}

func GetMissionResponse(controller *DiscordController, mission *models.Mission) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 39423 // #0099FF
	missionEmbed.Footer.IconURL = THUMBNAIL_SUCCESS

	response := MissionResponse(mission,
		"found",
		[]*discordgo.MessageEmbed{missionEmbed},
		[]discordgo.MessageComponent{getMissionActionRow(controller, mission)},
	)
	response.Data.Title = fmt.Sprintf("Mission %s", mission.Format())
	return response
}

func GetMissionError(err error) *discordgo.InteractionResponse {
	response := ErrorResponse("No mission could be found with the provided parameters", err, "get-mission")
	return response
}

func UpdateMissionResponse(mission *models.Mission) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 39423 // #0099FF
	missionEmbed.Footer.IconURL = THUMBNAIL_SUCCESS

	response := MissionResponse(mission, "successfully updated", []*discordgo.MessageEmbed{missionEmbed}, []discordgo.MessageComponent{})
	response.Data.Title = "Mission updated"
	return response
}

func CancelMissionResponse(mission *models.Mission) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 15219772 // #E83C3C
	missionEmbed.Footer.IconURL = THUMBNAIL_ERROR

	response := MissionResponse(mission, "successfully canceled", []*discordgo.MessageEmbed{missionEmbed}, []discordgo.MessageComponent{})
	response.Data.Title = "Mission canceled"
	return response
}
