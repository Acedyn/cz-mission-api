package discord

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

type DiscordModal struct {
	Modal   func(*models.Mission) []discordgo.MessageComponent
	Handler func(*discordgo.Session, *discordgo.InteractionCreate, string)
}

func getModals(controller *DiscordController) map[string]*DiscordModal {
	return map[string]*DiscordModal{
		"set-mission-parameters": {
			Modal: func(mission *models.Mission) []discordgo.MessageComponent {
				messageComponents := []discordgo.MessageComponent{}

				for _, parameter := range controller.DatabaseController.GetMissionClass(mission).Parameters {
					parameterValue := mission.GetParameterValue(parameter)
					messageComponents = append(messageComponents, discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       strings.Title(parameter),
								Style:       discordgo.TextInputShort,
								Placeholder: "",
								CustomID:    parameter,
								Value:       parameterValue,
							},
						},
					})
				}
				return messageComponents
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, id string) {
				utils.Log.Debug("Update mission parameter", id, "modal instruction received")

				mission, err := controller.DatabaseController.GetMissionFromString(id)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord modal\n\t", err)
					session.InteractionRespond(
						interaction.Interaction,
						ErrorResponse(
							fmt.Sprintf(
								"An error occured when processing getting mission with id %s",
								id,
							),
							err),
					)
					return
				}

				missionParameters, err := mission.GetParsedParameters()
				if err != nil {
					utils.Log.Error(err)
				}
				data := interaction.ModalSubmitData()
				for index, parameter := range controller.DatabaseController.GetMissionClass(mission).Parameters {
					parameterValue := data.Components[index].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
					missionParameters[parameter] = parameterValue
				}

				created := !mission.Initialized
				err = controller.DatabaseController.UpdateMissionParameters(mission, missionParameters)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord modal\n\t", err)
					session.InteractionRespond(
						interaction.Interaction,
						MissionResponseError(mission, err),
					)
					return
				}
				if created {
					err = session.InteractionRespond(
						interaction.Interaction,
						CreateMissionResponse(controller, mission),
					)
				} else {
					err = session.InteractionRespond(
						interaction.Interaction,
						UpdateMissionResponse(mission),
					)
				}
				if err != nil {
					utils.Log.Error("An error occured while responding to the interaction", interaction.ID, "\n\t", err)
					return
				}
			},
		},
		"update-mission": {
			Modal: func(mission *models.Mission) []discordgo.MessageComponent {
				return []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Name",
								Style:       discordgo.TextInputShort,
								Placeholder: "Insert new name",
								CustomID:    "name",
								Value:       mission.Name,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Short Description",
								Style:       discordgo.TextInputShort,
								Placeholder: "Insert new short description",
								CustomID:    "short-description",
								Value:       mission.ShortDescription,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Long Description",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Insert new long description",
								CustomID:    "long-description",
								Value:       mission.LongDescription,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Reward",
								Style:       discordgo.TextInputShort,
								Placeholder: "Insert new reward",
								CustomID:    "reward",
								Value:       strconv.Itoa(int(mission.Reward)),
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Close time",
								Style:       discordgo.TextInputShort,
								Placeholder: time.Now().Format(time.RFC822),
								CustomID:    "close-time",
								Value:       mission.CloseAt.Format(time.RFC822),
							},
						},
					},
				}
			},
			Handler: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, id string) {
				utils.Log.Debug("Update mission", id, "modal instruction received")

				mission, err := controller.DatabaseController.GetMissionFromString(id)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord modal\n\t", err)
					session.InteractionRespond(
						interaction.Interaction,
						ErrorResponse(
							fmt.Sprintf(
								"An error occured when processing getting mission with id %s",
								id,
							),
							err),
					)
					return
				}

				data := interaction.ModalSubmitData()
				name := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
				shortDescription := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
				longDescription := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
				rewardRaw := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
				reward, err := strconv.ParseInt(rewardRaw, 10, 32)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord modal\n\t", err)
					session.InteractionRespond(
						interaction.Interaction,
						MissionResponseError(mission, err),
					)
					return
				}
				closeTimeRaw := data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
				closeTime, err := time.Parse(time.RFC822, closeTimeRaw)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord modal\n\t", err)
					session.InteractionRespond(
						interaction.Interaction,
						MissionResponseError(mission, err),
					)
					return
				}
				created := !mission.Initialized
				err = controller.DatabaseController.UpdateMission(mission, name, shortDescription, longDescription, float64(reward), closeTime)
				if err != nil {
					utils.Log.Error("Could not handle mission update discord modal\n\t", err)
					session.InteractionRespond(
						interaction.Interaction,
						MissionResponseError(mission, err),
					)
					return
				}
				if created {
					err = session.InteractionRespond(
						interaction.Interaction,
						CreateMissionResponse(controller, mission),
					)
				} else {
					err = session.InteractionRespond(
						interaction.Interaction,
						UpdateMissionResponse(mission),
					)
				}
				if err != nil {
					utils.Log.Error("An error occured while responding to the interaction", interaction.ID, "\n\t", err)
					return
				}
			},
		},
	}
}

func getMissionModal(controller *DiscordController, mission *models.Mission) []discordgo.MessageComponent {
	idKey := "update-mission"
	modal, ok := getModals(controller)[idKey]
	if !ok {
		utils.Log.Error("Could not get the mission component", idKey, ": Modal does not exists")
		return []discordgo.MessageComponent{}
	}
	return modal.Modal(mission)
}

func getMissionParameterModal(controller *DiscordController, mission *models.Mission) []discordgo.MessageComponent {
	idKey := "set-mission-parameters"
	modal, ok := getModals(controller)[idKey]
	if !ok {
		utils.Log.Error("Could not get the mission component", idKey, ": Modal does not exists")
		return []discordgo.MessageComponent{}
	}
	return modal.Modal(mission)
}

func init() {
	controllerInitializers = append(
		controllerInitializers,
		func(controller *DiscordController) (err error) {
			controller.Modals = getModals(controller)
			return err
		},
	)

	controllerListeners = append(
		controllerListeners,
		func(controller *DiscordController) (err error) {
			controller.Session.AddHandler(
				func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
					if interaction.Type != discordgo.InteractionModalSubmit {
						return
					}
					idSplit := strings.Split(
						interaction.ModalSubmitData().CustomID,
						":",
					)
					modal, ok := controller.Modals[idSplit[0]]
					if ok {
						modal.Handler(session, interaction, idSplit[1])
					}
				},
			)
			utils.Log.Info("Discord components listeners started")
			return err
		},
	)
}
