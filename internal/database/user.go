package database

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/exp/slices"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

func (controller *DatabaseController) CreateUser(
	user *models.User,
) (err error) {
	user.Points = 0
	err = controller.DB.Create(user).Error
	if err != nil {
		return fmt.Errorf("Could not store user on database\n\t%s", err)
	}

	return err
}

func (controller *DatabaseController) GetUser(id uint32) *models.User {
	var user *models.User
	err := controller.DB.First(&user, &id).Error
	if err != nil {
		return nil
	}
	return user
}

func (controller *DatabaseController) GetUserParticipations(user *models.User) []models.Participation {
	var participations []models.Participation
	controller.DB.Model(&user).Association("Participations").Find(&participations)
	return participations
}

func (controller *DatabaseController) GetOrCreateUserFromString(id string) (*models.User, error) {
	userId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Could not get user by ID %s\n\t%s", id, err)
	}
	user := controller.GetUser(uint32(userId))
	if user == nil {
		user = &models.User{
			ID: uint32(userId),
		}
		err = controller.CreateUser(user)
		if err != nil {
			return nil, fmt.Errorf("Could not create user with ID %s\n\t%s", id, err)
		}
	}
	return user, nil
}

func (controller *DatabaseController) GetUsers(
	limit int,
	sort *string,
	ascending bool,
	filters map[string][]any,
) (users []models.User) {
	defaultSort := "updated_at"
	if sort == nil {
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

	request.Find(&users)
	return users
}

func (controller *DatabaseController) UpdateUser(user *models.User, points float64) (err error) {
	user.Points = points
	user.UpdatedAt = time.Now()
	controller.DB.Save(user)
	return err
}
