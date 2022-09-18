package database

import (
	"fmt"
	"reflect"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
	"golang.org/x/exp/slices"
)

func (controller *DatabaseController) CreateMission(
	mission *models.Mission,
) (err error) {
	err = mission.Initialize()
	if err != nil {
		return fmt.Errorf("Could not initialize mission data\n\t%s", err)
	}

	err = controller.DB.Create(&mission).Error
	if err != nil {
		return fmt.Errorf("Could not store mission on database\n\t%s", err)
	}

	return err
}

func (controller *DatabaseController) GetMission(
	id uint32,
) *models.Mission {
	var mission models.Mission
	controller.DB.First(&mission, id)
	return &mission
}

func (controller *DatabaseController) GetMissions(
	limit int,
	sort *string,
) (missions []models.Mission) {
	defaultSort := "updated_at"
	if sort == nil {
		sort = &defaultSort
	}
	if slices.Contains(GetMissionFieldNames(), *sort) {
		utils.Log.Warning(GetMissionFieldNames())
		utils.Log.Warning(
			"Invalid sort key on mission database fetching (",
			*sort,
			"): Using",
			defaultSort,
		)
		sort = &defaultSort
	}
	controller.DB.Limit(limit).
		Order(fmt.Sprintf("%s desc", *sort)).
		Find(&missions)
	return missions
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
