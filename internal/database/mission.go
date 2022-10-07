package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/datatypes"

	"github.com/cardboard-citizens/cz-mission-api/internal/missions"
	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

func (controller *DatabaseController) CreateMission(
	mission *models.Mission,
) (err error) {
	mission.Initialized = false
	mission.Canceled = false
	mission.CloseAt = time.Now()
	mission.CreatedAt = time.Now()
	mission.UpdatedAt = time.Now()
	mission.Category = controller.GetMissionClass(mission).Category
	mission.Logo = controller.GetMissionClass(mission).Logo
	mission.Parameters = datatypes.JSON([]byte(`{}`))
	err = controller.DB.Create(mission).Error
	if err != nil {
		return fmt.Errorf("Could not store mission on database\n\t%s", err)
	}

	return err
}

func (controller *DatabaseController) GetMission(id uint32) *models.Mission {
	var mission models.Mission
	controller.DB.First(&mission, id)
	return &mission
}

func (controller *DatabaseController) GetMissionFromString(id string) (*models.Mission, error) {
	missionId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Could not get mission by ID %s : %s", id, err)
	}
	mission := controller.GetMission(uint32(missionId))
	if mission == nil {
		return nil, fmt.Errorf("Could not get mission by ID %s : Mission not found", id)
	}
	return mission, nil
}

func (controller *DatabaseController) GetMissions(
	limit int,
	sort *string,
	ascending bool,
	filters map[string][]any,
) (missions []models.Mission) {
	defaultSort := "updated_at"
	if sort == nil {
		sort = &defaultSort
	}
	missionFieldNames := GetMissionFieldNames()
	if !slices.Contains(missionFieldNames, *sort) {
		utils.Log.Warning(
			"Invalid sort key on mission database fetching (",
			*sort,
			"): Using",
			defaultSort,
		)
		sort = &defaultSort
	}

	sortKey := fmt.Sprintf("%s", *sort)
	if !ascending {
		sortKey = fmt.Sprintf("%s desc", *sort)
	}
	request := controller.DB.
		Order(sortKey).
		Where("canceled = ?", false)

	if limit > 0 {
		request.Limit(limit)
	}

	for filter_key, filter_value := range filters {
		request.Where(filter_key, filter_value...)
	}

	request.Find(&missions)
	return missions
}

func (controller *DatabaseController) ValidateMission(mission *models.Mission, user *models.User) (paricipation *models.Participation, err error) {
	missionClass := controller.GetMissionClass(mission)
	validated, err := missionClass.Validation(mission, user)
	if err != nil {
		return nil, fmt.Errorf("Could not validate mission %s\n\t%s", mission.Format(), err)
	}
	if validated {
		participation := &models.Participation{
			Users:    []*models.User{user},
			Mission:  *mission,
			Progress: 1,
		}
		controller.CreateParicipation(participation)
		return participation, err
	}
	return nil, err
}

func (controller *DatabaseController) UpdateMission(mission *models.Mission, name, shortDescription, longDescription string, reward float64, closeTime time.Time) (err error) {
	mission.Initialized = true
	mission.Name = name
	mission.ShortDescription = shortDescription
	mission.LongDescription = longDescription
	mission.Reward = reward
	mission.CloseAt = closeTime
	mission.UpdatedAt = time.Now()
	err = controller.CheckMission(mission)
	if err != nil {
		return fmt.Errorf("Could not update mission %s\n\t%s", mission.Format(), err)
	}
	controller.DB.Save(mission)
	return err
}

func (controller *DatabaseController) UpdateMissionParameters(mission *models.Mission, parameters map[string]string) error {
	encodedParameters, err := json.Marshal(parameters)
	if err != nil {
		return err
	}

	mission.Initialized = true
	mission.Parameters = datatypes.JSON(encodedParameters)
	mission.UpdatedAt = time.Now()
	err = controller.CheckMission(mission)
	if err != nil {
		return fmt.Errorf("Could not update mission %s\n\t%s", mission.Format(), err)
	}
	controller.DB.Save(mission)
	return nil
}

func (controller *DatabaseController) CancelMission(mission *models.Mission) (err error) {
	mission.Canceled = true
	mission.UpdatedAt = time.Now()
	err = controller.CheckMission(mission)
	if err != nil {
		return fmt.Errorf("Could not update mission %s\n\t%s", mission.Format(), err)
	}
	controller.DB.Save(mission)
	return err
}

func GetMissionFieldNames() []string {
	missionType := reflect.TypeOf((*models.Mission)(nil)).Elem()
	fieldNames := make([]string, 0, missionType.NumField())

	for i := 0; i < missionType.NumField(); i++ {
		field := missionType.Field(i)
		if tag := field.Tag.Get("json"); tag != "" {
			fieldNames = append(fieldNames, tag)
		}
	}

	return fieldNames
}

func (controller *DatabaseController) GetMissionClass(mission *models.Mission) *missions.MissionClass {
	return missions.GetMissionsClasses()[mission.Class]
}

func (controller *DatabaseController) CheckMission(mission *models.Mission) (err error) {
	if mission.Name == "" {
		return errors.New("Invalid mission: No name provided")
	}
	missionClassKeys := missions.GetMissionClassKeys()
	if !slices.Contains(missionClassKeys, mission.Class) {
		return fmt.Errorf(
			"Invalid mission %s: Given class is not part of the available classes (%s)",
			mission.Name,
			mission.Class,
		)
	}

	mission.Name = html.EscapeString(strings.TrimSpace(mission.Name))
	mission.ShortDescription = html.EscapeString(
		strings.TrimSpace(mission.ShortDescription),
	)
	mission.LongDescription = html.EscapeString(
		strings.TrimSpace(mission.LongDescription),
	)
	return err
}
