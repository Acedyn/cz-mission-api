package main

import (
	"github.com/cardboard-citizens/cz-goodboard-api/internal/controllers"
	"github.com/spf13/cobra"
)

func migrateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Long:  "Migrate the database to match the given schema",
		Run: func(cmd *cobra.Command, args []string) {
			server := controllers.Server{}
			server.InitializeDb(cmd.Flag("dbdriver").Value.String(), cmd.Flag("dbdriver").Value.String())
			server.MigrateDb()
		},
	}

	command.Flags().StringP("dbname", "n", "cz-goodboard", "Name of the database")
	command.Flags().StringP("dbdriver", "d", "sqlite", "Driver that will be used to interact with the database (postgres, sqlite...)")
	return command
}

func init() {
	rootCmd.AddCommand(migrateCommand())
}
