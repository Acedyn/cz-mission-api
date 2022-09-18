package discord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

type DiscordButton struct {
	Button  discordgo.Button
	Handler func(*discordgo.Session, *discordgo.InteractionCreate, string)
}

func getButtons(controller *DiscordController) map[string]*DiscordButton {
	return map[string]*DiscordButton{
		"update-mission": {
			Button: discordgo.Button{
				Emoji: discordgo.ComponentEmoji{
					Name: "🔧",
				},
				Label: "Update",
				Style: discordgo.PrimaryButton,
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, id string) {
				utils.Log.Debug("Update mission", "", "component button instruction received")

				missionId, err := strconv.ParseInt(id, 10, 32)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord button", id, "Invalid ID")
					return
				}
				mission := controller.databaseController.GetMission(uint32(missionId))
				if mission == nil {
					utils.Log.Error("Could not handle mission update discord button", id, ": Mission not found")
					return
				}

				err = session.InteractionRespond(
					interaction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseModal,
						Data: &discordgo.InteractionResponseData{
							CustomID:   fmt.Sprintf("%s:%s", "update-mission", interaction.ID),
							Title:      "Update Mission",
							Components: getMissionModal(mission),
						},
					},
				)

				if err != nil {
					utils.Log.Error(
						fmt.Errorf(
							"An error occured while responding to the interaction %s\n\t%s",
							interaction.ID,
							err,
						),
					)
					return
				}
			},
		},
		"cancel-mission": {
			Button: discordgo.Button{
				Emoji: discordgo.ComponentEmoji{
					Name: "🗑️",
				},
				Label: "Cancel",
				Style: discordgo.DangerButton,
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, id string) {
				utils.Log.Debug("Cancel mission component button instruction received")
			},
		},
	}
}

func getMissionActionRow(
	controller *DiscordController,
	mission *models.Mission,
) discordgo.ActionsRow {
	actionsRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{},
	}
	for _, idKey := range []string{"update-mission", "cancel-mission"} {
		button, ok := controller.Buttons[idKey]
		if !ok {
			utils.Log.Error(
				"Could not get the mission component",
				idKey,
				": Component does not exists",
			)
			continue
		}
		button.Button.CustomID = fmt.Sprintf("%s:%d", idKey, mission.ID)
		actionsRow.Components = append(actionsRow.Components, button.Button)
	}
	return actionsRow
}

func init() {
	controllerInitializers = append(
		controllerInitializers,
		func(controller *DiscordController) (err error) {
			controller.Buttons = getButtons(controller)
			return err
		},
	)

	controllerListeners = append(
		controllerListeners,
		func(controller *DiscordController) (err error) {
			controller.Session.AddHandler(
				func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
					if interaction.Type != discordgo.InteractionMessageComponent {
						return
					}
					idSplit := strings.Split(
						interaction.MessageComponentData().CustomID,
						":",
					)
					button, ok := controller.Buttons[idSplit[0]]
					if ok {
						button.Handler(session, interaction, idSplit[1])
					}
				},
			)
			utils.Log.Info("Discord components listeners started")
			return err
		},
	)
}
