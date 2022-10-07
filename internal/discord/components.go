package discord

import (
	"fmt"
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
		"set-mission-parameters": {
			Button: discordgo.Button{
				Emoji: discordgo.ComponentEmoji{
					Name: "🔧",
				},
				Label: "Set Parameters",
				Style: discordgo.PrimaryButton,
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, id string) {
				utils.Log.Debug("Update mission parameters", id, "component instruction received")

				mission, err := controller.DatabaseController.GetMissionFromString(id)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord button\n\t", err)
					return
				}

				err = session.InteractionRespond(
					interaction.Interaction,
					MissionParametersModalResponse("Update Mission Parameters", controller, mission),
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
		"update-mission": {
			Button: discordgo.Button{
				Emoji: discordgo.ComponentEmoji{
					Name: "🔧",
				},
				Label: "Update",
				Style: discordgo.PrimaryButton,
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, id string) {
				utils.Log.Debug("Update mission", id, "component instruction received")

				mission, err := controller.DatabaseController.GetMissionFromString(id)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord button\n\t", err)
					return
				}

				err = session.InteractionRespond(
					interaction.Interaction,
					MissionModalResponse("Update Mission", controller, mission),
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
				utils.Log.Debug("Cancel mission", id, "component instruction received")

				mission, err := controller.DatabaseController.GetMissionFromString(id)
				if err != nil {
					utils.Log.Error("Could not handle mission cancel discord button\n\t", err)
					return
				}

				err = controller.DatabaseController.CancelMission(mission)
				if err != nil {
					utils.Log.Error("Could not handle mission cancel discord button\n\t", err)
					return
				}

				err = session.InteractionRespond(
					interaction.Interaction,
					CancelMissionResponse(mission),
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
	}
}

func getMissionActionRow(
	controller *DiscordController,
	mission *models.Mission,
) discordgo.ActionsRow {
	actionsRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{},
	}
	for _, idKey := range []string{"set-mission-parameters", "update-mission", "cancel-mission"} {
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
