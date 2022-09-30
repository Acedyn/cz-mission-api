package database

import (
	"fmt"
	"time"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

func (controller *DatabaseController) CreateParicipation(
	paricipation *models.Participation,
) (err error) {
	paricipation.CreatedAt = time.Now()
	paricipation.UpdatedAt = time.Now()
	err = controller.DB.Create(paricipation).Error
	if err != nil {
		return fmt.Errorf("Could not store paticipation on database\n\t%s", err)
	}

	return err
}
