package main

import (
	"github.com/cardboard-citizens/cz-mission-api/internal/server"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
	"github.com/spf13/cobra"
)

func serveCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  "Start the database connection, REST API and the discord API",
		Run: func(cmd *cobra.Command, args []string) {
			server := server.Server{}
			err := server.Initialize(cmd.Flag("dbdriver").Value.String(), cmd.Flag("dbname").Value.String(), cmd.Flag("discordbottoken").Value.String(), cmd.Flag("discordguildid").Value.String())
			if err != nil {
				utils.Log.Error(err)
				return
			}
			err = server.Run()
			if err != nil {
				utils.Log.Error(err)
				return
			}
		},
	}

	command.Flags().StringP("port", "p", "8080", "HTTP port to listen to")
	command.Flags().StringP("dbname", "n", "cz-mission", "Name of the database")
	command.Flags().StringP("dbdriver", "d", "sqlite", "Driver that will be used to interact with the database (postgres, sqlite...)")
	command.Flags().StringP("discordbottoken", "b", "", "Private token to posess the bot")
	command.Flags().StringP("discordguildid", "g", "", "ID of the discord server")
	command.MarkFlagRequired("bottoken")
	command.MarkFlagRequired("guildid")
	return command
}

func init() {
	rootCmd.AddCommand(serveCommand())
}
