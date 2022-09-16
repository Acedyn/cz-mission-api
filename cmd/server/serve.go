package main

import (
	"github.com/cardboard-citizens/cz-goodboard-api/internal/controllers"
	"github.com/spf13/cobra"
)

func serveCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  "Start the database connection, REST API and the discord API",
		Run: func(cmd *cobra.Command, args []string) {
			server := controllers.Server{}
			server.InitializeDb(cmd.Flag("dbdriver").Value.String(), cmd.Flag("dbdriver").Value.String())
		},
	}

	command.Flags().StringP("port", "p", "8080", "HTTP port to listen to")
	command.Flags().StringP("dbname", "n", "cz-goodboard", "Name of the database")
	command.Flags().StringP("dbdriver", "d", "sqlite", "Driver that will be used to interact with the database (postgres, sqlite...)")
	return command
}

func init() {
	rootCmd.AddCommand(serveCommand())
}
