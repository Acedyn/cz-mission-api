package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

const (
	THUMBNAIL_SUCCESS = "https://icons.iconarchive.com/icons/paomedia/small-n-flat/32/sign-check-icon.png"
	THUMBNAIL_ERROR   = "https://icons.iconarchive.com/icons/paomedia/small-n-flat/32/sign-error-icon.png"
)

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
				Value:  mission.Category,
				Inline: true,
			},
			{
				Name:   "Class",
				Value:  mission.Class,
				Inline: true,
			},
			{
				Name:   "Reward",
				Value:  fmt.Sprintf("%d", int(mission.Reward)),
				Inline: true,
			},
			{
				Name:   "Close date",
				Value:  mission.CloseAt.Format("2 jan 2006 15:04"),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: mission.Logo,
		},
		URL: "https://cardboardcitizen.com",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf(
				"Mission last updated on %s",
				mission.UpdatedAt.Format("2 jan 2006 15:04"),
			),
		},
	}
}

func MissionResponse(
	mission *models.Mission,
	message string,
	embeds []*discordgo.MessageEmbed,
	components []discordgo.MessageComponent,
) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Mission %s#%d %s",
				mission.Name,
				mission.ID,
				message,
			),
			Embeds:     embeds,
			Components: components,
		},
	}
}

func ErrorResponse(message string, err error) *discordgo.InteractionResponse {
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
		},
	}
}

func MissionResponseError(mission *models.Mission, err error,
) *discordgo.InteractionResponse {
	response := ErrorResponse(
		fmt.Sprintf(
			"An error occured when processing call for the mission %s#%d",
			mission.Name,
			mission.ID,
		),
		err,
	)
	return response
}

func MissionModalResponse(
	title string,
	controller *DiscordController,
	mission *models.Mission,
) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("%s:%d", "update-mission", mission.ID),
			Title:      title,
			Components: getMissionModal(controller, mission),
		},
	}
}

func MissionParametersModalResponse(
	title string,
	controller *DiscordController,
	mission *models.Mission,
) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("%s:%d", "set-mission-parameters", mission.ID),
			Title:      title,
			Components: getMissionParameterModal(controller, mission),
		},
	}
}

func CreateMissionResponse(
	controller *DiscordController,
	mission *models.Mission,
) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 39423 // #0099FF
	missionEmbed.Footer.IconURL = THUMBNAIL_SUCCESS

	response := MissionResponse(
		mission,
		"successfully created",
		[]*discordgo.MessageEmbed{missionEmbed},
		[]discordgo.MessageComponent{getMissionActionRow(controller, mission)},
	)
	response.Data.Title = fmt.Sprintf("Mission %s", mission.Format())
	return response
}

func GetMissionResponse(
	controller *DiscordController,
	missions []models.Mission,
) *discordgo.InteractionResponse {
	missionEmbeds := make([]*discordgo.MessageEmbed, 0, len(missions))
	for _, mission := range missions {
		missionEmbed := MissionEmbed(&mission)
		missionEmbed.Color = 39423 // #0099FF
		missionEmbed.Footer.IconURL = THUMBNAIL_SUCCESS
		missionEmbeds = append(missionEmbeds, missionEmbed)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf("%d Missions found", len(missions)),
			Embeds:     missionEmbeds,
			Components: []discordgo.MessageComponent{getMissionActionRow(controller, &missions[0])},
		},
	}
}

func UpdateMissionResponse(
	mission *models.Mission,
) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 39423 // #0099FF
	missionEmbed.Footer.IconURL = THUMBNAIL_SUCCESS

	response := MissionResponse(
		mission,
		"successfully updated",
		[]*discordgo.MessageEmbed{missionEmbed},
		[]discordgo.MessageComponent{},
	)
	return response
}

func CancelMissionResponse(
	mission *models.Mission,
) *discordgo.InteractionResponse {
	missionEmbed := MissionEmbed(mission)
	missionEmbed.Color = 15219772 // #E83C3C
	missionEmbed.Footer.IconURL = THUMBNAIL_ERROR

	response := MissionResponse(
		mission,
		"successfully canceled",
		[]*discordgo.MessageEmbed{missionEmbed},
		[]discordgo.MessageComponent{},
	)
	return response
}
