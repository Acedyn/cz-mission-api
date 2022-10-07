package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/cardboard-citizens/cz-mission-api/internal/database"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

func migrateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Long:  "Migrate the database to match the given schema",
		Run: func(cmd *cobra.Command, args []string) {
			databaseController := database.DatabaseController{
				DbDriver: cmd.Flag("dbdriver").Value.String(),
				DbName:   cmd.Flag("dbname").Value.String(),
			}
			err := databaseController.Initialize()
			if err != nil {
				utils.Log.Error(err)
				return
			}
			databaseController.Migrate()
		},
	}

	command.Flags().StringP("dbname", "n", os.Getenv("DB_NAME"), "Name of the database")
	command.Flags().
		StringP("dbdriver", "d", os.Getenv("DB_DRIVER"), "Driver that will be used to interact with the database (postgres, sqlite...)")
	return command
}

func init() {
	godotenv.Load()
	rootCmd.AddCommand(migrateCommand())
}
