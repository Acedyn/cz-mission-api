package database

import (
	"fmt"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

func (controller *DatabaseController) CreateMission(mission *models.Mission) (err error) {
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

func (controller *DatabaseController) GetMissions() (missions []models.Mission) {
	controller.DB.Find(&missions)
	return missions
}
