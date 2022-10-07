package database

import (
	"fmt"
    "os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

type DatabaseController struct {
	*gorm.DB
	DbDriver string
	DbName   string
}

func (controller *DatabaseController) Initialize() (err error) {
	if controller.DbDriver == "sqlite" {
		dbFile := fmt.Sprintf("%s.db", controller.DbName)
		controller.DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("Could connect to sqlite database\n\t%s", err)
		}
		utils.Log.Info("Openned sqlite database connection on", dbFile)
	} else if controller.DbDriver == "postgres" {
		controller.DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("Could connect to postgres database\n\t%s", err)
		}
	} else {
		return fmt.Errorf("Could not initialize database : Invalid or unsuported driver (%s)", controller.DbDriver)
	}

	return err
}

func (controller *DatabaseController) Migrate() (err error) {
	err = controller.DB.Debug().AutoMigrate(

		models.User{},
		models.Mission{},
		models.Participation{},
	)
	if err != nil {
		return fmt.Errorf(
			"An error occured during the database migration\n\t%s",
			err,
		)
	}

	utils.Log.Info("Database migration successfull")
	return err
}
