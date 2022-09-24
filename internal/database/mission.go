package database

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
	"golang.org/x/exp/slices"
)

func (controller *DatabaseController) CreateMission(
	mission *models.Mission,
) (err error) {
	mission.Initialized = false
	mission.Canceled = false
	mission.CloseAt = time.Now()
	mission.CreatedAt = time.Now()
	mission.UpdatedAt = time.Now()
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
	controller.DB.Limit(limit).
		Order(sortKey).
		Where("canceled = ?", false).
		Find(&missions)
	return missions
}

func (controller *DatabaseController) UpdateMission(mission *models.Mission, name, shortDescription, longDescription string, reward float64) (err error) {
	mission.Name = name
	mission.ShortDescription = shortDescription
	mission.LongDescription = longDescription
	mission.Reward = reward
	mission.Initialized = true
	mission.UpdatedAt = time.Now()
	err = mission.Validate()
	if err != nil {
		return fmt.Errorf("Could not update mission %s\n\t%s", mission.Format(), err)
	}
	controller.DB.Save(mission)
	return err
}

func (controller *DatabaseController) CancelMission(mission *models.Mission) (err error) {
	mission.Canceled = true
	mission.UpdatedAt = time.Now()
	err = mission.Validate()
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
